package main

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

func create(wg sync.WaitGroup, path string, size int) {
	defer wg.Done()
	log.Println("Create: " + path)
	handle, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	rand.Seed(time.Now().UnixNano())
	
	buf := new(bytes.Buffer)
	for i := 0; i < size; i++ {
		r := rand.Uint32
		binary.Write(buf, binary.LittleEndian, r)
		handle.Write(buf.Bytes())
	}
	handle.Close()
	log.Println("Created: " + path)
}

func read(wg sync.WaitGroup, path string, num int) {
	defer wg.Done()
	handle, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	
	for i := 0; i < num; i++ {
		
	}
}

func write(path string, num int) {
	defer wg.Done()
}

func main() {
	var wg sync.WaitGroup
	
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
	path := os.Args[1]
	create(path, size)
	wg.Add(readThr)
	for i := 0; i < readThr; i++ {
		go read(path, readNum)
	}
	wg.Add(writeThr)
	for i := 0; i < writeThr; i++ {
		go write(path, writeNum)
	}
	wg.Wait()
}
