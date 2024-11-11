package commands

import (
	"fmt"

	"github.com/didit-pub/rt-gcli/pkg/updater"
	"github.com/didit-pub/rt-gcli/pkg/version"
	"github.com/spf13/cobra"
)

func newUpdateMeCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "updateme",
		Short: "Verifica y aplica actualizaciones",
		Run: func(cmd *cobra.Command, args []string) {
			checkOnly, _ := cmd.Flags().GetBool("check")

			fmt.Printf("Versión actual: %s\n", version.GetVersion())
			fmt.Println("Buscando actualizaciones...")

			release, hasUpdate, err := updater.CheckForUpdates()
			if err != nil {
				fmt.Printf("Error al buscar actualizaciones: %v\n", err)
				return
			}

			if !hasUpdate {
				fmt.Println("Ya tienes la última versión.")
				return
			}

			fmt.Printf("Nueva versión disponible: %s\n", release.TagName)

			if checkOnly {
				return
			}

			fmt.Println("Iniciando actualización...")
			if err := updater.DoSelfUpdate(release); err != nil {
				fmt.Printf("Error durante la actualización: %v\n", err)
				return
			}

			fmt.Println("Actualización completada exitosamente!")
		},
	}

	return cmd
}
