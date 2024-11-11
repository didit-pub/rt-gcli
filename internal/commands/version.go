package commands

import (
	"fmt"

	"github.com/didit-pub/rt-gcli/pkg/version"
	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "version",
		Short: "Muestra la versión del cliente",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("Version: %s\n", version.Version)
			fmt.Printf("Commit SHA: %s\n", version.CommitSHA)
			fmt.Printf("Build Date: %s\n", version.BuildDate)
			return nil
		},
	}

	return cmd
}
