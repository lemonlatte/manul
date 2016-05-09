package main

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func getVendorSubmodules() (map[string]string, error) {
	output, err := execute(
		exec.Command("git", "submodule", "status"),
	)
	if err != nil {
		return nil, err
	}

	vendors := map[string]string{}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(strings.TrimLeft(line, " -"), " ")
		if len(parts) >= 2 {
			path := parts[1]
			commit := parts[0]
			if strings.HasPrefix(path, "vendor/") {
				path = strings.TrimPrefix(path, "vendor/")
				vendors[path] = commit
			}
		}
	}

	return vendors, nil
}

func addVendorSubmodule(importpath string) error {
	var (
		target   = "vendor/" + importpath
		prefixes = []string{"https://", "git://", "git+ssh://"}

		errs []string
	)

	for _, prefix := range prefixes {
		url := prefix + importpath

		_, err := execute(
			exec.Command("git", "submodule", "add", "-f", url, target),
		)
		if err == nil {
			return nil
		}

		errs = append(errs, err.Error())
	}

	return errors.New(strings.Join(errs, "\n"))
}

func removeVendorSubmodule(importpath string) error {
	path := "vendor/" + importpath

	_, err := execute(exec.Command("git", "submodule", "deinit", "-f", path))
	if err != nil {
		return err
	}

	_, err = execute(exec.Command("git", "rm", "-r", "-f", path))
	if err != nil {
		return err
	}

	return nil
}

func updateVendorSubmodule(importpath string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	cmd := exec.Command("git", "pull", "origin", "master")
	cmd.Dir = filepath.Join(cwd, "vendor", importpath)

	_, err = execute(cmd)

	return err
}