/*
 * Shift
 *
 * Copyright (c) 2016 QingCloud.com.
 * All rights reserved.
 * Code is licensed with BSD.
 *
 */

package cmd

import (
	"fmt"

	"build"

	"github.com/spf13/cobra"
)

// NewVersionCommand show the version number of radon client.
func NewVersionCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version number of radon client",
		Run:   versionCommandFn,
	}

	return cmd
}

func versionCommandFn(cmd *cobra.Command, args []string) {
	build := build.GetInfo()
	fmt.Printf("datatravelcli:[%+v]\n", build)
}
