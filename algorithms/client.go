package algorithms

import (
	"time"
)

type TorrentClient struct {
	Peers                        map[string]*Peer
	IsSeeder                     bool
	UnchokeInterval              time.Duration
	OptimisticInterval           time.Duration
	DownloadRateCheckingInterval time.Duration
	SnubbedCheckingInterval      time.Duration
	lastOptimisticUnchoke        *Peer
	// for piece selection
	TotalPieces  int
	OwnBitfield  []bool
	Pieces       []*Piece
	Downloading  map[int]bool
	PieceHashMap map[int][]byte
	Strategy     string // "rarest", "random", "strict", "endgame"
}
