build:
	docker-compose up

remove:
	docker-compose down --rmi local -v --remove-orphans