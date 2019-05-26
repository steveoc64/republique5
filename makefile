SUBDIRS := republique cmd/republique cmd/republique-ui

all: $(SUBDIRS)
$(SUBDIRS):
	$(MAKE) -C $@

.PHONY: all $(SUBDIRS)

help:
	@echo proto compile dump newgame serve info oob ui

proto:
	$(MAKE) -C republique/proto

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
	republique serve -port 1815 -web 8015 -game Jena-intro

info: cmd/republique
	republique info Jena-intro

oob: cmd/republique
	republique oob Jena-intro

ui: cmd/republique-ui
	$(MAKE) -C cmd/republique-ui
	republique-ui
