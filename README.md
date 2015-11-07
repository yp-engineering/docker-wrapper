# Docker Wrapper

Docker Wrapper is a small layer around the command-line docker client.  
It knows about the Docker CLI flags and options and can parse a `docker 
run` command for image name, allowing you to insert command line 
arguments and flags at run time.

We use the docker-wrapper with Mesos.  Mesos uses the docker CLI to run 
Docker Containerizer commands and this wrapper can intercept them and 
modify docker run arguments before passing on to the real docker CLI.

## Building `docker-wrapper`

This is a golang binary, so you need a working golang dev environment.

    go get github.com/yp-engineering/docker-wrapper

For developing new features, you will likely be working from a 
branch in your own fork, so the above line just gets the skeletal 
default repo, replace with your own docker-wrapper repo URL.

    cd ./docker-wrapper
    
    # local binary ./docker-wrapper
    make
    
    # install binary into $GOPATH/bin/docker-wrapper
    make install

Once you have built the binary, you can copy it to other machines as it 
is self-contained and should run on any comparable system (assuming no 
modules introduced shared libs).

## Run Modules

NOTE: adding modules is a compile time operation, there is no runtime 
module discovery nor dynamic loading of new modules taking place.

The whole point of this wrapper was to intercept calls to `docker run` 
and be able to add arguments to the command line flags.  To support 
that, there is a simple interface which a WrapperRunModule should 
implement.

* CLI: docker {DockerFlags} run {DockerRunCommandFlags}


    // Module interface for docker wrapper Run modules
    //   - Priority()  - a way to set order of operation - sorted in ascending order for execution
    //   - HandleRun(...) - handle any run-command context, setting global vars as needed and return new docker run args to inject
    type WrapperRunModule interface {
        Priority() int
        HandleRun(DockerFlags, DockerRunCommandFlags) []string
    }

The primary method is HandleRun which should examine the docker flags 
and run command flags and return any arguments you want to add.  Your 
args are inserted right after the "run" subcommand.

We also provide a default implementation which you can use as the base 
of your own Run Module struct: DefaultRunModule (Priority() ==> .priority ==> 0)

One way to implement is to have a struct which includes the 
DefaultRunModule and then implement your own HandleRun func.

    type MyRunModule struct {
        DefaultRunModule
    }
    
    func (m *MyRunModule) HandleRun(flags DockerFlags, runFlags DockerRunCommandFlags) []string {
        // ... inspect flags and return any options you want to add
    }

Once you have implemented your module, you will need to Register an 
instance of it with the main package's list of run modules using the 
`RegisterRunModule` func:

    // in your file my_run_module.go
    package main
    type MyRunModule struct {
        DefaultRunModule
        // ...
    }
    
    func init() {
        RegisterRunModule(&MyRunModule{priority: 10})
    }


## Package and Installation

There is a target to build a tpkg:

    make package

Version updates are manual at this time (be sure to change the 
`tpkg.yml` version and the `docker_wrapper.go` `VERSION` constant)

## Mesos Slave Config

To setup your mesos-slave to use the new docker-wrapper binary, just 
place the full path to the docker binary into `/etc/mesos-slave/docker` 
config file:

    # /etc/mesos-slave/docker
    /home/ops/bin/docker-wrapper

Or if your OS package does not support automatic /etc/mesos feature 
loading, use the `--docker $GOPATH/bin/docker-wrapper` command line 
arguments to start the Mesos agent.

### Package Installation

We use tpkg.github.io packaging.

The tpkg can be installed under any `TPKG_HOME`:

    sudo TPKG_HOME=/home/ops tpkg -i docker-wrapper

The Mesos Slave configuration (`/etc/mesos-slave/docker`) will then 
point to `/home/ops/bin/docker-wrapper` binary (uses `$TPKG_HOME`).

## Logging and Debug

Docker-wrapper tries to output a logfile to `/var/log/docker-wrapper.log`

If `DOCKER_WRAPPER_DEBUG=1` (or --debug docker flag) then the log is 
output to STDERR.
