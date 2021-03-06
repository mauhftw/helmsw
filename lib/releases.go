package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/cheggaaa/pb"
	color "github.com/logrusorgru/aurora"
)

// Version struct for github releases
type Version struct {
	Tag string `json:"tag_name"`
}

// Set global label variables
var (
	installed = color.Yellow("* Installed").Bold()
	selected  = color.Green("* Selected").Bold()
	tmpPath   = "/tmp"
)

// CheckOnlineReleases check helm's latest releases on github
func CheckOnlineReleases(url string) ([]string, error) {

	// Perform request to get helm releases
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		err := errors.New(string(body))
		return nil, err
	}
	defer resp.Body.Close()

	// Unmarshall JSON response
	var versions []Version
	json.Unmarshal(body, &versions)

	// Convert Struct to array and sort it
	githubReleases := []string{}
	for _, v := range versions {
		githubReleases = append(githubReleases, v.Tag)
	}
	return githubReleases, nil
}

// CheckLocalReleases check local helm releases
func CheckLocalReleases(helmVersionPath string) (string, error) {

	// List installed helm releases
	ls := &BashCmd{
		Cmd:      "ls",
		Args:     []string{"-1"},
		ExecPath: helmVersionPath,
	}
	localReleases, err := ExecBashCmd(ls)
	if err != nil {
		return "", err
	}
	return localReleases, nil
}

// LabelInstalledReleases marks which helm releases are installed
// TODO: Use structs
func LabelInstalledReleases(localReleases string, githubReleases []string, output []string) []string {

	// Parse helm semver releases
	lsToSlice := []string{}
	lsToSlice = strings.Split(localReleases, "\n")
	lsToSlice = lsToSlice[:len(lsToSlice)-1]
	for _, v := range lsToSlice {
		localReleases = strings.Trim(v, "helm-")

		// Erase release from list if it's already installed
		for k, j := range githubReleases {
			if localReleases == j {
				githubReleases = append(githubReleases[:k], githubReleases[k+1:]...)
			}
		}

		localReleases = fmt.Sprintf("%-15s %s", localReleases, installed)
		output = append(output, localReleases)
	}
	// Descending order of installed releases
	sort.Sort(sort.Reverse(sort.StringSlice(output)))
	return output
}

// InstallRelease installs a new helm release locally
// TODO: Check Architecture (x86,x64) for building up the link
// TODO: Check if destination dir exists
// TODO: Change CONSTANTS by argument variables, should be passed to the functions
func InstallRelease(result string, bin string, helmVersionPath string) error {

	// Check OS
	uname := &BashCmd{
		Cmd:  "uname",
		Args: []string{"-s"},
	}
	out, err := ExecBashCmd(uname)
	if err != nil {
		return err
	}

	// TODO: Fix removing the \n on the bashCmd function
	// Perform output clean up
	unameToSlice := []string{}
	unameToSlice = strings.Split(out, "\n")
	osType := strings.ToLower(fmt.Sprintf("%s-amd64", unameToSlice[0]))

	// Download file
	destinationPath := fmt.Sprintf("%s/helm-%s", helmVersionPath, result)
	err = DownloadRelease(result, osType)
	if err != nil {
		return err
	}

	// Untar helm version
	bin = fmt.Sprintf("helm-%s", result)
	v := fmt.Sprintf("%s/helm", osType)
	tar := &BashCmd{
		Cmd:      "tar",
		Args:     []string{"zxvf", bin, v, "--strip-components=1"},
		ExecPath: tmpPath,
	}
	_, err = ExecBashCmd(tar)
	if err != nil {
		return err
	}

	// Rename helm release to specific version
	mv := &BashCmd{
		Cmd:      "mv",
		Args:     []string{"helm", destinationPath},
		ExecPath: tmpPath,
	}
	_, err = ExecBashCmd(mv)
	if err != nil {
		return err
	}

	return nil
}

// DownloadRelease Downloads the selected helm release
func DownloadRelease(result string, osType string) error {

	// TODO: Figure out why we can download in a temporal folder and then move that file
	// Create the download file
	temporalPath := fmt.Sprintf("%s/helm-%s", tmpPath, result)
	tmpFile, err := os.Create(temporalPath)
	if err != nil {
		return err
	}

	defer tmpFile.Close()

	// Perform request to get helm releases
	url := "https://get.helm.sh/"
	helm := fmt.Sprintf("%shelm-%s-%s.tar.gz", url, result, osType)
	resp, err := http.Get(helm)
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		err := errors.New(string(fmt.Sprintf("Trying to download from %s", url)))
		return err
	}

	defer resp.Body.Close()

	// Set progress bar
	var progressBar *pb.ProgressBar
	contentLength, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		progressBar = pb.New(0)
	} else {
		progressBar = pb.New(int(contentLength))
	}

	// Set progress bar settings
	// TOOD: Use writers to log into logrus
	progressBar.ShowSpeed = true
	progressBar.SetWidth(80)
	progressBar.SetRefreshRate(time.Millisecond * 1000)
	progressBar.SetUnits(pb.U_BYTES)
	progressBar.Start()

	// Create Writer and read data transfered
	writer := io.MultiWriter(tmpFile, progressBar)
	_, err = io.Copy(writer, resp.Body)
	if err != nil {
		return err
	}
	progressBar.Finish()

	return nil
}

// SwitchRelease switches between helm releases
func SwitchRelease(binToSlice string, helmBinPath string, helmVersionPath string) error {

	// Delete actual symlink
	rmLn := &BashCmd{
		Cmd:      "find",
		Args:     []string{"-L", ".", "-xtype", "l", "-delete"},
		ExecPath: helmBinPath,
	}
	_, err := ExecBashCmd(rmLn)
	if err != nil {
		return err
	}

	// Create symlink to helm new version
	ln := &BashCmd{
		Cmd:  "ln",
		Args: []string{"-s", fmt.Sprintf("%s/helm-%s", helmVersionPath, binToSlice), fmt.Sprintf("%s/helm", helmBinPath)},
	}
	_, err = ExecBashCmd(ln)
	if err != nil {
		return err
	}
	return nil
}

// HighlightSelectedRelease labels the selected helm version to be used
func HighlightSelectedRelease(output []string, helmBinPath string) ([]string, error) {

	// Print value from a symlink
	readLink := &BashCmd{
		Cmd:      "readlink",
		Args:     []string{"-f", "helm"},
		ExecPath: helmBinPath,
	}
	out, err := ExecBashCmd(readLink)
	if err != nil {
		return nil, err
	}

	// Parse helm version selected
	// TODO: Optimize parsing
	readLinkToSlice := []string{}
	readLinkToSlice = strings.Split(out, "\n")
	readLinkToSlice = strings.SplitN(readLinkToSlice[0], "-", 2)

	// Put label on selected helm release
	currentVersion := fmt.Sprintf("%-15s %s", readLinkToSlice[1], installed)
	for k, v := range output {
		if v == currentVersion {
			selectedVersion := fmt.Sprintf("%-15s %s", readLinkToSlice[1], selected)
			output[k] = selectedVersion
		}
	}

	return output, nil
}
