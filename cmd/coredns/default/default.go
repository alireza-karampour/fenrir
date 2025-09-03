package coredns

import (
	. "github.com/alireza-karampour/fenrir/cmd/coredns"
	"github.com/alireza-karampour/fenrir/pkg/cli/subcmd/coredns"
	"github.com/spf13/cobra"
)

var (
	OptImage *string
)

// DefaultCmd represents the default command
var DefaultCmd = &cobra.Command{
	Use:   "default",
	Short: "sets persistant defaults for the cluster's main coredns",
	RunE: func(cmd *cobra.Command, args []string) error {
		if *OptImage != "" {
			c := coredns.New()
			err := c.ChangeImage(*OptImage)
			if err != nil {
				return err
			}
		}
		return nil
	},
}

func init() {
	CorednsCmd.AddCommand(DefaultCmd)
	OptImage = DefaultCmd.Flags().StringP("image", "i", "", "default image to use for coredns deployment")
}
