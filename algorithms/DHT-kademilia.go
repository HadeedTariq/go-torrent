package algorithms

import (
	"crypto/sha1"
	"encoding/hex"
	"math/big"
	"sync"
)

const IDLength = 160
const BucketSize = 20

type NodeId [IDLength / 8]byte

func NewNodeID(data string) NodeId {
	hash := sha1.Sum([]byte(data))

	var id NodeId
	copy(id[:], hash[:])
	return id
}

// just for the purpose of human readability
func (id NodeId) String() string {
	return hex.EncodeToString(id[:])
}

func (id NodeId) XOR(other NodeId) *big.Int {
	var result [IDLength / 8]byte

	for i := 0; i < len(id); i++ {
		result[i] = id[i] ^ other[i]
	}

	return new(big.Int).SetBytes(result[:])
}

type Contact struct {
	ID      NodeId
	Address string
}

type RoutingTable struct {
	contacts []Contact
	selfId   NodeId
	mu       sync.RWMutex
}

func NewRoutingTable(selfId NodeId) *RoutingTable {
	return &RoutingTable{
		contacts: make([]Contact, 0),
		selfId:   selfId,
	}
}

func (rt *RoutingTable) Add(contact Contact) {
	//~ lets consider the multiple go routines are running concurrently that's why I have to create the lock

	rt.mu.Lock()
	defer rt.mu.Unlock()

	// ~ also have to avoid the duplicate contacts

	for _, con := range rt.contacts {
		if con.ID == contact.ID {
			return
		}
	}

	// so if it is full we have to evict

	if len(rt.contacts) < BucketSize {
		rt.contacts = append(rt.contacts, contact)
	} else {
		rt.contacts = append(rt.contacts[:1], contact)
	}
}
