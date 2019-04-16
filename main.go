package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Altemista/altemista-billing/pkg/costs"
	"github.com/Altemista/altemista-billing/pkg/s3store"
	flag "github.com/spf13/pflag"
)

const (
	iso8601 = "2006-01-02"
)

func handleCosts(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("target")
	month := r.URL.Query().Get("month")

	// Validate passed month
	month += "-01"
	parsedMonth, err := time.Parse(iso8601, month)
	if err != nil {
		log.Println("Error parsing passed month argument")
		w.Write([]byte("Error parsing passed month argument"))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Select appropriate API
	var costapi = costs.Default()
	switch target {
	case "aws":
		costapi = costs.AWS()
	case "azure":
		costapi = costs.Azure()
	case "on-premise":
		costapi = costs.OnPremise()
	default:
		//stays costs.Default()
	}

	output, err := costapi(parsedMonth)

	if err != nil {
		log.Println("GetCostAndUsageRequest failed", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Upload to S3
	_, err = s3store.Upload(strings.NewReader(output.CsvFileContent), "bills/test_costs.csv")
	if err != nil {
		log.Println("Writing to s3 failed: ", err)
	}

	// Tell web browsers to download the csv representation of the results
	w.Header().Set("Content-Disposition", "attachment; filename=costs.csv")
	w.Write([]byte(output.CsvFileContent))

}

func main() {
	// Parse command line flags
	var serve *bool = flag.Bool("serve", false, "If set, the program will respond to http requests at :8080 instead of just running once for a specific month")
	var month *string = flag.String("month", "", "Specifies the month the program should generate billing data for in iso8601 (YYYY-MM). Ignored if serve is set.")
	var apiflag *string = flag.String("api", "", "Specifies the API to be queried. Possible values are aws, azure, on-premise")
	flag.Parse()

	// Serve on port if serve is set or no flags are passed
	if *serve || (*month == "" && *apiflag == "") {
		if *month != "" {
			log.Println("--month was passed but serve is set. --month will be ignored.")
		}
		if *apiflag != "" {
			log.Println("--target was passed but serve is set. --target will be ignored.")
		}

		// set up server
		port := ":8080"
		http.HandleFunc("/costs", handleCosts)
		log.Printf("Serving on %v ...", port)
		log.Fatal(http.ListenAndServe(port, nil))

	} else {
		if *month != "" {

			// Select appropriate API
			costapi := costs.Default()
			switch *apiflag {
			case "aws":
				costapi = costs.AWS()
			case "azure":
				costapi = costs.Azure()
			case "on-premise":
				costapi = costs.OnPremise()
			default:
				// stays costs.Default()
			}

			// Validate the string and parse into time.Time struct
			*month += "-01"
			targetmonth, err := time.Parse(iso8601, *month)
			if err != nil {
				log.Println("Error parsing passed month argument")
				// 22 signifies invalid argument
				os.Exit(22)
			}

			output, err := costapi(targetmonth)

			if err != nil {
				log.Println("GetCostAndUsageRequest failed", err)
				os.Exit(1)
			}

			// Upload to S3
			_, err = s3store.Upload(strings.NewReader(output.CsvFileContent), "bills/test_costs.csv")
			if err != nil {
				log.Println("Writing to s3 failed: ", err)
			}

			log.Println("results:")
			log.Println(output.CsvFileContent)
		} else {
			// warn if insufficient arguments were passed
			log.Println("Exiting: You need to pass either the --serve flag a --month to be analyzed")
		}

	}
}
