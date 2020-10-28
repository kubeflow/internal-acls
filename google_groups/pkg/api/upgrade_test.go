package api

import (
	"github.com/google/go-cmp/cmp"
	"github.com/kubeflow/internal-acls/google_groups/pkg/api/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestUpgrade(t *testing.T) {
	type testCase struct {
		filename string
		contents string
		expected v1alpha1.GoogleGroup
	}

	cases := []testCase {
		{
			filename: "calendar-admins.members.txt",
			contents: `

abe@acme.com
joe@gmail.net

`,
			expected: v1alpha1.GoogleGroup{
				ObjectMeta: v1.ObjectMeta{
					Name: "calendar-admins@kubeflow.org",
				},
				Spec:       v1alpha1.GoogleGroupSpec{
					Members: []v1alpha1.Member{
						{
							Principal: v1alpha1.Principal{
								User:"abe@acme.com",
							},
						},
						{
							Principal: v1alpha1.Principal{
								User:"joe@gmail.net",
							},
						},
					},
				},
			},
		},
	}

	for _, c := range cases {
		actual, err := ConvertTextToGroup(c.filename, c.contents)

		if err != nil {
			t.Errorf("Failed to convert text %v", err)
			continue
		}

		if diff := cmp.Diff(c.expected, *actual); diff != "" {
			t.Errorf("ConvertTextToGroup() mismatch (-want +got):\n%s", diff)
		}
	}
}
