TARGET=pegasus_engine
REGISTRY=innohubregister
NAME=pegasus-engine
VERSION=0.1.0
BUILD=`date +%FT%T%Z`

PACKAGES=`go list ./... | grep -v /vendor/`
VETPACKAGES=`go list ./... | grep -v /vendor/ | grep -v /examples/`
GOFILES=`find . -name "*.go" -type f -not -path "./vendor/*"`

default:
	@go build -o=${TARGET} -tags=jsoniter cmd/pegasus_engine.go

list:
	@echo ${PACKAGES}
	@echo ${VETPACKAGES}
	@echo ${GOFILES}

fmt:
	@gofmt -s -w ${GOFILES}

fmt-check:
	@diff=$$(gofmt -s -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

install:
	@govendor sync -v

test:
	@go test -cpu=1,2,4 -v -tags integration ./...

vet:
	@go vet $(VETPACKAGES)

genetest:
	@go test -run TestGeneticAlgorithm ./pkg/genetic -v -count 1


docker: docker-version

docker-version:
	docker build -t ${REGISTRY}/${NAME}:${VERSION} .

docker-tag:
	docker tag ${REGISTRY}/${NAME}:${VERSION} ${REGISTRY}/${NAME}:latest

push: docker-version docker-tag
	docker push ${REGISTRY}/${NAME}:${VERSION}; docker push ${REGISTRY}/${NAME}:latest

clean:
	@if [ -f ${TARGET} ] ; then rm ${TARGET} ; fi

.PHONY: default fmt fmt-check install test vet clean
