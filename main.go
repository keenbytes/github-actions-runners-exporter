package main

import (
	"encoding/json"
	"fmt"
	"gopkg.pl/phings/broccli/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	cli := broccli.NewCLI("github-actions-runners-exporter", "GitHub Actions' runners exporter for Prometheus", "infra-team@cardinals")
	cmdRun := cli.AddCmd("run", "Runs the daemon, requires GITHUB_TOKEN environment variable", runHandler)
	cmdRun.AddFlag("organization", "o", "", "GitHub Organization owner of the runners", broccli.TypeString, broccli.IsRequired)
	cmdRun.AddFlag("sleep", "s", "", "Seconds between each request to GitHub API", broccli.TypeInt, broccli.IsRequired)
	cmdRun.AddFlag("port", "p", "", "Port to expose /metrics endpoint on", broccli.TypeInt, broccli.IsRequired)
	_ = cli.AddCmd("version", "Prints version", versionHandler)
	if len(os.Args) == 2 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		os.Args = []string{"App", "version"}
	}
	os.Exit(cli.Run())
}

func versionHandler(c *broccli.CLI) int {
	fmt.Fprintf(os.Stdout, VERSION+"\n")
	return 0
}

func runHandler(cli *broccli.CLI) int {
	if os.Getenv("GITHUB_TOKEN") == "" {
		fmt.Fprint(os.Stderr, "!!! GITHUB_TOKEN environment variable is missing\n")
		return 1
	}

	// TODO: This one doesn't seem to be too useful as GitHub does not clean up runners
	// from the list straight away.  It's done after 14 days or so.
	totalCount := promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "github_actions",
		Subsystem: "runners",
		Name:      "total_count",
		Help:      "Number of running GitHub Actions' runners.",
	})

	totalCountOnline := promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "github_actions",
		Subsystem: "runners",
		Name:      "total_count_online",
		Help:      "Number of running GitHub Actions' runners with 'online' status.",
	})
	var response Response

	// Periodally call GitHub API for number of runners, parse the response into Response struct
	// and set prometheus' Gauge value
	go func() {
		for {
			// flag is already validated by cli lib
			sleepInt, _ := strconv.Atoi(cli.Flag("sleep"))

			// We are never going to have more than 100 runners so getting page 1 only
			req, err := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/orgs/%s/actions/runners?per_page=100&page=1", cli.Flag("organization")), strings.NewReader(""))
			if err != nil {
				fmt.Fprintf(os.Stderr, "!!! Error getting request object to GitHub API\n")
				continue
			}

			req.Header.Set("Accept", "application/vnd.github+json")
			req.Header.Set("Authorization", "Bearer "+os.Getenv("GITHUB_TOKEN"))

			c := &http.Client{}
			resp, err := c.Do(req)
			if err != nil {
				fmt.Fprint(os.Stderr, "!!! Error making request to GitHub API\n")
				time.Sleep(time.Duration(sleepInt*10) * time.Second)
				continue
			}

			defer resp.Body.Close()
			b, _ := ioutil.ReadAll(resp.Body)

			err = json.Unmarshal(b, &response)
			if err != nil {
				fmt.Fprintf(os.Stderr, "!!! Error unmarshalling GitHub API response\n")
			}

			totalCount.Set(float64(response.TotalCount))
			totalCountOnline.Set(float64(response.TotalCountOnline()))

			time.Sleep(time.Duration(sleepInt) * time.Second)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(fmt.Sprintf(":%s", cli.Flag("port")), nil)
	return 0
}
