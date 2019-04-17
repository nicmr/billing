package cmd

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/Altemista/altemista-billing/pkg/s3store"
)

// createBillCmd represents the createBill command
var (
	month   string
	api     string
	costCmd = &cobra.Command{
		Use:   "cost",
		Short: "Analyzes costs and creates billing documents for a single month",
		Long:  `Analyzes Altemista cloud costs and creates billing documents for a single month`,
		Run: func(cmd *cobra.Command, args []string) {
			cost()
		},
	}
)

func init() {

	costCmd.Flags().StringVarP(&month, "month", "m", "current", "Specifies the month: current, last, or 'YYYY-MM'")
	costCmd.Flags().StringVar(&api, "api", "aws", "Specifies the API to work with: aws, azure or onpremise")
	rootCmd.AddCommand(costCmd)
}

func cost() {
	// Select appropriate API
	costapi := parseCostAPI(api)

	// Validate the string and parse into time.Time struct
	parsedMonth, err := parseMonth(month)
	if err != nil {
		log.Println("Error parsing passed month argument")
		// 22 signifies invalid argument
		os.Exit(22)
	}

	// Execute the request
	output, err := costapi(parsedMonth)

	if err != nil {
		log.Println("GetCostAndUsageRequest failed", err)
		os.Exit(1)
	}

	// Upload to S3
	filename := "bills/test_costs_" + time.Now().Format("2006-01-02_15:04:05") + ".csv"
	_, err = s3store.Upload(strings.NewReader(output.CsvFileContent), filename)
	if err != nil {
		log.Println("Writing to s3 failed: ", err)
	}

	log.Println("results:")
	log.Println(output.CsvFileContent)
}
