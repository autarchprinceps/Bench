package main

/**
 * TODO
 * metrics
 */

import (
	"log"
	"os"
	"math/rand"
	"time"
	"encoding/binary"
	"bytes"
	"strconv"
	"sync"
)

func create(path string, size int) {
	log.Println("Create: " + path)
	handle, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	rand.Seed(time.Now().UnixNano())
	buf := new(bytes.Buffer)
	for i := 0; i < size; i++ {
		r := rand.Uint32()
		binary.Write(buf, binary.LittleEndian, r)
		handle.Write(buf.Bytes())
	}
	handle.Close()
	log.Println("Created: " + path)
}

func read(wg sync.WaitGroup, path string, num int, blockSize int) {
	defer wg.Done()
	handle, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	data := make([]byte, blockSize)
	for i := 0; i < num; i++ {
		// pos := rand.Int31n(len(handle)) TODO replace len(handle) with file length
		_, err := handle.ReadAt(data, pos)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func write(wg sync.WaitGroup, path string, num int, blockSize int) {
	defer wg.Done()
	handle, err := os.OpenFile(path, os.O_wronly, 0777) // TODO os.ModePerm
	if err != nil {
		log.Fatal(err)
	}
	data := make([]byte, blockSize)
	for i := 0; i < num; i++ {
		// TODO fill data with rand
		// TODO pos
		_, err := handle.WriteAt(data, pos)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	var wg sync.WaitGroup
	
	// param parsing
	size, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	readThr, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Fatal(err)
	}
	readNum, err := strconv.Atoi(os.Args[4])
	if err != nil {
		log.Fatal(err)
	}
	writeThr, err := strconv.Atoi(os.Args[5])
	if err != nil {
		log.Fatal(err)
	}
	writeNum, err := strconv.Atoi(os.Args[6])
	if err != nil {
		log.Fatal(err)
	}
	blockSize, err := strconv.Atoi(os.Args[7])
	if err != nil {
		log.Fatal(err)
	}
	path := os.Args[1]
	
	// fileops
	if size > 0 {
		create(path, size)
	}
	wg.Add(readThr)
	for i := 0; i < readThr; i++ {
		go read(wg, path, readNum, blockSize)
	}
	wg.Add(writeThr)
	for i := 0; i < writeThr; i++ {
		go write(wg, path, writeNum, blockSize)
	}
	wg.Wait()
}
