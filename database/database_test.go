package database_test

import (
	"testing"

	"github.com/sanjayJ369/LangApp/database"
	"github.com/sanjayJ369/LangApp/testhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestDatabase(t *testing.T) {
	t.Parallel()	
	
	tmpFile := testhelper.GetTempFileLoc()
	
	h, err := database.NewSqlite(tmpFile)
	require.NoError(t, err, "creating handler")

	err = h.Insert("abiser", "Ivory black; animal charcoal.")
	require.NoError(t, err, "inserting values")

	err = h.Insert("abiser", "Ivory black; animal charcoal.")
	require.NoError(t, err, "inserting values second time")

	meaning, err := h.Get("abiser")
	require.NoError(t, err, "getting meaning")

	assert.Equal(t, "Ivory black; animal charcoal.", meaning)
}

func TestBadgerDatabase(t *testing.T) {
	t.Parallel()
	
	tmpFile := testhelper.GetTempFileLoc()
	
	h, err := database.NewBadger(tmpFile + "/")
	t.Cleanup(func() {
		require.NoError(t, h.Close(), "closing badger")
	})
	require.NoError(t, err, "creating handler")

	err = h.Insert("abiser", "Ivory black; animal charcoal.")
	require.NoError(t, err, "inserting values")

	meaning, err := h.Get("abiser")
	require.NoError(t, err, "getting meaning")
	
	assert.Equal(t, "Ivory black; animal charcoal.", meaning)
}
