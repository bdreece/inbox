$(srcdir)/bin:
	mkdir -p $(abspath $@)

$(srcdir)/bin/inbox: $(srcdir)/vendor/modules.txt | $(srcdir)/bin
	go build -v -o $(abspath ./bin) ./...

$(srcdir)/vendor/modules.txt:
	go mod vendor

$(sysconfdir)/inbox:
	mkdir -p $@

$(sysconfdir)/inbox/%: $(srcdir)/configs/% | $(sysconfdir)/inbox
	install -m 0644 -T $< $@

$(bindir)/inbox: $(srcdir)/bin/inbox
	install -m 0755 -T $(abspath $<) $@

