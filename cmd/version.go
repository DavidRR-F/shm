package cmd

import (
	"fmt"

	"github.com/DavidRR-F/shm/internal/version"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Get shm version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("shm-%s", version.Version)
	},
}
