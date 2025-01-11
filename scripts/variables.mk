GO = go
GOFLAGS = -v

AIR = $(GO) run github.com/air-verse/air@latest
AIRFLAGS = -c configs/air.toml

prefix = /usr/local

sysconfdir = $(prefix)/etc

datarootdir = $(prefix)/share
datadir = $(datarootdir)

exec_prefix = $(prefix)
bindir = $(exec_prefix)/bin

srcdir = $(abspath .)
