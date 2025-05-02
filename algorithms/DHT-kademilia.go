package algorithms

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math/big"
	"net"
	"sort"
	"time"
)

// ~ so I want to build a distributed messaging system

// ~ so the first thing that comes up in the peer is peer
const IdLength = 8
const IdBits = 16

type NodeID [IdBits / IdLength]byte

func NewNodeId(data []byte) NodeID {
	hash := sha1.Sum(data[:])
	var newId NodeID
	copy(newId[:], hash[:len(newId)])

	return newId
}

func (id NodeID) String() string {
	return hex.EncodeToString(id[:])
}

func (id NodeID) XOR(otherId NodeID) *big.Int {
	var result NodeID

	for i := 0; i < len(id); i++ {
		result[i] = id[i] ^ otherId[i]
	}

	return new(big.Int).SetBytes(result[:])
}

// ~ so the work of routing table is that it maintain the list of the closest nodes
// ~ and the routing table contain the  k buckets and each bucket contain the list of peers info and there should be let say 2 contacts info in each bucket and bucket size should be depend on the size of the id like for our case 16

const contactSize = 2

type Contacts struct {
	Id           NodeID
	Address      string
	last_seen_at time.Time
}

type KBucket struct {
	contacts []Contacts
}

// ~ so in k buckets I have to add the contacts and with in contacts I have to add the peer info
func (kb *KBucket) Add(contact Contacts) bool {
	if len(kb.contacts) < contactSize {
		kb.contacts = append(kb.contacts, contact)
		return true
	} else {
		evicted := kb.Evict()
		if evicted {
			kb.contacts = append(kb.contacts, contact)
			return true
		} else {
			fmt.Println("There is no space in this bucket")
			return false
		}
	}
}
func (kb *KBucket) Evict() bool {
	sort.Slice(kb.contacts, func(i, j int) bool {
		return kb.contacts[i].last_seen_at.Before(kb.contacts[j].last_seen_at)
	})
	contact := kb.contacts[0]
	_, err := net.Dial("tcp", contact.Address)
	if err != nil {
		kb.contacts = append(kb.contacts[:0], kb.contacts[0+1:]...)
		return true
	}

	return false

}

type RoutingTable struct {
	buckets [4]KBucket
	selfId  NodeID
}

// ~ let say this is for storing messages in memory
type Messages struct {
	SenderId       NodeID
	ReceiverId     NodeID
	MessageContent string
	MessageId      string
}
type MessagingPeer struct {
	ID           NodeID
	Messages     []Messages
	Port         int
	RoutingTable *RoutingTable
}

// ~ so now as the structure is build now my main focus is towards creating the node-id and functionalities around routing table
