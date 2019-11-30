package assignment02IBC

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
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

			go verifyAndBroadcast(trans, peers, true)

			conn.Close()

		} else if trans == "Broadcast" {
			conn.Close()
			conn, err := n.Accept()
			if err != nil {
				log.Println(err)
			}
			go broadCastMeInTheAir(conn, peers)
		}
	}
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
			trans := "bad"
			gobEncoder := gob.NewEncoder(conn)
			err := gobEncoder.Encode(&trans)
			if err != nil {
				log.Println(err)
			}
		}
	}

	if blk.HashValue != tempChain.HashValue {
		verifyAndBroadcast(blk.Transaction, peers, false)
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
	balance := GetBalance(split[0], tempChain)
	println("\nAvailable Balance for transacton: ", balance)
	req, err := strconv.Atoi(split[1])
	if err != nil {
		fmt.Println("\nInvalid Required Balance\nProgram Terminating")
		req = balance + 50
	}
	if req <= balance { //Commit block and broadcast
		if bo {
			trans = trans + " \n 50 -> " + name
		}
		tempChain = InsertBlock(trans, tempChain)
		//add threads wala function
		var blk Block
		blk.HashValue = tempChain.HashValue
		blk.Transaction = tempChain.Transaction
		blk.PreviousBlock = nil

		println("\n")
		ListBlocks(tempChain)
		println("\n")

		go SendBroadCast(peers, blk)

	} else {
		println("\nInsufficient funds")
	}
}
