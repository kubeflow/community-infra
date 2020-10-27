package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

const (
	BulkProjectDeleteKind = "BulkProjectDelete"
	BulkProjectMoveKind   = "BulkProjectMove"
)

// BulkProjectDelete delete a bunch of GCP projects
type BulkProjectDelete struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec BulkProjectDeleteSpec `json:"spec,omitempty"`
}

type BulkProjectDeleteSpec struct {
	OrganizationId string   `json:"organizationId,omitempty"`
	Projects       []string `json:"projects,omitempty"`
	Folders        []string `json:"folders,omitempty"`
}

// BulkProjectMove moves GCP projects into folders
type BulkProjectMove struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec BulkProjectMoveSpec `json:"spec,omitempty"`
}

type BulkProjectMoveSpec struct {
	OrganizationId string        `json:"organizationId,omitempty"`
	Moves          []Move `json:"moves,omitempty"`
}

type Move struct {
	Project string `json:"project,omitempty"`
	Folder  string `json:"folder,omitempty"`
}
