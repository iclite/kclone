package main

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var defaultURL = "https://github.com/git-fixtures/basic.git"

var args = map[string]string{
	"clone": defaultURL,
}

var targetFolder = []string{}

func TestKclone(t *testing.T) {

	getTargetFolder("github.com/git-fixtures/basic")
	if _, err := os.Stat(targetFolder[0]); !os.IsNotExist(err) {
		err := os.RemoveAll(targetFolder[0])
		CheckIfError(err)
	}

	t.Run("kclone", func(t *testing.T) {
		arguments := append([]string{"run", "."}, "-t", args["clone"])
		cmd := exec.Command("go", arguments...)

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			t.Errorf("error running cmd %q", err)
		}
	})
}

func TestInfo(t *testing.T) {
	Info("Info Test")
}

func TestCheckIfError(t *testing.T) {
	CheckIfError(nil)

	oldOsExit := osExit
	var exitNumber int
	defer func() { osExit = oldOsExit }()
	exitCode := func(code int) {
		exitNumber = code
	}
	osExit = exitCode

	CheckIfError(errors.New("Some errors"))

	if exp := 1; exitNumber != exp {
		t.Errorf("Expected exit code: %d, got: %d", exp, exitNumber)
	}

}

func TestGetArgs(t *testing.T) {
	os.Args = []string{""}
	GetArgs()

	os.Args = append(os.Args, "-t", defaultURL)
	test, url := GetArgs()

	if test == false {
		t.Fatalf("GetArgs test is false")
	}

	if url != defaultURL {
		t.Fatalf("GetArgs url is %s", url)
	}
}

func TestGetUserPath(t *testing.T) {
	userPath, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("process ran with err %v, want exit status 1", err)
	}

	test := false
	path := GetUserPath(test)
	if path != userPath {
		t.Fatalf("path != userPath " + path + " " + userPath)
	}

	test = true
	path = GetUserPath(test)
	if path != "." {
		t.Fatalf("TestGetUserPath path != userPath " + path + ".")
	}
}

func TestGetClonePath(t *testing.T) {
	path := GetClonePath(defaultURL, ".")
	clonePath := filepath.Join(".", "gitworks", "github.com", "git-fixtures", "basic")
	if path != clonePath {
		t.Fatalf("TestGetClonePath path != clonePath " + path + " " + clonePath)
	}
}

func TestClone(t *testing.T) {
	path := "gitworks/github.com/git-fixtures/basic"

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		err := os.RemoveAll(path)
		CheckIfError(err)
	}

	Clone(defaultURL, path)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatalf("TestGetClonePath path != clonePath")
	}
}

func TestShowComplete(t *testing.T) {
	ShowComplete("")
}

func TestRunMain(t *testing.T) {
	path := "gitworks/github.com/git-fixtures/basic"
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		err := os.RemoveAll(path)
		CheckIfError(err)
	}

	os.Args = append(os.Args, "-t", defaultURL)

	main()
}

func getTargetFolder(dir string) string {
	path := "."
	path = filepath.Join(path, "gitworks", dir)

	targetFolder = append(targetFolder, path)
	return path
}
