package cmd

import (
	"free-hls.go/server"
	"github.com/spf13/cobra"
)

var (
	addr string
)

var (
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "运行服务端",
		Run: func(cmd *cobra.Command, args []string) {
			server.Start(addr)
		},
	}
)

func init() {
	serverCmd.Flags().StringVarP(&addr, "addr", "a", ":8787", "Addr")
	rootCmd.AddCommand(serverCmd)
}
