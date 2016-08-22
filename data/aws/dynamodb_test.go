package aws

import (
	"reflect"
	"testing"

	"github.com/titan-x/titan/data"
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
	su1 := data.User1

	// todo: for/range loop for a number of users
	u1, ok := db.GetByID(su1.ID)
	if !ok {
		t.Fatal("coulnd't get user")
	}

	if u1.ID != su1.ID ||
		u1.Registered != su1.Registered ||
		u1.Email != su1.Email ||
		u1.PhoneNumber != su1.PhoneNumber ||
		u1.GCMRegID != su1.GCMRegID ||
		u1.APNSDeviceToken != su1.APNSDeviceToken ||
		u1.Name != su1.Name ||
		!reflect.DeepEqual(u1.Picture, su1.Picture) ||
		u1.JWTToken != su1.JWTToken {
		t.Fatal("user fields are invalid")
	}
}

func TestGetByMail(t *testing.T) {

}

func TestSaveUser(t *testing.T) {
	// create then update
}
