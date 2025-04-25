package algorithms

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

const maxHash = 256 // Keep the ID space small for demo

func hashKey(input string) int {
	hasher := sha1.New()
	hasher.Write([]byte(input))
	hash := hasher.Sum(nil)
	val, _ := hex.DecodeString(fmt.Sprintf("%x", hash))
	return int(val[0]) % maxHash
}

type Node struct {
	Name      string
	ID        int
	Data      map[string]string
	Successor *Node
}

func NewNode(name string) *Node {
	return &Node{
		Name: name,
		ID:   hashKey(name),
		Data: make(map[string]string),
	}
}

func (n *Node) SetSuccessor(succ *Node) {
	n.Successor = succ
}

func (n *Node) IsResponsible(keyHash int) bool {
	if n.Successor == nil {
		return true
	}
	if n.ID < n.Successor.ID {
		return keyHash > n.ID && keyHash <= n.Successor.ID
	}
	return keyHash > n.ID || keyHash <= n.Successor.ID
}

func (n *Node) Store(key string, value string) {
	keyHash := hashKey(key)
	if n.IsResponsible(keyHash) {
		fmt.Printf("%s is storing key '%s' => %s\n", n.Name, key, value)
		n.Data[key] = value
	} else {
		n.Successor.Store(key, value)
	}
}

func (n *Node) Find(key string) string {
	keyHash := hashKey(key)
	if n.IsResponsible(keyHash) {
		return n.Data[key]
	}
	return n.Successor.Find(key)
}
