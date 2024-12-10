package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

type fakeConfig struct {
	Log Log `koanf:"log" json:"log" yaml:"log"`
}

func (c *ConfigTestSuite) TestInit() {
	// -- Given
	//
	fileContent := `{"log": {"level": "debug"}}`
	DefaultConfigDir, _ = os.MkdirTemp(os.TempDir(), "*")
	defer func() {
		_ = os.Remove(DefaultConfigDir)
	}()
	_ = os.WriteFile(filepath.Join(DefaultConfigDir, "skiff.json"), []byte(fileContent), 0o600)

	expected := &fakeConfig{}
	expected.Log.Level = "debug"
	expected.Log.Outputs = "stderr"

	c.T().Setenv("SKIFF_LOG_OUTPUTS", "stderr")
	k := InitKoanf("skiff", &fakeConfig{})
	actual := new(fakeConfig)

	// -- When
	//
	err := k.Unmarshal("", actual)

	// -- Then
	//
	if c.NoError(err) {
		c.Equal(expected, actual)
	}
}

func (c *ConfigTestSuite) TestNewLogger() {
	// -- Given
	//
	dir, _ := os.MkdirTemp(os.TempDir(), "*")
	defer func(name string) {
		_ = os.Remove(name)
	}(dir)
	logFile := filepath.Join(dir, "skiff.log")
	l := Log{Outputs: NewList(logFile), Level: "info"}

	// -- When
	//
	z, err := NewLogger(l)

	// -- Then
	//
	if c.NoError(err) {
		z.Info("hi")
		content, _ := os.ReadFile(logFile)
		type statement struct {
			Msg string `json:"msg"`
		}
		stmt := new(statement)
		_ = json.Unmarshal(content, &stmt)

		c.Equal("hi", stmt.Msg)
	}
}

func (c *ConfigTestSuite) TestToEnvVar() {
	// -- Given
	//
	type Config struct {
		Str   string  `koanf:"str"`
		Int   int     `koanf:"int"`
		Bool  bool    `koanf:"bool"`
		Float float32 `koanf:"float"`
		List  List    `koanf:"list"`
	}

	con := Config{
		Str:   "derp",
		Int:   1,
		Bool:  true,
		Float: 2,
	}

	expected := map[string]string{
		"TEST_STR":   "derp",
		"TEST_INT":   "1",
		"TEST_BOOL":  "true",
		"TEST_FLOAT": "2",
		"TEST_LIST":  "1,2,3",
	}

	c.T().Setenv("TEST_LIST", "1,2,3")

	// -- When
	//
	actual := ToEnvVars("test", con)

	// -- Then
	//
	c.Equal(expected, actual)
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
