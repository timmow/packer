package packer

import (
	"bytes"
	"cgl.tideland.biz/asserts"
	"os"
	"strings"
	"testing"
)

func testEnvironment() Environment {
	config := DefaultEnvironmentConfig()
	config.Ui = &ReaderWriterUi{
		new(bytes.Buffer),
		new(bytes.Buffer),
	}

	env, err := NewEnvironment(config)
	if err != nil {
		panic(err)
	}

	return env
}

func TestEnvironment_DefaultConfig_Commands(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	config := DefaultEnvironmentConfig()
	assert.Empty(config.Commands, "should have no commands")
}

func TestEnvironment_DefaultConfig_Ui(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	config := DefaultEnvironmentConfig()
	assert.NotNil(config.Ui, "default UI should not be nil")

	rwUi, ok := config.Ui.(*ReaderWriterUi)
	assert.True(ok, "default UI should be ReaderWriterUi")
	assert.Equal(rwUi.Writer, os.Stdout, "default UI should go to stdout")
	assert.Equal(rwUi.Reader, os.Stdin, "default UI should read from stdin")
}

func TestNewEnvironment_NoConfig(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	env, err := NewEnvironment(nil)
	assert.Nil(env, "env should be nil")
	assert.NotNil(err, "should be an error")
}

func TestEnvironment_Cli_CallsRun(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	command := &TestCommand{}
	commands := make(map[string]Command)
	commands["foo"] = command

	config := &EnvironmentConfig{}
	config.Commands = []string{"foo"}
	config.CommandFunc = func(n string) Command { return commands[n] }

	env, _ := NewEnvironment(config)
	assert.Equal(env.Cli([]string{"foo", "bar", "baz"}), 0, "runs foo command")
	assert.True(command.runCalled, "run should've been called")
	assert.Equal(command.runEnv, env, "should've ran with env")
	assert.Equal(command.runArgs, []string{"bar", "baz"}, "should have right args")
}

func TestEnvironment_DefaultCli_Empty(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	defaultEnv := testEnvironment()

	assert.Equal(defaultEnv.Cli([]string{}), 1, "CLI with no args")
}

func TestEnvironment_DefaultCli_Help(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	defaultEnv := testEnvironment()

	// A little lambda to help us test the output actually contains help
	testOutput := func() {
		buffer := defaultEnv.Ui().(*ReaderWriterUi).Writer.(*bytes.Buffer)
		output := buffer.String()
		buffer.Reset()
		assert.True(strings.Contains(output, "usage: packer"), "should print help")
	}

	// Test "--help"
	assert.Equal(defaultEnv.Cli([]string{"--help"}), 1, "--help should print")
	testOutput()

	// Test "-h"
	assert.Equal(defaultEnv.Cli([]string{"-h"}), 1, "--help should print")
	testOutput()
}

func TestEnvironment_DefaultCli_Version(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	defaultEnv := testEnvironment()

	// Test the basic version options
	assert.Equal(defaultEnv.Cli([]string{"version"}), 0, "version should work")
	assert.Equal(defaultEnv.Cli([]string{"--version"}), 0, "--version should work")
	assert.Equal(defaultEnv.Cli([]string{"-v"}), 0, "-v should work")

	// Test the --version and -v can appear anywhere
	assert.Equal(defaultEnv.Cli([]string{"bad", "-v"}), 0, "-v should work anywhere")
	assert.Equal(defaultEnv.Cli([]string{"bad", "--version"}), 0, "--version should work anywhere")

	// Test that "version" can't appear anywhere
	assert.Equal(defaultEnv.Cli([]string{"bad", "version"}), 1, "version should NOT work anywhere")
}

func TestEnvironment_SettingUi(t *testing.T) {
	assert := asserts.NewTestingAsserts(t, true)

	ui := &ReaderWriterUi{new(bytes.Buffer), new(bytes.Buffer)}

	config := &EnvironmentConfig{}
	config.Ui = ui

	env, _ := NewEnvironment(config)

	assert.Equal(env.Ui(), ui, "UIs should be equal")
}