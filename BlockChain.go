package assignment02IBC

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
)

type Block struct {
	PreviousBlock *Block
	HashValue     string
	Transaction   string
}

func hashBlock(block string) string {

	// Create a new HMAC by defining the hash type and the key (as byte array)
	h := hmac.New(sha256.New, []byte("secret"))

	// Write Data to it
	h.Write([]byte(block))

	// Get result and encode as hexadecimal string
	sha := hex.EncodeToString(h.Sum(nil))

	return sha

}

func GetBalance(name string, chainHead *Block) int {
	if chainHead.PreviousBlock == nil { //genesis Node
		if strings.Contains(chainHead.Transaction, name) {
			bal := 0
			split := strings.Split(chainHead.Transaction, " ")
			if split[0] == name { //payee
				bal, _ = strconv.Atoi(split[1])
				bal = 0 - bal
			} else if split[2] == name { //miner
				bal, _ = strconv.Atoi(split[0])
			} else if split[3] == name { //reciver
				bal, _ = strconv.Atoi(split[1])
			} else if split[7] == name { //miner
				bal, _ = strconv.Atoi(split[5])
			}
			return bal
		}

		return 0

	} else {

		bal := GetBalance(name, chainHead.PreviousBlock)
		bal1 := 0
		if strings.Contains(chainHead.Transaction, name) {
			split := strings.Split(chainHead.Transaction, " ")
			if split[0] == name { //payee
				bal1, _ = strconv.Atoi(split[1])
				bal1 = 0 - bal1
			} else if split[2] == name { //miner
				bal1, _ = strconv.Atoi(split[0])
			} else if split[3] == name { //reciver
				bal1, _ = strconv.Atoi(split[1])
			} else if split[7] == name { //miner
				bal1, _ = strconv.Atoi(split[5])
			}
		}

		return bal + bal1
	}
}

func InsertBlock(transaction string, chainHead *Block) *Block {

	var newBlock *Block = new(Block)

	if chainHead == nil {

		newBlock.Transaction = transaction
		newBlock.PreviousBlock = chainHead
		newBlock.HashValue = ""
		println("Genesis Block Added")
	} else {

		newBlock.Transaction = transaction
		newBlock.PreviousBlock = chainHead
		newBlock.HashValue = hashBlock(chainHead.Transaction + chainHead.HashValue)
		println("New Block Added")
	}

	return newBlock
}

func VerifyChain(chainHead *Block) string {

	if chainHead.PreviousBlock == nil { //genesis Node

		return hashBlock(chainHead.Transaction)

	} else {
		blockHash := VerifyChain(chainHead.PreviousBlock)

		if blockHash == chainHead.HashValue {
			println("Hash Matches")

		} else {
			println("Hash Does Not Match")

		}

		return hashBlock(chainHead.Transaction + blockHash)

	}
}

func ChangeBlock(oldTrans string, newTrans string, chainHead *Block) {

	if chainHead.PreviousBlock == nil { //Genesis Block

		if chainHead.Transaction == oldTrans { //If required block is found

			chainHead.Transaction = newTrans
			println("Block Changed")
		}

	} else { //recursivly iterate to the required block

		ChangeBlock(oldTrans, newTrans, chainHead.PreviousBlock)

		if chainHead.Transaction == oldTrans { //If required block is found

			chainHead.Transaction = newTrans
			println("Block Changed")
		}
	}
}

func ListBlocks(chainHead *Block) {

	if chainHead.PreviousBlock == nil { //genesis Node

		println("Transaction: " + chainHead.Transaction)
		println("Genesis Block.\n")

	} else {
		println("Transaction " + chainHead.Transaction)
		println("Hash Value: " + chainHead.HashValue)

		ListBlocks(chainHead.PreviousBlock)
	}
}
