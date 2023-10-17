// nolint
package cmd

import (
	zkclient "github.com/scalog/scalog/zkclient"

	"github.com/spf13/cobra"
)

// zookeeper represents the client command
var zkClientCmd = &cobra.Command{
	Use:   "zkclient",
	Short: "Zookeeper Client",
	Long:  "Zookeeper Client",
	Run: func(cmd *cobra.Command, args []string) {
		zkclient.Start()
	},
}

func init() {
	RootCmd.AddCommand(zkClientCmd)
}
