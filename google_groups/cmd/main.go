// Package main implements a tool to sync google groups
package main

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/gogo/protobuf/proto"
	"github.com/kubeflow/internal-acls/google_groups/pkg/api"
	"github.com/kubeflow/internal-acls/google_groups/pkg/api/v1alpha1"
	"github.com/kubeflow/internal-acls/google_groups/pkg/gcp"
	"github.com/kubeflow/internal-acls/google_groups/pkg/groups"
	"github.com/kubeflow/internal-acls/google_groups/pkg/util"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	admin "google.golang.org/api/admin/directory/v1"
	settingsSdk "google.golang.org/api/groupssettings/v1"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)


type RunOptions struct{
	Input string
	CredentialsFile string
	Secret string
	Continuous bool
	SyncPeriod time.Duration
	ForcedResyncPreiod time.Duration
}

type ImportOptions struct{
	Output string
	Domain string
}

var (
	opts = RunOptions{}
	iOpts = ImportOptions{}

	rootCmd    = &cobra.Command{}

	runCmd     = &cobra.Command{
		Use:   "run",
		Short: "Run a sync",
		Run: func(cmd *cobra.Command, args []string) {
			run()
		},
	}

	importCmd  = &cobra.Command{
		Use:   "import",
		Short: "Import existing group definitions to YAML files.",
		Run: func(cmd *cobra.Command, args []string) {
			runImport()
		},
	}

	upgradeCmd     = &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade/Fix all groups. Bulk applies one or more transformations to all groups.",
		Run: func(cmd *cobra.Command, args []string) {
			upgrade()
		},
	}

	log logr.Logger

	scopes = []string {
		admin.AdminDirectoryGroupMemberScope,
		admin.AdminDirectoryGroupScope,
		admin.CloudPlatformScope,
		settingsSdk.AppsGroupsSettingsScope,
	}
)

func init() {
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(upgradeCmd)
	rootCmd.AddCommand(importCmd)

	upgradeCmd.Flags().StringVarP(&opts.Input, "input", "", "", "A glob to match config files to upgrade.")
	upgradeCmd.Flags().StringVarP(&iOpts.Output, "output", "", "", "The directory to write the Group specs to")
	upgradeCmd.MarkFlagRequired("input")
	upgradeCmd.MarkFlagRequired("output")

	runCmd.Flags().StringVarP(&opts.Input, "input", "", "", "A glob to match config files to apply.")

	runCmd.Flags().StringVarP(&opts.CredentialsFile, "credentials-file", "", "", "JSON File containing OAuth2Client credentials as downloaded from APIConsole. Can be a GCS file.")
	runCmd.Flags().StringVarP(&opts.Secret, "secret", "", "", "The name of a secret in GCP secret manager where the OAuth2 token should be cached. Should be in the form {project}/{secret}")
	runCmd.Flags().BoolVarP(&opts.Continuous, "continuous", "", false, "If true runs forever; resyncing the groups whenever a change is detected")
	runCmd.Flags().DurationVarP(&opts.SyncPeriod, "sync-period", "", 30 * time.Second, "How often to check for changes. This should be O(seconds)")
	runCmd.Flags().DurationVarP(&opts.ForcedResyncPreiod, "forced-sync-period", "", 4 * time.Hour, "How often to resync even when no changes have been detected. Should be on the order of hours")

	importCmd.Flags().StringVarP(&opts.CredentialsFile, "credentials-file", "", "", "JSON File containing OAuth2Client credentials as downloaded from APIConsole.")
	importCmd.Flags().StringVarP(&iOpts.Domain, "domain", "", "kubeflow.org", "The domain containing the Google groups to import")
	importCmd.Flags().StringVarP(&iOpts.Output, "output", "", "", "The directory to write the results to")
	importCmd.MarkFlagRequired("input")
	importCmd.MarkFlagRequired("output")
	importCmd.MarkFlagRequired("domain")
}

