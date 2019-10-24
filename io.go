package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func collectGitRepositories(gitRoots []string) (gits []string) {
	for _, root := range gitRoots {
		if root == "." {
			continue
		}
		if strings.TrimSpace(root) == "" {
			continue
		}
		listing, err := ioutil.ReadDir(root)
		if err != nil {
			log.Println("Couldn't resolve path (skipping):", err)
			continue
		}
		for _, item := range listing {
			path := filepath.Join(root, item.Name())
			if isGitRepository(path, item) {
				gits = append(gits, path)
			}
		}
	}

	return gits
}
func filterGitRepositories(paths []string) (gits []string) {
	for _, path := range paths {
		stat, err := os.Stat(path)
		if err != nil {
			log.Println("Couldn't resolve path (skipping):", err)
			continue
		}
		if isGitRepository(path, stat) {
			gits = append(gits, path)
		}
	}
	return gits
}
func isGitRepository(path string, item os.FileInfo) bool {
	if !item.IsDir() {
		return false
	}

	_, err := os.Stat(filepath.Join(path, ".git"))
	if os.IsNotExist(err) {
		return false
	}

	return true
}

func execute(dir, command string) (string, error) {
	args := strings.Fields(command)
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

func prompt(message string) string {
	log.Println(message)
	s := bufio.NewScanner(os.Stdin)
	s.Scan()
	return s.Text()
}
