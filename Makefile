# cuee (POSIX) Makefile
.POSIX:
.SUFFIXES:

# macros
BIN = cuee
GOCMD = go
GOBUILD = $(GOCMD) build
GOGET = $(GOCMD) get

# by default, build
all: build

# build the binary
build: 
	$(GOBUILD) -o $(BIN) -v 

# clean up
clean:
	$(GOCMD) clean
	rm -f $(BIN)
