include .env
export

.PHONY: swagger-track tidy refresh print

swagger-track:
	swag init \
		--generalInfo ./track-service/cmd/main.go \
		--output ./track-service/docs \
		--parseDependency \
		--parseInternal

tidy:
	go mod tidy

refresh:
	cp .env.example .env || true
	git pull origin master
	docker-compose build track db migrate
	docker-compose up -d db
	until docker-compose exec db pg_isready -U $(POSTGRES_USER); do sleep 1; done
	sleep 2
	docker-compose run --rm migrate
	docker-compose up -d track nginx

print:
	echo $$DATABASE_URL