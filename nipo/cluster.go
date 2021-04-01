package main

import (
	"github.com/NipoDB/nipolib"
	"strconv"
	"time"
)

type Slave struct {
	Node              *Node
	Status, CheckedAt string
	Database          *Database
}

type Cluster struct {
	Slaves []Slave
	Status string
}

/*
sync database to slave
*/
func (database *Database) SyncSlave(config *Config, slave *Slave) {
	nipoConfig := nipo.CreateConfig(slave.Node.Token, slave.Node.Ip, slave.Node.Port)
	database.Foreach(func(key, value string) {
		_, ok := nipo.Set(nipoConfig, key, value)
		if !ok {
			config.logger("Set command on slave does not work correctly", 2)
		}
	})
}

/*
returs the status of cluster and slaves in json format
*/
func (cluster *Cluster) GetStatus() string {
	result := "{ "
	for index, slave := range cluster.Slaves {
		tempStr := "\"" + strconv.Itoa(slave.Node.Id) + "\" : " + "{ \"ip\" : \"" + slave.Node.Ip + "\" , \"status\" : \"" + slave.Status + "\" , \"checkedat\" : \"" + slave.CheckedAt + "\" }"
		if !(index == len(cluster.Slaves)-1) {
			tempStr += ","
		}
		result += tempStr
	}
	result += " }"
	return result
}

/*
create the cluster with config
*/
func (config *Config) CreateCluster() *Cluster {
	cluster := Cluster{}
	for _, slave := range config.Slaves {
		tempSlave := Slave{}
		tempSlave.Node = slave
		tempSlave.Status = "none"
		tempSlave.CheckedAt = "none"
		cluster.Slaves = append(cluster.Slaves, tempSlave)
	}
	cluster.Status = "none"
	return &cluster
}

/*
this function is the main of cluster, the health check and update state of cluster and slaves and 
also syncing the slaves is controlled with this function
*/
func (database *Database) RunCluster(config *Config, cluster *Cluster) {
	for {
		for index, slave := range cluster.Slaves {
			nipoconfig := nipo.CreateConfig(slave.Node.Token, slave.Node.Ip, slave.Node.Port)
			result, _ := nipo.Ping(nipoconfig)
			if result == "pong\n" {
				if slave.Status == "unhealthy" {
					slave.Status = "recover"
					config.logger("slave by id : "+strconv.Itoa(slave.Node.Id)+" becomes healthy", 1)
					config.logger("slave by id : "+strconv.Itoa(slave.Node.Id)+" is in recovery", 1)
					Lock.Lock()
					database.SyncSlave(config, &slave)
					Lock.Unlock()
					slave.Status = "healthy"
					config.logger("slave by id : "+strconv.Itoa(slave.Node.Id)+" recovery compleated", 1)
				}
				cluster.Slaves[index].Status = "healthy"
				cluster.Slaves[index].CheckedAt = time.Now().Format("2006-01-02 15:04:05.000")
				cluster.Status = "healthy"
			} else {
				if slave.Status == "healthy" {
					config.logger("slave by id : "+strconv.Itoa(slave.Node.Id)+" becomes unhealthy", 1)
				}
				cluster.Slaves[index].Status = "unhealthy"
				cluster.Slaves[index].CheckedAt = time.Now().Format("2006-01-02 15:04:05.000")
				cluster.Status = "unhealthy"
				config.logger("slave by id : "+strconv.Itoa(slave.Node.Id)+" is not healthy", 2)
			}
		}
		time.Sleep(time.Duration(config.Global.CheckInterval) * time.Millisecond)
	}
}

/*
sets the key and value on slaves of cluster
called from cmdSet function on each set command execution
*/
func (cluster *Cluster) SetOnSlaves(config *Config, key, value string) bool {
	for _, slave := range cluster.Slaves {
		if slave.Status == "healthy" {
			nipoconfig := nipo.CreateConfig(slave.Node.Token, slave.Node.Ip, slave.Node.Port)
			_, pingOK := nipo.Ping(nipoconfig)
			if pingOK {
				_, ok := nipo.Set(nipoconfig, key, value)
				if !ok {
					config.logger("Set command on slave does not work correctly", 2)
				}
			}
		}
		if slave.Status == "unhealthy" {
			slave.Database.Set(key, value)
		}
	}
	return true
}
