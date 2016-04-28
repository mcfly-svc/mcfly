package db

import (
	"fmt"
	"os/exec"
)

func RunHelperScript(sh string) {
	cmd := exec.Command("/bin/sh", sh)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + string(output))
		return
	} else {
		fmt.Println(string(output))
	}
}
