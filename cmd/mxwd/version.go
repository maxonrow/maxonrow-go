package main

import (
	"fmt"

	"github.com/spf13/cobra"
	ver "github.com/maxonrow/maxonrow-go/version"
)

// VersionCmd ...
var Version = &cobra.Command{
	Use:   "version",
	Short: "Print version info of maxonrow",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Maxonrow:", ver.Version)
	},
}
