package tests

import (
	"testing"

	"github.com/stretchr/testify/require"

	. "treehole_next/models"
)

func TestUserHiddenDivisionList(t *testing.T) {
	// In test mode the API uses user 1. The restriction must not affect user 2.
	setting := UserHiddenDivision{UserID: 1, DivisionID: 2}
	require.NoError(t, DB.Create(&setting).Error)
	t.Cleanup(func() { DB.Delete(&setting) })

	response := testAPIArray(t, "get", "/api/divisions", 200)
	for _, division := range response {
		require.NotEqualValues(t, 2, division["id"])
	}
	require.NotEmpty(t, response)

	var otherUser Divisions
	require.NoError(t, VisibleDivisions(DB, 2).Find(&otherUser).Error)
	found := false
	for _, division := range otherUser {
		if division.ID == 2 {
			found = true
		}
	}
	require.True(t, found, "unconfigured users should still see division 2")
}
