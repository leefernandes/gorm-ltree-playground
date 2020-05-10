package main

import (
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/kelseyhightower/envconfig"
)

type environment struct {
	PostgresHost     string `envconfig:"POSTGRES_HOST"`
	PostgresPort     string `envconfig:"POSTGRES_PORT"`
	PostgresDatabase string `envconfig:"POSTGRES_DATABASE"`
	PostgresUser     string `envconfig:"POSTGRES_USER"`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD"`
}

type thing struct {
	gorm.Model
	UUID            string
	Owner           string
	Name            string
	NameModifiedBy  string
	Color           string
	ColorModifiedBy string
	Code            string
	CodeModifiedBy  string
	Path            string `gorm:"type:ltree;"`
}

func main() {

	var env environment
	err := envconfig.Process("", &env)
	if err != nil {
		log.Fatal("unable to parse environment", err)
	}

	connectionStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		env.PostgresHost,
		env.PostgresPort,
		env.PostgresUser,
		env.PostgresPassword,
		env.PostgresDatabase,
	)

	db, err := gorm.Open("postgres", connectionStr)
	if err != nil {
		log.Fatal("failed to connect database", err)
	}
	defer db.Close()

	db.Exec("create extension ltree;")
	db.Exec("create extension btree_gin;")
	db.Exec("create extension btree_gist;")

	db.AutoMigrate(&thing{})

	// check if we have things
	var count int
	if err := db.Table("things").Count(&count).Error; err != nil {
		log.Fatal(err)
	}

	if 0 == count {
		createThings(db)
	}

	var results []thing
	// get A.AA, and all descendents
	if err := db.Debug().Where("path <@ ?", "A.AA").Find(&results).Error; err != nil {
		log.Fatal(err)
	}
	spew.Dump("A.AA results:", results)

	// get A.AB.ABC, and all descendents
	results = []thing{}
	if err := db.Debug().Where("path <@ ?", "A.AB.ABC").Find(&results).Error; err != nil {
		log.Fatal(err)
	}
	spew.Dump("A.AB.ABC results:", results)
}

func createThings(db *gorm.DB) {

	things := []*thing{
		// org A
		{
			Color: "orange",
			Path:  "A",
			Owner: "A",
		},
		// first sub org
		{
			Color: "yellow",
			Path:  "A.AA",
			Owner: "AA",
		},
		{
			Color: "red",
			Path:  "A.AA.AAA",
			Owner: "AAA",
		},
		{
			Color: "teal",
			Path:  "A.AA.AAA.AAAA",
			Owner: "AAAA",
		},
		{
			Color: "teal",
			Path:  "A.AA.AAA.AAAB",
			Owner: "AAAB",
		},
		// second sub org
		{
			Color: "yellow",
			Path:  "A.AB",
			Owner: "AB",
		},
		{
			Color: "yellow",
			Path:  "A.AB.ABA",
			Owner: "ABA",
		},
		{
			Color: "yellow",
			Path:  "A.AB.ABB",
			Owner: "ABB",
		},
		{
			Color: "yellow",
			Path:  "A.AB.ABC",
			Owner: "ABC",
		},
		{
			Color: "yellow",
			Path:  "A.AB.ABC.ABCA",
			Owner: "ABCA",
		},
		{
			Color: "yellow",
			Path:  "A.AB.ABC.ABCB",
			Owner: "ABCB",
		},
		{
			Color: "yellow",
			Path:  "A.AB.ABC.ABCC",
			Owner: "ABCC",
		},
		// third sub org
		{
			Color: "yellow",
			Path:  "A.AC",
			Owner: "AC",
		},
	}

	for i := range things {
		thing := things[i]
		thing.UUID = uuid.New().String()
		db.Create(thing)
	}

}
