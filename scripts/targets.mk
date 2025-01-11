all: build

## help: print this help message
.PHONY: help
help:
	@echo 'Usage: make [SCRIPT]'
	@echo 'Scripts:'
	@sed -n 's/^## //p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/  /'

## clean: remove dependencies and build artifacts
.PHONY: clean
clean:
	rm -rf $(srcdir)/vendor
	rm -rf $(srcdir)/bin
	rm -rf $(srcdir)/tmp

## restore: install dependencies
.PHONY: restore
restore: $(srcdir)/vendor/modules.txt

## build: compile the application
.PHONY: build
build: $(srcdir)/bin/inbox

## watch: watch the application
.PHONY: watch
watch:
	$(AIR) $(AIRFLAGS)

## install: install the application
.PHONY: install
install: $(sysconfdir)/inbox/settings.yml $(sysconfdir)/inbox/settings.prod.yml $(bindir)/inbox

## uninstall: uninstall the application
.PHONY: uninstall
uninstall:
	rm $(sysconfdir)/inbox/settings.yml $(sysconfdir)/inbox/settings.prod.yml
	rmdir $(sysconfdir)/inbox
	rm $(bindir)/inbox
