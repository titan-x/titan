package aws

import "testing"

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
	db.GetByID(id)
}

func TestGetByMail(t *testing.T) {

}

func TestSaveUser(t *testing.T) {
	// create then update
}
