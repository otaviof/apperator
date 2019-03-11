ARG BASEIMAGE
FROM ${BASEIMAGE}
ADD build/_output/bin/apperator-test /usr/local/bin/apperator-test
ARG NAMESPACEDMAN
ADD $NAMESPACEDMAN /namespaced.yaml
ADD build/test-framework/go-test.sh /go-test.sh
