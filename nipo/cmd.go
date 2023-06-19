package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

/*
change proc config of database
*/
func (database *Database) setProcCmd(cmd string) *Database {
	var err error
	cmds := strings.Split(cmd, " ")

	core, err := strconv.Atoi(cmds[0])
	if err != nil {
		fmt.Println(err)
	}

	thread, err := strconv.Atoi(cmds[1])
	if err != nil {
		fmt.Println(err)
	}

	database.SetProc(core, thread)
	return database
}

/*
validates that user have permission to run the command or not
*/
func validateCmd(cmd string, user *User) bool {
	cmds := strings.Split(user.Cmds, "||")
	for count := range cmds {
		if cmds[count] == "all" {
			return true
		}
		if cmds[count] == cmd {
			return true
		}
	}
	return false
}

/*
validates that user have access to this key (regex) or not
*/
func validateKey(key string, user *User) bool {
	keys := strings.Split(user.Keys, "||")
	for count := range keys {
		matched, err := regexp.MatchString(keys[count], key)
		if err != nil {
			fmt.Println(err)
		}
		if matched {
			return true
		}
	}
	return false
}

/*
sets the key and value into database
*/
func (database *Database) cmdSet(cmd string) (*Database, bool) {
	cmdFields := strings.Fields(cmd)
	key := cmdFields[1]
	db := CreateTempDatabase()
	ok := false
	if len(cmdFields) >= 3 {
		value := cmdFields[2]
		for n := 3; n < len(cmdFields); n++ {
			value += " " + cmdFields[n]
		}
		slaveSynced := database.cluster.SetOnSlaves(database.config, key, value)
		for !slaveSynced {
			slaveSynced = database.cluster.SetOnSlaves(database.config, key, value)
		}
		ok = database.Set(key, value)
		if !ok {
			return db, false
		} else {
			db.items[key] = value
			return db, true
		}
	}
	return db, ok
}

/*
gets the value of given keys, it will create a new databse as result
*/
func (database *Database) cmdGet(cmd string) *Database {
	cmdFields := strings.Fields(cmd)
	db := CreateTempDatabase()
	for _, key := range cmdFields {
		if cmdFields[0] != key {
			value, ok := database.Get(key)
			if ok {
				db.items[key] = value
			}
		}
	}
	return db
}

/*
selects the keys and values which are matched in given regex
it will create a new databse as result
*/
func (database *Database) cmdSelect(cmd string) *Database {
	cmdFields := strings.Fields(cmd)
	key := cmdFields[1]
	db, err := database.Select(key)
	if err != nil {
		fmt.Println(err)
	}
	return db
}

/*
summeries the values which mathed in given regex
it will ignore non number values
*/
func (database *Database) cmdSum(cmd string) *Database {
	cmdFields := strings.Fields(cmd)
	key := cmdFields[1]
	db, err := database.Select(key)
	returnDB := CreateTempDatabase()
	var sum float64 = 0
	if err != nil {
		fmt.Println(err)
	}
	db.Foreach(func(key, value string) {
		valFloat, _ := strconv.ParseFloat(value, 64)
		sum = sum + valFloat
	})
	returnDB.items[key] = fmt.Sprintf("%f", sum)
	return returnDB
}

/*
calculate the average of the values which mathed in given regex
it will ignore non number values
*/
func (database *Database) cmdAvg(cmd string) *Database {
	cmdFields := strings.Fields(cmd)
	key := cmdFields[1]
	db, err := database.Select(key)
	returnDB := CreateTempDatabase()
	var sum float64 = 0
	count := 0
	if err != nil {
		fmt.Println(err)
	}
	db.Foreach(func(key, value string) {
		valFloat, _ := strconv.ParseFloat(value, 64)
		sum = sum + valFloat
		count++
	})
	avg := sum / (float64(count))
	returnDB.items[key] = fmt.Sprintf("%f", avg)
	return returnDB
}

