package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

func Info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

func Warning(format string, args ...interface{}) {
	fmt.Printf("\x1b[36;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

func CheckArgs(arg ...string) {
	if len(os.Args) < len(arg)+1 {
		Warning("Usage: %s %s", os.Args[0], strings.Join(arg, " "))
		os.Exit(1)
	}
}

func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

func main() {
	CheckArgs("<url>")
	url := os.Args[1]
	userPath, err := os.UserHomeDir()
	CheckIfError(err)

	reg := regexp.MustCompile(`(http(s)?:\/\/|git@)([0-9a-zA-Z\.]+)(\/|:)(.*)(.git)`)
	if reg == nil {
		fmt.Println("regex error")
	}

	res := reg.FindAllStringSubmatch(url, -1)
	clonePath := filepath.Join(userPath + "\\gitworks\\" + res[0][3] + "\\" + res[0][5])

	Info("git clone %s %s --recursive", url, clonePath)

	cmd := exec.Command("git", "clone", url, clonePath, "--recursive")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	CheckIfError(err)

	Info("\nClone complete!")
	Info("\nOpen in explorer:")
	Info("\n    explorer %s", clonePath)
	Info("\nOpen in VS Code:")
	Info("\n    code %s\n", clonePath)
}
