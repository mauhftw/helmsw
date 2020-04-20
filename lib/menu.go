package lib

import (
	"strings"

	"github.com/manifoldco/promptui"
)

// DisplayMenu display interactive promptui menu
// TODO: Use structs
func DisplayMenu(output []string) (string, string, error) {

	// Display interactive menu
	prompt := promptui.Select{
		Label:        "Select Helm version",
		HideSelected: true,
		Items:        output,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", "", err
	}

	// Parse helm version from label (* Installed)
	binToSlice := []string{}
	binToSlice = strings.Split(result, " ")
	version := binToSlice[0]

	return version, result, nil
}
