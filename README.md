# Kubeflow Community Infrastructure

This is a repository for using declarative configs and GitOps to managed shared community Kubeflow GCP infrastructure.

The management cluster is setup using the Kubeflow [management blueprint](kpt pkg get https://github.com/kubeflow/gcp-blueprints.git/management@master).

## Creating GCP Resources.

If you need to create GCP resources for Kubeflow or gain access to GCP resoruces you do so by creating
PRs against this repository.

* We use [ACM](https://cloud.google.com/anthos-config-management/docs/concepts/repo) to sync
  the [Cloud Config Connector(CNRM)](https://cloud.google.com/config-connector/docs/overview)
  to GKE cluster that will apply those resouces.

* ACM has an oppinionated layout to the repository which is rooted at "/prod"

  * See the [docs](https://cloud.google.com/anthos-config-management/docs/concepts/repo) for 
    how this repository should be layed out

  * There should be a namespace for every GCP project that is managed

* To create a new project

  1. Create subfolder `/prod/namespaces/${PROJECt}`
  1. Create `/pord/namespaces/${PROJECT}/namespace.yaml` defining a kubernetes namespace
  1. Create `/pord/namespaces/${PROJECT}/project.yaml` containing a [Project](https://cloud.google.com/config-connector/docs/reference/resources#project)
     resource defining your project
  1. Create `/pord/namespaces/${PROJECT}/iam-policy-members.yaml` containing a [IAMPolicyMember](https://cloud.google.com/config-connector/docs/reference/resources#iampolicymember)
     resource granting IAM permissions to access the project is necessary

* Wait for the PR to be approved
* Once the PR is merged the resources should be created automatically.

## Setup

1. Follow the [management blueprint](kpt pkg get https://github.com/kubeflow/gcp-blueprints.git/management@master)
   
   * Do not install CNRM; we will use ConfigSync to install CNRM

1. Follow the [ACM installation guide](https://cloud.google.com/anthos-config-management/docs/how-to/installing)


   * Create the service account 'cnrm-system' in project `kf-kcc-admin`
   * **Note** It looks like when using ACM to install and manage CNRM you can't use workload identity and need to provide
     a GCP service account key.


1. Make sure the CNRM service account has roles `roles/owner` and project creator on the community folder