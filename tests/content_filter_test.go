package tests

import (
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"net/http"
	"net/http/httptest"

	"testing"

	"github.com/stretchr/testify/require"

	. "treehole_next/models"
)

func TestContentFilterEndpoints(t *testing.T) {
	// Isolate the suite's shared fixtures from collection edits and new posts.
	originalDB := DB
	DB = originalDB.Begin()
	require.NoError(t, DB.Error)
	t.Cleanup(func() { DB.Rollback(); DB = originalDB })

	// The test-mode API uses user 1. User 2 remains unrestricted.
	require.NoError(t, DB.Create(&UserContentFilter{UserID: 1}).Error)
	require.NoError(t, DB.Create(&UserHiddenDivision{UserID: 1, DivisionID: 2}).Error)
	require.NoError(t, DB.Create(&UserHiddenTag{UserID: 1, Name: "*LGBT+讨论"}).Error)

	tag := Tag{Name: "*lGbT+讨论"}
	require.NoError(t, DB.Create(&tag).Error)
	unlisted := Tag{Name: "提醒*lgbt+讨论"}
	require.NoError(t, DB.Create(&unlisted).Error)

	safeTag := Tag{Name: "filtertestsafe"}
	require.NoError(t, DB.Create(&safeTag).Error)
	blockedByTag := Hole{DivisionID: 1, Tags: Tags{&tag, &safeTag}}
	roundtableTag := Tag{Name: "roundtableonly"}
	require.NoError(t, DB.Create(&roundtableTag).Error)
	blockedByDivision := Hole{DivisionID: 2, Tags: Tags{&roundtableTag}}
	allowed := Hole{DivisionID: 1, Tags: Tags{&safeTag, &unlisted}}
	floorIDs := map[int]int{}
	for _, hole := range []*Hole{&blockedByTag, &blockedByDivision, &allowed} {
		require.NoError(t, DB.Create(hole).Error)
		floor := Floor{HoleID: hole.ID, Content: "contentfiltersearch", Ranking: 0, UserID: 1}
		require.NoError(t, DB.Create(&floor).Error)
		floorIDs[hole.ID] = floor.ID
		require.NoError(t, DB.Create(&UserFavorite{UserID: 1, HoleID: hole.ID}).Error)
		require.NoError(t, DB.Create(&UserSubscription{UserID: 1, HoleID: hole.ID}).Error)
	}
	testCommon(t, "get", fmt.Sprintf("/api/tags/%d", unlisted.ID), 200)
	for _, name := range []string{tag.Name, roundtableTag.Name} {
		testCommon(t, "post", "/api/tags", 404, Map{"name": name})
		require.Empty(t, testAPIArray(t, "get", "/api/tags?s="+name, 200))
	}
	testCommon(t, "get", fmt.Sprintf("/api/tags/%d", roundtableTag.ID), 404)
	testCommon(t, "post", "/api/tags", 200, Map{"name": safeTag.Name})

	// Use synthetic messages and verify both inbox and delivery for mixed recipients.
	for _, id := range []int{1, 2} {
		var u User
		require.NoError(t, DB.FirstOrCreate(&u, User{ID: id}).Error)
	}
	for _, h := range []*Hole{&blockedByTag, &blockedByDivision, &allowed} {
		fid := floorIDs[h.ID]
		report := Report{FloorID: fid, Reason: "fixture"}
		require.NoError(t, DB.Create(&report).Error)
		msg := Message{Title: fmt.Sprintf("inbox-%d", h.ID), Type: MessageTypeReply, RelatedFloorID: &fid, Recipients: []int{1}, Data: Map{"snapshot": "old"}}
		require.NoError(t, DB.Create(&msg).Error)
		legacy := Message{Title: fmt.Sprintf("legacy-%d", h.ID), URL: fmt.Sprintf("/api/reports/%d", report.ID), Type: MessageTypeReportDealt, Recipients: []int{1}}
		require.NoError(t, DB.Create(&legacy).Error)
		punishment := Punishment{UserID: 1, FloorID: &fid, DivisionID: h.DivisionID}
		require.NoError(t, DB.Create(&punishment).Error)
		status := 404
		if h == &allowed {
			status = 200
		}
		testCommon(t, "get", fmt.Sprintf("/api/reports/%d", report.ID), status)
		push := Notification{Title: fmt.Sprintf("delivery-%d", h.ID), Type: MessageTypeReply, RelatedFloorID: &fid, Recipients: []int{1, 2}, Data: Map{"snapshot": "old"}}
		_, err := push.Send()
		require.NoError(t, err)
		var received []int
		require.NoError(t, DB.Table("message_user").Joins("JOIN message ON message.id=message_user.message_id").Where("message.title = ?", push.Title).Pluck("user_id", &received).Error)
		if h == &allowed {
			require.Contains(t, received, 1)
		} else {
			require.NotContains(t, received, 1)
		}
		require.Contains(t, received, 2)
	}
	body := testCommon(t, "get", "/api/messages", 200)
	require.Contains(t, string(body), fmt.Sprintf("legacy-%d", allowed.ID))
	punishments := testAPIArray(t, "get", "/api/users/me/punishments", 200)
	var pf []int
	for _, p := range punishments {
		if p["floor_id"] != nil {
			pf = append(pf, int(p["floor_id"].(float64)))
		}
	}
	require.Contains(t, pf, floorIDs[allowed.ID])
	for _, h := range []*Hole{&blockedByTag, &blockedByDivision} {
		require.NotContains(t, string(body), fmt.Sprintf("legacy-%d", h.ID))
		require.NotContains(t, pf, floorIDs[h.ID])
	}

	require.Contains(t, string(body), fmt.Sprintf("inbox-%d", allowed.ID))
	require.NotContains(t, string(body), `"snapshot":"old"`)
	for _, h := range []*Hole{&blockedByTag, &blockedByDivision} {
		require.NotContains(t, string(body), fmt.Sprintf("inbox-%d", h.ID))
		for _, path := range []string{"/api/user/favorites", "/api/users/subscriptions"} {
			testCommon(t, "post", path, 404, Map{"hole_id": h.ID})
		}
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

	for _, path := range []string{"/api/user/favorites", "/api/user/favorites?favorite_group_id=0", "/api/users/subscriptions"} {
		items := testAPIArray(t, "get", path, 200)
		ids := map[int]bool{}
		for _, item := range items {
			ids[int(item["id"].(float64))] = true
		}
		require.False(t, ids[blockedByTag.ID])
		require.False(t, ids[blockedByDivision.ID])
		require.True(t, ids[allowed.ID])
	}
	for _, path := range []string{"/api/user/favorites?plain=true", "/api/user/favorites?plain=true&favorite_group_id=0", "/api/users/subscriptions?plain=true"} {
		result := testAPI(t, "get", path, 200)
		ids := map[int]bool{}
		for _, id := range result["data"].([]any) {
			ids[int(id.(float64))] = true
		}
		require.False(t, ids[blockedByTag.ID])
		require.False(t, ids[blockedByDivision.ID])
		require.True(t, ids[allowed.ID])
	}
	for _, hole := range []*Hole{&blockedByTag, &blockedByDivision} {
		testCommon(t, "get", fmt.Sprintf("/api/holes/%d/summary", hole.ID), 404)
		testCommon(t, "get", fmt.Sprintf("/api/floors/%d/history", floorIDs[hole.ID]), 404)
		require.Empty(t, testAPIArray(t, "get", fmt.Sprintf("/api/floors?hole_id=%d", hole.ID), 200))
	}

	// Writes resolve embedded mentions with the same visibility as reads.
	text := fmt.Sprintf("##%d ##%d ##%d", floorIDs[blockedByTag.ID], floorIDs[blockedByDivision.ID], floorIDs[allowed.ID])
	response := testAPI(t, "post", fmt.Sprintf("/api/holes/%d/floors", allowed.ID), 201, Map{"content": text})
	embedded := response["mention"].([]any)
	require.Len(t, embedded, 1)
	require.EqualValues(t, floorIDs[allowed.ID], embedded[0].(map[string]any)["id"])

	// A successful collection edit also returns the current collection.
	for _, path := range []string{"/api/user/favorites", "/api/users/subscriptions"} {
		response := testAPI(t, "post", path, 201, Map{"hole_id": allowed.ID})
		ids := response["data"].([]any)
		require.Contains(t, ids, float64(allowed.ID))
		require.NotContains(t, ids, float64(blockedByTag.ID))
		require.NotContains(t, ids, float64(blockedByDivision.ID))
		deleted := testAPI(t, "delete", path, 200, Map{"hole_id": allowed.ID})
		remaining := deleted["data"].([]any)
		for _, h := range []*Hole{&allowed, &blockedByTag, &blockedByDivision} {
			require.NotContains(t, remaining, float64(h.ID))
		}
	}

	// Return a mixed Elasticsearch page, including hidden hits and highlights.
	// The API must use the database visibility rules before emitting either.
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Elastic-Product", "Elasticsearch")
		hits := []Map{}
		for _, h := range []*Hole{&blockedByTag, &allowed, &blockedByDivision} {
			hits = append(hits, Map{"_index": "floors", "_id": fmt.Sprint(floorIDs[h.ID]), "highlight": Map{"content": []string{fmt.Sprintf("snippet-%d-end", h.ID)}}})
		}
		json.NewEncoder(w).Encode(Map{"hits": Map{"hits": hits}})
	}))
	defer server.Close()
	previous := ES
	defer func() { ES = previous }()
	var err error
	ES, err = elasticsearch.NewTypedClient(elasticsearch.Config{Addresses: []string{server.URL}})
	require.NoError(t, err)
	for _, endpoint := range []struct{ method, path string }{
		{"get", "/api/floors/search?search=contentfiltersearch"},
		{"post", "/api/floors/search?search=contentfiltersearch&accurate=true"},
		{"get", "/api/floors?s=contentfiltersearch"},
	} {
		body := testCommon(t, endpoint.method, endpoint.path, 200)
		var hits []Map
		require.NoError(t, json.Unmarshal(body, &hits))
		require.Len(t, hits, 1)
		require.EqualValues(t, allowed.ID, hits[0]["hole_id"])
		require.Contains(t, string(body), fmt.Sprintf("snippet-%d-end", allowed.ID))
		for _, h := range []*Hole{&blockedByTag, &blockedByDivision} {
			require.NotContains(t, string(body), fmt.Sprintf("snippet-%d-end", h.ID))
		}
	}
	var visible []int
	require.NoError(t, FilterHoles(DB.Model(&Hole{}), 2).Where("hole.id IN ?", []int{blockedByTag.ID, blockedByDivision.ID, allowed.ID}).Pluck("hole.id", &visible).Error)
	require.Len(t, visible, 3)

	// List-only configuration must not enable the extended restrictions.
	require.NoError(t, DB.Delete(&UserContentFilter{}, "user_id = ?", 1).Error)
	testCommon(t, "get", fmt.Sprintf("/api/holes/%d", blockedByTag.ID), 200)
	testCommon(t, "get", fmt.Sprintf("/api/holes/%d", blockedByDivision.ID), 200)
	testCommon(t, "get", fmt.Sprintf("/api/tags/%d", tag.ID), 200)
	for _, path := range []string{"/api/user/favorites?plain=true", "/api/users/subscriptions?plain=true"} {
		result := testAPI(t, "get", path, 200)
		require.Contains(t, result["data"], float64(blockedByTag.ID))
		require.Contains(t, result["data"], float64(blockedByDivision.ID))
	}
	controlHits := testAPIArray(t, "get", "/api/floors/search?search=contentfiltersearch", 200)
	require.Len(t, controlHits, 3)

}
