SHELL:=/bin/bash
PACKAGE:=github.com/AppleGamer22/rake
VERSION:=$(shell git describe --tags --abbrev=0 || echo '$(PACKAGE)/shared.Version')
HASH:=$(shell git rev-list -1 HEAD)
LDFLAGS:=-ldflags="-X '$(PACKAGE)/shared.Version=$(subst v,,$(VERSION))' -X '$(PACKAGE)/shared.Hash=$(HASH)'"

build: server cli

server:
	go build -race $(LDFLAGS) -o rakeserver ./server

cli:
	go build -race $(LDFLAGS) -o rake ./cli

test:
	go clean -testcache
	go test -v -race -cover ./shared/... ./server/...

completion:
	go run ./cli completion bash > rake.bash
	go run ./cli completion fish > rake.fish
	go run ./cli completion zsh > rake.zsh
	go run ./cli completion powershell > rake.ps1

manual:
	go run ./utils/replace rake.1 -b "vVERSION" -a "$(VERSION)"
	go run ./utils/replace rake.1 -b "DATE" -a "$(shell go run ./utils/date)"

clean:
	rm -rf rake rakeserver bin dist rake.bash rake.fish rake.zsh rake.ps1
	go clean -testcache -cache

.PHONY: server cli debug test clean completion manual