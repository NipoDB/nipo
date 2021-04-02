package main

import (
	"testing"
)

func TestCreateDatabase(t *testing.T) {
	db := CreateDatabase()
	if len(db.items) != 0 {
		t.Errorf("The database returned by CreateDatabase() is not empty ")
	}
}

func TestGet(t *testing.T) {
	db := &Database{items: map[string]string{"hello": "world"}}

	value, ok := db.Get("hello")
	if !ok {
		t.Errorf("Could not get value from db using Get()")
	}
	if value != "world" {
		t.Errorf("The value returned from Get() is not correct")
	}
}

func TestSet(t *testing.T) {
	db := &Database{items: map[string]string{}}

	db.Set("hello", "world")
	if db.items["hello"] != "world" {
		t.Errorf("The value is not set correct using Set()")
	}
}

func TestSelect(t *testing.T) {
	// TODO: needs more cases to be checked
	dbItems := map[string]string{
		"hello": "world1",
		"helli": "world2",
	}
	db := &Database{items: dbItems}

	newDb, err := db.Select("hell.")
	if err != nil {
		t.Errorf("Could not select a sub-db using Select()")
	}
	if len(newDb.items) != 2 {
		t.Errorf("The returned sub-db from Select() is not correct")
	}
}
