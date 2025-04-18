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
	buf[0] = 19 // protocol name length
	copy(buf[1:], "BitTorrent protocol")
	// 8 reserved bytes (zero)
	copy(buf[28:], infoHash[:])    // info_hash (20 bytes)
	copy(buf[48:], []byte(peerID)) // peer_id (20 bytes)

	// Send handshake
	_, err := conn.Write(buf)
	if err != nil {
		return fmt.Errorf("failed to send handshake: %w", err)
	}

	// Read handshake response (68 bytes)
	resp := make([]byte, 68)
	_, err = io.ReadFull(conn, resp)
	if err != nil {
		return fmt.Errorf("failed to read handshake: %w", err)
	}

	// Validate info hash
	if !bytes.Equal(resp[28:48], infoHash[:]) {
		return fmt.Errorf("info hash mismatch")
	}

	log.Println("🤝 Handshake successful")
	return nil
}

func (tc *TorrentClient) PeerLoop(conn net.Conn, peer *Peer, wg *sync.WaitGroup) {
	wg.Done()
	// read the bit field like which pieces a peer have
	for _, val := range peer.Bitfield {
		if val {
			_, err := conn.Write(InterestedMessage())
			if err != nil {
				log.Printf("❌ Failed to send interested: %v", err)
				return
			}
			log.Println("📨 Sent INTERESTED to peer")

			msg, err := ReadMessage(conn)

			if err != nil {
				log.Printf("Peer disconnected: %v", err)
				return
			}

			switch msg.ID {
			case 1:
				log.Println("✅ Peer unchoked us")
				peer.Choked = false

			case 0: // CHOKE
				log.Println("🚫 Peer choked us")
				peer.Choked = true
			case 4: // HAVE
				// update peer bitfield
			case 5: // BITFIELD
				// set peer.Bitfield
			}

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
