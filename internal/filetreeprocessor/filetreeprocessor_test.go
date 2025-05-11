package filetreeprocessor

import (
	"encoding/json"
	"github.com/apex/log"
	"github.com/gchiesa/ska/pkg/templateprovider"
	"github.com/gruntwork-io/terratest/modules/git"
	compare "github.com/kilianpaquier/compare/pkg"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

const (
	blueprintFolder   = "blueprint"
	resultFolder      = "update-destination-file-tree"
	loadMultipartsDir = "load-multiparts"
	stagingTreeFolder = "render-staging-tree"
	workingDirFolder  = "build-staging-file-tree"
	variablesFile     = "variables.json"
)

func TestBuildStagingFileTree(t *testing.T) {
	t.Parallel()
	tplSvc := templateprovider.ByType(templateprovider.SprigTemplateType, "sprig")
	logger := log.WithFields(log.Fields{})

	testCases := []struct {
		name        string
		fixtureName string
	}{
		{"build staging file tree -> single-file-flat-dir", "single-file-flat-dir"},
		{"build staging file tree -> multi-dir-multi-file-w-empty-file-names", "multi-dir-multi-file-w-empty-file-names"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fixturePath := getFixtureBasePath(t, tc.fixtureName)
			tempWorkindDirFolder := t.TempDir()
			tempDestinationFolder := t.TempDir()
			tp := &FileTreeProcessor{
				sourcePath:          filepath.Join(fixturePath, blueprintFolder),
				destinationPathRoot: tempDestinationFolder,
				workingDir:          tempWorkindDirFolder,
				template:            tplSvc,
				log:                 logger,
			}
			variables, err := loadVariables(t, fixturePath, variablesFile)
			assert.NoError(t, err)
			err = tp.buildStagingFileTree(variables)
			assert.NoError(t, err)
			// compare the fixture output with the produced output
			err = compare.Dirs(filepath.Join(fixturePath, workingDirFolder), tp.WorkingDir())
			assert.NoError(t, err)
		})
	}
}

func TestRenderStagingFiles(t *testing.T) {
	t.Parallel()
	tplSvc := templateprovider.ByType(templateprovider.SprigTemplateType, "sprig")
	logger := log.WithFields(log.Fields{})

	testCases := []struct {
		name        string
		fixtureName string
	}{
		{"rendering staging files -> single-file-flat-dir", "single-file-flat-dir"},
		{"rendering staging files -> multi-dir-multi-file-w-empty-file-names", "multi-dir-multi-file-w-empty-file-names"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fixturePath := getFixtureBasePath(t, tc.fixtureName)
			tempWorkindDirFolder := t.TempDir()
			tempDestinationFolder := t.TempDir()
			tp := &FileTreeProcessor{
				sourcePath:          filepath.Join(fixturePath, blueprintFolder),
				destinationPathRoot: tempDestinationFolder,
				workingDir:          tempWorkindDirFolder,
				template:            tplSvc,
				log:                 logger,
			}
			variables, err := loadVariables(t, fixturePath, variablesFile)
			assert.NoError(t, err)
			err = tp.buildStagingFileTree(variables)
			assert.NoError(t, err)
			assert.NoError(t, tp.loadMultiparts())

			err = tp.renderStagingFileTree(variables)
			assert.NoError(t, err)
			// compare the fixture output with the produced output
			err = compare.Dirs(filepath.Join(fixturePath, stagingTreeFolder), tp.WorkingDir())
			assert.NoError(t, err)
		})
	}
}

func getFixtureBasePath(t *testing.T, fixtureName string) string {
	return filepath.Join(git.GetRepoRoot(t), "tests", "fixtures", fixtureName)
}

func loadVariables(t *testing.T, path, variablesFile string) (map[string]interface{}, error) {
	fileData, err := os.ReadFile(filepath.Join(path, variablesFile))
	if err != nil {
		return nil, err
	}
	variables := make(map[string]interface{})
	if err := json.Unmarshal(fileData, &variables); err != nil {
		return nil, err
	}
	return variables, nil
}
