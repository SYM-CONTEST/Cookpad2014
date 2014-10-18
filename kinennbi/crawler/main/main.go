package main

import (
	"github.com/ChimeraCoder/anaconda"
	"github.com/SYM-CONTEST/Cookpad2014/kinennbi/crawler"
)

func main() {
	anaconda.SetConsumerKey("n567b7sH6HrPIWBZyhHM2QiaK")
	anaconda.SetConsumerSecret("ygYGJ7aXEh2UQgLI7pOOWU5cixK6o7pDWYVY4MmvRaerJjqLwT")
	api := anaconda.NewTwitterApi("298482612-AzpnvM6K8TfLw1kbOVnJTlwlQjEWGEGPgXdd7Viz", "HMvID4dg5K3WF6jo3urHYmsVk2MkAmY4V43kLBRE190DH")

	c := crawler.Crawler{Api: api}
	c.AnalyzeAnniversary()
}
