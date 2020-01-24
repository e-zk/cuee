# cuee (POSIX) Makefile
.POSIX:
.SUFFIXES:

.include <config.mk>

# by default, build
all: build

# build the binary
build: 
	$(GOBUILD) -o $(BIN) -v 

install:build
	cp -v $(BIN) $(BINPATH)
	chmod +x $(BINPATH)

# clean up
clean:
	$(GOCMD) clean
	rm -f $(BIN)
