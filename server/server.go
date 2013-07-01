package main

import (
	"fmt"
	"github.com/nicolaspaton/rhube"
	"io"
	"math"
	"net"
	"os"
	"strconv"
	"time"
)

func main() {
	fmt.Println("Rhube Server started.")
	ln, err := net.Listen("tcp", ":6380")
	if err != nil {
		panic(err)
	}
	db := rhube.NewDB()
	db.Set("toto", []byte("hey"))
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		go handleConn(conn, db)
	}
}

func handleConn(c net.Conn, db *rhube.DB) {
	defer c.Close()
	r := rhube.NewWireReader(c)
	w := rhube.NewWireWriter(c)
	for {
		args, err := r.ReadCommand()
		if err != nil {
			if err == io.EOF {
				fmt.Println("Connection closed.")
				return
			} else {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		handleCommand(w, db, args)
	}
}

func handleCommand(w *rhube.WireWriter, db *rhube.DB, args []string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			if r.(error) != nil {
				w.WriteErrorReply(r.(error).Error())
			}
		}
	}()

	switch args[0] {
	case "info":
		w.WriteBulkReply([]byte("hello:world"))

	// Keys

	case "del":
		w.WriteIntReply(db.Del(args[1:len(args)]...))
	case "keys":
		w.WriteStringMultiBulkReply(db.Keys(args[1]))
	case "rename":
		err := db.Rename(args[1], args[2])
		if hasNoError(w, err) {
			w.WriteStatusReply("OK")
		}
	case "renamenx":
		i, err := db.Renamenx(args[1], args[2])
		if hasNoError(w, err) {
			w.WriteIntReply(i)
		}

	case "persist":
		w.WriteBoolReply(db.Persist(args[1]))
	case "expire":
		intReq, err := strconv.Atoi(args[2])
		if hasNoError(w, err) {
			w.WriteBoolReply(db.Expire(args[1], intReq))
		}
	case "pexpire":
		intReq, err := strconv.Atoi(args[2])
		if hasNoError(w, err) {
			w.WriteBoolReply(db.Pexpire(args[1], intReq))
		}
	case "expireat":
		epoch, err := strconv.Atoi(args[2])
		if hasNoError(w, err) {
			t := time.Unix(int64(epoch), 0)
			w.WriteBoolReply(db.Expireat(args[1], t))
		}
	case "pexpireat":
		epoch, err := strconv.Atoi(args[2])
		if hasNoError(w, err) {
			seconds := math.Floor(float64(epoch) / 1000.0)
			milliseconds := float64(epoch) - (seconds * 1000.0)
			t := time.Unix(int64(seconds), int64(milliseconds)*1000000)
			w.WriteBoolReply(db.Expireat(args[1], t))
		}
	case "ttl":
		w.WriteIntReply(db.TTL(args[1]))
	case "pttl":
		w.WriteIntReply(db.PTTL(args[1]))

	// Strings

	case "get":
		w.WriteBulkReply(db.Get(args[1]))
	case "set":
		db.Set(args[1], []byte(args[2]))
		w.WriteStatusReply("OK")
	case "getset":
		resp, err := db.Getset(args[1], []byte(args[2]))
		if hasNoError(w, err) {
			w.WriteBulkReply(resp)
		}
	case "decr":
		resp, err := db.Decr(args[1])
		if hasNoError(w, err) {
			w.WriteIntReply(resp)
		}
	case "decrby":
		intReq, err := strconv.Atoi(args[2])
		if hasNoError(w, err) {
			intResp, err := db.Decrby(args[1], intReq)
			if hasNoError(w, err) {
				w.WriteIntReply(intResp)
			}
		}
	case "incr":
		resp, err := db.Incr(args[1])
		if hasNoError(w, err) {
			w.WriteIntReply(resp)
		}
	case "incrby":
		intReq, err := strconv.Atoi(args[2])
		if hasNoError(w, err) {
			intResp, err := db.Incrby(args[1], intReq)
			if hasNoError(w, err) {
				w.WriteIntReply(intResp)
			}
		}
	case "incrbyfloat":
		floatReq, err := strconv.ParseFloat(args[2], 64)
		if hasNoError(w, err) {
			strResp, err := db.Incrbyfloat(args[1], floatReq)
			if hasNoError(w, err) {
				w.WriteBulkReply([]byte(strResp))
			}
		}
	case "append":
		w.WriteIntReply(db.Append(args[1], []byte(args[2])))
	case "strlen":
		w.WriteIntReply(db.Strlen(args[1]))
	case "getrange":
		start, err := strconv.Atoi(args[2])
		if hasNoError(w, err) {
			stop, err := strconv.Atoi(args[3])
			if hasNoError(w, err) {
				resp := db.Getrange(args[1], start, stop)
				if hasNoError(w, err) {
					w.WriteBulkReply(resp)
				}
			}
		}
	case "setrange":
		start, err := strconv.Atoi(args[2])
		if hasNoError(w, err) {
			intResp := db.Setrange(args[1], start, []byte(args[3]))
			if hasNoError(w, err) {
				w.WriteIntReply(intResp)
			}
		}
	case "mget":
		w.WriteMultiBulkReply(db.Mget(args[1:len(args)]...))

	// Sets

	case "mset":
		db.Mset(args[1:len(args)]...)
		w.WriteStatusReply("OK")
	case "msetnx":
		w.WriteBoolReply(db.Msetnx(args[1:len(args)]...))
	case "sadd":
		w.WriteIntReply(db.Sadd(args[1], args[2:len(args)]...))
	case "srem":
		w.WriteIntReply(db.Srem(args[1], args[2:len(args)]...))
	case "sismember":
		w.WriteBoolReply(db.Sismember(args[1], args[2]))
	case "scard":
		w.WriteIntReply(db.Scard(args[1]))
	case "sinter":
		w.WriteStringMultiBulkReply(db.Sinter(args[1:len(args)]...))
	case "sunion":
		w.WriteStringMultiBulkReply(db.Sunion(args[1:len(args)]...))
	case "sinterstore":
		w.WriteIntReply(db.Sinterstore(args[1], args[2:len(args)]...))
	case "sunionstore":
		w.WriteIntReply(db.Sunionstore(args[1], args[2:len(args)]...))
	case "smembers":
		w.WriteStringMultiBulkReply(db.Smembers(args[1]))
	case "smove":
		w.WriteBoolReply(db.Smove(args[1], args[2], args[3]))
	case "srandmember":
		w.WriteBulkReply([]byte(db.Srandmember(args[1])))
	case "spop":
		w.WriteBulkReply([]byte(db.Spop(args[1])))

		// Hashes

	case "hget":
		w.WriteBulkReply([]byte(db.Hget(args[1], args[2])))
	case "hset":
		w.WriteBoolReply(db.Hset(args[1], args[2], args[3]))
	case "hsetnx":
		w.WriteBoolReply(db.Hsetnx(args[1], args[2], args[3]))
	case "hmset":
		w.WriteBoolReply(db.Hmset(args[1], args[2:len(args)]...))
	case "hmget":
		w.WriteStringMultiBulkReply(db.Hmget(args[1], args[2:len(args)]...))
	case "hexist":
		w.WriteBoolReply(db.Hexist(args[1], args[2]))
	case "hdel":
		w.WriteIntReply(db.Hdel(args[1], args[2:len(args)]...))
	case "hgetall":
		w.WriteHashMultiBulkReply(db.Hgetall(args[1]))
	case "hincrby":
		intReq, err := strconv.Atoi(args[3])
		if hasNoError(w, err) {
			intResp, err := db.Hincrby(args[1], args[2], intReq)
			if hasNoError(w, err) {
				w.WriteIntReply(intResp)
			}
		}
	case "hincrbyfloat":
		floatReq, err := strconv.ParseFloat(args[3], 64)
		if hasNoError(w, err) {
			strResp, err := db.Hincrbyfloat(args[1], args[2], floatReq)
			if hasNoError(w, err) {
				w.WriteBulkReply([]byte(strResp))
			}
		}
	case "hkeys":
		w.WriteStringMultiBulkReply(db.Hkeys(args[1]))
	case "hlen":
		w.WriteIntReply(db.Hlen(string(args[1])))
	case "hvals":
		w.WriteStringMultiBulkReply(db.Hvals(args[1]))

		// Lists

	case "lset":
		intReq, err := strconv.Atoi(args[2])
		if hasNoError(w, err) {
			err = db.Lset(args[1], intReq, []byte(args[3]))
			if hasNoError(w, err) {
				w.WriteStatusReply("OK")
			}
		}
	case "lindex":
		intReq, err := strconv.Atoi(args[2])
		if hasNoError(w, err) {
			w.WriteBulkReply(db.Lindex(args[1], intReq))
		}
	case "linsert":
		intResp, err := db.Linsert(args[1], args[2], []byte(args[3]), []byte(args[4]))
		if hasNoError(w, err) {
			w.WriteIntReply(intResp)
		}
	}

}

func hasNoError(w *rhube.WireWriter, err error) bool {
	if err == nil {
		return true
	}

	w.WriteErrorReply(err.Error())
	return false
}
