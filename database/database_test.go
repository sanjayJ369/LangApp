package database_test

import (
	"testing"

	"github.com/sanjayJ369/LangApp/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestDatabase(t *testing.T) {
	h, err := database.New("test.db")
	require.NoError(t, err, "creating handler")

	err = h.Insert("abiser", "Ivory black; animal charcoal.")
	require.NoError(t, err, "inserting values")
	err = h.Insert("abiser", "Ivory black; animal charcoal2.")
	require.NoError(t, err, "inserting values second time")

	meaning, err := h.Get("abiser")
	require.NoError(t, err, "getting meaning")
	assert.Equal(t, "Ivory black; animal charcoal2.", meaning)
}
