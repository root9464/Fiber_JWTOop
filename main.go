package main

import (
	"fmt"
	"log"
	"root/database"
	"root/services/root"
)



func main() {
    fmt.Println("hello")
	method := database.MethodDB()
    db, err := method.Connect()
    if err != nil {
        log.Fatal("Ошибка подключения к базе данных. \n", err)
    }
	root.Root(db)
	
}