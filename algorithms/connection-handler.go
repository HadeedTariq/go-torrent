package algorithms

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

func (tc *TorrentClient) ConnectToPeer(peer *Peer, address string, infoHash [20]byte, peerID string, wg *sync.WaitGroup) {
	defer wg.Done()
	conn, err := net.DialTimeout("tcp", address, 5*time.Second)
	if err != nil {
		log.Printf("❌ Failed to connect to peer %s: %v", address, err)
		return
	}

	peer.Conn = conn
	log.Printf("🔗 Connected from %s to peer: %s", peerID, address)

	err = tc.PerformHandshake(conn, infoHash, peerID)
	if err != nil {
		log.Printf("❌ Handshake failed with %s: %v", address, err)
		conn.Close()
		return
	}

	peer.HandshakeDone = true

	// Start message loop
	wg.Add(1)
	go tc.PeerLoop(conn, peer, wg)
}

func (tc *TorrentClient) PerformHandshake(conn net.Conn, infoHash [20]byte, peerID string) error {
	// Construct the handshake

	buf := make([]byte, 68)

	// first  I have to intiate the memory for protocol
	buf[0] = 19
	copy(buf[1:], "BitTorrent protocol")

	copy(buf[28:], infoHash[:])
	copy(buf[48:], []byte(peerID))

	_, err := conn.Write(buf)

	if err != nil {
		return fmt.Errorf("failed to send handshake: %w", err)
	}

	resp := make([]byte, 68)

	// now I have to read the response and write in to memory and check the info hash
	_, err = io.ReadFull(conn, resp)

	if err != nil {
		return fmt.Errorf("failed to read handshake: %w", err)
	}

	if !bytes.Equal(resp[28:48], infoHash[:]) {
		return fmt.Errorf("info hash mismatch")
	}

	log.Println("🤝 Handshake successful")

	return nil
}
func (tc *TorrentClient) PeerLoop(conn net.Conn, peer *Peer, wg *sync.WaitGroup) {
	defer wg.Done()

	// Send INTERESTED once (you may enhance with bitfield logic later)
	_, err := conn.Write(InterestedMessage())
	if err != nil {
		log.Printf("❌ Failed to send interested: %v", err)
		return
	}
	log.Println("📨 Sent INTERESTED to peer")

	// Start loop to listen for messages
	for {
		msg, err := ReadMessage(conn)
		if err != nil {
			log.Printf("❌ Peer disconnected or error: %v", err)
			return
		}

		switch msg.ID {
		case 0: // CHOKE
			log.Println("🚫 Peer choked us")
			peer.Choked = true

		case 1: // UNCHOKE
			log.Println("✅ Peer unchoked us")
			peer.Choked = false

		case 4: // HAVE
			// update peer.Bitfield[msg.Payload[0]] = true
			log.Println("📦 Peer sent HAVE (implement logic here)")

		case 5: // BITFIELD
			peer.Bitfield = ParseBitfield(msg.Payload)
			log.Printf("📊 Received bitfield: %v", peer.Bitfield)

		case 7: // PIECE
			log.Println("📥 Received PIECE (implement logic here)")

		default:
			log.Printf("🔎 Unknown message ID: %d", msg.ID)
		}
	}
}

func InterestedMessage() []byte {
	buf := make([]byte, 5)

	binary.BigEndian.PutUint32(buf[0:4], 1)
	buf[4] = 2
	return buf
}

type Message struct {
	Length  int
	ID      byte
	Payload []byte
}

func ReadMessage(conn net.Conn) (Message, error) {
	var lengthBuf [4]byte
	_, err := io.ReadFull(conn, lengthBuf[:])

	if err != nil {
		return Message{}, err
	}

	length := binary.BigEndian.Uint32(lengthBuf[:])
	if length == 0 {
		// keep-alive
		return Message{Length: 0}, nil
	}
	msg := Message{
		Length: int(length),
	}

	payload := make([]byte, length)
	_, err = io.ReadFull(conn, payload)
	if err != nil {
		return Message{}, err
	}
	msg.ID = payload[0]
	msg.Payload = payload[1:]

	return msg, nil

}

func ParseBitfield(payload []byte) []bool {
	bits := []bool{}
	for _, b := range payload {
		for i := 7; i >= 0; i-- {
			bits = append(bits, (b>>i)&1 == 1)
		}
	}
	return bits
}
