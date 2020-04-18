package lib

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
)

// DisplayMenu display interactive promptui menu
// TODO: Use structs
func DisplayMenu(output []string) (string, string, string, error) {

	// Display interactive menu
	prompt := promptui.Select{
		Label: "Select Helm version",
		Items: output,
	}

	_, result, err := prompt.Run()
	if err != nil {
		//log.Errorf("Prompt failed %v\n", err)
		return "", "", "", err
	}

	// Parse helm version from label (* Installed)
	binToSlice := []string{}
	binToSlice = strings.Split(result, " ")
	version := binToSlice[0]

	// Send message of selected option
	msg := fmt.Sprintf("You choose %q\n", version)
	return version, result, msg, nil
}
