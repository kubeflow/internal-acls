package groups

import (
	"context"
	"github.com/go-logr/logr"
	"github.com/kubeflow/internal-acls/google_groups/pkg/api/v1alpha1"
	admin "google.golang.org/api/admin/directory/v1"
	settingsSdk "google.golang.org/api/groupssettings/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

// GroupImporter is used to import existing groups to YAML files.
type GroupImporter struct {
	Client *http.Client
	Log logr.Logger
}

// Import group definitions
func (s *GroupImporter) Import(org string) ([]*v1alpha1.GoogleGroup, error) {
	log := s.Log
	results := []*v1alpha1.GoogleGroup{}

	// TODO(jlewi): Using admin.NewService was giving me auth problems. I think it might be an OAuthScope
	// issue because the credential didn't have all the scopes but admin.NewService appears to request all scopes
	// so the client was requesting a token for scopes it wasn't authorized for.
	service, err := admin.New(s.Client)

	if err != nil {
		return results, err
	}

	settingsService, err := settingsSdk.New(s.Client)

	if err != nil {
		return results, err
	}

	groups := []*admin.Group{}
	pageFunc := func(g *admin.Groups) error {
		groups = append(groups, g.Groups...)
		return nil
	}

	err = service.Groups.List().Domain(org).Pages(context.Background(), pageFunc)

	if err != nil {
		return results, err
	}

	for _, g := range groups {
		newGroup  := &v1alpha1.GoogleGroup{
			ObjectMeta: metav1.ObjectMeta{
				Name: g.Email,
			},
			Spec: v1alpha1.GoogleGroupSpec{
				Name: g.Name,
				Email: g.Email,
				Description: g.Description,
				Members: []v1alpha1.Member{},
			},
		}

		// Update the group settings.
		gSettings, err := settingsService.Groups.Get(g.Email).Do()

		if err != nil {
			log.Error(err, "Error getting group settings", "group", g.Email)
			continue
		}

		newGroup.Spec.AllowExternalMembers = gSettings.AllowExternalMembers
		newGroup.Spec.WhoCanPostMessage = gSettings.WhoCanPostMessage
		newGroup.Spec.WhoCanJoin = gSettings.WhoCanJoin

		appendMembers := func(page *admin.Members) error {

			for _, m := range page.Members {
				newGroup.Spec.Members = append(newGroup.Spec.Members, v1alpha1.Member{
					Email: m.Email,
					Role: m.Role,
				})
			}
			return nil
		}

		err = service.Members.List(g.Email).Pages(context.Background(), appendMembers)

		if err != nil {
			log.Error(err, "Error getting group members", "group", g.Email)
			continue
		}

		results = append(results, newGroup)
	}
	return results, nil
}
