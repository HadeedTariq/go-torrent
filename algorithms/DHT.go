package algorithms

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

// so the first thing required in the dht is the node
const MAX_SIZE = 256

type Node struct {
	Name      string
	ID        int
	Data      map[string]string
	Successor *Node
}

// the other functionality that comes to the dht is the creating hash of the key

func NewNode(name string) *Node {
	return &Node{
		Name: name,
		ID:   keyHasher(name),
		Data: make(map[string]string),
	}
}
func keyHasher(key string) int {

	hasher := sha1.New()
	hasher.Write([]byte(key))
	hash := hasher.Sum(nil)

	val, _ := hex.DecodeString(fmt.Sprintf("%x", hash))

	return int(val[0]) % MAX_SIZE
}

// now the next thing that comes out is to store in the node

func (node *Node) SetSuccessor(newNode *Node) {
	node.Successor = newNode
}

func (node *Node) Store(key string, value string) {
	keyHash := keyHasher(key)
	if node.isResponsible(keyHash) {
		node.Data[key] = value
	} else {
		node.Successor.Store(key, value)
	}

}

func (node *Node) isResponsible(keyHash int) bool {
	// now what I have to do in is responsible is to run some checks
	// so how successor is decide

	if node.Successor == nil {
		return true
	}
	if node.ID < node.Successor.ID {
		return node.ID < keyHash && keyHash <= node.Successor.ID
	}

	return keyHash > node.ID || keyHash <= node.Successor.ID
}

// also have to implement the find functionality there

func (node *Node) Find(key string) string {
	keyHash := keyHasher(key)

	if node.isResponsible(keyHash) {
		return node.Data[key]
	} else {
		return node.Successor.Find(key)
	}
}
