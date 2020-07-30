package core

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadDefaultParameters(t *testing.T) {
	var params cliParams

	params.setup(nil)
	if err := params.parse([]string{}); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "", params.configFile)
	assert.False(t, params.commit)
	assert.False(t, params.version)
}

func TestReadValidParams(t *testing.T) {
	var params cliParams
	configFilename := "/path/to/configfile.yaml"
	args := []string{
		"-config", configFilename,
		"-commit",
		"-version",
	}

	params.setup(nil)
	if err := params.parse(args); err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, configFilename, params.configFile)
	assert.True(t, params.commit)
	assert.True(t, params.version)
}

func TestReadInvalidParams(t *testing.T) {
	var params cliParams
	var buf bytes.Buffer
	args := []string{"-configfile", "asdfasdf", "-commit", "true"}

	params.setup(&buf)
	err := params.parse(args)
	if err == nil {
		t.Fatal("expected an error")
	}

	assert.NotEmpty(t, buf)
}

func TestValidation(t *testing.T) {
	assert.Nil(t, (&cliParams{configFile: "configfile.yaml"}).validate())
	assert.NotNil(t, (&cliParams{}).validate())
}
