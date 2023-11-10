package alias_test

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/thalesfsp/mole/alias"
)

const (
	FixtureDir = "./testdata"
)

var home string

func TestAddThenGetThenDeleteAlias(t *testing.T) {
	expectedAlias, err := addAlias()
	if err != nil {
		t.Errorf("error creating alias file %v", err)
	}

	expectedAliasFilePath := filepath.Join(home, ".mole", fmt.Sprintf("%s.toml", expectedAlias.Name))

	if _, err := os.Stat(expectedAliasFilePath); os.IsNotExist(err) {
		t.Errorf("alias file could not be found after the attempt to create it")
	}

	al, err := alias.Get(expectedAlias.Name)
	if err != nil {
		t.Errorf("%v", err)
	}

	if !reflect.DeepEqual(expectedAlias, al) {
		t.Errorf("expected: %s, actual: %s", expectedAlias, al)
	}

	err = alias.Delete(expectedAlias.Name)
	if err != nil {
		t.Errorf("error while deleting %s alias file: %v", expectedAlias.Name, err)
	}

	if _, err := os.Stat(expectedAliasFilePath); !os.IsNotExist(err) {
		t.Errorf("alias file found after the attempt to delete it")
	}

}

func TestShow(t *testing.T) {
	ids := []string{"test-env"}
	fx, err := filepath.Abs(FixtureDir)
	if err != nil {
		t.Errorf("error while loading data for TestShow: %v", err)
	}

	for _, id := range ids {
		fixturePath := filepath.Join(fx, fmt.Sprintf("show.alias.%s.fixture", id))
		expectedBytes, err := ioutil.ReadFile(fixturePath)
		if err != nil {
			t.Errorf("error while loading data for TestShow: %v", err)
		}

		expected := string(expectedBytes)

		output, err := alias.Show(id)
		if err != nil {
			t.Errorf("error showing alias %s: %v", id, err)
		}

		if output != expected {
			t.Errorf("output doesn't match. Failing the test.")
		}
	}
}

func TestShowAll(t *testing.T) {
	fx, err := filepath.Abs(FixtureDir)
	if err != nil {
		t.Errorf("error while loading data for TestShow: %v", err)
	}

	fixturePath := filepath.Join(fx, "show.alias.fixture")
	expectedBytes, err := ioutil.ReadFile(fixturePath)
	if err != nil {
		t.Errorf("error while loading data for TestShow: %v", err)
	}

	expected := string(expectedBytes)

	output, err := alias.ShowAll()
	if output != expected {
		t.Errorf("output doesn't match. Failing the test.")
	}
}

func TestMain(m *testing.M) {
	home, err := setup()
	if err != nil {
		fmt.Printf("error while loading data for TestShow: %v", err)
		os.RemoveAll(home)
		os.Exit(1)
	}

	code := m.Run()

	os.RemoveAll(home)

	os.Exit(code)
}

func addAlias() (*alias.Alias, error) {
	a := &alias.Alias{
		Name:              "alias",
		TunnelType:        "local",
		Verbose:           true,
		Insecure:          true,
		Detach:            true,
		Source:            []string{":1234"},
		Destination:       []string{"192.168.1.1:80"},
		Server:            "server.com",
		Key:               "path/to/key",
		KeepAliveInterval: "5s",
		ConnectionRetries: 3,
		WaitAndRetry:      "10s",
		SshAgent:          "path/to/agent",
		Timeout:           "1m",
		SshConfig:         "/home/user/.ssh/config",
	}

	err := alias.Add(a)
	if err != nil {
		return nil, err
	}

	return a, nil
}

// setup prepares the system environment to run the tests by:
// 1. Create temp dir and <dir>/.mole
// 2. Copy fixtures to <dir>/.mole
// 3. Set temp dir as the user testDir dir
func setup() (string, error) {
	testDir, err := ioutil.TempDir("", "mole-alias")
	if err != nil {
		return "", fmt.Errorf("error while setting up tests: %v", err)
	}

	moleAliasDir := filepath.Join(testDir, ".mole")
	err = os.Mkdir(moleAliasDir, 0755)
	if err != nil {
		return "", fmt.Errorf("error while setting up tests: %v", err)
	}

	err = os.Setenv("HOME", testDir)
	if err != nil {
		return "", fmt.Errorf("error while setting up tests: %v", err)
	}

	err = os.Setenv("USERPROFILE", testDir)
	if err != nil {
		return "", fmt.Errorf("error while setting up tests: %v", err)
	}

	fx, err := filepath.Abs(FixtureDir)
	if err != nil {
		return "", err
	}

	fixtures := []string{"test-env.toml", "example.toml"}
	for _, fixture := range fixtures {
		err = CopyFile(filepath.Join(fx, fixture), filepath.Join(moleAliasDir, fixture))
		if err != nil {
			return "", err
		}
	}

	home = testDir

	return moleAliasDir, nil
}

// GoLang: os.Rename() give error "invalid cross-device link" for Docker container with Volumes.
// CopyFile(source, destination) will work moving file between folders
// Source: https://gist.github.com/var23rav/23ae5d0d4d830aff886c3c970b8f6c6b
func CopyFile(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return fmt.Errorf("Couldn't open source file: %s", err)
	}

	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return fmt.Errorf("Couldn't open dest file: %s", err)
	}
	defer outputFile.Close()

	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return fmt.Errorf("Writing to output file failed: %s", err)
	}

	return nil
}
