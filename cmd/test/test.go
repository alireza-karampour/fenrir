package test

import (
	"github.com/spf13/cobra"
	"github.com/alireza-karampour/fenrir/cmd"
)

// TestCmd represents the test command
var TestCmd = &cobra.Command{
	Use:   "test",
	Short: "builds the test suite using ginkgo and deployes it inside the cluster to run",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	cmd.RootCmd.AddCommand(TestCmd)
}
