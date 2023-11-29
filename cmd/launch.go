package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var lauchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Start the dockerizer server",
	Run: func(cmd *cobra.Command, args []string) {
		port := cmd.Flag("port").Value.String()
		println(port)
		logfile := cmd.Flag("logfile").Value.String()
		println(logfile)
	},
}

func init() {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
    fmt.Println("Something went wrong")
		panic(err)
	}
	lauchCmd.Flags().IntP("port", "p", 3031, "The port on which dockerizer will listen for incoming HTTP requests. Must be higher than 1024")
	lauchCmd.Flags().StringP("logfile", "l", fmt.Sprintf("%s/.dockerizer.log", userHomeDir), "The log file for dockerizer where you can checkout errors and all executed commands")
	rootCmd.AddCommand(lauchCmd)
}
