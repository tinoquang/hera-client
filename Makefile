run: setup
	@docker-compose up -d --build hera-server

setup:
	@docker-compose up -d postgres influxdb grafana

teardown:
	@docker-compose down -v --remove-orphans

