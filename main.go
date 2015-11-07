// Copyright 2015 YP LLC.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"os"
	"sort"
)

var (
	// parsed out of DockerRunCommandFlags.Args.Image
	dockerFullImageName string // globals to store imagename:tag from {{docker run}} command line (see run_cmd.go)
	dockerImageName     string
	dockerImageTag      string

	// from -e ENV vars
	mesosTaskId   string
	marathonAppId string
)

// if debug is on logging goes to stdout/stderr, else setup logging to append
// to file ...
var logFile *os.File

// ********************

// Module interface for docker wrapper Run modules
//   - Priority()  - a way to set order of operation - sorted in ascending order for execution
//   - HandleRun(...) - handle any run-command context, setting global vars as needed and return new docker run args to inject

type WrapperRunModule interface {
	Priority() int
	HandleRun(DockerFlags, DockerRunCommandFlags) []string
}

// plural for sorting purposes
type WrapperRunModules []WrapperRunModule

// define sort.Interface using Priority() to sort module list
func (mods WrapperRunModules) Len() int      { return len(mods) }
func (mods WrapperRunModules) Swap(i, j int) { mods[i], mods[j] = mods[j], mods[i] }
func (mods WrapperRunModules) Less(i, j int) bool {
	return mods[i].Priority() < mods[j].Priority()
}

// the known list of modules for docker run
var registeredRunModules WrapperRunModules

// modules need to call this to register themselves
func RegisterRunModule(m WrapperRunModule) {
	if m != nil {
		registeredRunModules = append(registeredRunModules, m)
	}
}

// provide a simple pseudo Abstract example impl
type DefaultRunModule struct {
	Name     string
	priority int
}

func (d *DefaultRunModule) HandleRun(flags DockerFlags, runFlags DockerRunCommandFlags) []string {
	return []string{}
}

func (d *DefaultRunModule) Priority() int {
	return d.priority
}

// ********************
// ********************

//***********************************************************************

func main() {
	setupLogging()
	defer teardownLogging()

	// create new string slice without "docker-wrapper" first element, in case we need to add args
	newDockerArgs := os.Args[1:]

	// using a command-line parsing library can help grab IMAGE name
	// reliably but has it's own drawbacks: can be stale if new options are
	// added and users attempt to use those new options
	parseCommandlineArgs(newDockerArgs)

	if isDebugEnabled() {
		log.Printf("DEBUG: DOCKER IMAGE == %q", dockerImageName)
		log.Printf("DEBUG: DOCKER TAG == %q", dockerImageTag)
	}

	// if we have an image and a docker run command, we can add functionality here using modules
	if dockerImageName != "" && simpleIsDockerRunCommand(newDockerArgs) {

		// for each registered Run Module -- run them in Priority order
		sort.Sort(registeredRunModules)
		for _, mod := range registeredRunModules {
			// run the module and collect any new docker run params to inject
			modArgs := mod.HandleRun(dockerFlags, dockerRunFlags)
			if modArgs != nil && len(modArgs) > 0 {
				newDockerArgs = injectRunArgs(newDockerArgs, modArgs)
			}
		}

	}

	// now exec docker for real
	dockerExec(newDockerArgs)
}

//***************************************************************************
//***************************************************************************

// isDebugEnabled checks for --debug or DOCKER_WRAPPER_DEBUG env var.
// used in setupLogging() and elsewhere.
func isDebugEnabled() bool {
	return (os.Getenv("DOCKER_WRAPPER_DEBUG") == "1") || dockerFlags.Debug
}

// attempt to log to a file - otherwise stdout
func setupLogging() {
	if !isDebugEnabled() {
		// TODO: allow override - possibly via ENV DOCKER_WRAPPER_LOG=...
		logFileName := "/var/log/docker-wrapper.log"
		logFile, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			// check if we can write to directory - otherwise just log to stderr?
			if os.IsPermission(err) {
				log.Printf("WARN: Error opening logfile, fallback to STDERR: %v", err)
			} else {
				panic(err)
			}
		} else {
			// moved close to teardownLogging()
			//defer f.Close()
			log.SetOutput(logFile)
		}
	}

	// use date, time and filename for log output
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func teardownLogging() {
	// flush and close the file?
	if logFile != nil {
		logFile.Close()
	}
}

func printHelpText() {
	fmt.Println("Usage: docker-wrapper [OPTIONS] COMMAND [arg...]\n\nA Thin wrapper around docker\n")
}

func printVersionText() {
	fmt.Printf("docker-wrapper version: %v\n", VERSION)
}

////////////////
