package algorithms

import (
	"net"
	"sync"
	"time"
)

type Peer struct {
	IP              net.IP
	PORT            int
	Interested      bool
	Choked          bool
	DownloadRate    int
	LastUnchokedAt  time.Time
	BytesDownloaded int
	LastCheckedTime time.Time
	Snubbed         bool
	SnubbedUntil    time.Time
	mu              sync.Mutex
	Bitfield        []bool
}
