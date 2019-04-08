package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/Altemista/altemista-billing/pkg/costs"
	"github.com/Altemista/altemista-billing/pkg/s3store"
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

	_, err = s3store.Upload(strings.NewReader(output.Response), "bills/test.log")
	if err != nil {
		log.Println("Writing to s3 failed: ", err)
	}

	//w.Write([]byte(output.Response))
	w.Header().Set("Content-Disposition", "attachment; filename=costs.csv")
	w.Write([]byte(output.CsvFileContent))
}

func main() {
	http.HandleFunc("/costs", handleCosts)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
