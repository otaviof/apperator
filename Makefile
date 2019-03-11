OPERATOR = "apperator"
IMAGE_PREFIX = "otaviof"
NAMESPACE = ${KUBERNETES_NAMESPACE}

default: build

dep:
	dep ensure -v -vendor-only

# to run after modifying types, will generate items for k8s api
generate: FORCE
	operator-sdk generate k8s

build:
	operator-sdk build --enable-tests $(IMAGE_PREFIX)/$(OPERATOR):latest

push: build
	docker push $(IMAGE_PREFIX)/$(OPERATOR):latest

test: FORCE
	go test -v ./pkg/controller/apperatorapp

integration:
	operator-sdk test local ./test/e2e --namespace $(NAMESPACE)

FORCE: ;
