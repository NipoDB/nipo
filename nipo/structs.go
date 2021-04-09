package main

import (
	"sync"
	"net"
)

var Wait sync.WaitGroup
var Lock sync.Mutex

type Database struct {
	items map[string]string
}

type Client struct {
	Connection net.Conn
	User       User
	Authorized bool
}

type Slave struct {
	Node              *Node
	Status, CheckedAt string
	Database          *Database
}

type Cluster struct {
	Slaves []Slave
	Status string
}

type User struct {
	Name, Token, Keys, Cmds string
}

type Node struct {
	Id                             int
	Ip, Port, Authorization, Token string
}

type Config struct {
	Global struct {
		Authorization, Master string
		CheckInterval         int
	}

	Slaves []*Node

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

	Users []*User
}