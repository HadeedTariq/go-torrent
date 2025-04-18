package algorithms

import (
	"net"
	"sync"
	"time"
)

type Peer struct {
	IP              net.IP
	PORT            uint16
	Conn            net.Conn
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
	HandshakeDone   bool
}
