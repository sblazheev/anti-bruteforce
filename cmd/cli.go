/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/spf13/cobra" //nolint:depguard
)

var (
	cliAction string
	cliValue  string
	cliHost   string
	cliCmd    = &cobra.Command{
		Use:   "cli",
		Short: "Cli",
		Long:  `Cli`,
		Run: func(_ *cobra.Command, _ []string) {
			urlRequest := ""
			delBlackList := "%s/delete/black-list"
			addBlackList := "%s/add/black-list"
			delWhiteList := "%s/delete/white-list"
			addWhiteList := "%s/add/white-list"
			delBucket := "%s/delete/bucket"

			jsonData := ""
			switch cliAction {
			case "addbl":
				urlRequest = fmt.Sprintf(addBlackList, cliHost)
				jsonData = "{\"net\":\"" + cliValue + "\"}"
			case "delbl":
				urlRequest = fmt.Sprintf(delBlackList, cliHost)
				jsonData = "{\"net\":\"" + cliValue + "\"}"
			case "addwl":
				urlRequest = fmt.Sprintf(addWhiteList, cliHost)
				jsonData = "{\"net\":\"" + cliValue + "\"}"
			case "delwl":
				urlRequest = fmt.Sprintf(delWhiteList, cliHost)
				jsonData = "{\"net\":\"" + cliValue + "\"}"
			case "dell":
				urlRequest = fmt.Sprintf(delBucket, cliHost)
				jsonData = "{\"login\":\"" + cliValue + "\"}"
			case "delip":
				urlRequest = fmt.Sprintf(delBucket, cliHost)
				jsonData = "{\"password\":\"" + cliValue + "\"}"
			}

			req, err := http.NewRequest("POST", urlRequest, bytes.NewBuffer([]byte(jsonData)))
			if err != nil {
				log.Fatalf("Error creating request: %s", err)
			}
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Fatalf("Error sending request: %s", err)
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			// Check the response status code
			if resp.StatusCode != http.StatusOK {
				fmt.Printf("\x1b[31mError - %s\x1b[0m\n", body)
			} else {
				fmt.Println("\x1b[32mSuccessfully\x1b[0m")
			}
		},
	}
)

func init() {
	RootCmd.AddCommand(cliCmd)
	cliCmd.PersistentFlags().StringVarP(&cliAction, "action", "a", "",
		"Action addbl/delbl/addwl/delwl/dell/delip")
	cliCmd.PersistentFlags().StringVarP(&cliValue, "value", "v", "",
		"Action value (example 192.1.1.0/25)")
	cliCmd.PersistentFlags().StringVarP(&cliHost, "host", "r", "http://127.0.0.1",
		"Host request (example http://127.0.0.1)")
}
