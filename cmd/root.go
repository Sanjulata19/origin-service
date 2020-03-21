package cmd

import (
	"github.com/nullserve/static-host/static_host"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "static-host",
	Short: "Static Host is a cloud storage web server and static file router",
	Run:   root,
}

func root(cmd *cobra.Command, args []string) {
	static_host.Main()
}
