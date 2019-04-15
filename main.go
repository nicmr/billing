package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Altemista/altemista-billing/pkg/costs"
	"github.com/Altemista/altemista-billing/pkg/s3store"
	flag "github.com/spf13/pflag"
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
	var serve *bool = flag.Bool("serve", false, "If set, the program will respond to http requests at :8080 instead of just running once")
	var month *string = flag.String("month", "", "Specifies the month the program should generate billing date for in iso8601 (YYYY-MM). Ignored if serve is set.")
	var target *string = flag.String("target", "", "Specifies the API to be queried. Possible values are aws, azure, on-premise")
	flag.Parse()

	port := ":8080"

	// Serve on port if serve is set or no flags are passed
	if *serve || (*month == "" && *target == "") {
		if *month != "" {
			log.Println("--month was passed but serve is set. --month will be ignored.")
		}
		if *target != "" {
			log.Println("--target was passed but serve is set. --target will be ignored.")
		}
		http.HandleFunc("/costs", handleCosts)
		log.Printf("Serving on %v ...", port)
		log.Fatal(http.ListenAndServe(port, nil))
	} else {
		if *month != "" {
			costapi := costs.Default()
			if *target == "aws" {
				costapi = costs.AWS()
			} else if *target == "azure" {
				costapi = costs.Azure()
			} else if *target == "on-premise" {
				costapi = costs.OnPremise()
			}

			output, err := costapi(*month)

			if err != nil {
				log.Println("GetCostAndUsageRequest failed", err)
				os.Exit(1)
			}

			_, err = s3store.Upload(strings.NewReader(output.CsvFileContent), "bills/test_costs.csv")
			if err != nil {
				log.Println("Writing to s3 failed: ", err)
			}
		} else {
			log.Println("You need to pass either the --serve flag a --month to be analyzed")
		}

	}
}
