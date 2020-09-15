# Kubeflow Community Infrastructure

This is a repository for using declarative configs and GitOps to managed shared community Kubeflow GCP infrastructure.

The management cluster is setup using the Kubeflow [management blueprint](https://github.com/kubeflow/gcp-blueprints).

## Creating GCP Resources.

If you need to create GCP resources for Kubeflow or gain access to GCP resources you do so by creating
PRs against this repository.

* We use [ACM](https://cloud.google.com/anthos-config-management/docs/concepts/repo) to sync
  the [Cloud Config Connector(CNRM)](https://cloud.google.com/config-connector/docs/overview)
  to GKE cluster that will apply those resouces.

* ACM has an oppinionated layout to the repository which is rooted at "/prod"

  * See the [docs](https://cloud.google.com/anthos-config-management/docs/concepts/repo) for 
    how this repository should be layed out

  * There should be a namespace for every GCP project that is managed

* Follow these steps to create new project. Note that `${PROJECT}` name must be globally unique across all GCP projects.

  1. Create subfolder `/prod/namespaces/${PROJECT}`.

  1. Create `/prod/namespaces/${PROJECT}/namespace.yaml` defining a Kubernetes namespace.
  Namespace name should be equal to `${PROJECT}` name.

  1. Create `/prod/namespaces/${PROJECT}/project.yaml` containing a
  [`Project`](https://cloud.google.com/config-connector/docs/reference/resource-docs/resourcemanager/project) resource defining your project.

  1. Create `/prod/namespaces/${PROJECT}/iam-policy-members.yaml` containing a [`IAMPolicyMember`](https://cloud.google.com/config-connector/docs/reference/resource-docs/iam/iampolicymember) resource list with necessary IAM permissions to access the project.
  Each `IAMPolicyMember` should have unique name.

      You can give `roles/editor` to your GCP user account to view created project.

      If you want to integrate your project with `kubeflow-ci`,
      you have to give access to this service account: `serviceAccount:kubeflow-testing@kubeflow-ci.iam.gserviceaccount.com`.


      `kubeflow-testing` service account should have these permissions:
        - `roles/editor` to modify GCP resources.
        - `roles/cloudbuild.builds.editor` to create Cloud Builds.
        - `roles/container.admin` to manage Kubernetes clusters.


* Wait for the PR to be approved

* Once the PR is merged the resources should be created automatically and you can access created GCP project.
You can run `kubectl describe` on appropriate resource in `kf-community-admin` cluster to check status.

## Setup

1. Follow the [management blueprint](https://github.com/kubeflow/gcp-blueprints)
   
   * Do not install CNRM; we will use ConfigSync to install CNRM

1. Follow the [ACM installation guide](https://cloud.google.com/anthos-config-management/docs/how-to/installing)


   * Create the service account 'cnrm-system' in project `kf-kcc-admin`
   * **Note** It looks like when using ACM to install and manage CNRM you can't use workload identity and need to provide
     a GCP service account key.


1. Make sure the CNRM service account has roles `roles/owner` and project creator on the community folder
