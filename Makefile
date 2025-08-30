.PHONY: server cli debug test clean completion manual
SHELL:=$(shell which bash)
PACKAGE:=github.com/AppleGamer22/raker
VERSION:=$(shell git describe --tags --abbrev=0 || echo '$(PACKAGE)/shared.Version')
HASH:=$(shell git rev-list -1 HEAD)
LDFLAGS:=-ldflags="-X '$(PACKAGE)/shared.Version=$(subst v,,$(VERSION))' -X '$(PACKAGE)/shared.Hash=$(HASH)'"

build: server cli

server:
	# go build -race $(LDFLAGS) -o raker-server ./server
	docker compose up -d database
	-go run ./server || true
	docker stop database

cli:
	go build -race $(LDFLAGS) -o raker ./cli

test:
	go clean -testcache
	go test -v -race -cover ./shared/... ./server/...

debug:
	stalk watch -c "go run ./server" server/** shared/** templates/*

completion:
	go run ./cli completion bash > raker.bash
	go run ./cli completion fish > raker.fish
	go run ./cli completion zsh > raker.zsh
	go run ./cli completion powershell > raker.ps1

manual:
	go run ./utils/replace raker.1 -b "vVERSION" -a "$(VERSION)"
	go run ./utils/replace raker.1 -b "DATE" -a "$(shell go run ./utils/date)"

clean:
	rm -rf raker raker-server bin dist raker.bash raker.fish raker.zsh raker.ps1
	go clean -testcache -cache
