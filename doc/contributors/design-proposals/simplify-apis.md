---
title: simplify-olm-apis
authors:
  - "@njhale"
reviewers:
  - "@ecordell"
approvers:
  - TBD
creation-date: 2019-09-05
last-updated: 2019-09-19
status: provisional
see-also:
  - "http://bit.ly/rh-epic_simplify-olm-apis"
---

# simplify-olm-apis

## Release Signoff Checklist

- [ ] Enhancement is `implementable`
- [ ] Design details are appropriately documented from clear requirements
- [ ] Test plan is defined
- [ ] Graduation criteria for dev preview, tech preview, GA
- [ ] User-facing documentation is created in [openshift/docs]

## Summary

This enhancement iterates on OLM APIs to reduce the number of resource types and provide a smaller surface by which to manage an operator's lifecycle.

TODO: finish summary

## Motivation

Operator authors perceive OLM/Marketplace v1 APIs as difficult to use. This is thought to stem from three primary sources of complexity:

1. Too many types
2. Redundancies (e.g. OperatorSource and CatalogSource)
3. Effort to munge native k8s resources into OLM resources (e.g. Deployment/RBAC to CSV)

Negative perceptions stunt community adoption and cause a loss of mindshare, while complexity impedes product stability and feature delivery. Reducing OLM's API surface will help to avert these scenarios by making direct interaction with OLM more straightforward. A simplified user experience will encourage operator development and testing.

### Goals

- Define an API resource that aggregates the complete state of an operator
- Define a single API resource that allows an authorized user to:
  - install an operator without ancillary setup (e.g `OperatorGroup` not required)
  - optionally subscribe to installation of operator updates

- Remain backwards compatible with operators installed by previous versions of OLM
- Retain all of OLM's current features

### Non-Goals

- Define the implementation of an operator bundle
- Replace OLM's existing APIs
- Deprecate `OperatorSource`

## Proposal

### User Stories

#### As an __OLM Admin__, I want to

- restrict the resources OLM will apply from a bundle image/index when installing an operator
- restrict the `ServiceAccounts` a user can install operators with
- restrict the bundle indexes a given user can use to resolve operator bundles

#### As an __Operator Installer__, I want to

- view the state of an operator by inspecting the status of a single API resource
- deploy an operator by applying its manifests directly
- deploy an operator using a bundle image
- deploy and upgrade an operator by referencing a bundle index

#### As an _Operator User__, I want to

- view the operators I can use on a cluster

### Implementation Details/Notes/Constraints

#### Viewing On-cluster Operators

To enable viewing operators on a cluster:

- Add a new cluster-scoped resource that represents a cluster
  - call this resource `Operator`
- On the `Operator` resource, surface:
  - references to the resources that compose an operator
  - operator metadata

__Details:__

Add the `operators.coreos.com/v2alpha1` API group version to OLM

