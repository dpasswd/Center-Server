package version

import (
	"dh-passwd/global"
	"fmt"
	"github.com/spf13/cobra"
)

var (
	configYml string
	port      string
	mode      string
	StartCmd  = &cobra.Command{
		Use:     "version",
		Short:   "Get version info",
		Example: "dh-passwd version",
		PreRun: func(cmd *cobra.Command, args []string) {

		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func run() error {
	fmt.Println(global.Version)
	return nil
}
