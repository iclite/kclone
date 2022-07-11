package main

import (
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var kcloneTest = flag.Bool("kclone", false, "run the kclone tests")

var defaultURL = "https://github.com/git-fixtures/basic.git"

var args = map[string]string{
	"clone": defaultURL,
}

var targetFolder = []string{}

func TestKclone(t *testing.T) {

	flag.Parse()

	getTargetFolder("git-fixtures/basic")
	defer deleteTargetFolder()

	t.Run("kclone", func(t *testing.T) {
		arguments := append([]string{"run", "."}, args["clone"])
		cmd := exec.Command("go", arguments...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			t.Errorf("error running cmd %q", err)
		}
	})
}

func getTargetFolder(dir string) string {
	userPath, errUserPath := os.UserHomeDir()
	CheckIfError(errUserPath)

	path := userPath
	path = filepath.Join(path, "gitworks", dir)

	targetFolder = append(targetFolder, path)
	return path
}

func deleteTargetFolder() {
	for _, folder := range targetFolder {
		err := os.RemoveAll(folder)
		CheckIfError(err)
	}
}
