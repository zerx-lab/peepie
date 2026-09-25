package controller

import (
	"testing"

	"pentagi/pkg/config"

	"github.com/stretchr/testify/assert"
)

func TestConfiguredFlowImage(t *testing.T) {
	t.Parallel()

	cfg := &config.Config{
		DockerFlowImage:      "env/flow",
		DockerAssistantImage: " env/assistant ",
		Overrides:            config.NewOverrides(),
	}

	// Each mode reads its own env value until the Web UI overrides it.
	assert.Equal(t, "env/flow", configuredFlowImage(cfg, false))
	assert.Equal(t, "env/assistant", configuredFlowImage(cfg, true))

	cfg.Overrides.Set(config.CategoryExecution, config.KeyDockerAssistantImage, "ui/assistant")
	assert.Equal(t, "env/flow", configuredFlowImage(cfg, false))
	assert.Equal(t, "ui/assistant", configuredFlowImage(cfg, true))

	// A cleared field is stored as "" and hands the choice back to the LLM
	// even though the process environment still carries an image.
	cfg.Overrides.Set(config.CategoryExecution, config.KeyDockerFlowImage, "")
	assert.Empty(t, configuredFlowImage(cfg, false))
}