func initLogger() {
	// TODO(jlewi): Make the verbosity level configurable.

	// Start with a production logger config.
	config := zap.NewProductionConfig()

	// TODO(jlewi): In development mode we should use the console encoder as opposed to json formatted logs.

	// Increment the logging level.
	// TODO(jlewi): Make this a flag.
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	zapLog, err := config.Build()
	if err != nil {
		panic(fmt.Sprintf("Could not create zap instance (%v)?", err))
	}
	log = zapr.NewLogger(zapLog)

	zap.ReplaceGlobals(zapLog)
}

func getWebFlowLocal() *gcp.CachedCredentialHelper {
	webFlow, err := gcp.NewWebFlowHelper(opts.CredentialsFile, scopes)

	if err != nil {
		log.Error(err, "Failed to create a WebFlowHelper credential helper")
		return nil
	}

	home, err := os.UserHomeDir()

	if err != nil {
		log.Error(err,"Could not get home directory")
		return nil
	}

	cacheFile := filepath.Join(home, ".cache", "kubeflow", "groups.sync.token")
	h := &gcp.CachedCredentialHelper {
		CredentialHelper: webFlow,
		TokenCache: &gcp.FileTokenCache{
			CacheFile: cacheFile,
			Log:       log,
		},
		Log: log,
	}

	return h
}

func getWebFlowSecretManager() *gcp.CachedCredentialHelper {
	webFlow, err := gcp.NewWebFlowHelper(opts.CredentialsFile, scopes)

	if err != nil {
		log.Error(err, "Failed to create a WebFlowHelper credential helper")
		return nil
	}

	pieces := strings.Split(opts.Secret, "/")

	if len(pieces) != 2 {
		log.Error(fmt.Errorf("Secret %v not in form {project}/{secret}", opts.Secret), "Incorrectly specified secret", "secret", opts.Secret)
		return nil
	}

	cache, err := gcp.NewSecretCache(pieces[0], pieces[1], "latest")

	if err != nil {
		log.Error(err, "Could not create cache for secret manager")
		return nil
	}

	cache.Log = log

	h := &gcp.CachedCredentialHelper {
		CredentialHelper: webFlow,
		TokenCache: cache,
		Log: log,
	}

	return h
}

// getAdminClient initializes an admin client using a local credential cahce
func getAdminClient(h gcp.CredentialHelper) *http.Client{
	ctx := context.Background()
	ts, err := h.GetTokenSource(ctx)

	if err != nil {
		log.Error(err, "Failed to obtain token source")
		return nil
	}

	client := oauth2.NewClient(ctx, ts)
	return client
}

func run() {
	initLogger()

	var credsHelper gcp.CredentialHelper

	if opts.Secret != "" {
		log.Info("Getting OAuth2 credential via webflow and secret manager")
		credsHelper = getWebFlowSecretManager()
	} else {
		log.Info("Getting OAuth2 credential via webflow and local cache")
		credsHelper = getWebFlowLocal()
	}

	if credsHelper == nil {
		return
	}

	client := getAdminClient(credsHelper)

	if client == nil {
		return
	}

	s := &groups.GroupSyncer{
		Client: client,
		Log: log,
	}


	runSync := func () error {
		defs := api.ReadGroups(opts.Input)

		if len(defs) == 0 {
			log.Info("No groups matched glob", "glob", opts.Input)
			return nil
		}

		err := s.Sync(defs)

		if err != nil {
			log.Error(err, "Failed to sync")
		}

		return err
	}

	lastHash := ""

	// Set the resync time in the past to force a resync immediately
	nextResyncTime := time.Now().Add(-10 *time.Minute)

	for ;; {
		// Get the current content hash so we can see if its changed.
		log.Info("Reading glob", "directory", opts.Input)
		matches, err := filepath.Glob(opts.Input)
		if err != nil {
			log.Error(err, "Error matching glob path", "glob", opts.Input)
			return
		}

		log.Info("Found files", "files", matches)

		hashBytes, err := util.ContentHash(matches)

		if err != nil {
			log.Error(err, "Could not hash the file contents")
			return
		}

		newHash := string(hashBytes)

		if newHash != lastHash || time.Now().After(nextResyncTime) {
			log.Info("Sync needed", "lastHash", lastHash, "newHash", newHash, "nextResyncTime", nextResyncTime)
			err := runSync()

			if err == nil {
				// Only update the hash if the sync succeeded otherwise we want to try again.
				// ToDO(jlewi): Should we do exponential backoff in the event of errors.
				lastHash = newHash
				nextResyncTime = time.Now().Add(opts.ForcedResyncPreiod)

				log.Info("Updated content hash and resync time", "lasthash", lastHash, "nextResyncTime", nextResyncTime)
			}
		} else {
			log.Info("No sync needed", "lastHash", lastHash, "newHash", newHash, "nextResyncTime", nextResyncTime)
		}

		if !opts.Continuous {
			return
		}

		time.Sleep(opts.SyncPeriod)
	}
}

