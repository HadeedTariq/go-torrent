package main

import (
	"fmt"
	"os"
	"torrent-client/algorithms"

	"github.com/jackpal/bencode-go"
)

// why info dict is necessary because it have to get the exact info of the torrent file and parse that out

type InfoDict struct {
	PieceLength int    `bencode:"piece length"`
	Pieces      string `bencode:"pieces"`
	Name        string `bencode:"name"`
	Length      int    `bencode:"length"`
}

type TorrentMeta struct {
	Announce string   `bencode:"announce"`
	Info     InfoDict `bencode:"info"`
}

func readTorrentFile(filePath string) (*TorrentMeta, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var meta TorrentMeta
	err = bencode.Unmarshal(file, &meta)
	if err != nil {
		return nil, err
	}

	return &meta, nil
}

func extractHashes(piecseString string) [][]byte {
	hashLen := 20
	data := []byte(piecseString)

	numHashes := len(piecseString) / hashLen

	hashes := make([][]byte, numHashes)

	// so what is the procedure of extracting the pieces
	for i := 0; i < len(data); i += hashLen {
		hashes = append(hashes, data[i:i+hashLen])
	}

	return hashes
}

func main() {
	// port := 6881
	// address := fmt.Sprintf("0.0.0.0:%d", port)

	// listener, err := net.Listen("tcp", address)
	// if err != nil {
	// 	panic(fmt.Sprintf("Failed to listen on port %d: %v", port, err))
	// }
	// defer listener.Close()

	// fmt.Printf("🌐 Torrent client listening on %s\n", address)

	// // This part is reading  the torrent file and after reading the file it extract the piece hashes and return the piece hashes and after that I instantiate the client and pass pieces in to the client
	// filePath := "debian-12.10.0-amd64-netinst.iso.torrent"

	// meta, err := readTorrentFile(filePath)
	// if err != nil {
	// 	fmt.Println("Error reading torrent file:", err)
	// 	return
	// }

	// hashes := extractHashes(meta.Info.Pieces)

	// client := algorithms.TorrentClient{}

	// client.InitPieces(hashes)
	// fmt.Printf("✅ Torrent '%s' loaded with %d pieces.\n", meta.Info.Name, client.TotalPieces)

	// // so now the pieces initiation task is done now I have to connect with peers list
	// // so the peers parsers are written now my main concern is to connect to the tracker and store the peer info in an application memory
	// baseURL := meta.Announce

	// var buf bytes.Buffer
	// err = bencode.Marshal(&buf, meta.Info)
	// if err != nil {
	// 	log.Fatal("failed to bencode info dict")
	// }
	// infoHash := sha1.Sum(buf.Bytes())
	// encodedInfoHash := utils.EncodeInfoHash(infoHash)

	// // so the error basically exist due to the double encoding
	// peer_id := utils.GeneratePeerID()
	// params := url.Values{}
	// params.Add("peer_id", peer_id) // 20-byte peer ID
	// params.Add("port", strconv.Itoa(port))
	// params.Add("uploaded", "0")
	// params.Add("downloaded", "0")
	// params.Add("left", strconv.Itoa(meta.Info.Length))
	// params.Add("compact", "1")
	// params.Add("event", "started")

	// query := fmt.Sprintf("info_hash=%s&%s", encodedInfoHash, params.Encode())
	// fullURL := fmt.Sprintf("%s?%s", baseURL, query)
	// fmt.Println(fullURL)

	// // Send the request
	// resp, err := http.Get(fullURL)
	// if err != nil {
	// 	log.Fatalf("❌ Tracker request failed: %v", err) // stop the program
	// }
	// defer resp.Body.Close()

	// type TrackResponseStruct struct {
	// 	Interval int         `bencode:"interval"`
	// 	Peers    interface{} `bencode:"peers"`
	// }

	// // Decode tracker response
	// var trackerResp TrackResponseStruct
	// err = bencode.Unmarshal(resp.Body, &trackerResp)
	// if err != nil {
	// 	log.Fatalf("failed to parse tracker response: %v", err)
	// }

	// // Handle compact peer list
	// peersRaw := trackerResp.Peers

	// fmt.Println("passes")
	// var parsedPeers map[string]*algorithms.Peer
	// switch peers := peersRaw.(type) {
	// case string: // compact format
	// 	parsedPeers, err = utils.ParsePeers([]byte(peers))
	// case []interface{}: // dictionary format
	// 	parsedPeers, err = utils.ParsePeersFromDict(peers)
	// default:
	// 	log.Fatalf("unsupported peers format: %T", peersRaw)
	// }
	// if err != nil {
	// 	log.Fatalf("failed to parse peers: %v", err)
	// }
	// // ~ I think so parsing peers part is done

	// client.Peers = parsedPeers
	// var wg sync.WaitGroup

	// for key, peer := range client.Peers {
	// 	wg.Add(1)
	// 	go client.ConnectToPeer(peer, key, infoHash, peer_id, &wg)
	// }

	// wg.Wait()
	// fmt.Println("Finished connecting to peers")
	nodeA := algorithms.NewNode("NodeA")
	nodeB := algorithms.NewNode("NodeB")
	nodeC := algorithms.NewNode("NodeC")

	// Connect them in a ring
	nodeA.SetSuccessor(nodeB)
	nodeB.SetSuccessor(nodeC)
	nodeC.SetSuccessor(nodeA)

	// Store some keys
	nodeA.Store("apple", "a")
	nodeB.Store("banana", "🍌")
	nodeC.Store("cherry", "c")

	// Lookup keys
	fmt.Println("Find 'banana':", nodeA.Find("banana"))
	fmt.Println("Find 'cherry':", nodeB.Find("cherry"))
	fmt.Println("Find 'apple':", nodeC.Find("apple"))
}
