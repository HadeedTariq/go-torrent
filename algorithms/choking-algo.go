package algorithms

import (
	"fmt"
	"math/rand"
	"sort"
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
}

type TorrentClient struct {
	Peers                        []*Peer
	IsSeeder                     bool
	UnchokeInterval              time.Duration
	OptimisticInterval           time.Duration
	DownloadRateCheckingInterval time.Duration
	SnubbedCheckingInterval      time.Duration
	lastOptimisticUnchoke        *Peer
}

func (tc *TorrentClient) SnubberChecker() {
	snubbedTimer := time.NewTicker(tc.SnubbedCheckingInterval)

	go func() {
		for range snubbedTimer.C {
			for _, peer := range tc.Peers {
				peer.mu.Lock()

				if peer.snubbed && time.Now().After(peer.snubbedUntil) {
					peer.snubbed = false
				}

				if !peer.snubbed && peer.DownloadRate < 10 && time.Since(peer.LastUnchokedAt) < 10*time.Second {
					peer.snubbed = true
					peer.snubbedUntil = time.Now().Add(20 * time.Minute)
				}

				peer.mu.Unlock()
			}
		}
	}()
}

func (tc *TorrentClient) UpdateDownloadRateOfPeers() {
	downloadRateTimer := time.NewTicker(tc.DownloadRateCheckingInterval)
	for range downloadRateTimer.C {
		select {
		case <-downloadRateTimer.C:
			for _, peer := range tc.Peers {
				duration := time.Now().Sub(peer.lastCheckedTime)
				rate := float64(peer.bytesDownloaded) / duration.Seconds()
				peer.DownloadRate = int(rate)
				peer.bytesDownloaded = 0
				peer.lastCheckedTime = time.Now()
			}
		}
	}
}

func (tc *TorrentClient) RunCheckLoop() {
	ticker := time.NewTicker(tc.UnchokeInterval)
	optimisticTicker := time.NewTicker(tc.OptimisticInterval)

	for {
		select {
		case <-ticker.C:
			if tc.IsSeeder {
				tc.runSeederChoke()
			} else {
				tc.runLeecherChoke()
			}
		case <-optimisticTicker.C:
			tc.runOptimisticLeecher()
		}
	}
}

func (tc *TorrentClient) runSeederChoke() {
	interestedPeers := []*Peer{}
	for _, peer := range tc.Peers {
		if peer.Interested && !peer.snubbed {
			interestedPeers = append(interestedPeers, peer)
		}
	}

	sort.Slice(interestedPeers, func(i, j int) bool {
		return interestedPeers[i].LastUnchokedAt.Before(interestedPeers[j].LastUnchokedAt)
	})

	for i, peer := range interestedPeers {
		if i < 3 {
			peer.Choked = false
			peer.LastUnchokedAt = time.Now()
		} else {
			peer.Choked = true
		}
	}

	if len(interestedPeers) > 3 {
		index := rand.Intn(len(interestedPeers)-3) + 3
		interestedPeers[index].Choked = false
		interestedPeers[index].LastUnchokedAt = time.Now()
	}
}

func (tc *TorrentClient) runLeecherChoke() {
	interestedPeers := []*Peer{}
	for _, peer := range tc.Peers {
		if peer.Interested && !peer.snubbed {
			interestedPeers = append(interestedPeers, peer)
		}
	}

	sort.Slice(interestedPeers, func(i, j int) bool {
		return interestedPeers[i].DownloadRate > interestedPeers[j].DownloadRate
	})

	for i, peer := range interestedPeers {
		if i < 3 {
			peer.Choked = false
			peer.LastUnchokedAt = time.Now()
		} else {
			peer.Choked = true
		}
	}
}

func (tc *TorrentClient) runOptimisticLeecher() {
	fmt.Println("🎲 [Leecher] Running optimistic unchoke...")

	chokedInterestedPeers := []*Peer{}
	for _, peer := range tc.Peers {
		if peer.Interested && peer.Choked && !peer.snubbed {
			chokedInterestedPeers = append(chokedInterestedPeers, peer)
		}
	}

	if len(chokedInterestedPeers) == 0 {
		fmt.Println("😢 No choked + interested peers found for optimistic unchoke.")
		return
	}

	randomIndex := rand.Intn(len(chokedInterestedPeers))
	selectedPeer := chokedInterestedPeers[randomIndex]

	selectedPeer.Choked = false
	selectedPeer.LastUnchokedAt = time.Now()
	fmt.Printf("🌟 Optimistically unchoked peer: %s\n", selectedPeer.ID)
}
