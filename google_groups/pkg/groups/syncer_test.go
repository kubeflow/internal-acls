package groups

import (
	"github.com/google/go-cmp/cmp"
	"github.com/kubeflow/internal-acls/google_groups/pkg/api/v1alpha1"
	admin "google.golang.org/api/admin/directory/v1"
	"testing"
)

func TestMembersDiff(t *testing.T) {

	type testCase struct {
		current []string
		desired []string
		expected memberDiff
	}

	testCases := []testCase {
		{
			current : []string {"a", "b", "c"},
			desired: []string{"b", "c", "d"},
			expected: memberDiff{
				ToAdd:    []v1alpha1.Member{
					{
						Email: "d",
					},
				},
				ToRemove: []string{"a"},
			},
		},
	}

	for _, c := range testCases {
		current := []*admin.Member{}

		for _, e := range c.current {
			current = append(current, &admin.Member{Email: e})
		}

		desired := []v1alpha1.Member{}

		for _, e := range c.desired {
			desired = append(desired, v1alpha1.Member{
				Email: e,
			})
		}

		memberDiff := diffCurrentDesiredMembers(current, desired)

		if diff := cmp.Diff(c.expected, memberDiff); diff != "" {
			t.Errorf("diffCurrentDesiredMembers() mismatch (-want +got):\n%s", diff)
		}
	}
}
