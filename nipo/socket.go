package main

import (
	"bufio"
	"encoding/json"
	"net"
	"os"
	"os/signal"
	"syscall"
	"runtime"
	"strings"
	"reflect"
)

func CreateClient() *Client {
	return &Client{}
}

/*
validate the client with given token
*/
func (client *Client) Validate(token string, config *Config) bool {
	client.Authorized = false
	for _, tempuser := range config.Users {
		if token == tempuser.Token {
			client.Authorized = true
			client.User = *tempuser
			return client.Authorized
		}
	}
	return client.Authorized
}

/*
handles opened socket. after checking the authorization field at config, validates
given token, checks the command fields count, executes the command, converts to json
and finally writes on opened socket
*/
func (database *Database) HandleSocket(client *Client) {
	defer client.Connection.Close()
	strRemoteAddr := client.Connection.RemoteAddr().String()
	input, err := bufio.NewReader(client.Connection).ReadString('\n')
	if err != nil {
		database.config.logger("Read from socket error : "+err.Error(), 2)
		return
	}
	inputFields := strings.Fields(input)
	if len(inputFields) >= 2 {
		if inputFields[1] == "ping" {
			_, err := client.Connection.Write([]byte("pong\n"))
			if err != nil {
				database.config.logger("Error on writing to connection", 1)
			}
			return
		}
		if inputFields[1] == "status" {
			status := ""
			if database.config.Cluster.Master == "true" {
				status = database.cluster.GetStatus()
			} else {
				status = "Not Clustered"
			}
			_, err := client.Connection.Write([]byte(status+"\n"))
			if err != nil {
				database.config.logger("Error on writing to connection", 1)
			}
			return
		}
		if inputFields[1] == "exit" {
			database.config.logger("Client closed the connection from "+strRemoteAddr, 2)
			return
		}
		if inputFields[1] == "EOF" {
			database.config.logger("Client terminated the connection from "+strRemoteAddr, 2)
			return
		}
	}
	if database.config.Global.Authorization == "true" {
		if client.Validate(inputFields[0], database.config) {
			cmd := ""
			if len(inputFields) >= 3 {
				cmd = inputFields[1]
				for n := 2; n < len(inputFields); n++ {
					cmd += " " + inputFields[n]
				}
			}
			returneddb, message := database.cmd(cmd, &client.User)
			jsondb, err := json.Marshal(returneddb.items)
			if message != "" {
				_, err := client.Connection.Write([]byte(message+"\n"))
				if err != nil {
					database.config.logger("Error on writing to connection", 1)
				}
			}
			if err != nil {
				database.config.logger("Error in converting to json : "+err.Error(), 1)
			}
			if len(jsondb) > 2 {
				_, err := client.Connection.Write([]byte(string(jsondb)+"\n"))
				if err != nil {
					database.config.logger("Error on writing to connection", 1)
				}
			}
		} else {
			database.config.logger("Wrong token "+strRemoteAddr, 1)
			_ = client.Connection.Close()
		}
	} else {
		cmd := ""
		if len(inputFields) >= 3 {
			cmd = inputFields[1]
			for n := 2; n < len(inputFields); n++ {
				cmd += " " + inputFields[n]
			}
		}
		returneddb, message := database.cmd(cmd, &client.User)
		jsondb, err := json.Marshal(returneddb.items)
		if message != "" {
			_, err := client.Connection.Write([]byte(message+"\n"))
			if err != nil {
				database.config.logger("Error on writing to connection", 1)
			}
		}
		if err != nil {
			database.config.logger("Error in converting to json : "+err.Error(), 1)
		}
		if len(jsondb) > 2 {
			_, err := client.Connection.Write([]byte(string(jsondb)+"\n"))
			if err != nil {
				database.config.logger("Error on writing to connection", 1)
			}
		}
	}
}

/*
handles sighup from OS
*/
func (database *Database) HandleSigHup(){
	for {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGHUP)
		for range c {
			ok := false
			tempConfig, ok = ReloadConfig(os.Args[1])
			if ok {
				if !reflect.DeepEqual(tempConfig, database.config) {
					database.config = tempConfig
					database.cluster = database.config.CreateCluster()
					database.config.logger("Nipo reloaded", 1)
					database.Run()
				} else {
					database.config.logger("Nipo NOT reloaded, config file does not changed", 1)
				}
			}
		}
	}
}

/*
initialize the socket
*/
func (database *Database) InitSocket() {
	var err error
	database.config.logger("Opening Socket on "+database.config.Listen.Ip+":"+database.config.Listen.Port+"/"+database.config.Listen.Protocol, 1)
	database.socket, err = net.Listen(database.config.Listen.Protocol, database.config.Listen.Ip+":"+database.config.Listen.Port)
	if err != nil {
		database.config.logger("Error listening: "+err.Error(), 1)
		os.Exit(1)
	}
}

/*
called from main function, runs the service, multi-thread and multi-process handles here
calls the HandleSocket function
*/
func (database *Database) Run() {
	tempConfig = database.config
	go database.HandleSigHup()
	go database.RunCluster()
	defer database.socket.Close()
	runtime.GOMAXPROCS(database.config.Proc.Cores)
	for thread := 0; thread < database.config.Proc.Threads; thread++ {
		Wait.Add(1)
		go func() {
			defer Wait.Done()
			for {
				client := CreateClient()
				var err error
				client.Connection, err = database.socket.Accept()
				if err != nil {
					database.config.logger("Error accepting socket : "+err.Error(), 2)
				}
				database.HandleSocket(client)
			}
		}()
		Wait.Wait()
	}
}
