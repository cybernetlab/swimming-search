package domain_test

import (
	"testing"

	"github.com/cybernetlab/swimming-search/internal/domain"
	"github.com/stretchr/testify/require"
)

func TestAuthorize_Anonymous(t *testing.T) {
	require.False(t, domain.User{}.Authorize("help:get"))
}

func TestAuthorize_RegularUser(t *testing.T) {
	require.True(t, domain.User{Name: "test"}.Authorize("help:get"))
	require.False(t, domain.User{Name: "test"}.Authorize("user:create"))
	require.False(t, domain.User{Name: "test"}.Authorize("user:get"))
}

func TestAuthorize_AdminUser(t *testing.T) {
	require.True(t, domain.User{Name: "test", IsAdmin: true}.Authorize("help:get"))
	require.True(t, domain.User{Name: "test", IsAdmin: true}.Authorize("user:create"))
	require.True(t, domain.User{Name: "test", IsAdmin: true}.Authorize("user:get"))
}
