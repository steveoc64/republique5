SUBDIRS := republique cmd/republique cmd/republique-ui

all: $(SUBDIRS)
$(SUBDIRS):
	$(MAKE) -C $@

.PHONY: all $(SUBDIRS)

compile:
	republique compile oob/jena*/fr*/*.oob

dump:
	cat oob/jena*/fr*/*.json | jq

test:
	#republique compile scenarios/jena-auerstadt-1806/Jena.scenario
	republique compile games/Jena-1.game
	republique info Jena-1

serve: 
	republique serve -port 1815 -web 8015 -game Jena-1

info:
	republique info Jena-1

oob:
	republique oob Jena-1

ui:
	$(MAKE) -C cmd/republique-ui
	republique-ui
