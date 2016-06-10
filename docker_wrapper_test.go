// Copyright 2015 YP LLC.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// tests for wrapper code and docker binary

package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleIsDockerRunCommand(t *testing.T) {
	var testArgs []string
	var result bool

	testArgs = []string{"docker", "run", "theimage/name"}
	result = simpleIsDockerRunCommand(testArgs)
	assert.True(t, result, "docker run should return true")

	testArgs = []string{"docker", "pull", "theimage/name"}
	result = simpleIsDockerRunCommand(testArgs)
	assert.False(t, result, "docker pull is not docker run")

	testArgs = []string{"docker", "runthis", "theimage/name"}
	result = simpleIsDockerRunCommand(testArgs)
	assert.False(t, result, "'docker runthis' should not return true: got %v", result)
}

func TestIsDebugEnabled(t *testing.T) {

	var result bool

	os.Unsetenv("DOCKER_WRAPPER_DEBUG")
	result = isDebugEnabled()
	if result {
		t.Errorf("isDebugEnabled should return false when DOCKER_WRAPPER_DEBUG is not set. got %v", result)
	}

	os.Setenv("DOCKER_WRAPPER_DEBUG", "1")
	result = isDebugEnabled()
	if !result {
		t.Errorf("isDebugEnabled should return true when DOCKER_WRAPPER_DEBUG is set. got %v", result)
	}

	os.Unsetenv("DOCKER_WRAPPER_DEBUG")
}

var exampleRun1Args = []string{"run",
	"-d", "-c", "256", "-m", "33554432",
	"-e", "MARATHON_APP_VERSION=2015-06-16T19:01:46.290Z",
	"-e", "HOST=mesosdev5.np.wc1.yellowpages.com",
	"-e", "PORT_10000=31782",
	"-e", "MESOS_TASK_ID=container-echo-test.237350f2-145a-11e5-a886-56847afe9799",
	"-e", "PORT=31782",
	"-e", "PORTS=31782",
	"-e", "MARATHON_APP_ID=/container-echo-test",
	"-e", "PORT0=31782",
	"-e", "MESOS_SANDBOX=/mnt/mesos/sandbox",
	"-v", "/tmp/mesos/slaves/20150612-194908-1313800458-5050-5553-S0/frameworks/20150529-012325-1313800458-5050-26433-0000/executors/container-echo-test.237350f2-145a-11e5-a886-56847afe9799/runs/ffeaf330-c579-4780-98f7-53877eecea99:/mnt/mesos/sandbox",
	"--net", "bridge",
	"--entrypoint", "/bin/sh",
	"--name", "mesos-ffeaf330-c579-4780-98f7-53877eecea99",
				"centos:centos6.6",
				"-c", "while : ; do uptime; sleep 10 ; done"}
var exampleRun1Image = "centos" // ":centos6.6"

//********************
// An extreme example of what we could be working with:
// - note the image name embedded without a tag
// - note the last arg is part of CMD and has a colon ':'
var exampleRun2Args = []string{"run",
	"-d", "--restart", "always", "--link", "nsqlookupd1:nsqlookupd",
	"-v", "/var/run/docker.sock:/var/run/docker.sock",
	"-v", "/usr/local/bin/docker:/usr/local/bin/docker",
	"-v", "/tmp:/tmp",
	"-e", "DOCKER_HOST=\"unix:///var/run/docker.sock\"",
	"--privileged",
	"-v", "/path/to/script.sh:/path/to/script.sh",
	"--name", "nsqexec",
	"jess/nsqexec", "--", "-d", "-exec=/path/to/script.sh",
	"-topic", "hooks-docker", "-channel", "hook",
	"-lookupd-addr", "nsqlookupd:4161"}
var exampleRun2Image = "jess/nsqexec"

var exampleRun2bArgs = []string{"run",
	"--name", "nsqexec",
	"jess/nsqexec", "-exec=/path/to/script.sh",
	"-topic", "hooks-docker", "-channel", "hook",
	"-lookupd-addr", "nsqlookupd:4161"}

func TestParseCommandlineArgs(t *testing.T) {
	// sample mesos command line - make sure we can pull image name
	parseCommandlineArgs(exampleRun1Args)
	if dockerImageName != exampleRun1Image {
		t.Errorf("parseCommandLineArgs expected %q, got %q", exampleRun1Image, dockerImageName)
	}

	// example complex command line from github.com/docker
	parseCommandlineArgs(exampleRun2Args)
	if dockerImageName != exampleRun2Image {
		t.Errorf("parseCommandLineArgs expected %q, got %q", exampleRun2Image, dockerImageName)
	}

	// example without extra -- (double dash) crutch
	parseCommandlineArgs(exampleRun2bArgs)
	if dockerImageName != exampleRun2Image {
		t.Errorf("parseCommandLineArgs expected %q, got %q", exampleRun2Image, dockerImageName)
	}
}

