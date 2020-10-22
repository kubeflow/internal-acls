package api

import (
	"github.com/ghodss/yaml"
	"github.com/go-logr/zapr"
	"github.com/kubeflow/internal-acls/google_groups/pkg/api/v1alpha1"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// ReadGroups reads in all group specs from a directory
func ReadGroups(inputGlob string) ([]*v1alpha1.GoogleGroup) {
	log := zapr.NewLogger(zap.L())
	results := []*v1alpha1.GoogleGroup{}
	log.Info("Reading glob", "directory", inputGlob)
	matches, err := filepath.Glob(inputGlob)
	if err != nil {
		log.Error(err, "Error matching glob path", "glob", inputGlob)
		return results
	}

	for _, f := range matches {
		log.Info("Reading file", "input", f)
		b, err := ioutil.ReadFile(f)

		if err != nil {
			log.Error(err, "Error reading file.", "file", f)
			continue
		}

		g := &v1alpha1.GoogleGroup{}
		err = yaml.Unmarshal(b, g)

		if err != nil {
			log.Error(err, "Error parsing GoogleGroup from file.", "file", f)
			continue
		}
		results = append(results, g)
	}

	return results
}

func ensureDirExists(dir string) error {
	log := zapr.NewLogger(zap.L())
	_, err := os.Stat(dir)

	if err != nil {
		if os.IsNotExist(err) {
			log.Info("Create cache directory", "dir", dir)
			err := os.MkdirAll(dir, 0770)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

// WriterGroups serializes the groups as YAML to the specified directory
func WriteGroups(groups []*v1alpha1.GoogleGroup, output string) error {
	log := zapr.NewLogger(zap.L())

	err := ensureDirExists(output)

	if err != nil {
		log.Error(err, "Could not ensure output directory exists", "output", output)
		return err
	}


	for _, g := range groups {
		gBytes, err := yaml.Marshal(g)

		if err != nil {
			log.Error(err, "Error marshling group", "group", g)
			continue
		}

		fileName := strings.Split(g.Spec.Email, "@")[0]
		yamlFile := filepath.Join(output, fileName+".yaml")

		err = ioutil.WriteFile(yamlFile, gBytes, 0644)

		if err != nil {
			log.Error(err, "Error writing file", "target", yamlFile)
			continue
		}

		log.Info("Converted group file.", "output", yamlFile)
	}
	return nil
}