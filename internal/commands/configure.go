package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func dumpConfig() error {
	settings := viper.AllSettings()
	cleanSettings := make(map[string]interface{})

	// Filtrar valores vacíos
	for k, v := range settings {
		if str, ok := v.(string); ok {
			if str != "" {
				cleanSettings[k] = v
			}
		} else if v != nil {
			cleanSettings[k] = v
		}
	}

	configPath := viper.ConfigFileUsed()
	if configPath == "" {
		configPath = filepath.Join(os.Getenv("HOME"), ".rt-gcli.yaml")
	}
	// Escribir solo los valores no vacíos
	newViper := viper.New()
	newViper.MergeConfigMap(cleanSettings)
	return newViper.WriteConfigAs(configPath)
}

func newConfigureSetCmd() *cobra.Command {
	var (
		key   string
		value string
	)
	cmd := &cobra.Command{
		Use:          "set",
		Short:        "Set a value",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validar que la clave existe en la configuración
			if !viper.InConfig(key) && !viper.IsSet(key) {
				return fmt.Errorf("key '%s' not found in config", key)
			}
			viper.Set(key, value)
			return dumpConfig()
		},
	}
	cmd.Flags().StringVarP(&key, "key", "k", "", "Key to set")
	cmd.Flags().StringVarP(&value, "value", "v", "", "Value to set")
	cmd.MarkFlagRequired("key")
	cmd.MarkFlagRequired("value")
	return cmd
}

func newConfigureGetCmd() *cobra.Command {
	var (
		key string
	)
	cmd := &cobra.Command{
		Use:          "get",
		Short:        "Get a value",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !viper.InConfig(key) && !viper.IsSet(key) {
				return fmt.Errorf("key '%s' not found in config", key)
			}
			fmt.Println(viper.Get(key))
			return nil
		},
	}
	cmd.Flags().StringVarP(&key, "key", "k", "", "Key to get")
	cmd.MarkFlagRequired("key")
	return cmd
}

func newConfigureDumpCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "dump",
		Short:        "Show current configuration",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Mostrar configuración
			for k, v := range viper.AllSettings() {
				fmt.Printf("%s: %v\n", k, v)
			}
			return nil
		},
	}
	return cmd
}

func newConfigureCmd() *cobra.Command {

	cmd := &cobra.Command{
		Use:          "configure",
		Short:        "Configure the client",
		SilenceUsage: true,
	}
	cmd.AddCommand(newConfigureSetCmd())
	cmd.AddCommand(newConfigureGetCmd())
	cmd.AddCommand(newConfigureDumpCmd())
	return cmd

}
