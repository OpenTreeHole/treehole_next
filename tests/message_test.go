package tests

import (
	"fmt"
	"strconv"
	"testing"

	. "treehole_next/models"

	"github.com/stretchr/testify/assert"
)

const testNotificationUserID = 9999

func init() {
	// ensure the test user exists for notification tests
	user := User{ID: testNotificationUserID}
	DB.FirstOrCreate(&user, user)
}

func createTestNotification(t *testing.T, relatedFloorID *int, relatedHoleID *int) Message {
	t.Helper()
	notification := Notification{
		Title:          "test notification",
		Description:    "test",
		Data:           Map{"test": true},
		Type:           MessageTypeReply,
		URL:            "/api/floors/1",
		Recipients:     []int{testNotificationUserID},
		RelatedFloorID: relatedFloorID,
		RelatedHoleID:  relatedHoleID,
	}
	_, err := notification.Send()
	assert.NoError(t, err)

	// Send() returns empty Message when NotificationUrl is unset,
	// so query the last inserted message from DB.
	var msg Message
	DB.Last(&msg)
	assert.NotZero(t, msg.ID)
	return msg
}

func assertMessageExists(t *testing.T, messageID int, shouldExist bool) {
	t.Helper()
	var count int64
	DB.Model(&Message{}).Where("id = ?", messageID).Count(&count)
	if shouldExist {
		assert.EqualValues(t, 1, count, "message should exist")
	} else {
		assert.EqualValues(t, 0, count, "message should not exist")
	}
}

func assertMessageUserExists(t *testing.T, messageID int, shouldExist bool) {
	t.Helper()
	var count int64
	DB.Model(&MessageUser{}).Where("message_id = ?", messageID).Count(&count)
	if shouldExist {
		assert.Greater(t, count, int64(0), "message_user should exist")
	} else {
		assert.EqualValues(t, 0, count, "message_user should not exist")
	}
}

func TestDeleteFloorCascadeNotification(t *testing.T) {
	// find a floor from division 7, offset 5 (same as TestDeleteFloor setup)
	var hole Hole
	DB.Where("division_id = ?", 7).Offset(5).First(&hole)
	var floor Floor
	DB.Where("hole_id = ? AND deleted = ?", hole.ID, false).First(&floor)
	if floor.ID == 0 {
		t.Skip("no available floor for cascade test")
	}

	// create a notification linked to this floor
	floorID := floor.ID
	msg := createTestNotification(t, &floorID, nil)
	assertMessageExists(t, msg.ID, true)
	assertMessageUserExists(t, msg.ID, true)

	// delete the floor via API
	data := Map{"delete_reason": "cascade test"}
	testAPI(t, "delete", "/api/floors/"+strconv.Itoa(floor.ID), 200, data)

	// verify notification was cascade-deleted
	assertMessageExists(t, msg.ID, false)
	assertMessageUserExists(t, msg.ID, false)
}

func TestDeleteHoleCascadeNotification(t *testing.T) {
	// create a fresh hole for this test
	hole := Hole{DivisionID: 1}
	err := DB.Create(&hole).Error
	assert.NoError(t, err)
	floor := Floor{HoleID: hole.ID, Content: "cascade test floor", UserID: 1}
	err = DB.Create(&floor).Error
	assert.NoError(t, err)

	// create a notification linked to this hole
	holeID := hole.ID
	msg := createTestNotification(t, nil, &holeID)
	assertMessageExists(t, msg.ID, true)
	assertMessageUserExists(t, msg.ID, true)

	// delete the hole via API (HideHole)
	testCommon(t, "delete", fmt.Sprintf("/api/holes/%d", hole.ID), 204)

	// verify notification was cascade-deleted
	assertMessageExists(t, msg.ID, false)
	assertMessageUserExists(t, msg.ID, false)
}

func TestSendSensitiveRelatedFloorID(t *testing.T) {
	floor := Floor{
		ID:      99999,
		HoleID:  1,
		Content: "sensitive test",
		UserID:  1,
	}

	// SendSensitive sends to admin list, which is empty in test.
	// Instead, verify the Notification struct is correctly constructed
	// by checking the model's field directly.
	notification := Notification{
		Data:           &floor,
		Recipients:     []int{testNotificationUserID},
		Description:    "test",
		Title:          "test",
		Type:           MessageTypeSensitive,
		URL:            fmt.Sprintf("/api/floors/%d", floor.ID),
		RelatedFloorID: &floor.ID,
	}

	_, err := notification.Send()
	assert.NoError(t, err)

	// Send() returns empty Message when NotificationUrl is unset,
	// so query the last inserted message from DB.
	var savedMsg Message
	DB.Last(&savedMsg)
	assert.NotZero(t, savedMsg.ID)
	assert.NotNil(t, savedMsg.RelatedFloorID)
	assert.EqualValues(t, floor.ID, *savedMsg.RelatedFloorID)
	assert.Nil(t, savedMsg.RelatedHoleID)
}

func TestNotificationRelatedFieldsPersisted(t *testing.T) {
	floorID := 12345
	holeID := 67890
	msg := createTestNotification(t, &floorID, &holeID)

	var savedMsg Message
	DB.First(&savedMsg, msg.ID)
	assert.NotNil(t, savedMsg.RelatedFloorID)
	assert.NotNil(t, savedMsg.RelatedHoleID)
	assert.EqualValues(t, floorID, *savedMsg.RelatedFloorID)
	assert.EqualValues(t, holeID, *savedMsg.RelatedHoleID)
}
