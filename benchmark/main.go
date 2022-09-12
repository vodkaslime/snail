package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/vodkaslime/snail"
	"github.com/vodkaslime/snail/utils"
)

const BenchTestWordCount int = 1000000
const BenchTestFilePath string = "./test.db"

func benchTestCase(tableSize int) error {
	err := utils.ClearFile(BenchTestFilePath)
	if err != nil {
		return err
	}

	cfg := snail.Config{
		TableSize: tableSize,
		FilePath:  BenchTestFilePath,
	}
	fs := snail.NewLocalFS(&cfg)

	wordsPool := []string{
		"yo this is really good",
		"the quick brown fox jumps over the lazy dog",
		"how about this",
	}

	// Write to the FS
	println("writing to fs")
	start := time.Now()
	for i := 0; i < BenchTestWordCount; i++ {
		w := wordsPool[rand.Intn(len(wordsPool))]

		b := []byte(w)
		_, err := fs.Write(&b)
		if err != nil {
			return err
		}
	}

	println("write took time " + time.Since(start).String())

	err = utils.ClearFile(BenchTestFilePath)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	tableSizes := []int{100, 200, 400, 1000, 2000, 4000, 10000, 20000, 40000, 100000}
	for _, s := range tableSizes {
		println(fmt.Sprintf("current test case with table size: %d", s))
		err := benchTestCase(s)
		if err != nil {
			println(err.Error())
			return
		}
	}
}
