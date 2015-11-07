// Copyright 2015 YP LLC.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

// collection of utility methods used in docker-wrapper

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

const (
	// docker-wrapper is expected to be installed *outside* of below safe PATH so
	// we can find the real docker binary in dockerDo
	// e.g. /go/bin/docker-wrapper
	SafeDockerSearchPath = "/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin"
)

// parseJsonFromString uses the generic interface{} to read arbitrary data
func parseJsonFromString(jsonstring string) (interface{}, error) {
	var f interface{}
	err := json.Unmarshal([]byte(jsonstring), &f)
	return f, err
}

// ---------------------------------------------------------------------------

// check os.Args for the single word "run"
func simpleIsDockerRunCommand(args []string) bool {
	for i := range args {
		if args[i] == "run" {
			return true
		}
	}
	return false
}

// findBinary uses a restricted PATH to find an executable
func findBinary(name string) (string, error) {
	os.Setenv("PATH", SafeDockerSearchPath)
	return exec.LookPath(name)
}

// sh runs a subprocess and returns the String output
func sh(name string, argv ...string) (string, error) {
	// try to find the binary first
	binary, err := findBinary(name)
	if err != nil {
		panic(err)
	}
	cmd := exec.Command(binary, argv...)
	if isDebugEnabled() {
		log.Printf("DEBUG: sh CMD: %q", cmd)
	}
	out, err := cmd.Output()
	return strings.Trim(string(out), " \n"), err
}

// dockerExec execs out to real `docker` binary with argv arguments, replacing
// current process
func dockerExec(argv []string) {
	// grab the (pre-update) environment for our Exec later
	env := os.Environ()

	dockerBinary, err := findBinary("docker")
	if err != nil {
		panic(err)
	}

	// setup new args slice for real docker binary cli
	// syscall.Exec requires "docker" as first arg also
	newArgs := []string{"docker"}
	newArgs = append(newArgs, argv...)

	if isDebugEnabled() {
		log.Printf("DEBUG: docker binary: %s", dockerBinary)
		log.Printf("DEBUG: docker args: %v", newArgs)
	}

	dockerError := syscall.Exec(dockerBinary, newArgs, env)
	if dockerError != nil {
		panic(dockerError)
	}
}

// injectRunArgs takes the current argument list and another argument list to
// inject after the "docker run" portions of the command arguments
func injectRunArgs(args []string, inject_args []string) []string {
	if inject_args == nil || len(inject_args) == 0 {
		return args
	}
	newArgs := []string{}

	if isDebugEnabled() {
		log.Printf("DEBUG: request to inject args: %q", inject_args)
	}

	// find the run command argument index
	runIndex := -1
	for i := range args {
		if args[i] == "run" {
			runIndex = i
		}
	}
	if runIndex == -1 {
		// something funny - not docker run? just append all args
		newArgs = args
	} else {
		// split the difference and merge in new args after "run"
		newArgs = append(newArgs, args[0:runIndex+1]...)
		newArgs = append(newArgs, inject_args...)
		newArgs = append(newArgs, args[runIndex+1:]...)
	}

	return newArgs
}

////////////////////////////////////////
////////////////////////////////////////

// dockerInspect uses docker to inspect an image or a container
func dockerInspect(name string) (string, error) {
	return sh("docker", "inspect", name)
}
