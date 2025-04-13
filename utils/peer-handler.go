package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"sync"
	"time"
	"torrent-client/algorithms"
)

// so we are writting two parser to store the peers because some tracker support binary and some doesn't so we should have a fallback for that

func ParsePeers(peersBin []byte) ([]algorithms.Peer, error) {
	const peerSize = 6 // 4 bytes IP + 2 bytes port
	if len(peersBin)%peerSize != 0 {
		return nil, fmt.Errorf("malformed peers binary data")
	}

	numPeers := len(peersBin) / peerSize
	peers := make([]algorithms.Peer, 0, numPeers)

	for i := 0; i < len(peersBin); i += peerSize {
		ip := net.IP(peersBin[i : i+4])
		port := binary.BigEndian.Uint16(peersBin[i+4 : i+6])

		peer := algorithms.Peer{
			IP:              ip,
			PORT:            int(port),
			Interested:      false,
			Choked:          true,
			DownloadRate:    0,
			LastUnchokedAt:  time.Time{},
			BytesDownloaded: 0,
			LastCheckedTime: time.Now(),
			Snubbed:         false,
			SnubbedUntil:    time.Time{},
			Bitfield:        nil,
		}
		peers = append(peers, peer)
	}

	return peers, nil
}

func ParsePeersFromDict(peersDict []interface{}) ([]Peer, error) {
	peers := make([]Peer, 0, len(peersDict))
	for _, p := range peersDict {
		peerMap, ok := p.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid peer entry")
		}
		ipStr, ok := peerMap["ip"].(string)
		if !ok {
			return nil, fmt.Errorf("invalid peer IP")
		}
		port, ok := peerMap["port"].(int64)
		if !ok {
			return nil, fmt.Errorf("invalid peer port")
		}
		peers = append(peers, Peer{
			IP:   net.ParseIP(ipStr),
			Port: uint16(port),
		})
	}
	return peers, nil
}

func (p *Peer) DoHandshake(infoHash, peerID [20]byte) error {
	if p.Conn == nil {
		return fmt.Errorf("no connection established")
	}

	// Set deadline to avoid hanging
	p.Conn.SetDeadline(time.Now().Add(3 * time.Second))
	defer p.Conn.SetDeadline(time.Time{}) // Disable deadline

	// Prepare handshake buffer
	buf := make([]byte, 68) // 1 + 19 + 8 + 20 + 20
	buf[0] = 19             // pstrlen
	copy(buf[1:20], []byte("BitTorrent protocol"))
	copy(buf[28:48], infoHash[:])
	copy(buf[48:68], peerID[:])

	// Send handshake
	if _, err := p.Conn.Write(buf); err != nil {
		return fmt.Errorf("failed to write handshake: %v", err)
	}

	// Read response (should mirror the same structure)
	resp := make([]byte, 68)
	if _, err := io.ReadFull(p.Conn, resp); err != nil {
		return fmt.Errorf("failed to read handshake: %v", err)
	}

	// Validate response
	if resp[0] != 19 || string(resp[1:20]) != "BitTorrent protocol" {
		return fmt.Errorf("invalid handshake response")
	}
	if !bytes.Equal(resp[28:48], infoHash[:]) {
		return fmt.Errorf("mismatched info hash")
	}

	return nil
}

func ConnectToPeers(peers []Peer, infoHash, peerID [20]byte, maxConns int) ([]Peer, error) {
	var wg sync.WaitGroup
	connChan := make(chan Peer, len(peers))
	semaphore := make(chan struct{}, maxConns) // Limit concurrent connections

	for _, peer := range peers {
		wg.Add(1)
		go func(p Peer) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", p.IP, p.Port), 3*time.Second)
			if err != nil {
				return
			}
			p.Conn = conn

			if err := p.DoHandshake(infoHash, peerID); err != nil {
				conn.Close()
				return
			}

			connChan <- p
		}(peer)
	}

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(connChan)
	}()

	// Collect successful peers
	connectedPeers := make([]Peer, 0)
	for p := range connChan {
		connectedPeers = append(connectedPeers, p)
	}

	return connectedPeers, nil
}
