package database_test

import (
	"testing"

	"github.com/sanjayJ369/LangApp/database"
	"github.com/sanjayJ369/LangApp/testhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestSqlite(t *testing.T) {
	t.Parallel()

	tmpFile := testhelper.GetTempFileLoc()

	h, err := database.NewSqlite(tmpFile)
	require.NoError(t, err, "creating handler")

	t.Cleanup(func() {
		require.NoError(t, h.Close(), "closing db")
	})

	err = h.Insert("abiser", "Ivory black; animal charcoal.")
	require.NoError(t, err, "inserting value")

	err = h.Insert("abiser", "Ivory black; animal charcoal.")
	require.NoError(t, err, "inserting value second time")

	meaning, err := h.Get("abiser")
	require.NoError(t, err, "getting value")

	assert.Equal(t, meaning, "Ivory black; animal charcoal.")
}

func TestBadger(t *testing.T) {
	t.Parallel()

	tmpFile := testhelper.GetTempFileLoc()

	h, err := database.NewBadger(tmpFile + "/")
	require.NoError(t, err, "creating handler")

	t.Cleanup(func() {
		require.NoError(t, h.Close(), "closing db")
	})

	err = h.Insert("abiser", "Ivory black; animal charcoal.")
	require.NoError(t, err, "inserting value")

	err = h.Insert("abiser", "Ivory black; animal charcoal.")
	require.NoError(t, err, "inserting value second time")

	meaning, err := h.Get("abiser")
	require.NoError(t, err, "getting value")

	assert.Equal(t, meaning, "Ivory black; animal charcoal.")
}
