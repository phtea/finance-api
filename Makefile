# Start the application
run:
	docker-compose up --build

# Run migrations up
migrate-up:
	docker-compose run --rm app sh -c 'goose -dir ./migrations postgres "$$DATABASE_URL" up'

# Run migrations down
migrate-down:
	docker-compose run --rm app sh -c 'goose -dir ./migrations postgres "$$DATABASE_URL" down'

# Run migrations up and down
migrate:
	docker-compose run --rm app sh -c 'goose -dir ./migrations postgres "$$DATABASE_URL" up'
	docker-compose run --rm app sh -c 'goose -dir ./migrations postgres "$$DATABASE_URL" down'

# Stop and remove containers
stop:
	docker-compose down

# Clean up all resources (volumes and orphans)
clean:
	docker-compose down --volumes --remove-orphans
