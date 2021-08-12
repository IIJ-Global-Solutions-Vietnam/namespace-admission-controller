# namespace-admission-controller

## Install

- 1. Generate Rancher API Token
  - Select 'no scope'
- 2. Create value.yaml below.
```yaml
rancher:
  clusterID: "your clusterID"
  url: "your rancher API endpoint"
  apiToken: "your Rancher API bearer Token"
```
- 3. Install with Helm.
```bash
helm repo add higio https://raw.githubusercontent.com/IIJ-Global-Solutions-Vietnam/charts/gh-pages/
helm install higio/namespace-admission-controller --generate-name -f value.yaml
```

### Specify the namespace to install

You can install in the specified namespace by adding `-n` option
```
helm install higio/namespace-admission-controller --generate-name -n kube-system -f value.yaml
```

### Uninstall

- 1. Show the Helm relase name
```bash
$ helm list
NAME                                     	NAMESPACE	REVISION	UPDATED                             	STATUS  	CHART                               	APP VERSION
namespace-admission-controller-1622182468	default  	1       	2021-05-28 15:14:32.217555 +0900 JST	deployed	namespace-admission-controller-0.1.0	1.16.0
```

- 2. Uninstall with `helm uninstall` command
```bash
$ helm uninstall namespace-admission-controller-1622182468
```

## Require Permissions

The user that generate API Token must have below permissions.

| Reosurce | Verbs |
| -------- | ----- |
| GlobalRoles | Get <br> List |
| GlobalRoleBindings | Get <br> List |
| ClusterRoleTemplateBindings | Get <br> List |
| Projects | Get <br> List <br> Create |
| RoleTemplates | Get <br> List |
| ProjectRoleTemplateBindings | Get <br> List <br> Create |
