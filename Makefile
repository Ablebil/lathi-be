run:
	@docker compose up
build:
	@docker compose down
	@docker compose up --build
stop:
	@docker compose down
logs:
	@docker compose logs -f
restart:
	@docker compose down
	@docker compose up
down:
	@docker compose down -v
migrate-down:
	@docker compose exec app /app/server migrate -action down
migrate-up:
	@docker compose exec app /app/server migrate -action up
seed-all:
	@docker compose exec app /app/server seed
