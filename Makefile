PRJ_NAME=send_console-ng
MAJOR=$(shell ./verscripts/maj.sh)
MINOR=$(shell ./verscripts/min.sh)
CHASH=$(shell ./verscripts/hash.sh)
REV=$(shell ./verscripts/rev.sh)
ALL_DEPENDENCIES := $(PRJ_NAME)-$(MAJOR).$(MINOR)
SOURCES := $(wildcard *.go)

all: $(ALL_DEPENDENCIES)

$(PRJ_NAME)-$(MAJOR).$(MINOR): $(SOURCES)
	go build -ldflags "-w -X 'main.AppName=$(PRJ_NAME)' -X 'main.Version=$(MAJOR)' -X 'main.Build=$(MINOR)' -X 'main.Hash=$(CHASH)' -X 'main.Rev=$(REV)'" -o  $(PRJ_NAME)-$(MAJOR).$(MINOR)

clean:
	rm -rf  $(PRJ_NAME)-*
