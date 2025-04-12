export LOG_LEVEL=release
PROJECT_NAME=justinject
EXCLUDES_COVERAGE=domain|mock

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

helm-build:
	@ helm package helm -d helm-releases
	@ helm repo index helm-releases

clean: ## Builds binary
	@ printf "Cleaning application... "
	@ rm -rf ${PROJECT_NAME} coverage
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

minikube:
	@ if helm ls -n ${PROJECT_NAME} | grep ${PROJECT_NAME}; then \
		echo "Helm release ${PROJECT_NAME} already installed. Uninstalling..."; \
		helm uninstall ${PROJECT_NAME} -n ${PROJECT_NAME}; \
	fi
	@ helm install ${PROJECT_NAME} helm -n ${PROJECT_NAME}

test: clean
	@ printf "Running tests... "
	@ mkdir -p coverage
	@ go test ./... -coverprofile=coverage/cover.out.tmp -v | grep -Ev "${EXCLUDES_COVERAGE}"
	@ cat coverage/cover.out.tmp | grep -Ev "${EXCLUDES_COVERAGE}" > coverage/cover.out
	@ rm -f coverage/cover.out.tmp
	@	cat coverage/cover.out | \
			awk 'BEGIN {cov=0; stat=0;} \
			$$3!="" { cov+=($$3==1?$$2:0); stat+=$$2; } \
			END {printf("Total coverage: %.2f%% of statements\n", (cov/stat)*100);}'
	@ echo "done"

citest: test
	@ go tool cover -func coverage/cover.out -o coverage/cover.out

coverage: test
	@ go tool cover -html coverage/cover.out

lint:
	@ printf "Running lint... "
	@ golangci-lint run
	@ echo "done"

format:
	@ printf "Running format... "
	@ gofumpt -l -w .
	@ goimports-reviser -rm-unused -set-alias -format ./...
	@ echo "done"
