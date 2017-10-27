network:
	@-docker network create loopstock_network

docker-image:
	$(MAKE) -C integer-api $@
	$(MAKE) -C integer-gen $@
	$(MAKE) -C integeraverage-cal $@

run: network
	$(MAKE) -C broker $@
	$(MAKE) -C integer-api $@
	$(MAKE) -C integer-gen $@
	$(MAKE) -C integeraverage-cal $@

stop:
	$(MAKE) -C broker $@
	$(MAKE) -C integer-api $@
	$(MAKE) -C integer-gen $@
	$(MAKE) -C integeraverage-cal $@

cleanup:
	$(MAKE) -C broker $@
	$(MAKE) -C integer-api $@
	$(MAKE) -C integer-gen $@
	$(MAKE) -C integeraverage-cal $@
	docker network rm loopstock_network

ports:
	@echo nsqadmin: `$(MAKE) -C broker nsqadmin-port`
	@echo integer-api: `$(MAKE) -C integer-api integer-api-port`

.PHONY: network docker-image run stop cleanup