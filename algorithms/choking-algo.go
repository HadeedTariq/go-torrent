package algorithms

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

func (tc *TorrentClient) SnubberChecker() {
	snubbedTimer := time.NewTicker(tc.SnubbedCheckingInterval)

	go func() {
		for range snubbedTimer.C {
			for _, peer := range tc.Peers {
				peer.mu.Lock()

				if peer.Snubbed && time.Now().After(peer.SnubbedUntil) {
					peer.Snubbed = false
				}

				if !peer.Snubbed && peer.DownloadRate < 10 && time.Since(peer.LastUnchokedAt) < 10*time.Second {
					peer.Snubbed = true
					peer.SnubbedUntil = time.Now().Add(20 * time.Minute)
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
				duration := time.Now().Sub(peer.LastCheckedTime)
				rate := float64(peer.BytesDownloaded) / duration.Seconds()
				peer.DownloadRate = int(rate)
				peer.BytesDownloaded = 0
				peer.LastCheckedTime = time.Now()
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
		if peer.Interested && !peer.Snubbed {
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
		if peer.Interested && !peer.Snubbed {
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
		if peer.Interested && peer.Choked && !peer.Snubbed {
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
	fmt.Printf("🌟 Optimistically unchoked peer: %s\n", selectedPeer.IP)
}
