package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version number of sps",
	Long:  "All software has a version. This is my time",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sps v0.1")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
