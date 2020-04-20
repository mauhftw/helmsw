package main

import (
	"fmt"
	"os"
	"sort"

	lib "helmsw/lib"

	log "github.com/sirupsen/logrus"
)

// var (
// 	HELM_ROOT     = os.Getenv("HOME") + ".helmsw"
// 	HELM_VERSIONS = HELM_ROOT + "/versions"
// 	HELM_BINS     = HELM_ROOT + "/bin"
// )

// Set helm root directory
var helmRoot = os.Getenv("HOME")

// Helm struct defines directory tree
type Helm struct {
	Version string
	Bin     string
}

func main() {

	// Instance helm directory struct
	helm := &Helm{
		Version: helmRoot + "/versions",
		Bin:     helmRoot + "/bin",
	}

	// Check helmsw local dir
	err := lib.CheckHelmswDir(HELM_VERSIONS, HELM_BINS)
	if err != nil {
		log.Fatal(err)
	}

	// Check helm releases on github
	url := "https://api.github.com/repos/helm/helm/releases"
	githubReleases, err := lib.CheckOnlineReleases(url)
	if err != nil {
		log.Fatal(err)
	}

	output := []string{}

	// Check local helm releases
	localReleases, err := lib.CheckLocalReleases(HELM_VERSIONS)
	if err != nil {
		log.Fatal(err)
	}

	// Check if there are helm releases installed
	if localReleases != "" {
		output = lib.LabelInstalledReleases(localReleases, githubReleases, output)
	}

	// Merge installed and internet helm releases
	sort.Sort(sort.Reverse(sort.StringSlice(githubReleases)))
	output = append(output, githubReleases...)

	// TODO: Make a function to check version set
	ls := &lib.BashCmd{
		Cmd:      "ls",
		Args:     []string{"helm"},
		ExecPath: HELM_BINS,
	}
	_, err = lib.ExecBashCmd(ls)

	if err == nil {
		// Hightlight installed version
		output, err = lib.HighlightSelectedRelease(output, HELM_BINS)
		if err != nil {
			log.Error(err)
		}
	}

	// Display interactive menu
	version, result, msg, err := lib.DisplayMenu(output)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Info(msg)
	}

	// Checks if selected helm release exists locally
	bin := fmt.Sprintf("%s/helm-%s", HELM_VERSIONS, version)
	if _, err := os.Stat(bin); os.IsNotExist(err) {

		// Install a new Helm release
		log.Infof("%s is not installed in your system", bin)
		err := lib.InstallRelease(result, bin, HELM_VERSIONS)
		if err != nil {
			log.Error(err)
		}

		// Helm release is already installed
	} else {

		// Switch Helm release
		err = lib.SwitchRelease(version, HELM_BINS, HELM_VERSIONS)
		if err != nil {
			log.Error(err)
		}
	}
}
