// Copyright 2015 YP LLC.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.
package main

// collection of all of the Docker command-line flags and options

import (
	"log"
	"os"

	"github.com/jessevdk/go-flags"
)

///////////////////////////////////////////////////////////////////////////
// Docker 1.8.1 - flags to docker itself:
//  Options:
//    --config=~/.docker              Location of client config files
//    --tls=false                     Use TLS; implied by --tlsverify
//    --tlscacert=~/.docker/ca.pem    Trust certs signed only by this CA
//    --tlscert=~/.docker/cert.pem    Path to TLS certificate file
//    --tlskey=~/.docker/key.pem      Path to TLS key file
//    --tlsverify=false               Use TLS and verify the remote
//    -D, --debug=false               Enable debug mode
//    -H, --host=[]                   Daemon socket(s) to connect to
//    -h, --help=false                Print usage
//    -l, --log-level=info            Set the logging level
//    -v, --version=false             Print version information and quit
//
//    -d, --daemon=false                                        Enable daemon mode
//    -h, --help=false                                          Print usage

//***********************************************************************
// Setup for using go-flags library https://github.com/jessevdk/go-flags
//
// See also:
//   run_cmd.go
//   pull_cmd.go

type DockerFlags struct {
	// Example of verbosity with level
	Config        string         `long:"config" description:"Location of client config files"`
	Debug         bool           `short:"D" long:"debug" description:"Enable debug mode"`
	DisableLegReg bool           `long:"disable-legacy-registry" description:"Do not contact legacy registries (deprecated)"`
	Host          []string       `short:"H" long:"host" description:"Daemon socket(s) to connect to"`
	Help          bool           `short:"h" long:"help" description:"Print usage"`
	LogLevel      string         `short:"l" long:"log-level" description:"Set the logging level" default:"info"`
	Tls           bool           `long:"tls" description:"Trust certs signed only by this CA"`
	TlsCaCert     flags.Filename `long:"tlscacert" description:"Trust certs signed only by this CA"`
	TlsCert       flags.Filename `long:"tlscert" description:"Path to TLS certificate file"`
	TlsKey        flags.Filename `long:"tlskey" description:"Path to TLS key file"`
	TlsVerify     bool           `long:"tlsverify" description:"Use TLS and verify the remote"`
	Version       bool           `long:"version" description:"Print version information and quit"`
	// TODO: override help and version to output docker-wrapper specific info
	// NOTE: --daemon no longer an option as of docker 1.8
	//Daemon    bool           `short:"d" long:"daemon" description:"Enable daemon mode"`
}

