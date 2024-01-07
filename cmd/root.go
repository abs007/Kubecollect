/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/abs007/kcl/cmd/check"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kubecollect",
	Short: "A brief description of your application",
	Long:  ``,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(check.CheckCmd)
}
