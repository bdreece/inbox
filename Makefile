# inbox v0.1.0
# Copyright (C) 2025 Brian Reece

SHELL = /bin/sh
NPROCS = $(shell grep -c 'processor' /proc/cpuinfo)
MAKEFLAGS += -j$(NPROCS)

include scripts/variables.mk
include scripts/recipes.mk
include scripts/targets.mk
