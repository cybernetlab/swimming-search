package domain

import (
	"context"
)

type NodeID string

type nodeIDKeyType int

var nodeIDKey nodeIDKeyType

func ContextNodeID(ctx context.Context) (NodeID, error) {
	id, ok := ctx.Value(nodeIDKey).(NodeID)
	if !ok {
		return "", NewErrInvalidContext("nodeID")
	}
	return id, nil
}

func WithNodeID(ctx context.Context, id NodeID) context.Context {
	return context.WithValue(ctx, nodeIDKey, id)
}
