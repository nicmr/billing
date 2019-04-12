package main

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Altemista/altemista-billing/pkg/costs"
	"github.com/Altemista/altemista-billing/pkg/s3store"
)

// splitIfValid checks if month is a iso 8601 conforming string,
// then splits it into the first day of the month and the first day fo the following month
// These two tasks are combined in one function because it is more efficient, validating it is a side effect of splitting it
func splitIfValid(month string) (string, string, error) {
	const iso8601 = "2006-01-02"
	startstr := month + "-01"

	targetMonth, err := time.Parse(iso8601, startstr)
	if err != nil {
		return "", "", err
	}
	nextmonth := targetMonth.AddDate(0, 1, 0)
	y, m, _ := nextmonth.Date()
	end := time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)

	endstr := end.Format(iso8601)
	log.Println(startstr, endstr)
	return startstr, endstr, nil
}

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

func handleMonth(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")
	month := r.URL.Query().Get("month")

	// validate startStr and endStr

	var costapi = costs.Default()

	if target == "aws" {
		costapi = costs.AWS()
	} else if target == "azure" {
		costapi = costs.Azure()
	} else if target == "on-premise" {
		costapi = costs.OnPremise()
	}

	start, end, err := splitIfValid(month)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("400 Bad Request"))
		return
	}

	output, err := costapi(start, end)

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
	http.HandleFunc("/month", handleMonth)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
