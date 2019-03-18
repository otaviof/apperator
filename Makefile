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

push: build
	docker push $(IMAGE_PREFIX)/$(OPERATOR):latest

test: FORCE
	go test -cover -v ./pkg/controller/apperatorapp

integration-local:
	operator-sdk test local ./test/e2e --namespace $(NAMESPACE)

integration-cluster: build-test
	docker push $(IMAGE_PREFIX)/$(OPERATOR):test
	operator-sdk test cluster --namespace $(NAMESPACE) $(IMAGE_PREFIX)/$(OPERATOR):test

FORCE: ;
