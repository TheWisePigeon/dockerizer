package cmd

import (
	"os"
	"github.com/spf13/cobra"
)


var rootCmd = &cobra.Command{
	Use:   "dockerizer",
	Short: "docker commands runner",
	Long: `dockerizer is a docker commands runner that you can deploy on a server and send commands to via HTTP.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
