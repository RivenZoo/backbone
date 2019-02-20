package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

const gitCmd = "git"

func fetchRepo(gitRepo, gitVer string) (tmpPath string, err error) {
	if gitVer == "" {
		gitVer = "master"
	}
	tmpPath, err = ioutil.TempDir(os.TempDir(), "projcreator-")
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			os.RemoveAll(tmpPath) // remove tmp dir if error
		}
	}()

	// git clone -b master/tag gitRepo tmpPath
	cmd := exec.Command(gitCmd, "clone", "-b", gitVer, gitRepo, tmpPath)
	err = cmd.Run()
	if err == nil {
		// remove git dir
		os.RemoveAll(filepath.Join(tmpPath, ".git"))
	}
	return
}
