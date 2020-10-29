package controllers

import (
	"context"
	"fmt"
	"github.com/go-logr/zapr"
	"github.com/kubeflow/community-infra/pkg/api/v1alpha1"
	"go.uber.org/zap"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/servicemanagement/v1"
	"net/http"
)

func ApplyBulkDelete(bulkDelete *v1alpha1.BulkProjectDelete, orgHelper *OrgHelper) error {
	log := zapr.NewLogger(zap.L())
	folders := map[string]string{}
	if len(bulkDelete.Spec.Folders) > 0 {
		newFolders, err := orgHelper.GetFolders(bulkDelete.Spec.OrganizationId)
		if err != nil {
			log.Error(err, "Error listing folders")
			return err
		}
		folders = newFolders
	}
	log.Info("Got folders", "folders", folders)
	projects := map[string]bool{}
	// Check if the project has any service endpoints. These needed to be deleted first
	// before deleting the project as they will block project deletion
	for _, p := range bulkDelete.Spec.Projects {
		deleteEndpoint := func(r *servicemanagement.ListServicesResponse) error {
			for _, s := range r.Services {
				log.Info("Deleting service endpoint", "project", p, "serviceName", s.ServiceName)
				op, err := orgHelper.Sm.Services.Delete(s.ServiceName).Do()
				if err != nil {
					log.Error(err, "Error deleting service", "project", p, "serviceName", s.ServiceName)
				}
				log.Info("Deleting service endpoint", "project", p, "serviceName", s.ServiceName, "operation", op)
			}
			return nil
		}

		err := orgHelper.Sm.Services.List().ProducerProjectId(p).Pages(context.Background(), deleteEndpoint)

		if err != nil {
			gErr, ok := err.(*googleapi.Error)

			if ok && gErr.Code == http.StatusNotFound {
				log.Info("Project not found; either it doesn't exist or you don't have permissions", "project", p)
				continue
			}
			log.Error(err, "There was a problem listing services", "project", p)
		}

		projects[p] = true
	}

	// Now loop over the projects and delete them.
	for p, _ := range projects {
		log.Info("Deleting project", "project", p)

		_, err := orgHelper.Crm.Projects.Delete(p).Do()

		if err != nil {
			log.Error(err, "Error deleting project", "project", p)
		}
	}

	// Now loop over the folders and delete them
	for _, f := range bulkDelete.Spec.Folders {

		if _, ok := folders[f]; !ok {
			log.Error(fmt.Errorf("Could not locate id for folder"), "Folder with display name not found in org; it might have been deleted already", "organization", bulkDelete.Spec.OrganizationId, "folder", f)
			continue
		}
		log.Info("Deleting folder", "folder", f, "folderId", folders[f])

		_, err := orgHelper.CrmV2.Folders.Delete(folders[f]).Do()

		if err != nil {
			log.Error(err, "Error deleting folder", "folder", f, "folderId", folders[f])
		}
	}
	return nil
}
