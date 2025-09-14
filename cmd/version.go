package cmd

import (
	"fmt"

	"github.com/nauticale/yhm/internal/version"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get yhm version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("yhm-", version.Version)
	},
}
