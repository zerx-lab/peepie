package graph

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNormalizeDockerImage(t *testing.T) {
	t.Parallel()

	valid := map[string]string{
		"":                                     "",
		"   ":                                  "",
		" vxcontrol/kali-linux ":               "vxcontrol/kali-linux",
		"debian:latest":                        "debian:latest",
		"registry.local:5000/team/kali:2025.1": "registry.local:5000/team/kali:2025.1",
		"python@sha256:" + sha256Hex:           "python@sha256:" + sha256Hex,
	}
	for in, want := range valid {
		got, err := normalizeDockerImage(in)
		require.NoError(t, err, "input %q", in)
		assert.Equal(t, want, got, "input %q", in)
	}

	for _, in := range []string{"Kali/Linux", "kali linux", "kali:", "http://kali"} {
		_, err := normalizeDockerImage(in)
		assert.Error(t, err, "input %q", in)
	}
}

const sha256Hex = "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef"
