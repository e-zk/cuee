# binary
BIN = cuee

# install location
PREFIX=/usr/local
BINPATH=$(PREFIX)/bin/$(BIN)

# go
GOCMD = go
GOBUILD = $(GOCMD) build
GOGET = $(GOCMD) get
