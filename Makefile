# building the docker wrapper golang binary

.PHONY: all build install clean uninstall test

TMPDIR?=/tmp
INSTALL?=install

BINARY=docker-wrapper
PKG_SRC=main.go version.go util.go docker_flags.go run_cmd.go example_run_module.go
TEST_PKG_SRC=docker_wrapper_test.go

PACKAGE_DIR=$(TMPDIR)/docker-wrapper.tpkg.tmp
PACKAGE_BIN_DIR=$(PACKAGE_DIR)/reloc/bin
PACKAGE_ETC_DIR=$(PACKAGE_DIR)/reloc/etc
PACKAGE_SLAVE_CONF_DIR=$(PACKAGE_ETC_DIR)/mesos-slave
PACKAGE_LOG_CONFIG_DIR=$(PACKAGE_ETC_DIR)/logrotate.d
PACKAGE_LOG_CONFIG=logrotate.d/docker-wrapper_logrotate

all: build

# set VERSION from version.go, eval into Makefile for inclusion into tpkg.yml
version: version.go
	$(eval VERSION := $(shell grep "VERSION" version.go | cut -f2 -d'"'))

build: $(BINARY)

# this just builds local 'docker-wrapper' binary
$(BINARY): $(PKG_SRC)
	go build -v -x .

clean:
	go clean

uninstall:
	@$(RM) -iv `which docker-wrapper`

test: $(PKG_SRC) $(TEST_PKG_SRC)
	go test .

# this will install binary in your GOPATH
install: build test
	go install .

# NOTE: can only build tpkg for now
package: version build test
	$(RM) -r $(PACKAGE_DIR)
	mkdir -p $(PACKAGE_BIN_DIR) $(PACKAGE_SLAVE_CONF_DIR) $(PACKAGE_LOG_CONFIG_DIR)
	$(INSTALL) $(BINARY) $(PACKAGE_BIN_DIR)/.
	$(INSTALL) postinstall postremove $(PACKAGE_DIR)/.
	$(INSTALL) -m 0644  tpkg.yml $(PACKAGE_DIR)/.
	$(INSTALL) -m 0644 $(PACKAGE_LOG_CONFIG) $(PACKAGE_LOG_CONFIG_DIR)/.
	$(INSTALL) -m 0644 etc_mesos-slave_docker $(PACKAGE_SLAVE_CONF_DIR)/docker
	sed -i "s/version:.*/version: $(VERSION)/" $(PACKAGE_DIR)/tpkg.yml
	tpkg --make $(PACKAGE_DIR) --out $(CURDIR)
	$(RM) -r $(PACKAGE_DIR)
