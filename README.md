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
- 3. Install with helm.
```bash
helm repo add higio https://raw.githubusercontent.com/IIJ-Global-Solutions-Vietnam/charts/gh-pages/
helm install higio/namespace-admission-controller --generate-name -f value.yaml
```

## Require Permissions

- The user that generate API Token must have below permissions.

| Reosurce | Verbs |
| -------- | ----- |
| GlobalRoles | List <br> Get |
| GlobalRoleBindings | List <br> Get |
| ClusterRoleTemplates | List <br> Get |
| ClusterRoleTemplateBindings | List <br> Get |
| ProjectRoleTemplates | List <br> Get |
| ProjectRoleTemplateBindings | List <br> Get <br> Create |
| Projects | List <br> Get <br> Create|
