# Using Kustomize with CRDs

In order to use `CustomResourceDefintions` (CRDs) with [kustomize](https://github.com/kubernetes-sigs/kustomize), an OpenAPI schema containing their definitions is required. By default, kustomize looks for this schema in a file named `crds.json`.

After each addition/modification/removal of a CRD, this schema must be regenerated. A simple way to do this for OLM is to apply the CRDs to a kubernetes instance and extract only the definitions belonging to the `operators.coreos.com` API group:

```sh
kubectl get --raw /openapi/v2 | jq '{definitions: .definitions|with_entries(select(.key|test("com.coreos.operators")))}' > crds.json
```
