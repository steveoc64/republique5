SUBDIRS := gui/appwindow cmd/republique cmd/republique-ui

all: deps $(SUBDIRS)

deps:
	go get -u ./...

$(SUBDIRS):
	$(MAKE) -C $@

.PHONY: all $(SUBDIRS)

test:
	go test ./...

lint:
	golint ./...
	go vet ./...

help:
	@echo proto compile dump newgame serve info oob ui

protobuf:
	my-protoc.sh
	$(MAKE) -C proto
	neds-protoc.sh

compiler:
	$(MAKE) -C cmd/republique

compile:cmd/republique
	republique compile oob/jena*/fr*/*.oob

dump:
	cat oob/jena*/fr*/*.json | jq

introgame:cmd/republique
	republique compile games/Jena-intro.game
	republique info Jena-intro

newgame:cmd/republique
	#republique compile scenarios/jena-auerstadt-1806/Jena.scenario
	republique compile games/Jena-1.game
	republique info Jena-1

serve: cmd/republique
	republique serve -port 1815 -game Jena-intro

web: cmd/republique
	republique serve -port 1815 -web 8015 -game Jena-intro

info: cmd/republique
	republique info Jena-intro

oob: cmd/republique
	republique oob Jena-intro

ui: cmd/republique-ui
	$(MAKE) -C cmd/republique-ui
	republique-ui
