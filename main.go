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
	month := r.URL.Query().Get("month")

	// TODO: validate month

	var costapi = costs.Default()

	if target == "aws" {
		costapi = costs.AWS()
	} else if target == "azure" {
		costapi = costs.Azure()
	} else if target == "on-premise" {
		costapi = costs.OnPremise()
	}

	output, err := costapi(month)

	if err != nil {
		log.Println("GetCostAndUsageRequest failed", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = s3store.Upload(strings.NewReader(output.CsvFileContent), "bills/test_costs.csv")
	if err != nil {
		log.Println("Writing to s3 failed: ", err)
	}

	// Tell web browsers to "download" the response as "costs.csv".
	w.Header().Set("Content-Disposition", "attachment; filename=costs.csv")
	w.Write([]byte(output.CsvFileContent))

}

func main() {
	http.HandleFunc("/costs", handleCosts)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
