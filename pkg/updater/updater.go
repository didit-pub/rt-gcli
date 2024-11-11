package updater

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"

	"github.com/didit-pub/rt-gcli/pkg/version"
)

type GitHubRelease struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

func CheckForUpdates() (*GitHubRelease, bool, error) {
	// Ajusta esta URL a tu repositorio
	const githubAPI = "https://api.github.com/repos/didit-pub/rt-gcli/releases/latest"

	resp, err := http.Get(githubAPI)
	if err != nil {
		return nil, false, err
	}
	defer resp.Body.Close()

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, false, err
	}

	// Comparar versiones (eliminar 'v' del principio si existe)
	latestVersion := release.TagName
	if len(latestVersion) > 0 && latestVersion[0] == 'v' {
		latestVersion = latestVersion[1:]
	}

	hasUpdate := latestVersion != version.GetVersion()
	return &release, hasUpdate, nil
}

func DoSelfUpdate(release *GitHubRelease) error {
	// Construir el nombre del binario para la plataforma actual
	binaryName := fmt.Sprintf("app-%s-%s", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	// Encontrar el asset correcto
	var downloadURL string
	for _, asset := range release.Assets {
		if asset.Name == binaryName {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	if downloadURL == "" {
		return fmt.Errorf("no se encontró binario para %s-%s", runtime.GOOS, runtime.GOARCH)
	}

	// Descargar el nuevo binario
	resp, err := http.Get(downloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Crear archivo temporal
	tempFile, err := os.CreateTemp("", "app-update")
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())

	// Copiar contenido al archivo temporal
	if _, err := io.Copy(tempFile, resp.Body); err != nil {
		return err
	}
	tempFile.Close()

	// Hacer el archivo temporal ejecutable
	if err := os.Chmod(tempFile.Name(), 0755); err != nil {
		return err
	}

	// Obtener la ruta del ejecutable actual
	executable, err := os.Executable()
	if err != nil {
		return err
	}

	// En Windows, no podemos sobrescribir un archivo en uso
	// así que movemos el actual y luego copiamos el nuevo
	if runtime.GOOS == "windows" {
		oldPath := executable + ".old"
		os.Remove(oldPath) // Eliminar si existe
		os.Rename(executable, oldPath)
	}

	// Reemplazar el binario actual
	if err := os.Rename(tempFile.Name(), executable); err != nil {
		return err
	}

	return nil
}
