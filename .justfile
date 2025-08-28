#!/usr/bin/env -S just --justfile

default: help

alias bench := benchmark

[doc("list help")]
[group("help")]
help:
  @just --list --justfile {{justfile()}}

[doc("run ci")]
[group("ci")]
ci: vet && test

[doc("test code")]
[group("develop")]
test:
  @echo "Testting code..."
  @cd {{invocation_directory()}}; go test ./...

[doc("benchmark code")]
[group("develop")]
benchmark:
  @echo "Benchmark code..."
  @cd {{invocation_directory()}}; go test -benchmem -bench . .

vet:
  @echo "Vetting code..."
  @go vet ./...

[doc("formatting code")]
[group("develop")]
fmt:
  @echo "Formatting code..."
  @go fmt ./...
