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

- GlobalRoles
  - List
  - Get
- GlobalrRoleBindings
  - List
  - Get
- ClusterRoles
  - List
  - Get
- ClusterRoleTemplates
  - List
  - Get
- ProjectRoleTemplate
  - List
  - Get
- ProjectRoleTemplateBinding
  - List
  - Get
  - Create
- Projects
  - List
  - Get
  - Create
