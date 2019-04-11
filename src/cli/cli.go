/*
 * Shift
 *
 * Copyright (c) 2017 QingCloud.com.
 * All rights reserved.
 *
 */

package main

import (
	"fmt"
	"os"

	"cli/cmd"

	"github.com/spf13/cobra"
)

const (
	cliName        = "datatravelcli"
	cliDescription = "A simple command line client for datatravel"
)

var (
	rootCmd = &cobra.Command{
		Use:        cliName,
		Short:      cliDescription,
		SuggestFor: []string{"datatravelcli"},
	}
)

func init() {
	rootCmd.AddCommand(cmd.NewVersionCommand())
	rootCmd.AddCommand(cmd.NewDatatravelCommand())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
