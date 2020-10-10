package main

import (
	"context"
	"fmt"
	"log"

	"github.com/tebrizetayi/cockroackdbhw/dataservice"
	"github.com/tebrizetayi/cockroackdbhw/model"
)

func main() {

	insert()
}

func insert() {
	addr := "postgresql://root@localhost:26257/account?sslmode=disable"
	grmClient := &dataservice.GormClient{}
	grmClient.SetupDB(addr)
	ctx := context.Background()
	acc := model.AccountData{
		Name: "Tabriz",
	}
	createdAccount, err := grmClient.StoreAccount(ctx, acc)

	if err != nil {
		fmt.Errorf("Error happened %v \n", err)
	}
	fmt.Printf("Created Account with ID %s \n", createdAccount.ID)
	grmClient.Close()
	fmt.Println("Successfully Ended!")
}

func read(id string) {
	addr := "postgresql://root@localhost:26257/account?sslmode=disable"
	grmClient := &dataservice.GormClient{}
	grmClient.SetupDB(addr)
	ctx := context.Background()
	acc, err := grmClient.QueryAccount(ctx, id)
	if err != nil {
		fmt.Println("Reading error " + id)
	}
	log.Println(acc.ID + " is read!")

	grmClient.Close()
}
