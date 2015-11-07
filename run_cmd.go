// Copyright 2015 YP LLC.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"log"
	"strings"
)

const (
	MesosTaskEnv  = "MESOS_TASK_ID"
	MarathonAppId = "MARATHON_APP_ID"
)

// init will setup the RunCommand as part of the main go-flags option parser
func init() {
	optsParser.AddCommand("run",
		"Run a command in a new container",
		"Usage: docker run [OPTIONS] IMAGE [COMMAND] [ARG...]",
		&dockerRunFlags)
}

// DockerRunCommandFlags defined in docker_flags.go
func (x *DockerRunCommandFlags) Execute(args []string) error {
	if isDebugEnabled() {
		log.Printf("RunCommand Env=%q\n", x.Env)
	}

	// set docker image name and tag using splitter from full argument
	fullImageName := x.Args.Image
	setGlobalImageNameAndTag(fullImageName)

	// look for -e ENV vars for these
	setGlobalMesosTaskId(x.Env)
	setGlobalMarathonAppId(x.Env)

	// no error to return - don't halt exec to docker
	return nil
}

// set the Global vars dockerImageName and dockerImageTag
func setGlobalImageNameAndTag(fullImageName string) {
	var err error
	dockerFullImageName = fullImageName
	dockerImageName, dockerImageTag, err = splitFullImageNameWithTag(fullImageName)
	if err != nil {
		log.Printf("RunCommand split error: %q", err)
		// TODO: instead of error, should we just blank out the image and tag
		//return err
	}
	if isDebugEnabled() {
		log.Printf("RunCommand fullImageName=%q\n", fullImageName)
		log.Printf("RunCommand image=%q, tag=%q", dockerImageName, dockerImageTag)
	}
}

// splitFullImageNameWithTag separates out the docker imagename:tag
func splitFullImageNameWithTag(full string) (string, string, error) {
	// split Image string into dockerImageName:dockerImageTag
	var image string
	var tag string
	var splitError error
	imageParts := strings.Split(full, ":")
	if len(imageParts) >= 2 {
		// last piece is tag, all the rest is imageName (in case multiple :)
		image = strings.Join(imageParts[0:len(imageParts)-1], ":")
		tag = imageParts[len(imageParts)-1]
	} else if len(imageParts) == 1 {
		// no tag
		image = imageParts[0]
	} else {
		splitError = errors.New("docker-wrapper: Unable to split image name into parts")
	}
	return image, tag, splitError
}

// collectEnvValuesLike looks for env vars like string and returns list of values
func collectEnvValuesLike(env []string, like string) []string {
	var splitVal []string
	var results []string = []string{}
	for _, envItem := range env {
		if strings.HasPrefix(envItem, like) {
			splitVal = strings.SplitN(envItem, "=", 2)
			if splitVal != nil && len(splitVal) > 1 {
				results = append(results, splitVal[1])
			}
		}
	}
	return results
}

// singleEnvValueLike uses collect and returns just the first match value
func singleEnvValueLike(env []string, like string) string {
	possibles := collectEnvValuesLike(env, like)
	if possibles != nil && len(possibles) > 0 {
		return possibles[0]
	}
	return ""
}

// setGlobalMesosTaskId looks for MESOS_TASK_ID environment variable option.
// this assumes docker to be called from Mesos
func setGlobalMesosTaskId(env []string) {
	mesosTaskId = singleEnvValueLike(env, MesosTaskEnv)
}

// setGlobalMarathonAppId looks for MARATHON_APP_ID environment variable option.
// this assumes docker to be called from Mesos for Marathon
func setGlobalMarathonAppId(env []string) {
	marathonAppId = singleEnvValueLike(env, MarathonAppId)
}
