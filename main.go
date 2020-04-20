package main

import (
	"fmt"
	"os"
	"sort"

	lib "helmsw/lib"

	log "github.com/sirupsen/logrus"
)

// Set helm root directory
var helmRoot = os.Getenv("HOME") + "/.helmsw"

// HelmswPath struct defines directory tree
// TODO: Change to dynamic structures using pointers
type HelmswPath struct {
	Version string
	Bin     string
}

func main() {

	// Instance helm directory struct
	helmswPath := &HelmswPath{
		Version: helmRoot + "/versions",
		Bin:     helmRoot + "/bin",
	}

	// Check helmsw local dir
	err := lib.CheckHelmswDir(helmswPath.Version, helmswPath.Bin)
	if err != nil {
		log.Fatalf("Error when trying to create helmsw dir: %v", err)
	}

	// Check helm releases on github
	url := "https://api.github.com/repos/helm/helm/releases"
	githubReleases, err := lib.CheckOnlineReleases(url)
	if err != nil {
		log.Fatal(err)
	}

	output := []string{}

	// Check local helm releases
	localReleases, err := lib.CheckLocalReleases(helmswPath.Version)
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

	// Check if a helm version has been set
	ls := &lib.BashCmd{
		Cmd:      "ls",
		Args:     []string{"helm"},
		ExecPath: helmswPath.Bin,
	}
	_, err = lib.ExecBashCmd(ls)

	if err == nil {
		// Hightlight installed version
		output, err = lib.HighlightSelectedRelease(output, helmswPath.Bin)
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
	bin := fmt.Sprintf("%s/helm-%s", helmswPath.Version, version)
	if _, err := os.Stat(bin); os.IsNotExist(err) {

		// Install a new Helm release
		log.Infof("%s is not installed in your system", bin)
		err := lib.InstallRelease(result, bin, helmswPath.Version)
		if err != nil {
			log.Error(err)
		}

		// Helm release is already installed
	} else {

		// Switch Helm release
		err = lib.SwitchRelease(version, helmswPath.Bin, helmswPath.Version)
		if err != nil {
			log.Error(err)
		}
	}
}
