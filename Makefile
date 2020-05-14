TARGET=pegasus_engine
REGISTRY=innohubregister
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
	docker build -t ${REGISTRY}/${TARGET}:${VERSION} .

docker-tag:
	docker tag ${REGISTRY}/${TARGET}:${VERSION} ${REGISTRY}/${TARGET}:latest

push: docker-version docker-tag
	docker push ${REGISTRY}/${TARGET}:${VERSION}; docker push ${REGISTRY}/${TARGET}:latest

clean:
	@if [ -f ${TARGET} ] ; then rm ${TARGET} ; fi

.PHONY: default fmt fmt-check install test vet clean
