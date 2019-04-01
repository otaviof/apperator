# Vault-Handler Integration

Manifest of entities required from Vault. Allows to copy secrets over to Kubernetes, or to expose
secrets as files. Please consider `vault-handler` project, for further information.

Apperator supports a specially crafted `ConfigMap` to inform `vault-handler` manifest and
authorization information, as per the following example:

``` yaml
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: vault-secrets
spec:
  # authentication and authorization data
  authorization: |-
    # kubernetes secret name
    secretName: vault
    # kubernetes secret key names
    secretKeys:
      # when using approle, key name to identify role-id
      roleId: role-id
      # when using approle, key name to identify secret-id
      secretId: secret-id
      # when using token, key name to identify the token
      token: token
  # manifest listing the secrets to copy from Vault
  secrets: |-
    # group name, used to name files, or will be employed as a Kubernetes Secret name.
    name:
      # path in vault
      path: secrets/dir1/dir2
      # type of format the secret is stored, use "file" to save it to a file, or a valid Kubernetes
      # Secret type to copy it to the cluster.
      type: file # kubernetes.io/tls
      # data contained in respective vault path.
      data:
        # entity name, employed to name "file" or "key" for Kubernetes secret
        - name: keystore
          # file extension
          extension: jks # when using "file" as type
          # zip payload
          zip: true # only when having "file" as type
```

TODO: document how to set-up secret for vault-handler authorization;