func runImport() {
	initLogger()
	log.Info("Getting OAuth2 credential via webflow and local cache")
	credsHelper := getWebFlowLocal()

	if credsHelper == nil {
		return
	}

	client := getAdminClient(credsHelper)

	if client == nil {
		return
	}

	s := &groups.GroupImporter{
		Client: client,
		Log: log,
	}

	groups, err := s.Import(iOpts.Domain)

	if err != nil {
		log.Error(err, "Failed to import group specs")
		return
	}

	// This is a list of groups that are already in GitHub.
	// We don't want to import members for groups that aren't already in GitHub
	// because don't want to put emails into GitHub without people's consent.
	// So they should be the ones to open the PR adding themselves.
	// This should go away once we've completely migrated to using GitOps.
	groupsToImportMembersFor := map[string]bool {
		"calendar-admins@kubeflow.org": true,
		"ci-team@kubeflow.org": true,
		"ci-viewer@kubeflow.org": true,
		"code-search-team@kubeflow.org": true,
		"community-meeting-hosts@kubeflow.org":true,
		"devrel-team@kubeflow.org": true,
		"devstats@kubeflow.org": true,
		"drive-content-managers@kubeflow.org": true,
		"example-maintainers@kubeflow.org": true,
		"feast-team@kubeflow.org": true,
		"github-team@kubeflow.org": true,
		"google-codelab-projects-owners@kubeflow.org": true,
		"kf-demo-owners@kubeflow.org": true,
		"kf-kcc-admins@kubeflow.org": true,
		"release-team@kubeflow.org": true,
	}

	for _, g := range groups {
		if _, ok := groupsToImportMembersFor[g.Spec.Email]; !ok {
			log.Info("Removing members from group not in hardcoded allow list", "group", g.Spec.Email)
			g.Spec.Members = []v1alpha1.Member{}

			log.Info("Disabling autosync", "group", g.Spec.Email)
			g.Spec.AutoSync = proto.Bool(false)
		}
	}

	err = api.WriteGroups(groups, iOpts.Output)

	if err != nil {
		log.Error(err, "Error writing groups")
	}
}

func upgrade() {
	grps := api.ReadGroups(opts.Input)

	if len(grps) == 0 {
		log.Info("No groups matched glob", "glob", opts.Input)
		return
	}

	requiredUsers := []*v1alpha1.Member{}
	removeUsers := map[string]bool{
		"kf-autobot@kf-infra-gitops.iam.gserviceaccount.com": true,
	}

	err := api.Upgrade(grps, requiredUsers, removeUsers)

	if err != nil {
		log.Error(err, "Failed to upgrade specs")
		return
	}

	err = api.WriteGroups(grps, iOpts.Output)

	if err != nil {
		log.Error(err, "Failed to write specs", "output", iOpts.Output)
		return
	}
}

func main() {
	rootCmd.Execute()
}
