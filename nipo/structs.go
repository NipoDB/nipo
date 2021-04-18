package main

import (
	"sync"
	"net"
)

var Wait sync.WaitGroup
var Lock sync.Mutex

/*
defines a node
*/
type Node struct {
	Id                             int
	Ip, Port, Authorization, Token string
}

/*
defines user
*/
type User struct {
	Name, Token, Keys, Cmds string
}

/*
each client is a connection with user and authorization definition
*/
type Client struct {
	Connection net.Conn
	User       	User
	Authorized bool
}

/*
contains all directives at config file 
*/
type Config struct {
	Acl struct {
		Authorization string
		Users []*User
	}

	Cluster struct {
		Master 			string
		CheckInterval	int
		Slaves []*Node
	}

	Proc struct {
		Cores, Threads int
	}

	Listen struct {
		Ip, Port, Protocol string
	}

	Log struct {
		Level int
		Path  string
	}
}

var tempConfig *Config

/*
defines slaves from config object
*/
type Slave struct {
	Node		*Node
	Status, CheckedAt string
	Database	*Database
}

/*
defines the cluster with slaves and status of cluster
*/
type Cluster struct {
	Slaves []Slave
	Status string
}

/*
structure of main database. 
items used for key-value map
config contains the configuration file
cluster contains the cluster definition
socket contains the listener definition
*/
type Database struct {
	items map[string]string
	config 		*Config
	cluster		*Cluster
	socket 		net.Listener
	reloaded	bool
}
