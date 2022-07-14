package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

func Info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

func main() {
	test := flag.Bool("t", false, "Test mode.")
	flag.Parse()
	fmt.Println("-t:", *test)

	if len(flag.Args()) < 1 {
		fmt.Println("Usage: kclone <git url>")
	}

	url := flag.Args()[0]
	userPath, err := os.UserHomeDir()
	CheckIfError(err)

	if *test {
		userPath = "."
	}

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
