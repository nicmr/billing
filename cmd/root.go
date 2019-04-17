package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "altemista-billing",
	Short: "Analyzes Altemista cloud usage and generates billing documents",
	Long:  "altemista-billing is an application to calculate AWS costs per customer and generates billing documents accordingly.",
}

// Execute lets the cobra root Command parse the subcommands and params
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
