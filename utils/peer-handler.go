package utils

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
	"torrent-client/algorithms"
)

// so we are writting two parser to store the peers because some tracker support binary and some doesn't so we should have a fallback for that

// so we have to make sure that the unique peers are stored

func ParsePeers(peersBin []byte) (map[string]*algorithms.Peer, error) {
	const peerSize = 6 // 4 bytes IP + 2 bytes port
	if len(peersBin)%peerSize != 0 {
		return nil, fmt.Errorf("malformed peers binary data")
	}

	numPeers := len(peersBin) / peerSize
	peers := make(map[string]*algorithms.Peer, numPeers)
	for i := 0; i < len(peersBin); i += peerSize {
		ip := net.IP(peersBin[i : i+4])
		port := binary.BigEndian.Uint16(peersBin[i+4 : i+6])
		key := fmt.Sprintf("%s:%d", ip, port)
		peers[key] = &algorithms.Peer{
			IP:              ip,
			PORT:            port,
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
	}

	return peers, nil
}

func ParsePeersFromDict(peersDict []interface{}) (map[string]*algorithms.Peer, error) {
	peers := make(map[string]*algorithms.Peer, len(peersDict))

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
		key := fmt.Sprintf("%s:%d", ipStr, port)

		peers[key] = &algorithms.Peer{
			IP:              net.ParseIP(ipStr),
			PORT:            uint16(port),
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
	}
	return peers, nil
}
