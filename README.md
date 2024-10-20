<p align="center">
  <img src="docs/static/cassata_logo1.png" alt="Cassata Logo">
</p>

<h1 align="center">Cassata</h1>
<div align="center">

**Cassata provides an API interface for [Crossplane](https://www.crossplane.io/). It enables platform teams to offer crossplane as a service to product teams. Reduces the complexity of managing and using crossplane through easy-to-use & uniform APIs for all crossplane-managed resources. You dont have to worry about XRDs, Compositions and Claims anymore!**

</div>

## Introduction

Cassata has powerful APIs for managing cloud resources across all Crossplane Providers. 
- For Platform Teams
  - Rich admin APIs for managing users, workspaces, permissions & policies. 
  - Resource isolation & RBAC
- For Product Teams:
  - Declarative APIs for creating and managing your resources. 
  - Share resources with other teams.

## Getting Started

### Prerequisites

- Kubernetes cluster with kubectl configured
- A working Crossplane installation.
- Helm 3.x installed
- Optional: A database connection (Postgres is recommended). 


### Installation

To install Cassata using Helm, run the following command: 
```bash
helm repo add cassata https://cassata-io.github.io/cassata
helm install cassata cassata/cassata
```

***Check the [values.yaml](deployment/helm/values.yaml) file for more configuration options.***

Alternatively, you can use the `--set` flag to override the default values:
```bash
helm install cassata cassata/cassata --set database.external=true --set database.dsn="your_database_dsn"
```

