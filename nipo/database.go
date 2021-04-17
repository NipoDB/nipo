package main

import (
	"regexp"
	"os"
)

/*
creates the main database contains config and cluster definition
*/
func CreateDatabase() *Database {
	config := GetConfig(os.Args[1])
	cluster := config.CreateCluster()
	return &Database{items: map[string]string{}, config: config, cluster: cluster, reloaded: false}
}

/*
creates temporary database
*/
func CreateTempDatabase() *Database {
	return &Database{items: map[string]string{}}
}

/*
gets the value of given key
*/
func (database *Database) Get(key string) (string, bool) {
	value, ok := database.items[key]
	return value, ok
}

/*
sets the given key value
*/
func (database *Database) Set(key, value string) bool {
	database.items[key] = value
	rvalue, ok := database.Get(key)
	if rvalue != value {
		return false
	}
	return ok
}

/*
handles foreach operation on database object
*/
func (database *Database) Foreach(action func(string, string)) {
	for key, value := range database.items {
		action(key, value)
	}
}

/*
select the keys and values which are matched in given regex
*/
func (database *Database) Select(keyRegex string) (*Database, error) {
	selected := CreateTempDatabase()
	var err error
	for key, value := range database.items {
		matched, err := regexp.MatchString(keyRegex, key)
		if err != nil {
			return selected, err
		}
		if matched {
			selected.items[key] = value
		}
	}
	return selected, err
}
