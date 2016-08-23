package aws

import (
	"reflect"
	"testing"

	"github.com/titan-x/titan/data"
	"github.com/titan-x/titan/models"
)

const (
	endpoint = "http://localhost:8000"
	region   = "us-west-2"
)

func NewTestDynamoDB() *DynamoDB {
	return NewDynamoDB(region, endpoint)
}

func TestListTables(t *testing.T) {
	db := NewTestDynamoDB()
	tbl, err := db.listTables()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(tbl)
}

func TestSeed(t *testing.T) {
	db := NewTestDynamoDB()
	err := db.Seed(true)
	if err != nil {
		t.Fatal(err)
	}

	tbl, err := db.listTables()
	if err != nil {
		t.Fatal(err)
	}

	if len(tbl) < 1 {
		t.Fatal("tables not created")
	}
}

func TestGetByID(t *testing.T) {
	db := NewTestDynamoDB()
	users := []models.User{data.User1, data.User2}

	for _, user := range users {
		u, ok := db.GetByID(user.ID)
		if !ok {
			t.Fatal("coulnd't get user")
		}

		if u.ID != user.ID ||
			u.Registered != user.Registered ||
			u.Email != user.Email ||
			u.PhoneNumber != user.PhoneNumber ||
			u.GCMRegID != user.GCMRegID ||
			u.APNSDeviceToken != user.APNSDeviceToken ||
			u.Name != user.Name ||
			!reflect.DeepEqual(u.Picture, user.Picture) ||
			u.JWTToken != user.JWTToken {
			t.Fatal("user fields are invalid")
		}
	}
}

func TestGetByMail(t *testing.T) {

}

func TestSaveUser(t *testing.T) {
	// create then update
}
