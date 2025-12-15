/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra" //nolint:depguard
)

var RootCmd = &cobra.Command{
	Use:   "anti-bruteforce",
	Short: "Anti bruteforce",
	Long:  `Anti bruteforce App`,
}

var configFile string

func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	RootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