func TestParseJsonFromString_simple(t *testing.T) {
	// an array with a single map that has two key-values
	tester := "[{\"a\": 1, \"b\": \"ok\"}]"

	out, err := parseJsonFromString(tester)
	if err != nil {
		t.Errorf("Error parsing json: %v", err)
	}

	// expect a single value in a list
	test_list := out.([]interface{})
	if len(test_list) != 1 {
		t.Errorf("Error: expecting a single item in an array, got %d: %v", len(test_list), test_list)
	}

	// map values should match - cast to map
	test_hash := test_list[0].(map[string]interface{})
	assert.Equal(t, 1.0, test_hash["a"].(float64), "Expected float64 for 'a'")
	assert.Equal(t, "ok", test_hash["b"], "Expected string for 'b'")
}

func TestParseJsonFromString_label(t *testing.T) {
	// an array with a single map that has two key-values
	tester := "[{\"File\":\"/var/log/sample-access.log\",\"Topic\":\"sample_access_log\"},{\"File\":\"/var/log/sample-error.log\",\"Topic\":\"sample_error_log\"}]"

	out, err := parseJsonFromString(tester)
	if err != nil {
		t.Errorf("Error parsing json: %v", err)
	}

	// expect two values in a list
	test_list := out.([]interface{})
	if len(test_list) != 2 {
		t.Errorf("Error: expecting two items in an array, got %d: %v", len(test_list), test_list)
	}

	// mapped values should match - cast to map
	accesslog_hash := test_list[0].(map[string]interface{})
	assert.Equal(t, "sample_access_log", accesslog_hash["Topic"], "Expected correct Topic")
	assert.Equal(t, "/var/log/sample-access.log", accesslog_hash["File"], "Expected correct File")

	errorlog_hash := test_list[1].(map[string]interface{})
	assert.Equal(t, "sample_error_log", errorlog_hash["Topic"], "Expected correct Topic")
	assert.Equal(t, "/var/log/sample-error.log", errorlog_hash["File"], "Expected correct File")
}

const (
	DOCKER_INSPECT_JSON = `[
	{
		"Id": "fe60df6c5fab35af33f485ffea14f3cc913e832ffaeeba5293798eb976e55a98",
		"Parent": "d5a08de35f0e74bbc9d4593df0890d21852428b4593d8edd4e77a6f1ae2f6349",
		"Comment": "",
		"Created": "2015-10-09T22:51:21.766811982Z",
		"Container": "7c2131e2996137cb748e56e09b409e9ef3566e8f3ce72de1638d41d0f888f007",
		"ContainerConfig": {
			"Hostname": "3f43f5c20d23",
			"Domainname": "",
			"User": "",
			"AttachStdin": false,
			"AttachStdout": false,
			"AttachStderr": false,
			"PortSpecs": null,
			"ExposedPorts": null,
			"Tty": false,
			"OpenStdin": false,
			"StdinOnce": false,
			"Env": [
				"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
			],
			"Cmd": [
				"/bin/sh",
				"-c",
				"#(nop) LABEL com.example.logdata=[{\"File\":\"/var/log/sample-access.log\",\"Topic\":\"example_access_log\"},{\"File\":\"/var/log/sample-error.log\",\"Topic\":\"example_error_log\"}]"
			],
			"Image": "d5a08de35f0e74bbc9d4593df0890d21852428b4593d8edd4e77a6f1ae2f6349",
			"Volumes": null,
			"VolumeDriver": "",
			"WorkingDir": "",
			"Entrypoint": null,
			"NetworkDisabled": false,
			"MacAddress": "",
			"OnBuild": [],
			"Labels": {
				"com.example.logdata": "[{\"File\":\"/var/log/sample-access.log\",\"Topic\":\"example_access_log\"},{\"File\":\"/var/log/sample-error.log\",\"Topic\":\"example_error_log\"}]"
			}
		},
		"DockerVersion": "1.7.1",
		"Author": "",
		"Config": {
			"Hostname": "3f43f5c20d23",
			"Domainname": "",
			"User": "",
			"AttachStdin": false,
			"AttachStdout": false,
			"AttachStderr": false,
			"PortSpecs": null,
			"ExposedPorts": null,
			"Tty": false,
			"OpenStdin": false,
			"StdinOnce": false,
			"Env": [
				"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
			],
			"Cmd": [
				"/bin/sh",
				"-c",
				"/opt/run.sh"
			],
			"Image": "d5a08de35f0e74bbc9d4593df0890d21852428b4593d8edd4e77a6f1ae2f6349",
			"Volumes": null,
			"VolumeDriver": "",
			"WorkingDir": "",
			"Entrypoint": null,
			"NetworkDisabled": false,
			"MacAddress": "",
			"OnBuild": [],
			"Labels": {
				"com.example.logdata": "[{\"File\":\"/var/log/sample-access.log\",\"Topic\":\"example_access_log\"},{\"File\":\"/var/log/sample-error.log\",\"Topic\":\"example_error_log\"}]",
				"com.another": "Hello world"
			}
		},
		"Architecture": "amd64",
		"Os": "linux",
		"Size": 0,
		"VirtualSize": 215017816
	}
	]`
)
