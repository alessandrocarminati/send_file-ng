PRJ_NAME=send_console-ng
MAJOR=$(shell ./verscripts/maj.sh)
MINOR=$(shell ./verscripts/min.sh)
CHASH=$(shell ./verscripts/hash.sh)
DIRTY=$(shell ./verscripts/dirty.sh)
ALL_DEPENDENCIES := $(PRJ_NAME)-$(MAJOR).$(MINOR)
SOURCES := $(wildcard *.go)

all: $(ALL_DEPENDENCIES)

$(PRJ_NAME)-$(MAJOR).$(MINOR): $(SOURCES)
	go build -ldflags "-w -X 'main.AppName=$(PRJ_NAME)' -X 'main.Version=$(MAJOR)' -X 'main.Build=$(MINOR)' -X 'main.Hash=$(CHASH)' -X 'main.Dirty=$(DIRTY)'" -o  $(PRJ_NAME)-$(MAJOR).$(MINOR)

clean:
	rm -rf  $(PRJ_NAME)-*
