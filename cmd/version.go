/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra" //nolint:depguard
)

var (
	release   = "UNKNOWN"
	buildDate = "UNKNOWN"
	gitHash   = "UNKNOWN"
)

func printVersion() {
	if err := json.NewEncoder(os.Stdout).Encode(struct {
		Release   string
		BuildDate string
		GitHash   string
	}{
		Release:   release,
		BuildDate: buildDate,
		GitHash:   gitHash,
	}); err != nil {
		fmt.Printf("error while decode version info: %v\n", err)
	}
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Long:  `Print version`,
	Run: func(_ *cobra.Command, _ []string) {
		printVersion()
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
