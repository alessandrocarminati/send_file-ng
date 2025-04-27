PRJ_NAME=send_console-ng
MAJOR=$(shell ./verscripts/maj.sh)
MINOR=$(shell ./verscripts/min.sh)
CHASH=$(shell ./verscripts/hash.sh)
REV=$(shell ./verscripts/rev.sh)
ALL_DEPENDENCIES := $(PRJ_NAME)-$(MAJOR).$(MINOR)
SOURCES := $(wildcard *.go)
pwd := $(shell pwd)

all: $(ALL_DEPENDENCIES)

$(PRJ_NAME)-$(MAJOR).$(MINOR): $(SOURCES)
	go build -ldflags "-w -X 'main.AppName=$(PRJ_NAME)' -X 'main.Version=$(MAJOR)' -X 'main.Build=$(MINOR)' -X 'main.Hash=$(CHASH)' -X 'main.Rev=$(REV)'" -o  $(PRJ_NAME)-$(MAJOR).$(MINOR)

clean:
	rm -rf  $(PRJ_NAME)-*
	rm -f send_console-ng.spec
	rm -f send_console-ng-*.tar.gz

spec: send_console-ng.spec

send_console-ng.spec: send_console-ng.spec.in
	awk '{gsub(/%%PRJ_NAME%%/, "$(PRJ_NAME)"); gsub(/%%MAJOR%%/, "$(MAJOR)"); gsub(/%%MINOR%%/, "$(MINOR)"); print}' $< > $@

$(PRJ_NAME)-$(MAJOR).$(MINOR).tar.gz:
	tar --exclude='.git/*' --exclude='send_console-ng.spec' --exclude='$(PRJ_NAME)-*.tar.gz' -zcvf $(PRJ_NAME)-$(MAJOR).$(MINOR).tar.gz .

srpm: $(PRJ_NAME)-$(MAJOR).$(MINOR).tar.gz spec
	rpmbuild --define "_sourcedir $(pwd)" --define "_specdir $(pwd)" --define "_srcrpmdir $(pwd)" -bs send_console-ng.spec
