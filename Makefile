# Старт контейнеров и сборка
run:
	docker-compose up --build

# Миграции вперед
migrate-up:
	docker-compose run --rm app goose -dir ./migrations up

# Миграции назад
migrate-down:
	docker-compose run --rm app goose -dir ./migrations down

# Миграции вперед и назад (по умолчанию)
migrate:
	docker-compose run --rm app goose -dir ./migrations up
	docker-compose run --rm app goose -dir ./migrations down

# Остановка контейнеров
stop:
	docker-compose down

# Очистка контейнеров и сети
clean:
	docker-compose down --volumes --remove-orphans