// docker-wrapper: the docker run command options (the ones we care about).
// For now we just want to reliably get the image-name
//
type DockerRunCommandFlags struct {
	// We use positional args to reliably find image (first non-option arg is
	// image name, second is command).  This requires us to define all the
	// known docker-run options here to properly detect first non-option.
	//
	// Options from: Docker version 1.8.0-dev, build 8c7cd78, experimental
	//
	// NOTE: using strings here instead of uints, since wrapper doesn't care
	// (only want image name) and common CMD args is bash -c "string" and since
	// flag lib won't stop parsing until '--', assuming all flags belong to it
	//
	//CpuShares    uint           `short:"c" long:"cpu-shares" description:"CPU shares (relative weight)"`
	//
	Attach              []string       `short:"a" long:"attach" description:"Attach to STDIN, STDOUT or STDERR"`
	AddHost             []string       `long:"add-host" description:"Add a custom host-to-IP mapping (host:ip)"`
	BlkioWeight         string         `long:"blkio-weight" description:"Block IO (relative weight), between 10 and 1000"`
	BlkioWeightDevice   []string       `long:"blkio-weight-device" description:"Block IO weight (relative device weight)"`
	CpuShares           string         `short:"c" long:"cpu-shares" description:"CPU shares (relative weight)"`
	CapAdd              []string       `long:"cap-add" description:"Add Linux capabilities"`
	CapDrop             []string       `long:"cap-drop" description:"Drop Linux capabilities"`
	CgroupParent        string         `long:"cgroup-parent" description:"Optional parent cgroup for the container"`
	Cidfile             flags.Filename `long:"cidfile" description:"Write the container ID to the file"`
	CpuPeriod           string         `long:"cpu-period" description:"Limit CPU CFS (Completely Fair Scheduler) period"`
	CpuQuota            string         `long:"cpu-quota" description:"Limit CPU CFS (Completely Fair Scheduler) quota"`
	CpusetCpus          string         `long:"cpuset-cpus" description:"CPUs in which to allow execution (0-3, 0,1)"`
	CpusetMems          string         `long:"cpuset-mems" description:"MEMs in which to allow execution (0-3, 0,1)"`
	Detach              bool           `short:"d" long:"detach" description:"Run container in background and print container ID"`
	DetachKeys          string         `long:"detach-keys" description:"Override the key sequence for detaching a container"`
	Device              []string       `long:"device" description:"Add a host device to the container"`
	DeviceReadBps       []string       `long:"device-read-bps" description:"Limit read rate (bytes per second) from a device"`
	DeviceReadIops      []string       `long:"device-read-iops" description:"Limit read rate (IO per second) from a device"`
	DeviceWriteBps      []string       `long:"device-write-bps" description:"Limit write rate (bytes per second) to a device"`
	DisableWriteIops    []string       `long:"device-write-iops" description:"Limit write rate (IO per second) to a device"`
	DisableContentTrust bool           `long:"disable-content-trust" description:"Skip image verification"`
	Dns                 []string       `long:"dns" description:"Set custom DNS servers"`
	DnsOpt              []string       `long:"dns-opt" description:"Set DNS options"`
	DnsSearch           []string       `long:"dns-search" description:"Set custom DNS search domains"`
	Env                 []string       `short:"e" long:"env" description:"Set environment variables"`
	Entrypoint          string         `long:"entrypoint" description:"Overwrite the default ENTRYPOINT of the image"`
	EnvFile             []string       `long:"env-file" description:"Read in a file of environment variables"`
	Expose              []string       `long:"expose" description:"Expose a port or a range of ports"`
	GroupAdd            []string       `long:"group-add" description:"Add additional groups to join"`
	Hostname            string         `short:"h" long:"hostname" description:"Container host name"`
	Help                bool           `long:"help" description:"Print Usage"`
	Interactive         bool           `short:"i" long:"interactive" description:"Keep STDIN open even if not attached"`
	Ip                  string         `long:"ip" description:"Container IPv4 address (e.g. 172.30.100.104)"`
	Ip6                 string         `long:"ip6" description:"Container IPv6 address (e.g. 2001:db8::33)"`
	Ipc                 string         `long:"ipc" description:"IPC namespace to use"`
	Isolation           string         `long:"isolation" description:"Container isolation level"`
	KernelMemory        string         `long:"kernel-memory" description:"Kernel memory limit"`
	Label               []string       `short:"l" long:"label" description:"Set meta data on a container"`
	LabelFile           []string       `long:"label-file" description:"Read in a line delimited file of labels"`
	Link                []string       `long:"link" description:"Add link to another container"`
	LogDriver           string         `long:"log-driver" description:"Logging driver for container"`
	LogOpt              []string       `long:"log-opt" description:"Log driver options"`
	LxcConf             []string       `long:"lxc-conf" description:"Add custom lxc options (deprecated)"`
	Memory              string         `short:"m" long:"memory" description:"Memory limit"`
	MacAddress          string         `long:"mac-address" description:"Container MAC address (e.g. 92:d0:c6:0a:29:33)"`
	MemoryReservation   string         `long:"memory-reservation" description:"Memory soft limit"`
	MemorySwap          string         `long:"memory-swap" description:"Total memory (memory + swap), '-1' to disable swap"`
	MemorySwappiness    string         `long:"memory-swappiness" description:"Tuning container memory swappiness (0 to 100)"`
	Name                string         `long:"name" description:"Assign a name to the container"`
	Net                 string         `long:"net" description:"Set the Network mode for the container" default:"bridge"`
	NetAlias            []string       `long:"net-alias" description:"Add network-scoped alias for the container"`
	OomKillDisable      bool           `long:"oom-kill-disable" description:"Disable OOM Killer"`
	OomScoreAdj         string         `long:"oom-score-adj" description:"Tune host's OOM preferences (-1000 to 1000)"`
	PublishAll          bool           `short:"P" long:"publish-all" description:"Publish all exposed ports to random ports"`
	Publish             []string       `short:"p" long:"publish" description:"Publish a container's port(s) to the host"`
	PublishService      string         `long:"publish-service" description:"Publish this container as a service (deprecated)"`
	Pid                 string         `long:"pid" description:"PID namespace to use"`
	Privileged          bool           `long:"privileged" description:"Give extended privileges to this container"`
	ReadOnly            bool           `long:"read-only" description:"Mount the container's root filesystem as read only"`
	Restart             string         `long:"restart" description:"Restart policy to apply when a container exits" default:"no"`
	Rm                  bool           `long:"rm" description:"Automatically remove the container when it exits"`
	SecurityOpt         []string       `long:"security-opt" description:"Security Options"`
	ShmSize             string         `long:"shm-size" description:"Size of /dev/shm, default value is 64MB"`
	SigProxy            bool           `long:"sig-proxy" description:"Proxy received signals to the process"`
	StopSignal          string         `long:"stop-signal" description:"Signal to stop a container, SIGTERM by default"`
	Tty                 bool           `short:"t" long:"tty" description:"Allocate a pseudo-TTY"`
	Tmpfs               []string       `long:"tmpfs" description:"Mount a tmpfs directory"`
	User                string         `short:"u" long:"user" description:"Username or UID (format: <name|uid>[:<group|gid>])"`
	Ulimit              []string       `long:"ulimit" description:"Ulimit options"`
	Uts                 string         `long:"uts" description:"UTS namespace to use"`
	Volume              []string       `short:"v" long:"volume" description:"Bind mount a volume"`
	VolumeDriver        string         `long:"volume-driver" description:"Optional volume driver for the container"`
	VolumesFrom         []string       `long:"volumes-from" description:"Mount volumes from the specified container(s)"`
	Workdir             string         `short:"w" long:"workdir" description:"Working directory inside the container"`

	Args struct {
		Image   string
		CmdArgs []string
	} `positional-args:"yes" required:"yes"`
}

