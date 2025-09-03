package coredns

import (
	"github.com/spf13/cobra"
	"github.com/alireza-karampour/fenrir/cmd"
	"github.com/alireza-karampour/fenrir/pkg/cli/subcmd/coredns"
)

var (
	export *bool
)

// CorednsCmd represents the coredns command
var CorednsCmd = &cobra.Command{
	Use:   "coredns",
	Short: "used to customize cluster's main coredns deployment",
	RunE: func(cmd *cobra.Command, args []string) error {
		if export != nil && *export == true {
			c := coredns.New()
			err := c.Export()
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	cmd.RootCmd.AddCommand(CorednsCmd)

	export = CorednsCmd.Flags().BoolP("export", "e", false, "exports the active coredns config to kube-configs dir")
}
