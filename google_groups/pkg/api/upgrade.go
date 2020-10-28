package api

import (
	"github.com/go-logr/zapr"
	"github.com/kubeflow/internal-acls/google_groups/pkg/api/v1alpha1"
	"go.uber.org/zap"
	"sort"
)


const (
	expectedSuffix = ".members.txt"
)

// Upgrade all specs. addUsers a list of users to ensure are present in all groups
func Upgrade(groups []*v1alpha1.GoogleGroup, addUsers []*v1alpha1.Member, removeUsers map[string]bool) error {
	log := zapr.NewLogger(zap.L())

	initMissing := func() map[string]*v1alpha1.Member{
		r := map[string]*v1alpha1.Member {}
		for _, a := range addUsers{
			r[a.Email] = a
		}
		return r
	}

	for _, g := range groups {
		missing := initMissing()

		toKeep := []v1alpha1.Member{}
		for _, m := range g.Spec.Members {
			if _, ok := removeUsers[m.Email]; ok {
				log.Info("Removing member", "group" ,g.Spec.Email , "member", m.Email)
				continue
			}

			toKeep = append(toKeep, m)

			desired, ok := missing[m.Email]

			if !ok {
				continue
			}

			log.Info("Group already has member", "group",g.Spec.Email , "member", m.Email)

			m.Role = desired.Role

			delete(missing, m.Email)
		}

		g.Spec.Members = toKeep
		for _, m := range missing {
			log.Info("Adding member to group", "group", g.Spec.Email, "member", m.Email)

			g.Spec.Members = append(g.Spec.Members, *m)
		}

		// Sort the members
		sort.Slice(g.Spec.Members[:], func(i, j int) bool {
			return g.Spec.Members[i].Email < g.Spec.Members[j].Email
		})
	}

	return nil
}
