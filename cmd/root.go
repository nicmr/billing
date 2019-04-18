package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "altemista-billing",
	Short: "Analyzes Altemista cloud usage, generates invoices & uploads to S3",
	Long: `altemista-billing is an application to
	- calculate cloud usage costs per Altemista customer
	- generate invoices accordingly
	- upload invoices to S3 bucket`,
}

// Execute lets the cobra root Command parse the subcommands and params
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
