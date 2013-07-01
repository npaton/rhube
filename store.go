package rhube

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"log"
	"math/big"
	"os"
	"strconv"
)

var sharedBuffer []byte

func (db *DB) Load(fileName string) error {
	f, err := os.Open(fileName)
	defer f.Close()
	if err != nil {
		return err
	}

	sharedBuffer = make([]byte, 1000000) // ~1MB

	_, err = f.Read(sharedBuffer[:9])

	magickString := string(sharedBuffer)
	if err != nil {
		return err
	}

	if magickString[:5] != "REDIS" {
		return errors.New("File corrupt. Missing Magick string.")
	}

	if magickString[5:9] != "0006" {
		return errors.New("RDB version is not 6, only v6 supported atm. -" + magickString[5:9] + "- ")
	}

	// DB Number (FE+LengthEnc)
	_, err = f.Read(sharedBuffer[:1])

	dbNum := readLength(f)
	log.Println("Reading DB " + strconv.Itoa(int(dbNum)))

	// Process key value pairs
	readKeyValues(f)

	return nil
}

func readKeyValues(r io.Reader) int {
	for {
		r.Read(lengthBuf[:1])
		switch int(lengthBuf[:1][0]) {
		case 0:
			readStringKeyValue(r)
			// return 0
		case 12:
			readZipListKeyValue(r)
		default:
			// readStringKeyValue(r)
			log.Printf("VALUE: ", lengthBuf[:1], int(lengthBuf[:1][0]), hex.EncodeToString(lengthBuf[:1]))
		}
	}
}

func readZipListKeyValue(r io.Reader) (key string, value []byte, err error) {
	key = string(readRString(r))
	log.Println(key)
	value = readRString(r)
	log.Println(string(value))
	log.Println("============================")
	// 50c0b5bd760ee35582000dc5
	os.Exit(0)
	return "", []byte(""), nil
}

func readStringKeyValue(r io.Reader) (key string, value []byte, err error) {
	key = string(readRString(r))
	// log.Println(key)
	value = readRString(r)
	log.Println(key, string(value))
	return key, value, err
}

func readRString(r io.Reader) []byte {
	length := readLength(r)
	n, err := r.Read(sharedBuffer[:length])
	if uint64(n) != length || err != nil {
		log.Println("Error: readRedisString. (" + err.Error() + ") (" + strconv.Itoa(int(n)) + "not" + strconv.Itoa(int(length)) + ")")
	}
	return sharedBuffer[:length]
}

var lengthBuf []byte

func readLength(r io.Reader) uint64 {

	_, err := r.Read(lengthBuf[:1])
	if err != nil {
		panic(err)
	}

	num := 0 | 0 | 0 | (int64(lengthBuf[:1][0]) << 3)
	firstByte := big.NewInt(num)
	firstBit, secondBit := firstByte.Bit(1), firstByte.Bit(0)

	if firstBit == 1 && secondBit == 1 {
		firstByte.SetBit(firstByte, 0, 0)
		firstByte.SetBit(firstByte, 1, 0)
		log.Println("11")
		return 0
	}

	if firstBit == 1 && secondBit == 0 {
		firstByte.SetBit(firstByte, 0, 0)

		// _, err := r.Read(lengthBuf[:1])
		// log.Println("=>", lengthBuf[:1])
		fByte := lengthBuf[:1]
		_, err = r.Read(lengthBuf[:9])
		if err != nil {
			panic(err)
		}

		bits := append(fByte, lengthBuf[:9]...)
		log.Println("=>", string(bits))

		var i uint64
		buf := bytes.NewBuffer(bits)
		err = binary.Read(buf, binary.BigEndian, &i)
		if err != nil {
			log.Println("Error. 1 readLength. binary.Read failed:", err)
		}

		log.Println("10", bits, i, lengthBuf[:3])
		os.Exit(1)
		return i
	}

	if firstBit == 0 && secondBit == 1 {
		_, err := r.Read(lengthBuf[1:2])
		if err != nil {
			panic(err)
		}

		var extraByte uint8
		buf := bytes.NewBuffer(lengthBuf[:1])
		err = binary.Read(buf, binary.BigEndian, &extraByte)
		if err != nil {
			log.Println("Error. 3 readLength. binary.Read failed:", err)
		}

		firstByte.SetBit(firstByte, 1, 0)
		num := uint16(firstByte.Bytes()[0]) | (uint16(extraByte) << 8)
		log.Println("01", num)
		return uint64(num)
	}

	if firstBit == 0 && secondBit == 0 {
		num := uint64(lengthBuf[:1][0])
		log.Println("00 ->", num)
		return num
	}

	panic("fjdksljfklsdjfkl")
	return 0
}

func init() {
	lengthBuf = make([]byte, 100)
}
