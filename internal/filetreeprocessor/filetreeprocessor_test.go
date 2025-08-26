package filetreeprocessor

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/apex/log"
	"github.com/gchiesa/ska/pkg/templateprovider"
	"github.com/gruntwork-io/terratest/modules/git"
	compare "github.com/kilianpaquier/compare/pkg"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
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
		name              string
		fixtureName       string
		sourceIgnorePaths []string
	}{
		{"build staging file tree -> single-file-flat-dir", "single-file-flat-dir", []string{}},
		{"build staging file tree -> multi-dir-multi-file-w-empty-file-names", "multi-dir-multi-file-w-empty-file-names", []string{}},
		{"build staging file tree -> multi-dir-with-ignore-all-except-dirA", "multi-dir-with-ignore-all-except-dirA", []string{"*", "!dir-a/", "!dir-a/**"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fixturePath, err := provisionScenario(t, tc.fixtureName)
			assert.NoError(t, err)
			tempWorkindDirFolder := t.TempDir()
			tempDestinationFolder := t.TempDir()
			tp := &FileTreeProcessor{
				sourcePath:          filepath.Join(fixturePath, blueprintFolder),
				destinationPathRoot: tempDestinationFolder,
				workingDir:          tempWorkindDirFolder,
				template:            tplSvc,
				log:                 logger,
				sourceIgnorePaths:   tc.sourceIgnorePaths,
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
		name              string
		fixtureName       string
		sourceIgnorePaths []string
	}{
		{"rendering staging files -> single-file-flat-dir", "single-file-flat-dir", []string{}},
		{"rendering staging files -> multi-dir-multi-file-w-empty-file-names", "multi-dir-multi-file-w-empty-file-names", []string{}},
		{"rendering staging files -> multi-dir-with-ignore-all-except-dirA", "multi-dir-with-ignore-all-except-dirA", []string{"*", "!dir-a/", "!dir-a/**"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fixturePath, err := provisionScenario(t, tc.fixtureName)
			assert.NoError(t, err)
			tempWorkindDirFolder := t.TempDir()
			tempDestinationFolder := t.TempDir()
			tp := &FileTreeProcessor{
				sourcePath:          filepath.Join(fixturePath, blueprintFolder),
				destinationPathRoot: tempDestinationFolder,
				workingDir:          tempWorkindDirFolder,
				template:            tplSvc,
				log:                 logger,
				sourceIgnorePaths:   tc.sourceIgnorePaths,
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

// YAML-driven fixture provisioning

type fileNode struct {
	Content string              `yaml:"content,omitempty"`
	Files   map[string]fileNode `yaml:",inline"`
}

type fixtureYAML struct {
	Variables map[string]interface{}         `yaml:"variables"`
	Blueprint map[string]fileNode            `yaml:"blueprint"`
	Expected  map[string]map[string]fileNode `yaml:"expected"`
}

type scenariosYAML struct {
	Scenarios map[string]fixtureYAML `yaml:"scenarios"`
}

// writeTree writes a map-based tree into basePath
func writeTree(basePath string, tree map[string]fileNode) error {
	for name, node := range tree {
		p := filepath.Join(basePath, name)
		// if node has nested files, treat it as dir
		if len(node.Files) > 0 {
			if err := os.MkdirAll(p, 0o755); err != nil {
				return err
			}
			if err := writeTree(p, node.Files); err != nil {
				return err
			}
			continue
		}
		// it's a file (can be empty content)
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(p, []byte(node.Content), 0o644); err != nil {
			return err
		}
	}
	return nil
}

// provisionScenario loads internal/filetreeprocessor/fixtures/scenarios.yaml, picks the named scenario,
// and materializes it into a temporary directory with blueprint, variables.json and expected folders.
func provisionScenario(t *testing.T, scenarioName string) (string, error) {
	root := filepath.Join(git.GetRepoRoot(t), "internal", "filetreeprocessor", "fixtures")
	data, err := os.ReadFile(filepath.Join(root, "scenarios.yaml"))
	if err != nil {
		return "", err
	}
	var sc scenariosYAML
	if err := yaml.Unmarshal(data, &sc); err != nil {
		return "", err
	}
	fx, ok := sc.Scenarios[scenarioName]
	if !ok {
		return "", errors.New("scenario not found: " + scenarioName)
	}
	base := t.TempDir()
	// create blueprint
	if err := os.MkdirAll(filepath.Join(base, blueprintFolder), 0o755); err != nil {
		return "", err
	}
	if err := writeTree(filepath.Join(base, blueprintFolder), fx.Blueprint); err != nil {
		return "", err
	}
	// write variables.json
	vdata, _ := json.Marshal(fx.Variables)
	if err := os.WriteFile(filepath.Join(base, variablesFile), vdata, 0o644); err != nil {
		return "", err
	}
	// expected directories
	for key, tree := range fx.Expected {
		dirPath := filepath.Join(base, key)
		if err := os.MkdirAll(dirPath, 0o755); err != nil {
			return "", err
		}
		if err := writeTree(dirPath, tree); err != nil {
			return "", err
		}
	}
	return base, nil
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

func getFixtureBasePath(t *testing.T, fixtureName string) string {
	return filepath.Join(git.GetRepoRoot(t), "tests", "fixtures", fixtureName)
}
