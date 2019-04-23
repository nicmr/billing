package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "altemista-billing",
		Short: "Analyzes Altemista cloud usage, generates invoices & uploads to S3",
		Long: `altemista-billing is an application to
	- calculate cloud usage costs per Altemista customer
	- generate invoices accordingly
	- upload invoices to S3 bucket`,
	}
)

// Execute lets the cobra root Command parse the subcommands and params
func Execute() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "supplies the application with the specified config file")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}

}

func initConfig() {

	if cfgFile != "" {
		// Load the specified config
		viper.SetConfigFile(cfgFile)
	} else {
		// look for config file in working directory
		viper.SetConfigType("yaml")
		viper.SetConfigName("config")
		viper.AddConfigPath(".")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if cfgFile != "" {
			log.Println("No config file or unable to read it, using defaults")
		}
	}

}
