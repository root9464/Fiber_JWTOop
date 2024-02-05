package main

import (
	"fmt"
	"log"
	"root/database"
	"root/services/root"
	"root/services/users"
)



func main() {
    fmt.Println("hello")
	method := database.MethodDB()
    db, err := method.Connect()
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных. \n", err)
	}
	err = db.AutoMigrate(&users.User{}, &users.Token{})
	if err != nil {
		panic("failed to perform migration")
	}else {
		fmt.Println("migration successful")
	}
	root.Root(db)
	
}