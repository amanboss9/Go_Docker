package core

import (
	"Go_Docker/util"
	"encoding/csv"
	"io"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/mgo.v2"
)

// Mongo structure
type Mongo struct {
	Key   string `json:"key" bson:"key"`
	Value string `json:"value" bson:"value"`
}

//LoadCSV will load the csv
func LoadCSV() {

	session := util.MongoSession.Copy()
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	collection := session.DB(util.Config.DbName).C("csvload")
	absPath, _ := filepath.Abs("../Go_Docker/data/convertcsv.csv")
	file, err := os.Open(absPath)

	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		if record[0] != "key" {
			err = collection.Insert(&Mongo{Key: record[0], Value: record[1]})

			if err != nil {
				panic(err)
			}
			log.Printf("%#v", record)
		}

	}
}
