package cmd

import (
	"fmt"
	"goLoadRep/test"

	"github.com/spf13/cobra"
)

var url string
var requests int
var concurrency int

var rootCmd = &cobra.Command{
	Use:   "load-tester",
	Short: "Load testing for web services",
	Run: func(cmd *cobra.Command, args []string) {
		if err := test.RunLoadTest(url, requests, concurrency); err != nil {
			fmt.Println("Error in test:", err)
		}
	},
}

func init() {
	rootCmd.Flags().StringVar(&url, "url", "", "URL of the service to be tested")
	rootCmd.Flags().IntVar(&requests, "requests", 0, "Total number of requests")
	rootCmd.Flags().IntVar(&concurrency, "concurrency", 0, "Number of simultaneous calls")
}

func Execute() error {
	return rootCmd.Execute()
}
