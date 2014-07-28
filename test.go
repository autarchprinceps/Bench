package main

import (
	"log"
	"os"
	"math/rand"
	"time"
	"encoding/binary"
	"strconv"
	"sync"
)

var randData [67108864]byte // = 64 MiB
var leng int

func create(path string, blockCount int) {
	log.Println("Create: " + path)
	handle, err := os.Create(path)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	for j := 0; j < blockCount; j++ {
		handle.Write(randData[0:])
	}
	log.Println("Created: " + path)
}

func read(wg *sync.WaitGroup, path string, num int, blockSize int) {
	defer wg.Done()
	handle, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	statHandle, err := handle.Stat()
	if err != nil {
		log.Fatal(err)
	}
	size := statHandle.Size()
	
	data := make([]byte, blockSize)
	for i := 0; i < num; i++ {
		posFile := rand.Int63n(size - (int64)(blockSize + 20))
		_, err := handle.ReadAt(data, posFile)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func write(wg *sync.WaitGroup, path string, num int, blockSize int) {
	defer wg.Done()
	handle, err := os.OpenFile(path, os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	statHandle, err := handle.Stat()
	if err != nil {
		log.Fatal(err)
	}
	size := statHandle.Size()
	
	for i := 0; i < num; i++ {
		posFile := rand.Int63n(size - (int64)(blockSize + 20))
		posRand := rand.Intn(leng - blockSize)
		_, err := handle.WriteAt(randData[posRand:(posRand+blockSize)], posFile)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	leng = len(randData)
	for i := 0; i < (leng / 8) - 1; i++ {
		binary.PutVarint(randData[i*8:], rand.Int63())
	}
	
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
	// log.Print("preCreate")
	preCreate := time.Now()
	if size > 0 {
		create(path, size)
	}
	// log.Print("preRDWR")
	preRDWR := time.Now()
	wg.Add(readThr)
	for i := 0; i < readThr; i++ {
		go read(&wg, path, readNum, blockSize)
	}
	wg.Add(writeThr)
	for i := 0; i < writeThr; i++ {
		go write(&wg, path, writeNum, blockSize)
	}
	// log.Print("wait")
	wg.Wait()
	// log.Print("post")
	post := time.Now()
	
	// log.Print("Create start: ", preCreate)
	// log.Print("Read write start: ", preRDWR)
	// log.Print("End: ", post)
	
	log.Print("Creation rate [MiB/s]: ", (float64)(64000000000.0) * (float64)(size) / (float64)(preRDWR.Sub(preCreate).Nanoseconds()) ) // 64 MB * (Nanosec -> Sec)
	readAmount := readThr * readNum * blockSize
	log.Print("Read amount [MiB]: ", (float64)(readAmount) / (float64)(1024*1024))
	writeAmount := writeThr * writeNum * blockSize
	log.Print("Write amount [MiB]: ", (float64)(writeAmount) / (float64)(1024*1024))
	log.Print("Edit rate [MiB/s]: ", (float64)(1000000000) * (float64)(readAmount + writeAmount) / (float64)(post.Sub(preRDWR).Nanoseconds() * 1024 * 1024) )
}
