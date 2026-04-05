package sensitive

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"treehole_next/config"
)

func TestCheckSensitiveSkipsRemoteCheckWhenDisabled(t *testing.T) {
	oldOpenSensitiveCheck := config.Config.OpenSensitiveCheck
	config.Config.OpenSensitiveCheck = false
	t.Cleanup(func() {
		config.Config.OpenSensitiveCheck = oldOpenSensitiveCheck
	})

	var (
		resp *ResponseForCheck
		err  error
	)

	require.NotPanics(t, func() {
		resp, err = CheckSensitive(ParamsForCheck{
			Content:  "plain text",
			Id:       1,
			TypeName: TypeHole,
		})
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
	assert.True(t, resp.Pass)
	assert.Empty(t, resp.Detail)
	assert.Nil(t, resp.Labels)
}
