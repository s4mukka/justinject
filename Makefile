export LOG_LEVEL=release
PROJECT_NAME=justinject

build: clean
	@ printf "Building application... "
	@ go build \
		-trimpath  \
		-o build/${PROJECT_NAME}
	@ echo "done"

build-alpine: clean
	@ printf "Building application... "
	@ GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build \
		-trimpath  \
		-o build/${PROJECT_NAME}
	@ echo "done"

clean: ## Builds binary
	@ printf "Cleaning application... "
	@ rm -f ${PROJECT_NAME}
	@ echo "done"

run: build
	@ echo "Running application... "
	@ build/${PROJECT_NAME} broker
	@ echo "done"

docker: build-alpine
	@ printf "Building image... "
	@ docker build -t ${PROJECT_NAME} .
	@ echo "done"

up: docker
	@ printf "Running compose... "
	@ docker compose up -d
	@ echo "done"

down:
	@ printf "Stopping compose... "
	@ docker compose down
	@ echo "done"
