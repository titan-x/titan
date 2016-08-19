package aws

import "testing"

const (
	endpoint = "http://localhost:8000"
	region   = "us-west-2"
)

func TestListTables(t *testing.T) {
	db := NewDynamoDB(region, endpoint)
	tbl, err := db.listTables()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(tbl)
}

func TestSeed(t *testing.T) {
	db := NewDynamoDB(region, endpoint)
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
