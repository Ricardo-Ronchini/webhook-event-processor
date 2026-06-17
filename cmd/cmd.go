package cmd

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "app",
		Short: "Webhook event processor",
	}
)

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default: ./configs/config.yaml)")
	cobra.OnInitialize(initConfig)

	registerCommands(rootCmd)
}

func initConfig() {
	if err := godotenv.Load(); err != nil {
		log.Printf("⚠️  Cannot load .env, error: %v", err)
	}
}

func registerCommands(rootCmd *cobra.Command) {
	rootCmd.AddCommand(EventProcessor())
}