/*
counts the keys which mathed in given regex
it will ignore non number values
*/
func (database *Database) cmdCount(cmd string) *Database {
	cmdFields := strings.Fields(cmd)
	key := cmdFields[1]
	db, err := database.Select(key)
	returnDB := CreateTempDatabase()
	count := 0
	if err != nil {
		fmt.Println(err)
	}
	db.Foreach(func(key, value string) {
		count++
	})
	returnDB.items["count"] = strconv.Itoa(count)
	return returnDB
}

/*
the main function to handle the command
checks the validation and authorization of user to access the keys and commands
*/
func (database *Database) cmd(cmd string, user *User) (*Database, string) {
	database.config.logger("client executed command : "+cmd, 2)
	database.config.logger("cmd.go - func cmd - with cmd : "+cmd, 2)
	database.config.logger("cmd.go - func cmd - with user : "+user.Name, 2)
	cmdFields := strings.Fields(cmd)
	db := CreateTempDatabase()
	ok := false
	message := ""
	if len(cmdFields) >= 2 {
		switch cmdFields[0] {
		case "count":
			db = database.cmdCount(cmd)
		case "set":
			if database.config.Acl.Authorization == "true" {
				if validateCmd("set", user) {
					if validateKey(cmdFields[1], user) {
						Lock.Lock()
						db, ok = database.cmdSet(cmd)
						Lock.Unlock()
						if !ok {
							message = "set failed by user " + user.Name + " for command : " + cmd
							database.config.logger(message, 1)
						}
					} else {
						message = "User " + user.Name + " not allowed to use regex : " + cmdFields[1]
						database.config.logger(message, 1)
					}
				} else {
					message = "User " + user.Name + " not allowed to use command : " + cmd
					database.config.logger(message, 1)
				}
			} else {
				Lock.Lock()
				db, ok = database.cmdSet(cmd)
				Lock.Unlock()
				if !ok {
					message = "set failed by user " + user.Name + " for command : " + cmd
					database.config.logger(message, 1)
				}
			}
		case "get":
			if database.config.Acl.Authorization == "true" {
				if validateCmd("get", user) {
					if validateKey(cmdFields[1], user) {
						db = database.cmdGet(cmd)
					} else {
						message = "User " + user.Name + " not allowed to use regex : " + cmdFields[1]
						database.config.logger(message, 1)
					}
				} else {
					message = "User " + user.Name + " not allowed to use command : " + cmd
					database.config.logger(message, 1)
				}
			} else {
				db = database.cmdGet(cmd)
			}
		case "select":
			if database.config.Acl.Authorization == "true" {
				if validateCmd("select", user) {
					if validateKey(cmdFields[1], user) {
						db = database.cmdSelect(cmd)
					} else {
						message = "User " + user.Name + " not allowed to use regex : " + cmdFields[1]
						database.config.logger(message, 1)
					}
				} else {
					message = "User " + user.Name + " not allowed to use command : " + cmd
					database.config.logger(message, 1)
				}
			} else {
				db = database.cmdSelect(cmd)
			}
		case "sum":
			if database.config.Acl.Authorization == "true" {
				if validateCmd("sum", user) {
					if validateKey(cmdFields[1], user) {
						db = database.cmdSum(cmd)
					} else {
						message = "User " + user.Name + " not allowed to use regex : " + cmdFields[1]
						database.config.logger(message, 1)
					}
				} else {
					message = "User " + user.Name + " not allowed to use command : " + cmd
					database.config.logger(message, 1)
				}
			} else {
				db = database.cmdSum(cmd)
			}
		case "avg":
			if database.config.Acl.Authorization == "true" {
				if validateCmd("avg", user) {
					if validateKey(cmdFields[1], user) {
						db = database.cmdAvg(cmd)
					} else {
						message = "User " + user.Name + " not allowed to use regex : " + cmdFields[1]
						database.config.logger(message, 1)
					}
				} else {
					message = "User " + user.Name + " not allowed to use command : " + cmd
					database.config.logger(message, 1)
				}
			} else {
				db = database.cmdAvg(cmd)
			}
		case "setproc":
			db = database.setProcCmd(cmd)
		}
	}
	return db, message
}
