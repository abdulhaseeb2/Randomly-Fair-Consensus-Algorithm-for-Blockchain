package IBC_Project

import (
	"encoding/gob"
	"log"
	"math"
	"math/rand"
	"net"
	"time"
)

func SendBlockchain(port net.Conn, chainHead *Block) {
	gobEncoder := gob.NewEncoder(port)
	err := gobEncoder.Encode(&chainHead)
	if err != nil {
		log.Println(err)
	}
}

func SendPeers(port net.Conn, peers []string) {
	gobEncoder := gob.NewEncoder(port)
	err := gobEncoder.Encode(peers)
	if err != nil {
		log.Println(err)
	}
}

func SendTrans(trans string, miner string) {
	conn, err := net.Dial("tcp", miner)
	if err != nil {
		log.Println(err)
	}
	println("\nMiner: ", miner)
	trans1 := "Transaction"

	gobEncoder := gob.NewEncoder(conn)
	err = gobEncoder.Encode(&trans1)

	gobEncoder = gob.NewEncoder(conn)
	err = gobEncoder.Encode(&trans)

	println("\nTransaction Sent")

	if err != nil {
		log.Println(err)
	}

	conn.Close()
}

func SelectMiner(nodes int, peersPort []string, port string) string {
	miner := 0
	for {
		miner = rand.Intn(nodes)
		if port != peersPort[miner] {
			break
		}
	}
	return peersPort[miner]
}

func Gen_Storage(port string, sat_port string) int {

	//Telling the Node that he is the Selector of the random Storage Variable.
	conn, err := net.Dial("tcp", port)
	if err != nil {
		log.Println(err)
	}

	trans := "Selector"
	println("\nSelector ", port)

	gobEncoder := gob.NewEncoder(conn)
	err = gobEncoder.Encode(&trans)
	conn.Close()

	n, err := net.Listen("tcp", sat_port)
	if err != nil {
		log.Println(err)
	}
	conn1, err := n.Accept()
	if err != nil {
		log.Println(err)
	}

	message := 0

	dec := gob.NewDecoder(conn1)
	_ = dec.Decode(&message)

	conn1.Close()
	n.Close()
	println("Random Storage Value Selected: ", message)

	return message
}

func UniqueConsensus(nodes int, peersPort []string, port string, sat_port string) string {
	//randomly selecting a selector
	selector := 0
	for {
		selector = myrandom(nodes)
		if port != peersPort[selector] {
			break
		}
	}
	//Then Generating a random storage Value by the selector
	randomStorage := Gen_Storage(peersPort[selector], sat_port)

	//randomly selecting a Potential Miner
	pot_Miners := []string{}
	pot_Storage := []int{}
	for i := 0; i < 3; i++ {
		pm := 0
		for {
			pm = myrandom(nodes)
			found := false
			if port != peersPort[pm] && peersPort[pm] != peersPort[selector] {
				for j := 0; j < i; j++ {
					if peersPort[pm] == pot_Miners[j] {
						found = true
					}
				}
				if !(found) {
					pot_Miners = append(pot_Miners, peersPort[pm])

					println("Potential Miner: ", pot_Miners[len(pot_Miners)-1])

					pot_Storage = append(pot_Storage, int(math.Abs(float64(GetStorage(peersPort[pm], sat_port)-randomStorage))))

					break
				}
			}
		}
	}

	//Selecting the miner on the bases of minimum Storage
	min := pot_Storage[0]
	index := 0
	for i := 1; i < 3; i++ {
		if pot_Storage[i] < min {
			min = pot_Storage[i]
			index = i
		}
	}
	println("Miner Choosen: ", pot_Miners[index])

	return pot_Miners[index]

}

func GetStorage(port string, sat_port string) int {
	//Telling the Node that he is the Potential Miner and that he needs to give me his free storage.
	conn, err := net.Dial("tcp", port)
	if err != nil {
		log.Println(err)
	}

	trans := "Potential Miner"
	println("\nPotential Miner ", port)

	gobEncoder := gob.NewEncoder(conn)
	err = gobEncoder.Encode(&trans)
	conn.Close()

	n, err := net.Listen("tcp", sat_port)
	if err != nil {
		log.Println(err)
	}
	conn1, err := n.Accept()
	if err != nil {
		log.Println(err)
	}

	message := 0

	dec := gob.NewDecoder(conn1)
	_ = dec.Decode(&message)

	conn1.Close()
	n.Close()
	println("Free Storage of Potential Miner: ", message)

	return message
}

func myrandom(maxVal int) int {

	rand.Seed(time.Now().UTC().UnixNano())
	var x int
	var y int
	x = randInt(1, 499)
	y = randInt(500, 999)
	x = ((x/y)/(x*y)*(x+y) + (x * y))
	y = randInt(1000, 1500)
	x = (x + y) % maxVal //change 100 with max nodes here
	var z float64
	z = math.Abs(float64(x)) //make a neg value positive

	return int(z)
}

func randInt(min int, max int) int {
	//this num will be the input number that you will pass from your code
	num := 23
	return min + rand.Intn(max-min)*num
}
