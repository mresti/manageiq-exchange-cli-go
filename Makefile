# Project specific variables
PROJECT=manageiq-exchange
# --- the rest of the file should not need to be configured ---

# GO env
GOPATH=$(shell pwd)
GO=go
GOCMD=GOPATH=$(GOPATH) $(GO)
RELEASE_PATH := ${GOPATH}/release
DISTPATH := bin/$(PROJECT)

# Build versioning
COMMIT = $(shell git log -1 --format="%h" 2>/dev/null || echo "0")
VERSION=$(shell git describe --tags --always)
BUILD_DATE = $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
FLAGS = -ldflags "\
  -X $(PROJECT)/constants.COMMIT=$(COMMIT) \
  -X $(PROJECT)/constants.VERSION=$(VERSION) \
  -X $(PROJECT)/constants.BUILD_DATE=$(BUILD_DATE) \
  "

GOBUILD = $(GOCMD) build $(FLAGS)

.PHONY: all
all:	build

.PHONY: build
build:
	GOARCH=amd64 GOOS=darwin $(GOBUILD) -o $(DISTPATH).darwin $(PROJECT)
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(DISTPATH).linux $(PROJECT)

.PHONY: deploy
deploy: coverage build
	echo "Creating tar file"
	tar -zcf $(RELEASE_PATH)/$(PROJECT)-$(VERSION).tar.gz bin/$(PROJECT)*

.PHONY: run
run:
	- $(GOCMD) run src/$(PROJECT)/main.go

.PHONY: test
test:
	$(GOCMD) test ./src/$(PROJECT)/...

.PHONY: coverage
coverage:
	rm -fr coverage
	mkdir -p coverage
	$(GOCMD) list $(PROJECT)/... > coverage/packages
	@i=a ; \
	while read -r P; do \
		i=a$$i ; \
		$(GOCMD) test ./src/$$P -cover -coverpkg $$P -covermode=count -coverprofile=coverage/$$i.out; \
	done <coverage/packages
	echo "mode: count" > coverage/coverage
	cat coverage/*.out | grep -v "mode: count" >> coverage/coverage
	$(GOCMD) tool cover -html=coverage/coverage


.PHONY: clean
clean:
	rm -fR bin pkg
