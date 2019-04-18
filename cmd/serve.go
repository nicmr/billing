package cmd

import (
	"log"
	"net/http"
	"strings"

	"github.com/spf13/cobra"

	"github.com/Altemista/altemista-billing/pkg/s3store"
)

// serveCmd represents the serve command
var (
	port     string
	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "serve http requests",
		Long:  `Serve http requests. Specify the port with --port`,
		Run: func(cmd *cobra.Command, args []string) {
			http.HandleFunc("/cost", handleCosts)
			log.Printf("Serving on port %v ...", port)
			log.Fatal(http.ListenAndServe(":"+port, nil))
		},
	}
)

func init() {
	serveCmd.Flags().StringVarP(&port, "port", "p", "8000", "specifies port to serve on")
	rootCmd.AddCommand(serveCmd)
}

func handleCosts(w http.ResponseWriter, r *http.Request) {
	target := r.URL.Query().Get("api")
	month := r.URL.Query().Get("month")
	bucket := r.URL.Query().Get("bucket")

	// bucket is required
	if bucket == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Parameter missing: bucket."))
	}

	// try to parse month
	parsedMonth, err := parseMonth(month)
	if err != nil {
		log.Println("Error parsing passed month argument")
		w.Write([]byte("Error parsing passed month argument"))
		w.WriteHeader(http.StatusBadRequest)
	}

	// Select appropriate API
	costapi := parseCostAPI(target)

	// Execute the request
	output, err := costapi(parsedMonth)

	if err != nil {
		log.Println("GetCostAndUsageRequest failed", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Upload to S3
	_, err = s3store.Upload(strings.NewReader(output.CsvFileContent), bucket, "bills/cost", ".csv", true)
	if err != nil {
		log.Println("Writing to s3 failed: ", err)
	}

	// Tell web browsers to download the csv representation of the results
	w.Header().Set("Content-Disposition", "attachment; filename=costs.csv")
	w.Write([]byte(output.CsvFileContent))

}
