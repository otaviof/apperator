OPERATOR = "apperator"
IMAGE_PREFIX = "otaviof"
NAMESPACE = ${KUBERNETES_NAMESPACE}

default: build

bootstrap:
	dep ensure -v -vendor-only

# to run after modifying types, will generate items for k8s api
generate: FORCE
	operator-sdk generate k8s

build: FORCE
	operator-sdk build $(IMAGE_PREFIX)/$(OPERATOR):latest

build-test: FORCE
	operator-sdk build --enable-tests $(IMAGE_PREFIX)/$(OPERATOR):test

clean-vendor:
	rm -rf ./vendor > /dev/null

push: build
	docker push $(IMAGE_PREFIX)/$(OPERATOR):latest

push-test: build-test
	docker push $(IMAGE_PREFIX)/$(OPERATOR):test

test: FORCE
	go test -cover -v ./pkg/controller/apperatorapp

integration-local:
	operator-sdk test local ./test/e2e \
		--debug \
		--go-test-flags "-v" \
		--up-local \
		--namespace apperator

integration-cluster: push-test
	operator-sdk test cluster $(IMAGE_PREFIX)/$(OPERATOR):test \
		--namespace $(NAMESPACE) \
		--service-account $(NAMESPACE)

FORCE: ;
