package main

type Response struct {
	TotalCount int      `json:"total_count"`
	Runners    []Runner `json:"runners"`
}

type Runner struct {
	Id     int     `json:"id"`
	Name   string  `json:"name"`
	OS     string  `json:"os"`
	Status string  `json:"status"`
	Busy   bool    `json:"busy"`
	Labels []Label `json:"labels"`
}

type Label struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"read-only"`
}

func (r Response) TotalCountOnline() int {
	a := 0
	for _, r := range r.Runners {
		if r.Status == "online" {
			a++
		}
	}
	return a
}
