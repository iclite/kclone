package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

var test = flag.Bool("t", false, "Test mode.")
var osExit = os.Exit

func Info(format string, args ...interface{}) {
	fmt.Printf("\x1b[34;1m%s\x1b[0m\n", fmt.Sprintf(format, args...))
}

func CheckIfError(err error) {
	if err == nil {
		return
	} else {
		fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
		osExit(1)
	}
}

func GetArgs() (bool, string) {
	flag.Parse()
	var url string

	if len(flag.Args()) < 1 {
		Info("Usage: kclone <git url>")
		osExit(1)
	} else {
		url = flag.Args()[0]
	}

	return *test, url
}

func GetUserPath(test bool) string {
	userPath, err := os.UserHomeDir()
	CheckIfError(err)

	if test {
		userPath = "."
	}

	return userPath
}

func ShowComplete(path string) {
	Info("\nClone complete!")
	Info("\nOpen in explorer:")
	Info("\n    explorer %s", path)
	Info("\nOpen in VS Code:")
	Info("\n    code %s\n", path)
}

func GetClonePath(url string, userPath string) string {
	reg := regexp.MustCompile(`(http(s)?:\/\/|git@)([0-9a-zA-Z\.]+)(\/|:)(.*)(.git)`)

	res := reg.FindAllStringSubmatch(url, -1)
	clonePath := filepath.Join(userPath, "gitworks", res[0][3], res[0][5])

	return clonePath
}

func Clone(url string, path string) {

	Info("git clone %s %s --recursive", url, path)

	cmd := exec.Command("git", "clone", url, path, "--recursive")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	CheckIfError(err)
}

func main() {
	test, url := GetArgs()
	userPath := GetUserPath(test)
	clonePath := GetClonePath(url, userPath)

	Clone(url, clonePath)
	ShowComplete(clonePath)
}
