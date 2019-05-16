SUBDIRS := republique cmd/republique

all: $(SUBDIRS)
$(SUBDIRS):
	$(MAKE) -C $@

.PHONY: all $(SUBDIRS)

compile:
	republique compile oob/jena*/fr*/*.oob

dump:
	cat oob/jena*/fr*/*.json | jq
