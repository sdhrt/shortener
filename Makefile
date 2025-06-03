build:
	@go build -o ./bin/fs ./cmd/api/

run: build
	@./bin/fs

dev:
	cd www && pnpm run dev

web:
	make -j 2 run dev

c:
	curl -d '{"Url": "https://www.youtube.com/feed/history"}' \
		-H "Content-Type: application/json" \
		localhost:8080/create -vv

h:
	curl localhost:8080/healthcheck -vv
