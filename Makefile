include .env
export

.PHONY: swagger-track tidy refresh print build-front

swagger-track:
	swag init \
		--dir ./track-service \
		--generalInfo cmd/main.go \
		--output track-service/docs \
		--parseDependency \
		--parseInternal

tidy:
	go mod tidy

refresh:
	git pull origin master
	cp .env.example .env || true
	docker-compose build track db migrate
	docker-compose up -d db
	until docker-compose exec db pg_isready -U $(POSTGRES_USER); do sleep 1; done
	sleep 2
	docker-compose up -d track nginx

build-front:
	cd .. && \
	cd retry-front && \
	npm run build && \
	rm -rf ../retry/front-dist/* && \
	cp -r dist/* ../retry/front-dist/

print:
	echo $$DATABASE_URL