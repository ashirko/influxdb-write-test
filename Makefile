GOCMD=go
GOGET=$(GOCMD) get
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOINSTALL=$(GOCMD) install
PWD=$(shell pwd)
PREFIX=$(word 2, $(subst github.com, github.com, $(PWD)))/cmd/
SCRIPTNAMES=$(addprefix $(PREFIX), $(shell ls cmd))

get:
	$(GOGET) ./...

build:
	$(GOBUILD) ./...

test:
	$(GOTEST) ./...

install:
	$(GOINSTALL) $(SCRIPTNAMES)