- alpha versions make no [support guarantees](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api_changes.md#alpha-beta-and-stable-versions), so resources can be iterated quickly

Add a command option to `olm-operator` that enables `v1alpha2`

- leave disabled by default

```bash
olm-operator --enable-v2alpha1
```

Add the cluster-scoped `Operator` resource to the `v2alpha1`

- cluster-scoped resources can be viewed without granting namespace list
- they also require a `ClusterRole` to configure any access

Define a label convention for selecting resources that compose an operator

- `operators.coreos.com/v2/operators/<name>: ""`; where `name` is `metadata.name`
- using a unique label key helps avoid collisions
- using a deterministic label key lets users know the exact label used in advance

Define, as part of the `Operator` resource, a `status.resources` field that surfaces:

- `matchLabels`: label selector used to list operator components
- `components`: references to resources selected by `matchLabels`

Add OLM controller logic such that, for each `Operator` resource:

- `matchLabels` is resolved from the label convention and stored in the `status.resources.matchLabels` field

```yaml
status:
  resources:
    matchLabels:
      operators.coreos.com/v2/operators/plumbus: ""
```

Add OLM controller logic such that, for each `Operator` resource:

- cluster-scoped resources are selected using `status.resources.matchLabels`
- namespace-scoped resources are selected __across all namespaces__ using `status.resources.matchLabels`
- the union of both selections is stored in `status.resources.components`
- kinds selected in both cases may be predefined or queried from discovery

```yaml
status:
  resources:
    matchLabels:
      operators.coreos.com/v2/operator-name: plumbus
    components:
    - kind: Subscription
      namespace: my-ns
      name: plumbus
      uid: d655a13e-d06a-11e9-821f-9a2e3b9d8156
      apiVersion: operators.coreos.com/v1alpha1
      resourceVersion: 109719
    - kind: ClusterServiceVersion
      namespace: my-ns
      name: plumbus.v2.0.0-alpha
      uid: d70a53b5-d06a-11e9-821f-9a2e3b9d8156
      apiVersion: operators.coreos.com/v1alpha1
      resourceVersion: 109811
    # ...
```

Add OLM controller logic, such that for each `Operator` resource

- OLM resources matching `status.resources.matchLabels` have the matching labels projected onto all related generated resources
  - from `Subscriptions` onto resolved `Subscriptions` and `InstallPlans`
  - from `InstallPlans` onto `ClusterServiceVersions`, `CustomResourceDefinitions`, RBAC, etc.
  - from `ClusterServiceVersions` onto `Deployments` and `APIServices`

Define, as part of the `Operator` resource, a `status.metadata` field that

- contains the `displayName`, `description`, and `version` of an operator

```yaml
status:
  metadata:
    displayName: Plumbus
    description: Welcome to the exciting world of Plumbus ownership! A Plumbus will aid many things in life, making life easier. With proper maintenance, handling, storage, and urging, Plumbus will provide you with a lifetime of better living and happiness.
    version: 2.0.0-alpha
```

Define, as part of the `Operator` resource, a `status.apis` field that

- contains the `required` and `provided` APIs of an operator

```yaml
status:
  apis:
    provides:
    - group: how.theydoit.com
      version: v2alpha1
      kind: Plumbus
      plural: plumbai
      singular: plumbus
    requires:
    - group: how.theydoit.com
      version: v1
      kind: Grumbo
      plural: grumbos
      singular: grumbo
```

Define, as part of the `Operator` resource, a `status.permissions` field that

- contains the RBAC requirements of an operator

```yaml
status:
  cluster: # optional
  - rules:
    - apiGroups:
      - how.theydoit.com
      resources:
      - plumbus
      verbs:
      - get
      - list
      - watch
      - create
      - update
    serviceAccountName: mr-plumbus
  namespaced:
  - rules:
    - apiGroups:
      - how.theydoit.com
      resources:
      - grumbos
      verbs:
      - get
      - list
      - watch
    serviceAccountName: mr-plumbus
```

Add OLM controller logic, such that for each `Operator` resource

- a single CSV in the operator's component selection is picked
  - the _newest_ with respect to the `spec.version` field
- the pick's `displayName`, `description`, and `version` are projected onto the operator's `spec.metadata` field
- the pick's `required` and `provided` APIs are projected onto the operator's `spec.apis` field
- the pick's `permissions` are projected onto the operator's `spec.permissions` field

__Notes:__

Open Questions:

- Why cluster-scoped?
- Why create a new resource?

#### Aggregate Operator Status

In order to aggregate the status of an operator, the `Operator` resource will:

- surface status conditions for each selected component
- surface top-level status conditions that describe overall abnormal operator states

__Details:__

Define, as part of the `Operator` resource, a `status.resources.components[*].conditions` field that

- enriches each element with conditions representing the state of the referenced component
- should be relevant in the context of an operator
- need not be surfaced for all component kinds
- should follow [k8s status condition conventions](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties)

```yaml
status:
  resources:
    matchLabels:
        operators.coreos.com/v2/operators/plumbus: ""
    components:
    - kind: ClusterServiceVersion
      namespace: my-ns
      name: plumbus.v2.0.0-alpha
      uid: d70a53b5-d06a-11e9-821f-9a2e3b9d8156
      apiVersion: operators.coreos.com/v1alpha1
      resourceVersion: 109811
      conditions:
      - type: Installing
        status: True
        reason: AllPreconditionsMet
        message: deployment rolling out
        lastTransitionTime: "2019-09-16T22:26:29Z"
```

Phase-in OLM controller logic such that, for each `Operator` resource:

- component conditions from OLM kinds are updated
  - `ClusterServiceVersion`
  - `Subscription`
  - `InstallPlan`

Define, as part of the `Operator` resource, a `status.conditions` field that

- surfaces status conditions which describe the overall abnormal state of the operator
- should follow [k8s status condition conventions](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-architecture/api-conventions.md#typical-status-properties)

```yaml
status:
   conditions:
   - kind: Upgrading
     status: True
     reason: PivotingVersions
     message: pivoting between versions
     lastTransitionTime: "2019-09-16T22:26:29Z"
```

Phase-in OLM controller logic such that, for each `Operator` resource:

- status conditions track the following states
  - upgrading
  - degraded
    - API conflicts between operators

#### Installing from a bundle index

TODO: section

### Risks and Mitigations

TODO: section

## Design Details

### Test Plan

TODO: section

### Graduation Criteria

TODO: section

#### Examples

TODO: section

##### Dev Preview -> Tech Preview

TODO: section

##### Tech Preview -> GA

TODO: section

##### Removing a deprecated feature

TODO: section

### Upgrade / Downgrade Strategy

TODO: section

### Version Skew Strategy

TODO: section

## Implementation History

TODO: section

## Drawbacks

TODO: section

## Alternatives

TODO: section
