// nolint
package cmd

import (
	zk "github.com/scalog/scalog/zookeeper"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	zkCmd.PersistentFlags().IntP("zid", "z", 0, "zookeeper index")
	viper.BindPFlag("zid", zkCmd.PersistentFlags().Lookup("zid"))
}
