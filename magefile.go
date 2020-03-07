// +build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
	"os"
	"os/exec"
)

// Default target to run when none is specified
// If not set, running mage will list available targets
// var Default = Build

// A build step that requires additional params, or platform specific steps for example
func Build() error {
	mg.Deps(InstallDeps)
	fmt.Println("Building...")
	os.Mkdir("out", os.ModePerm)
	cmd := exec.Command("go", "build", "-o", "out", "-v", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// A custom install step if you need your bin someplace other than go/bin
func Install() error {
	mg.Deps(Build)
	fmt.Println("Installing...")
	return os.Rename("out/coffeebot", "/usr/local/bin/coffeebot")
}

// Manage your deps, or running package managers.
func InstallDeps() error {
	fmt.Println("Installing Deps...")
	cmd := exec.Command("go", "get", "-v", "all", "./...")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Clean up after yourself
func Clean() {
	fmt.Println("Cleaning...")
	os.RemoveAll("out")
}

func Run() error {
	mg.Deps(Build)
	cmd := exec.Command("out/coffeebot")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
