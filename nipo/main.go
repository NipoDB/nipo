package main

import (
    "github.com/NipoDB/nipo/nipo/database"
    "github.com/NipoDB/nipo/nipo/socket"
)

func main() {
    database := CreateDatabase()
    database.InitSocket()
    database.Run(0)
}