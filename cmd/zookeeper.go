// nolint
package cmd

import (
	zk "github.com/scalog/scalog/zookeeper"

	"github.com/spf13/cobra"
)

// zookeeper represents the client command
var zkCmd = &cobra.Command{
	Use:   "zookeeper",
	Short: "Zookeeper",
	Long:  "Zookeeper",
	Run: func(cmd *cobra.Command, args []string) {
		zk.ZKInit()
	},
}

func init() {
	RootCmd.AddCommand(zkCmd)
}
