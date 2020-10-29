package controllers

import (

	"github.com/go-logr/zapr"
	"github.com/kubeflow/community-infra/pkg/api/v1alpha1"
	"go.uber.org/zap"
	"strings"
)

func ApplyBulkMove(bulkMove *v1alpha1.BulkProjectMove, orgHelper *OrgHelper) error {
	log := zapr.NewLogger(zap.L())

	folders, err := orgHelper.GetFolders(bulkMove.Spec.OrganizationId)
	if err != nil {
		log.Error(err, "Error listing folders")
		return err
	}

	log.Info("Got folders", "folders", folders)

	// Loop over the moves and apply them
	for _, m := range bulkMove.Spec.Moves  {
		folderId, ok := folders[m.Folder]

		if !ok {
			log.Info("Can't move folder; folder not found with given name", "folder", m.Folder, "project", m.Project)
			continue
		}

		p, err := orgHelper.Crm.Projects.Get(m.Project).Do()

		if err != nil {
			log.Error(err, "Could not get project info", "project", m.Project)
			continue
		}

		p.Parent.Type = "folder"

		pieces := strings.Split(folderId, "/")
		p.Parent.Id = pieces[len(pieces) -1]


		_, err = orgHelper.Crm.Projects.Update(m.Project, p).Do()

		if err != nil {
			log.Error(err, "Failed to move project", "folder", m.Folder, "folderId", folderId, "project", p, "parent", p.Parent)
		} else {
			log.Info("Moved project", "folder", m.Folder, "folderId", folderId, "project", p, "parent", p.Parent)
		}
	}
	//
	//// Now loop over the folders and delete them
	//for _, f := range bulkDelete.Spec.Folders {
	//
	//	if _, ok := folders[f]; !ok {
	//		log.Error(fmt.Errorf("Could not locate id for folder"), "Folder with display name not found in org; it might have been deleted already", "organization", bulkDelete.Spec.OrganizationId, "folder", f)
	//		continue
	//	}
	//	log.Info("Deleting folder", "folder", f, "folderId", folders[f])
	//
	//	_, err := orgHelper.CrmV2.Folders.Delete(folders[f]).Do()
	//
	//	if err != nil {
	//		log.Error(err, "Error deleting folder", "folder", f, "folderId", folders[f])
	//	}
	//}
	return nil
}
