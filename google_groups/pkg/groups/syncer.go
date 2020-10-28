package groups

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/kubeflow/internal-acls/google_groups/pkg/api/v1alpha1"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/googleapi"
	settingsSdk "google.golang.org/api/groupssettings/v1"
	"net/http"
	"strings"
)

type GroupSyncer struct {
	Client *http.Client
	Log logr.Logger
}

func (s *GroupSyncer) Sync(groupSpecs []*v1alpha1.GoogleGroup) error {
	log := s.Log
	// TODO(jlewi): Using admin.NewService was giving me auth problems. I think it might be an OAuthScope
	// issue because the credential didn't have all the scopes but admin.NewService appears to request all scopes
	// so the client was requesting a token for scopes it wasn't authorized for.
	service, err := admin.New(s.Client)

	if err != nil {
		return err
	}

	settingsService, err := settingsSdk.New(s.Client)

	if err != nil {
		return err
	}

	for _, gDef := range groupSpecs {
		if gDef.Spec.AutoSync != nil && !*gDef.Spec.AutoSync {
			log.Info("AutoSync not enabled for group", "group", gDef.Spec.Email)
			continue
		}

		// Ensure each group exists and settings are up to date
		s.syncGroupSettings(gDef, service, settingsService)

		// Sync members
		s.syncMembers(gDef, service)
	}

	// Sync members
	return nil
}

// syncGroupSettings synchronizes a groups setting (but not the membership).
// This includes;
// * creating the group if it doesn't exist
// * setting description and properties on the group
func (s *GroupSyncer) syncGroupSettings(gDef *v1alpha1.GoogleGroup,service *admin.Service, settingsService *settingsSdk.Service) error {
	_, err := service.Groups.Get(gDef.Spec.Email).Do()

	log := s.Log
	if err != nil {
		gErr, ok := err.(*googleapi.Error)

		if !ok {
			log.Error(err, "Error getting group.", "group", gDef.Spec.Email)
			return err
		}

		if gErr.Code != http.StatusNotFound {
			log.Error(err, "Error getting group.", "group", gDef.Spec.Email)
			return err
		}

		// Group doesn't exist so create it
		log.Info("Creating group", "group", gDef.Spec.Email)
		// TODO(jlewi): How do we set who can join? How do we allow external members?
		pieces := strings.Split(gDef.Spec.Email, "@")
		newGroup := &admin.Group {
			Email: gDef.Spec.Email,
			Name: pieces[0],
			Description: gDef.Spec.Description,
		}
		_, err := service.Groups.Insert(newGroup).Do()

		if err != nil {
			log.Error(err, "Error creating group.", "group", gDef.Spec.Email)
			return err
		}
	}

	// Update the group settings.
	gSettings, err := settingsService.Groups.Get(gDef.Spec.Email).Do()

	if err != nil {
		s.Log.Error(err, "Error getting group settings", "group", gDef.Spec.Email)
		return err
	}

	// Ref: https://developers.google.com/admin-sdk/groups-settings/v1/reference/groups#json
	gSettings.Description = gDef.Spec.Description

	// The only mechanism for joining these groups should be via the configs and the sync
	// program
	gSettings.WhoCanJoin = gDef.Spec.WhoCanJoin
	// Most members will be joining with their non kubeflow accounts
	gSettings.AllowExternalMembers = gSettings.AllowExternalMembers

	// TODO(jlewi): How should we validate its set to a correct value?
	gSettings.WhoCanPostMessage = gDef.Spec.WhoCanPostMessage

	if gSettings.WhoCanPostMessage == "" {
		gSettings.WhoCanPostMessage = "ANYONE_CAN_POST"
	}

	log.Info("Updating group settings", "group", gDef.Spec.Email)
	_, err = settingsService.Groups.Update(gDef.Spec.Email, gSettings).Do()
	if err != nil {
		s.Log.Error(err, "Error updating group settings", "group", gDef.Spec.Email)
	}

	return nil
}

func (s *GroupSyncer) syncMembers(gDef *v1alpha1.GoogleGroup, service *admin.Service) error {
	log := s.Log
	currentMembers := []*admin.Member{}

	appendMembers := func(page *admin.Members) error {
		currentMembers = append(currentMembers, page.Members...)
		return nil
	}

	err := service.Members.List(gDef.Spec.Email).Pages(context.Background(), appendMembers)

	if err != nil {
		log.Error(err, "Error getting group members", "group", gDef.Spec.Email)
		return err
	}

	diff := diffCurrentDesiredMembers(currentMembers, gDef.Spec.Members)

	log.Info("Diff Group Membership", "group", gDef.Spec.Email, "diff", diff)

	// Add missing members
	for _, m := range diff.ToAdd {
		if !isValidGroupRole(GroupRole(m.Role)) {
			log.Error(fmt.Errorf("Member has invalid role"), "Member has invalid role", "group", gDef.Spec.Email, "member", m)
			continue
		}
		newMember := admin.Member{
			Email: m.Email,
			Role: m.Role,
		}
		result, err := service.Members.Insert(gDef.Spec.Email, &newMember).Do()

		if err != nil {
			log.Error(err, "Could not insert member", "group", gDef.Spec.Email, "member", newMember)
		} else {
			log.Info( "Inserted member", "group", gDef.Spec.Email, "member", result)
		}
	}

	// Delete removed members
	for _, m := range diff.ToRemove {
		err := service.Members.Delete(gDef.Spec.Email, m).Do()

		if err != nil {
			log.Error(err, "Could not delete member", "group", gDef.Spec.Email, "member", m)
		} else {
			log.Info( "Delete member", "group", gDef.Spec.Email, "member", m)
		}
	}
	return nil
}

// TODO(jlewi): We should also figure out if there are any members that need to be updated e.g. have their
// role changed.
type memberDiff struct{
	// List of members that are missing from the group
	ToAdd []v1alpha1.Member

	// List of members to remove from the group
	ToRemove []string
}

func diffCurrentDesiredMembers(current []*admin.Member, desired []v1alpha1.Member) memberDiff {
	cSet := map[string] bool {}

	for  _, m := range current {
		cSet[m.Email] = true
	}

	diff := memberDiff{
		ToAdd:  []v1alpha1.Member{},
		ToRemove: []string{},
	}

	dSet := map[string] bool {}

	// generate missing members
	for _, m := range desired {
		dSet[m.Email] = true

		if _, ok := cSet[m.Email]; !ok {
			diff.ToAdd = append(diff.ToAdd, m)
		}
	}

	// generate members to delete
	for _, m := range current {
		if _, ok := dSet[m.Email]; !ok {
			diff.ToRemove = append(diff.ToRemove, m.Email)
		}
	}

	return diff
}

func isValidGroupRole(role GroupRole) bool {
	for _, v := range []GroupRole{OwnerRole, ManagerRole, MemberRole} {
		if v == role {
			return true
		}
	}
	return false
}