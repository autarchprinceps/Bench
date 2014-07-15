package main

import (
	"log"
	"os"
	"math/rand"
	"time"
	"encoding/binary"
	"bytes"
	"strconv"
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
		r := rand.Uint32
		binary.Write(buf, binary.LittleEndian, r)
		handle.Write(buf.Bytes())
	}
	handle.Close()
	log.Println("Created: " + path)
}

func main() {
	size, err := strconv.Atoi(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}
	thr, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Fatal(err)
	}
	create(os.Args[1], size)
}
