package main

func main() {
    database := CreateDatabase()
    database.InitSocket()
    database.Run()
}