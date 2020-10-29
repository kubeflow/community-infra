package controllers

import (
	"context"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"google.golang.org/api/cloudresourcemanager/v1"
	crmV2 "google.golang.org/api/cloudresourcemanager/v2"
	"google.golang.org/api/servicemanagement/v1"
)

type OrgHelper struct {
	Crm   *cloudresourcemanager.Service
	CrmV2 *crmV2.Service
	Sm     *servicemanagement.APIService
}

// GetFolders returns a map from display name to folder id.
func (o *OrgHelper) GetFolders(orgId string) (map[string]string, error) {
	log := zapr.NewLogger(zap.L())
	// Map from the name of the folder to an id in the form folders/{folder_id} which is what delete expects
	folders := map[string]string{}

	addToMap := func(r *crmV2.ListFoldersResponse) error {
		for _, f := range r.Folders {
			folders[f.DisplayName] = f.Name
		}
		return nil
	}
	err := o.CrmV2.Folders.List().Parent(orgId).Pages(context.Background(), addToMap)

	if err != nil {
		log.Error(err, "Error listing folders")
		return nil, err
	}
	return folders, nil
}
