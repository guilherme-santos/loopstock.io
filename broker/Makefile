run:
	docker-compose up -d

logs:
	docker-compose logs -f ${SERVICE}

stop:
	docker-compose stop

cleanup:
	docker-compose down -v

nsqadmin-port:
	@docker-compose port nsqadmin 4171