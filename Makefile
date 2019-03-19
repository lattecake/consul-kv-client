GOCMD = /usr/local/go/bin/go
GOTEST = $(GOCMD) test

all: test

test:
	$(GOTEST) -v ./
