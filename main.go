package main

import (
	"fmt"

	"einer.io/routes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	EmailAddress string `gorm:"not null;unique"`
	Username     string `gorm:"not null;unique"`
	Password     string `gorm:"not null"`
	Firstname    string `gorm:"not null;unique"`
	Lastname     string `gorm:"not null"`
}

type parametersDatabase struct {
	Hostname string
	Username string
	Password string
	Database string
	Port     uint16
	SSL      bool
	Timezone string
}

func main() {
	p := &parametersDatabase{
		Hostname: "localhost",
		Username: "postgres",
		Password: "password",
		Database: "argon2",
		Port:     5432,
		SSL:      false,
		Timezone: "Asia/Bangkok",
	}

	pg, err := databaseConnnection(p)
	if err != nil {
		panic(err)
	}

	pg.AutoMigrate(&Users{})

	r := routes.NewRoutes(pg)

	r.Run()
}

func databaseConnnection(p *parametersDatabase) (*gorm.DB, error) {
	sslmode := "disable"
	if p.SSL {
		sslmode = "enable"
	}

	dataSource := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s", p.Hostname, p.Username, p.Password, p.Database, p.Port, sslmode, p.Timezone)
	pg, err := gorm.Open(postgres.Open(dataSource), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return pg, nil
}