// primary pre-command option flags
var dockerFlags DockerFlags

// docker run command flags
var dockerRunFlags DockerRunCommandFlags

// global parser so run_cmd can init subcommand.  ignore unknown and pass all options after double dash --
var optsParser = flags.NewParser(&dockerFlags, flags.PassDoubleDash|flags.IgnoreUnknown|flags.PassAfterNonOption)

// parseCommandlineArgs will run the option parser for the passed docker args
func parseCommandlineArgs(args []string) {
	// ParseArgs will call any registered subcommands (e.g. run_cmd, which should set image name)
	otherArgs, err := optsParser.ParseArgs(args)

	if isDebugEnabled() {
		log.Printf("DEBUG: os.Environ() = %q\n", os.Environ())
		log.Printf("DEBUG: os.Args = %q\n", os.Args)
		log.Printf("DEBUG: Docker FLAGS = %q\n", dockerFlags)
		log.Printf("DEBUG: remaining ARGS = %q\n", otherArgs)
	}

	// Version & Help: print our docker-wrapper specific info before calling docker
	if dockerFlags.Help {
		printHelpText()
	}
	if dockerFlags.Version {
		printVersionText()
	}

	// we only care to output parse errors if this was a docker-run command ...
	// otherwise let docker itself error on the args
	if err != nil && simpleIsDockerRunCommand(otherArgs) {
		// don't panic - we still want to exec `docker`
		log.Printf("WARN: %q\n", err)
	}
}
