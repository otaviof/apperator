OPERATOR = "apperator"
IMAGE_PREFIX = "otaviof"

default: build

dep:
	dep ensure -v -vendor-only

# to run after modifying types, will generate items for k8s api
generate: FORCE
	operator-sdk generate k8s

build: FORCE
	operator-sdk build $(IMAGE_PREFIX)/$(OPERATOR)

FORCE: ;