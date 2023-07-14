# github-actions-runners-exporter
Exports number of GitHub Actions' runners.  Calls GitHub API every X seconds to 
get the value.


## Building
Run `go build -o github-actions-runners-exporter` to compile the binary.

### Building docker image
To build the docker image, use the following command.

    docker build -t github-actions-runners-exporter .


## Running
Check below help message for `run` command:

    Usage:  github-actions-runners-exporter run [FLAGS]
    
    Runs the daemon, requires GITHUB_TOKEN environment variable
    
    Required flags: 
      -o,    --organization     GitHub Organization owner of the runners
      -p,    --port         Port to expose /metrics endpoint on
      -s,    --sleep        Seconds between each request to GitHub API


### Example usage

Set `GITHUB_TOKEN` environment variable with github token.  The token needs to
have read-only access to **self-hosted runners** within the 
**organizational permissions** (when using fine-grained token).

Run the program in the background with the following command:

    ./github-actions-runners-exporter run -o my-github-org-name -s 30 -p 8081

Open the `/metrics` endpoint to get the prometheus metrics that would include 
github actions' runners total count:

    % curl -sSL http://127.0.0.1:8081/metrics          
    # HELP github_actions_runners_total_count Number of running GitHub Actions' runners.
    # TYPE github_actions_runners_total_count gauge
    github_actions_runners_total_count 1
    # HELP github_actions_runners_total_count_online Number of running GitHub Actions' runners with 'online' status.
    # TYPE github_actions_runners_total_count_online gauge
    github_actions_runners_total_count_online 0
    # HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
    # TYPE go_gc_duration_seconds summary
    go_gc_duration_seconds{quantile="0"} 0
    ...
    # TYPE go_threads gauge
    go_threads 5
    # HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
    # TYPE promhttp_metric_handler_requests_in_flight gauge
    promhttp_metric_handler_requests_in_flight 1
    # HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
    # TYPE promhttp_metric_handler_requests_total counter
    promhttp_metric_handler_requests_total{code="200"} 0
    promhttp_metric_handler_requests_total{code="500"} 0
    promhttp_metric_handler_requests_total{code="503"} 0

#### Running in a docker container

Remember to export the `GITHUB_TOKEN` environment variable.

    docker run -e GITHUB_TOKEN github-actions-runners-exporter run -o Cardinal-Cryptography -p 8081 -s 60
