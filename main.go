package main

import (
	"errors"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/Altemista/altemista-billing/pkg/costs"
	"github.com/Altemista/altemista-billing/pkg/s3store"
)

func handleCosts(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")
	startStr := r.URL.Query().Get("start")
	endStr := r.URL.Query().Get("end")

	// validate startStr and endStr
	const iso8601 = "2006-01-02"
	start, err := time.Parse(iso8601, startStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 Bad Request"))
		return
	}
	end, err := time.Parse(iso8601, endStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 Bad Request"))
		return
	}
	if !end.After(start) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 Bad Request"))
		return
	}

	var costapi = costs.Default()

	if target == "aws" {
		costapi = costs.AWS()
	} else if target == "azure" {
		costapi = costs.Azure()
	} else if target == "on-premise" {
		costapi = costs.OnPremise()
	}

	output, err := costapi(startStr, endStr)

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
