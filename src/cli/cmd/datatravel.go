package cmd

import (
	"config"
	"fmt"

	"github.com/spf13/cobra"
)

// NewDatatravelCommand return a datatravel command
func NewDatatravelCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "datatravel <subcommand>",
		Short: "datatravel related commands",
	}
	cmd.AddCommand(NewDatatravelProgressRate())
	return cmd
}

// NewDatatravelProgressRate return the travel progress rate cmd
func NewDatatravelProgressRate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "progressrate",
		Short: "show datatravel progress rate",
		Run:   datatravelProgressRate,
	}

	return cmd
}

func datatravelProgressRate(cmd *cobra.Command, args []string) {
	rateMsg, _ := config.ReadTravelProgress("./datatravel-meta")
	fmt.Printf("%s\n", rateMsg)
}
