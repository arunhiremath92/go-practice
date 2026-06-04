package main

import (
	"fmt"
	"math/rand/v2"

	"github.com/aruhiremath92/go-pg-test/models"
)

func main() {
	dbHandler, err := models.ModelsInstance()
	if err != nil {
		panic("failed to get a db connection")
	}

	err = dbHandler.AddUserDetails(rand.IntN(100), "")
	if err != nil {
		fmt.Println("failed to insert data in to the db")
	}
	userDetails := dbHandler.GetUserDetails()
	for _, user := range userDetails {
		fmt.Println("row :", user.Age, user.Name)
	}
}
