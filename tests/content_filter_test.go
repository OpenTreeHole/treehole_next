package tests

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	. "treehole_next/models"
)

func TestContentFilterEndpoints(t *testing.T) {
	// The test-mode API uses user 1. User 2 remains unrestricted.
	require.NoError(t, DB.Create(&UserContentFilter{UserID: 1}).Error)
	require.NoError(t, DB.Create(&UserHiddenDivision{UserID: 1, DivisionID: 2}).Error)
	require.NoError(t, DB.Create(&UserHiddenTag{UserID: 1, Name: "LGBT"}).Error)
	t.Cleanup(func() {
		DB.Where("user_id = ?", 1).Delete(&UserContentFilter{})
		DB.Where("user_id = ?", 1).Delete(&UserHiddenDivision{})
		DB.Where("user_id = ?", 1).Delete(&UserHiddenTag{})
	})

	tag := Tag{Name: "lgbt"}
	require.NoError(t, DB.Create(&tag).Error)
	safeTag := Tag{Name: "filtertestsafe"}
	require.NoError(t, DB.Create(&safeTag).Error)
	blockedByTag := Hole{DivisionID: 1, Tags: Tags{&tag, &safeTag}}
	blockedByDivision := Hole{DivisionID: 2}
	allowed := Hole{DivisionID: 1, Tags: Tags{&safeTag}}
	floorIDs := map[int]int{}
	for _, hole := range []*Hole{&blockedByTag, &blockedByDivision, &allowed} {
		require.NoError(t, DB.Create(hole).Error)
		floor := Floor{HoleID: hole.ID, Content: "contentfiltersearch", Ranking: 0, UserID: 1}
		require.NoError(t, DB.Create(&floor).Error)
		floorIDs[hole.ID] = floor.ID
		require.NoError(t, DB.Create(&UserFavorite{UserID: 1, HoleID: hole.ID}).Error)
		require.NoError(t, DB.Create(&UserSubscription{UserID: 1, HoleID: hole.ID}).Error)
	}
	for _, hole := range []*Hole{&blockedByTag, &blockedByDivision} {
		testCommon(t, "get", fmt.Sprintf("/api/holes/%d", hole.ID), 404)
		testCommon(t, "get", fmt.Sprintf("/api/holes/%d/floors", hole.ID), 404)
		testCommon(t, "get", fmt.Sprintf("/api/floors/%d", floorIDs[hole.ID]), 404)
	}
	testCommon(t, "get", fmt.Sprintf("/api/holes/%d", allowed.ID), 200)
	testCommon(t, "get", fmt.Sprintf("/api/floors/%d", floorIDs[allowed.ID]), 200)
	testCommon(t, "get", "/api/divisions/2", 404)
	testCommon(t, "get", fmt.Sprintf("/api/tags/%d", tag.ID), 404)
	testCommon(t, "get", "/api/tags/lgbt/holes", 404)
	byTag := testAPIArray(t, "get", "/api/tags/filtertestsafe/holes", 200)
	for _, item := range byTag {
		require.NotEqualValues(t, blockedByTag.ID, item["id"])
	}
	tags := testAPIArray(t, "get", "/api/tags", 200)
	for _, item := range tags {
		require.NotEqual(t, "lgbt", item["name"])
	}
	search := testAPIArray(t, "get", "/api/floors/search?search=contentfiltersearch", 200)
	searchIDs := map[int]bool{}
	for _, floor := range search {
		searchIDs[int(floor["hole_id"].(float64))] = true
	}
	require.False(t, searchIDs[blockedByTag.ID])
	require.False(t, searchIDs[blockedByDivision.ID])
	require.True(t, searchIDs[allowed.ID])

	for _, path := range []string{"/api/user/favorites", "/api/users/subscriptions"} {
		items := testAPIArray(t, "get", path, 200)
		ids := map[int]bool{}
		for _, item := range items {
			ids[int(item["id"].(float64))] = true
		}
		require.False(t, ids[blockedByTag.ID])
		require.False(t, ids[blockedByDivision.ID])
		require.True(t, ids[allowed.ID])
	}
	for _, path := range []string{"/api/user/favorites?plain=true", "/api/users/subscriptions?plain=true"} {
		result := testAPI(t, "get", path, 200)
		ids := map[int]bool{}
		for _, id := range result["data"].([]any) {
			ids[int(id.(float64))] = true
		}
		require.False(t, ids[blockedByTag.ID])
		require.False(t, ids[blockedByDivision.ID])
		require.True(t, ids[allowed.ID])
	}
	var visible []int
	require.NoError(t, FilterHoles(DB.Model(&Hole{}), 2).Where("hole.id IN ?", []int{blockedByTag.ID, blockedByDivision.ID, allowed.ID}).Pluck("hole.id", &visible).Error)
	require.Len(t, visible, 3)
}
