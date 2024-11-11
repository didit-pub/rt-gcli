package main

import (
	"os"

	"github.com/didit-pub/rt-gcli/internal/commands"
)

func main() {
	// Ejecutar comandos
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
