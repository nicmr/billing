package main

import (
	"github.com/Altemista/altemista-billing/pkg/query"
	"log"
	"net/http"
)

func costs(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")

	// TODO: sanitize parameters

	var client query.CostsQuery = query.DefaultClient()
	if target == "aws" {
		client = query.NewAWS()
	} else if target == "azure" {
		client = query.NewAzure()
	} else if target == "on-premise" {
		client = query.NewOnPremise()
	}

	output, err := client.CostsBetween(start, end)

	if err != nil {
		log.Println("GetCostAndUsageRequest failed", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Write([]byte(output.Response))
}

func main() {
	http.HandleFunc("/costs", costs)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
