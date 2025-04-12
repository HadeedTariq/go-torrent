package algorithms

import (
	"sync"
	"time"
)

type Peer struct {
	ID              string
	Interested      bool
	Choked          bool
	DownloadRate    int
	LastUnchokedAt  time.Time
	bytesDownloaded int
	lastCheckedTime time.Time
	snubbed         bool
	snubbedUntil    time.Time
	mu              sync.Mutex
	Bitfield        []bool
}
