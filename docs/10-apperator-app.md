# `ApperatorApp`

Extends the behavior of Deployment Kubernetes resource, by composing how a application should look
like. The idea is to create other CRDs to define Environment variables, init-containers, sidecars and
functional monitoring in the same context.

``` yaml
---
apiVersion: apperator.otaviof.github.io/v1alpha1
kind: ApperatorApp
metadata:
  name: app
spec:
  # Kubernetes Deployment.
  deployment:
    spec: {}
  # Kubernetes envVar construction.
  envs:
    - environment-variables-example
  # Kubernetes initContainer spec.
  initContainers:
    - example-init-container
  # Kubernetes container spec.
  sidecars:
    - jmx-exporter-sidecar
    - example-sidecar
  # Functional monitoring using Prometheus alerts as input to restart application pods. The alerts
  # can be configured in Prometheus itself, by using it's default configuration file, or Prometheus
  # Operator.
  probes:
    - functional-monitoring-example
  # A special type of configuration to specify secrets to download from Vault, before the container
  # starts, therefore it's placed as a init-container, however with more steps involved.
  vault:
    - vault-secrets-manifest
```
