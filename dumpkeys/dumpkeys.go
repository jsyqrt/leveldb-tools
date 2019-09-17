package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/syndtr/goleveldb/leveldb"
)

var DB *leveldb.DB
var File *os.File
var WithValueSize bool

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

func initFile(filepath string) {
	f, err := os.Create(filepath)
	checkErr(err)
	File = f
}

func initWithValueSize(withValueSize bool) {
	WithValueSize = withValueSize
}

func loopAll(db *leveldb.DB, mapFunc func(key []byte, value []byte)) error {
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()

		mapFunc(key, value)
	}
	iter.Release()
	err := iter.Error()
	return err
}

func writeKeyToFile(key []byte, value []byte) {
	_, err := File.Write(key)
	checkErr(err)

	if WithValueSize {
		_, err = File.Write([]byte("\t" + strconv.Itoa(len(value)) + "\n"))
		checkErr(err)
	} else {
		_, err = File.Write([]byte("\n"))
		checkErr(err)
	}
}

func main() {
	usage := `
	usage: ./dumpkeys /path/to/db/dir /path/to/save/keyfile boolwithvaluesize

	example: ./dumpkeys /home/ubuntu/.cpchain/cpchain/chaindata ./keyfile-with-value-size true
	example: ./dumpkeys /home/ubuntu/.cpchain/cpchain/chaindata ./keyfile-without-value-size false
	`

	if len(os.Args) != 4 {
		fmt.Println(usage)
		panic("error args")
	}

	dbPath := os.Args[1]
	outputPath := os.Args[2]
	withValueSize := func() bool {
		if os.Args[3] == "true" {
			return true
		}
		return false
	}()

	initDB(dbPath)
	initFile(outputPath)
	initWithValueSize(withValueSize)

	err := loopAll(DB, writeKeyToFile)
	fmt.Println("final err", err)

	err = File.Sync()
	checkErr(err)

	DB.Close()
	File.Close()
}
