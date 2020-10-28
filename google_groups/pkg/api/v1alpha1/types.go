package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// GoogleGroup defines a google group.
type GoogleGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GoogleGroupSpec   `json:"spec,omitempty"`
}

type GoogleGroupSpec struct {
	// Whether to AutoSync the group
	AutoSync *bool `json:"autoSync,omitempty"`
	Email string `json:"email,omitempty"`
	Name string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	// https://developers.google.com/admin-sdk/groups-settings/v1/reference/groups#json
	WhoCanPostMessage string `json:"whoCanPostMessage",omitempty`
	WhoCanJoin string `json:"whoCanJoin",omitempty`
	AllowExternalMembers string `json:"allowExternalMembers",omitempty`
	Members []Member `json:"members,omitempty"`
}

type Member struct {
	// Principal is the identity of the member
	Email string `json:"email,omitempty"`

	// Role of the member
	// see https://developers.google.com/admin-sdk/directory/v1/reference/members/insert
	Role string `json:"role,omitempty"`
}
