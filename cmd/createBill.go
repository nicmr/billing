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
	month         string
	api           string
	createBillCmd = &cobra.Command{
		Use:   "createBill",
		Short: "Creates billing documents only for the specified month",
		Long:  `Analyzes Altemista cloud usage and creates billing documents only for the specified month`,
		Run: func(cmd *cobra.Command, args []string) {
			createBill()
		},
	}
)

func init() {

	createBillCmd.Flags().StringVarP(&month, "month", "m", "current", "Specifies the month: current, last, or 'YYYY-MM'")
	createBillCmd.Flags().StringVar(&api, "api", "aws", "Specifies the API to work with: aws, azure or onpremise")
	rootCmd.AddCommand(createBillCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createBillCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createBillCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func createBill() {
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
