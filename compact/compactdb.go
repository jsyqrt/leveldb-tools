package main

import (
	"fmt"
	"os"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var DB *leveldb.DB

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func initDB(dbPath string) {
	db, err := leveldb.OpenFile(dbPath, nil)
	checkErr(err)

	DB = db
}

func main() {
	usage := "usage: ./compact /path/to/db/dir"
	if len(os.Args) != 2 {
		fmt.Println(usage)
		panic("error args")
	}

	path := os.Args[1]

	initDB(path)

	err := DB.CompactRange(util.Range{nil, nil})

	fmt.Println("final err", err)

	DB.Close()
}
