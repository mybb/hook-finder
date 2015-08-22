# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build -o
GOCLEAN=$(GOCMD) clean
GOINSTALL=$(GOCMD) install

TOPLEVEL_PKG := github.com/mybb/hook-finder/src
BIN_NAME := bin/hook-finder

# TODO: build multiple platforms (Mac, Linux, Windows x86 & x64)
build:
	$(GOBUILD) $(BIN_NAME) $(TOPLEVEL_PKG)/$*
clean:
	$(GOCLEAN) $(TOPLEVEL_PKG)/$*
install:
	$(GOINSTALL) $(TOPLEVEL_PKG)/$*