package commands

import (
	"fmt"
	"time"

	"github.com/didit-pub/rt-gcli/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	cfg     config.Config
	silent  bool

	rootCmd = &cobra.Command{
		Use:   "rtg",
		Short: "Cliente CLI para Request Tracker",
		Long: `Un cliente de linea de comandos completo para interactuar 
               con el sistema Request Tracker (RT).`,
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.PersistentFlags().StringVar(&cfg.URL, "url", "", "URL del servidor RT")
	rootCmd.PersistentFlags().StringVar(&cfg.APIURL, "apiurl", "", "URL de la API del servidor RT")
	rootCmd.PersistentFlags().StringVar(&cfg.Username, "username", "", "Nombre de usuario")
	rootCmd.PersistentFlags().StringVar(&cfg.Password, "password", "", "Contraseña")
	rootCmd.PersistentFlags().StringVar(&cfg.Token, "token", "", "Token")
	rootCmd.PersistentFlags().DurationVar(&cfg.Timeout, "timeout", 10*time.Second, "Tiempo de espera para las solicitudes")
	rootCmd.PersistentFlags().BoolVar(&cfg.Debug, "debug", false, "Habilitar modo de depuración")
	rootCmd.PersistentFlags().StringVar(&cfg.Me, "me", "", "Usuario")
	rootCmd.PersistentFlags().BoolVar(&silent, "silent", false, "Silenciar la salida")
	viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))
	viper.BindPFlag("apiurl", rootCmd.PersistentFlags().Lookup("apiurl"))
	viper.BindPFlag("username", rootCmd.PersistentFlags().Lookup("username"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("token", rootCmd.PersistentFlags().Lookup("token"))
	viper.BindPFlag("timeout", rootCmd.PersistentFlags().Lookup("timeout"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	viper.SetDefault("timeout", 30*time.Second)
	viper.BindPFlag("me", rootCmd.PersistentFlags().Lookup("me"))
	// Añadir comandos
	rootCmd.AddCommand(
		newCreateCmd(),
		newGetCmd(),
		newUpdateCmd(),
		newCommentCmd(),
		newVersionCmd(),
		newUpdateMeCmd(),
		newConfigureCmd(),
	)
}

func initConfig() {
	viper.SetConfigType("yaml")
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Configurar los posibles nombres y ubicaciones del archivo de configuracion
		viper.SetConfigName(".rt-gcli")
		// Añadir rutas de busqueda en orden de preferencia
		viper.AddConfigPath(".")
		viper.AddConfigPath("$HOME")
	}
	viper.SetEnvPrefix("RT_GCLI")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		if cfg.Debug {
			fmt.Println("Usando archivo de configuracion:", viper.ConfigFileUsed())
		}
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		fmt.Println("error unmarshaling config: %w", err)
	}
}
