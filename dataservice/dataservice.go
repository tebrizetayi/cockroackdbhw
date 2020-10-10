package dataservice

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/tebrizetayi/cockroackdbhw/model"
	"github.com/twinj/uuid"
)

type IGormClient interface {
	UpdateAccount(ctx context.Context, accountData model.AccountData) (model.AccountData, error)
	StoreAccount(ctx context.Context, accountData model.AccountData) (model.AccountData, error)
	QueryAccount(ctx context.Context, accountId string) (model.AccountData, error)
	SetupDB(addr string)
	SeedAccounts() error
	Check() bool
	Close()
}

type GormClient struct {
	crDB *gorm.DB
}

func (gc *GormClient) Check() bool {
	return gc.crDB != nil
}

func (gc *GormClient) Close() {
	log.Println("Closing connection to CockroachDB")
	gc.crDB.Close()
}

// StoreAccount uses ACID TX.
func (gc *GormClient) StoreAccount(ctx context.Context, accountData model.AccountData) (model.AccountData, error) {

	if gc.crDB == nil {
		return model.AccountData{}, fmt.Errorf("Connection to DB not established!")
	}
	accountData.ID = uuid.NewV4().String()

	tx := gc.crDB.Begin()
	tx = tx.Create(&accountData)
	if tx.Error != nil {

		log.Fatalf("Error creating AccountData: %v \n", tx.Error.Error())
		return model.AccountData{}, tx.Error
	}
	tx = tx.Commit()
	if tx.Error != nil {
		log.Fatalf("Error committing AccountData: %v", tx.Error.Error())
		return model.AccountData{}, tx.Error
	}
	log.Println("Successfully created AccountData instance")
	return accountData, nil
}

// UpdateAccount uses ACID TX.
func (gc *GormClient) UpdateAccount(ctx context.Context, accountData model.AccountData) (model.AccountData, error) {

	if gc.crDB == nil {
		return model.AccountData{}, fmt.Errorf("Connection to DB not established!")
	}
	tx := gc.crDB.Begin()
	tx = tx.Save(&accountData)
	if tx.Error != nil {
		log.Fatalf("Error updating AccountData: %v", tx.Error.Error())
		return model.AccountData{}, tx.Error
	}
	tx.Commit()
	if tx.Error != nil {
		log.Fatalf("Error committing AccountData: %v", tx.Error.Error())
		return model.AccountData{}, tx.Error
	}
	log.Println("Successfully updated AccountData instance")

	// Read object from DB before return.
	accountData, _ = gc.QueryAccount(ctx, accountData.ID)
	return accountData, nil
}

func (gc *GormClient) QueryAccount(ctx context.Context, accountId string) (model.AccountData, error) {

	if gc.crDB == nil {
		return model.AccountData{}, fmt.Errorf("connection to DB not established!")
	}
	tx := gc.crDB.Begin()
	acc := model.AccountData{}
	tx = tx.First(&acc, "ID = ?", accountId)
	if tx.Error != nil {
		return acc, tx.Error
	}
	if acc.ID == "" {
		return acc, fmt.Errorf("no account found having ID %v", accountId)
	}
	tx.Commit()
	return acc, nil
}

func (gc *GormClient) SetupDB(addr string) {
	log.Println("Connecting with connection string: '%v'", addr)
	var err error
	gc.crDB, err = gorm.Open("postgres", addr)
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	gc.crDB.DB().SetMaxOpenConns(1000)

	gc.crDB.DB().SetMaxOpenConns(25)
	gc.crDB.DB().SetMaxIdleConns(25)
	gc.crDB.DB().SetConnMaxLifetime(5 * time.Minute)
	// Migrate the schema
	gc.crDB.AutoMigrate(&model.AccountData{})
}

func (gc *GormClient) SeedAccounts() error {
	if gc.crDB == nil {
		return fmt.Errorf("connection to DB not established")
	}
	gc.crDB.Delete(&model.AccountData{})
	total := 1000
	for i := 0; i < total; i++ {

		// Generate a key 10000 or larger
		key := strconv.Itoa(10000 + i)

		// Create an instance of our Account struct
		acc := model.AccountData{
			ID:   key,
			Name: "Person_" + strconv.Itoa(i),
		}

		gc.crDB.Create(&acc)
	}
	log.Println("Successfully created %v account instances.", 100)
	return nil
}
