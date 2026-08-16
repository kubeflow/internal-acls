package groups

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/go-logr/zapr"
	"github.com/kubeflow/internal-acls/google_groups/pkg/api/v1alpha1"
	"go.uber.org/zap"
	settingsSdk "google.golang.org/api/groupssettings/v1"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestSyncAppliesAllowExternalMembers(t *testing.T) {
	testCases := []struct {
		name    string
		current string
		desired string
	}{
		{
			name:    "enable external members",
			current: "false",
			desired: "true",
		},
		{
			name:    "disable external members",
			current: "true",
			desired: "false",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var updated settingsSdk.Groups
			client := &http.Client{Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
				body := "{}"
				if strings.Contains(req.URL.Path, "/groups/v1/") && req.Method == http.MethodGet {
					body = `{"allowExternalMembers":"` + tc.current + `"}`
				}
				if strings.Contains(req.URL.Path, "/groups/v1/") && req.Method == http.MethodPut {
					b, err := ioutil.ReadAll(req.Body)
					if err != nil {
						t.Fatal(err)
					}
					if err := json.Unmarshal(b, &updated); err != nil {
						t.Fatal(err)
					}
				}

				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     http.Header{"Content-Type": []string{"application/json"}},
					Body:       ioutil.NopCloser(bytes.NewBufferString(body)),
					Request:    req,
				}, nil
			})}

			syncer := &GroupSyncer{
				Client: client,
				Log:    zapr.NewLogger(zap.NewNop()),
			}
			err := syncer.Sync([]*v1alpha1.GoogleGroup{{
				Spec: v1alpha1.GoogleGroupSpec{
					Email:                "example@kubeflow.org",
					AllowExternalMembers: tc.desired,
				},
			}})
			if err != nil {
				t.Fatal(err)
			}

			if updated.AllowExternalMembers != tc.desired {
				t.Errorf("allowExternalMembers sent to Groups Settings API = %q, want %q", updated.AllowExternalMembers, tc.desired)
			}
		})
	}
}
