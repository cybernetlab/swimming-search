package domain_test

import (
	"testing"

	"github.com/cybernetlab/swimming-search/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestContextNodeID(t *testing.T) {
	var ctxErr domain.ErrInvalidContext

	id, err := domain.ContextNodeID(t.Context())
	require.ErrorAs(t, err, &ctxErr)
	require.Equal(t, "nodeID", ctxErr.Field)

	ctx := domain.WithNodeID(t.Context(), domain.NodeID("test"))
	id, err = domain.ContextNodeID(ctx)
	require.NoError(t, err)
	require.Equal(t, domain.NodeID("test"), id)
}
