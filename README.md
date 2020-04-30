# Kubeflow Community Infrastructure

This is a repository for using declarative configs and GitOps to managed shared community Kubeflow GCP infrastructure.

The management cluster is setup using the Kubeflow [management blueprint](kpt pkg get https://github.com/kubeflow/gcp-blueprints.git/management@master).

## Setup

1. Follow the [management blueprint](kpt pkg get https://github.com/kubeflow/gcp-blueprints.git/management@master)
   
   * Do not install CNRM; we will use ConfigSync to install CNRM

1. Follow the [ACM installation guide](https://cloud.google.com/anthos-config-management/docs/how-to/installing)


   * Create the service account 'cnrm-system' in project `kf-kcc-admin`
   * **Note** It looks like when using ACM to install and manage CNRM you can't use workload identity and need to provide
     a GCP service account key.


1. Make sure the CNRM service account has roles `roles/owner` and project creator on the community folder