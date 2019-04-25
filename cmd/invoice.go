package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Altemista/altemista-billing/pkg/s3store"
)

var (
	month    string
	provider string
	bucket   string
	// invoiceCmd represents the createBill command
	invoiceCmd = &cobra.Command{
		Use:   "invoice",
		Short: "Analyzes costs and creates billing documents for a single month",
		Long:  `Analyzes Altemista cloud costs and creates billing documents for a single month`,
		Run: func(cmd *cobra.Command, args []string) {
			cost()
		},
	}
)

func init() {

	// month flag & config
	invoiceCmd.Flags().StringVarP(&month, "month", "m", "current", "Specifies the month: current, last, or 'YYYY-MM'")
	if err := viper.BindPFlag("month", invoiceCmd.Flags().Lookup("month")); err != nil {
		log.Fatal("Unable to bind viper to flag:", err)
	}

	// provider flag & config
	invoiceCmd.Flags().StringVar(&provider, "provider", "aws", "Specifies the API to work with: aws, azure or onpremise")
	if err := viper.BindPFlag("provider", invoiceCmd.Flags().Lookup("provider")); err != nil {
		log.Fatal("Unable to bind viper to flag:", err)
	}

	invoiceCmd.Flags().StringVarP(&bucket, "bucket", "b", "", "S3 bucket for output documents (required) ")
	invoiceCmd.MarkFlagRequired("bucket")
	if err := viper.BindPFlag("bucket", invoiceCmd.Flags().Lookup("bucket")); err != nil {
		log.Fatal("Unable to bind viper to flag:", err)
	}

	rootCmd.AddCommand(invoiceCmd)
}

func cost() {
	// Select appropriate API
	costapi := parseCostProvider(viper.GetString("provider"))

	// Validate the string and parse into time.Time struct
	parsedMonth, err := parseMonth(viper.GetString("month"))
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
	filename := "bills/test_costs_"
	_, err = s3store.Upload(strings.NewReader(output.CsvFileContent), viper.GetString("bucket"), filename, ".csv", true)
	if err != nil {
		log.Println("Writing to s3 failed: ", err)
	}

	log.Println("results:")
	log.Println(output.CsvFileContent)
}
