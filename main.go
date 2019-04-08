package main

import (
	"log"
	"net/http"

	"github.com/Altemista/altemista-billing/pkg/costs"
)

func handleCosts(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")

	// TODO: sanitize parameters

	var client = costs.DefaultClient()

	if target == "aws" {
		client = costs.NewAWS()
	} else if target == "azure" {
		client = costs.NewAzure()
	} else if target == "on-premise" {
		client = costs.NewOnPremise()
	}

	output, err := client.CostsBetween(start, end)

	if err != nil {
		log.Println("GetCostAndUsageRequest failed", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Disposition", "attachment; filename=costs.csv")
	w.Write([]byte(output.CsvFileContent))
}

func main() {
	http.HandleFunc("/costs", handleCosts)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
