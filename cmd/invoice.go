package cmd

import (
	"log"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/Altemista/altemista-billing/pkg/billing"
	"github.com/Altemista/altemista-billing/pkg/cmd"
)

var (
	month    string
	provider string
	bucket   string
	margin   float64
	local    bool
	// invoiceCmd represents the createBill command
	invoiceCmd = &cobra.Command{
		Use:   "invoice",
		Short: "Analyzes costs and creates billing documents for a single month",
		Long:  `Analyzes Altemista cloud costs and creates billing documents for a single month`,
		Run: func(cmd *cobra.Command, args []string) {
			invoice()
		},
	}
)

func init() {

	invoiceCmd.Flags().StringVarP(&month, "month", "m", "current", "Specifies the month: current, last, or 'YYYY-MM'")

	invoiceCmd.Flags().StringVar(&provider, "provider", "aws", "Specifies the API to work with: aws, azure or onpremise")

	invoiceCmd.Flags().StringVarP(&bucket, "bucket", "b", "", "S3 bucket for output documents (required) ")
	invoiceCmd.MarkFlagRequired("bucket")

	invoiceCmd.Flags().BoolVarP(&local, "local", "l", false, "Saves files to local disk instead of uploading to S3. Overrides bucket param.")

	invoiceCmd.Flags().Float64Var(&margin, "margin", 0.00, "The relative margin that should be added on top of resource costs as ops compensation")

	rootCmd.AddCommand(invoiceCmd)
}

func invoice() {
	// Get bucket parameter

	if bucket == "" {
		log.Fatal("Required --bucket parameter missing. Please supply it via flag or config file.")
	}

	// Select appropriate API
	parsedProvider := parseCostProvider(provider)

	// Validate the month tring by parsing into time.Time struct
	parsedMonth, err := parseMonth(month)
	if err != nil {
		log.Println("Error parsing passed month argument")
		// 22 signifies invalid argument
		os.Exit(22)
	}

	// Flag and arg parsing complete, pass to application code
	err = cmd.Invoice(parsedProvider, parsedMonth, margin, bucket, local)
	if err != nil {
		log.Fatal(err)
	}
}

const (
	iso8601 = "2006-01-02"
)

func parseCostProvider(s string) (costapi billing.CloudProvider) {
	switch s {
	case "aws":
		costapi = billing.AWS()
	case "azure":
		costapi = billing.Azure()
	case "on-premise":
		costapi = billing.OnPremise()
	default:
		costapi = billing.Default()
	}
	return
}

func parseMonth(s string) (time.Time, error) {
	var parsedMonth time.Time
	switch s {
	case "current":
		parsedMonth = time.Now()
	case "last":
		y, m, _ := time.Now().Date()
		parsedMonth = time.Date(y, m, 1, 0, 0, 0, 0, time.UTC).AddDate(0, -1, 0)
	default:
		// try to parse as iso
		s += "-01"
		var err error
		parsedMonth, err = time.Parse(iso8601, s)
		if err != nil {
			return time.Time{}, err
		}
	}
	return parsedMonth, nil
}
