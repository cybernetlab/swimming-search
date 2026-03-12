package domain_test

import (
	"testing"

	"github.com/cybernetlab/swimming-search/internal/domain"
	"github.com/stretchr/testify/require"
)

func Test_ErrInvalidContext(t *testing.T) {
	err := domain.NewErrInvalidContext("test-field")
	require.Equal(t, "test-field", err.Field)
	require.Equal(t, "test-field not found in context", err.Error())
}

func Test_ErrAlreadyExists(t *testing.T) {
	err := domain.NewErrAlreadyExists("test-subj")
	require.Equal(t, "test-subj", err.Subject)
	require.Equal(t, "test-subj already exists", err.Error())
}
