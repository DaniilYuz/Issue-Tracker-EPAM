package project

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func checkGRPCError(t *testing.T, err error, wantErr bool, code codes.Code, msg string) {
	t.Helper()

	if !wantErr {
		assert.NoError(t, err)
		return
	}

	require.Error(t, err)
	st, ok := status.FromError(err)
	require.True(t, ok, "error should be a gRPC status error")

	assert.Equal(t, code, st.Code())
	assert.Equal(t, msg, st.Message())
}
