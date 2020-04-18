package lib

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// DownloadFile Downloads the selected helm release
func DownloadFile(result string, path string, osType string) error {

	// Create the file
	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer file.Close()

	// Perform request to get helm releases
	// TODO: Replace URL by variable
	helm := fmt.Sprintf("https://get.helm.sh/helm-%s-%s.tar.gz", result, osType)
	resp, err := http.Get(helm)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
