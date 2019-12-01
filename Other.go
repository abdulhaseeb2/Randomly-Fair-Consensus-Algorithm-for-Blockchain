package IBC_Project

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

func Me(conn net.Conn, port string, name string) {
	gobEncoder := gob.NewEncoder(conn)
	err := gobEncoder.Encode(&port)
	if err != nil {
		log.Println(err)
	}

	gobEncoder = gob.NewEncoder(conn)
	err = gobEncoder.Encode(&name)
	if err != nil {
		log.Println(err)
	}
}

var TempChain *Block
var Name string
var Storage []int

func ListenToPeers(peers []string, port string, host string) {

	n, err := net.Listen("tcp", port)
	if err != nil {
		log.Println(err)
	}

	for {
		conn, err := n.Accept()
		if err != nil {
			log.Println(err)
		}
		println("\nGot connected")
		trans := ""

		dec := gob.NewDecoder(conn)
		_ = dec.Decode(&trans)

		if trans == "Transaction" {
			trans = ""
			println("\nRecieved Tranaction")
			dec := gob.NewDecoder(conn)
			err = dec.Decode(&trans)
			if err != nil {
				log.Println(err)
			}

			println("Choosen as Miner.")

			go VerifyAndBroadcast(trans, peers, true)

			conn.Close()

		} else if trans == "Broadcast" {
			conn.Close()
			conn, err := n.Accept()
			if err != nil {
				log.Println(err)
			}
			go BroadCastMeInTheAir(conn, peers)
		} else if trans == "Selector" {
			conn.Close()
			conn, err := net.Dial("tcp", host)
			if err != nil {
				log.Println(err)
			}
			go GenStorageValue(conn)
		} else if trans == "Potential Miner" {
			conn.Close()
			conn, err := net.Dial("tcp", host)
			if err != nil {
				log.Println(err)
			}
			go GetFreeStorage(conn)
		}
	}
}

func GenStorageValue(conn net.Conn) {

	println("\nSelector")
	z := GenValue()

	gobEncoder := gob.NewEncoder(conn)
	_ = gobEncoder.Encode(&z)

	println("Random Value: ", z)
	conn.Close()
}

func GetFreeStorage(conn net.Conn) {
	println("\nPotential Miner")
	z := Storage[GenValue()]

	gobEncoder := gob.NewEncoder(conn)
	_ = gobEncoder.Encode(&z)

	println("Free Storage: ", z)
	conn.Close()
}

func GenValue() int {

	rand.Seed(time.Now().UTC().UnixNano())
	var x int
	var y int
	x = RandInt1(1, 499)
	y = RandInt1(500, 999)
	x = ((x/y)/(x*y)*(x+y) + (x * y))
	y = RandInt1(1000, 1500)
	x = (x + y) % 1000 //change 100 with max nodes here
	var z float64
	z = math.Abs(float64(x)) //make a neg value positive

	return int(z)
}

func RandInt1(min int, max int) int {
	//this num will be the input number that you will pass from your code
	num := 23
	return min + rand.Intn(max-min)*num
}

func BroadCastMeInTheAir(conn net.Conn, peers []string) {
	var blk Block
	println("\nRecieved Broadcast")
	for {
		//println("Mhere")
		dec := gob.NewDecoder(conn)
		err := dec.Decode(&blk)
		if err != nil {
			log.Println(err)
		}
		//println("Mhere")
		if blk.Transaction != "" {
			trans := "ok"
			println("\nok\n")
			gobEncoder := gob.NewEncoder(conn)
			err := gobEncoder.Encode(&trans)
			if err != nil {
				log.Println(err)
			}
			break
		} else {
			println("\nbad\n")
			trans := "ok"
			gobEncoder := gob.NewEncoder(conn)
			err := gobEncoder.Encode(&trans)
			if err != nil {
				log.Println(err)
			}
			break
		}
	}

	if blk.HashValue != TempChain.HashValue {
		VerifyAndBroadcast(blk.Transaction, peers, false)
	} else {
		println("\nDropping blk", blk.HashValue)
	}
	conn.Close()
}

func CmdINPUT(host string, port string) {
	for {
		print("Enter Transaction(e.g. Alice 50 -> Bob): ")
		reader := bufio.NewReader(os.Stdin) //create new reader, assuming bufio imported
		var storageString string
		storageString, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err)
		}
		conn, err := net.Dial("tcp", host)

		if err != nil {
			log.Println(err)
		}

		gobEncoder := gob.NewEncoder(conn)
		err = gobEncoder.Encode(&storageString) //transaction

		gobEncoder = gob.NewEncoder(conn)
		err = gobEncoder.Encode(&port) //port number
		if err != nil {
			log.Println(err)
		}
		println("\nTransaction Sent.")
		conn.Close()
	}
}

func DecryptBC(chainHead *Block) []string {

	if chainHead.PreviousBlock == nil { //genesis Node

		var str = []string{chainHead.Transaction}
		return str

	}
	return append(DecryptBC(chainHead.PreviousBlock), chainHead.Transaction)
}

func SendBroadCast(peers []string, blk Block) {
	for _, broadcast := range peers {
		conn, err := net.Dial("tcp", broadcast)
		if err != nil {
			log.Println(err)
		}

		//rand.Seed(time.Now().UnixNano())
		//n := rand.Intn(10) // n will be between 0 and 10
		//time.Sleep(time.Duration(n) * time.Second)

		println("Broadcasting To: ", broadcast)
		trans1 := "Broadcast"
		gobEncoder := gob.NewEncoder(conn)
		err = gobEncoder.Encode(&trans1)
		if err != nil {
			log.Println(err)
		}
		println("BroadcastingString sent")
		conn.Close()
		conn, err = net.Dial("tcp", broadcast)
		if err != nil {
			log.Println(err)
		}
		for {
			//println("MhereB")
			gobEncoder = gob.NewEncoder(conn)
			err = gobEncoder.Encode(&blk)
			if err != nil {
				log.Println(err)
			}
			//println("MhereB")
			dec := gob.NewDecoder(conn)
			err = dec.Decode(&trans1)
			if err != nil {
				log.Println(err)
			}
			//println("MhereB")
			println("\nBroadcasting blk: ", blk.HashValue)

			if trans1 == "ok" {
				break
			}
		}

		conn.Close()
	}
}

func VerifyAndBroadcast(trans string, peers []string, bo bool) {
	split := strings.Split(trans, " ")
	balance := GetBalance(split[0], TempChain)
	println("\nAvailable Balance for transacton: ", balance)
	req, err := strconv.Atoi(split[1])
	if err != nil {
		fmt.Println("\nInvalid Required Balance\nProgram Terminating")
		req = balance + 50
	}
	if req <= balance { //Commit block and broadcast
		if bo {
			trans = trans + " \n 50 -> " + Name
		}
		TempChain = InsertBlock(trans, TempChain)
		//add threads wala function
		var blk Block
		blk.HashValue = TempChain.HashValue
		blk.Transaction = TempChain.Transaction
		blk.PreviousBlock = nil

		println("\n", TempChain.Transaction)
		ListBlocks(TempChain)
		println("\n")

		go SendBroadCast(peers, blk)

	} else {
		println("\nInsufficient funds")
	}
}
