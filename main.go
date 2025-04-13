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

	// This part is reading  the torrent file and after reading the file it extract the piece hashes and return the piece hashes and after that I instantiate the client and pass pieces in to the client
	filePath := "debian-12.10.0-amd64-netinst.iso.torrent"

	meta, err := readTorrentFile(filePath)
	if err != nil {
		fmt.Println("Error reading torrent file:", err)
		return
	}

	hashes := extractHashes(meta.Info.Pieces)

	client := algorithms.TorrentClient{}

	client.InitPieces(hashes)
	fmt.Printf("✅ Torrent '%s' loaded with %d pieces.\n", meta.Info.Name, client.TotalPieces)

	// so now the pieces initiation task is done now I have to connect with peers list

}
