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

	tag := Tag{Name: "*lGbT+讨论"}
	require.NoError(t, DB.Create(&tag).Error)
	exactRule := UserHiddenTag{UserID: 1, Name: "开盒", ExactMatch: true}
	require.NoError(t, DB.Create(&exactRule).Error)
	exactTag := Tag{Name: exactRule.Name}
	require.NoError(t, DB.Create(&exactTag).Error)
	testCommon(t, "get", fmt.Sprintf("/api/tags/%d", exactTag.ID), 404)
	reminderTag := Tag{Name: "别来开盒"}
	require.NoError(t, DB.Create(&reminderTag).Error)
	testCommon(t, "get", fmt.Sprintf("/api/tags/%d", reminderTag.ID), 200)
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
	testCommon(t, "get", "/api/tags/*lGbT+讨论/holes", 404)
	byTag := testAPIArray(t, "get", "/api/tags/filtertestsafe/holes", 200)
	foundAllowed := false
	for _, item := range byTag {
		require.NotEqualValues(t, blockedByTag.ID, item["id"])
		if int(item["id"].(float64)) == allowed.ID {
			foundAllowed = true
		}
	}
	require.True(t, foundAllowed)
	tags := testAPIArray(t, "get", "/api/tags", 200)
	for _, item := range tags {
		require.NotEqual(t, tag.Name, item["name"])
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

	// List-only configuration must not enable the extended restrictions.
	require.NoError(t, DB.Delete(&UserContentFilter{}, "user_id = ?", 1).Error)
	testCommon(t, "get", fmt.Sprintf("/api/holes/%d", blockedByTag.ID), 200)
	testCommon(t, "get", fmt.Sprintf("/api/holes/%d", blockedByDivision.ID), 200)
	testCommon(t, "get", fmt.Sprintf("/api/tags/%d", tag.ID), 200)
}
