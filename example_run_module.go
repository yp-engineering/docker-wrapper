package main

// Example Docker Wrapper Run Module
//
// looks at some environment variables and flags and adds some options
//

import (
	"fmt"
	"log"
)

type SampleRunModule struct {
	DefaultRunModule
}

// HandleRun implements the WrapperRunModule interface
func (m *SampleRunModule) HandleRun(flags DockerFlags, runFlags DockerRunCommandFlags) []string {
	log.Println("INFO: SampleRunModule.HandleRun(...)")

	// look for a few 'standard' vars and craft our own
	// (run_cmd.go already looks for MESOS_TASK_ID and pals)

	ports := singleEnvValueLike(runFlags.Env, "PORTS")

	// we are combining a few pieces of data into a new env var flag
	newflags := []string{"-e", fmt.Sprintf("SAMPLE_RUN_MODULE=%s-%s", mesosTaskId, ports)}

	return newflags
}

// init calls RegisterRunModule
func init() {
	RegisterRunModule(&SampleRunModule{})
}
