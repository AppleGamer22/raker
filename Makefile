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
	go run . completion bash > rake.bash
	go run . completion fish > rake.fish
	go run . completion zsh > rake.zsh
	go run . completion powershell > rake.ps1

manual:
	if [[ "$$OSTYPE" == "linux-gnu"* ]]; then \
		sed -i "s/vVERSION/$(VERSION)/" rake.1; \
		sed -i "s/DATE/$$(date -Idate)/" rake.1; \
	elif [[ "$$OSTYPE" == "darwin"* ]]; then \
		sed -I '' "s/vVERSION/$(VERSION)/" rake.1; \
		sed -I '' "s/DATE/$$(date -Idate)/" rake.1; \
	fi

clean:
	rm -rf rake rakeserver bin dist rake.bash rake.fish rake.zsh rake.ps1
	go clean -testcache -cache

.PHONY: server cli debug test clean completion manual