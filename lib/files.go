package lib

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/cheggaaa/pb"
)

// DownloadFile Downloads the selected helm release
func DownloadFile(result string, path string, osType string) error {

	// Create the destination file
	destination, err := os.Create(path)
	if err != nil {
		return err
	}

	defer destination.Close()

	// Perform request to get helm releases
	// TODO: Replace URL by variable
	// TODO: Performs if response is != 200
	helm := fmt.Sprintf("https://get.helm.sh/helm-%s-%s.tar.gz", result, osType)
	resp, err := http.Get(helm)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// Set progress bar
	// TODO: Check errors. Read documentation
	var progressBar *pb.ProgressBar
	contentLength, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		progressBar = pb.New(0)
	} else {
		progressBar = pb.New(int(contentLength))
	}
	defer progressBar.Finish()

	// Set progress bar settings
	progressBar.ShowSpeed = true
	progressBar.SetWidth(80)
	progressBar.SetRefreshRate(time.Millisecond * 1000)
	progressBar.SetUnits(pb.U_BYTES)
	progressBar.Start()

	// Create Writer and read data transfered
	writer := io.MultiWriter(destination, progressBar)
	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
