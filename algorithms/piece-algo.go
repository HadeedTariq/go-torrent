package algorithms

// ! Rarest piece
// we have to ensure that the rarest piece are distributed b/w the peer so that each peer have the rarest piece and each can get the good dowload speed
// ! Random Policy
// But we also have to make sure that the newly joined peers get the whatever piece available so that they can even get started for that we have to make sure that for first 4 pieces the random ploicy works out
// ! Strict Policy
// So it is something like we also have to ensure that until unless a user doesn't have a complete piece who just starting out we will only give him the block of a first piece he is downloading to ensure asap he start contributing
// ! Endgame mode
// We also have to implement that out like when we are at the end of the downloading the file instead of waiting that we only request to the one peer to download the block of the piece and whenever we recieved that block we abort the request to the other peers to make sure fast download
// & Others
// So the piece duplication is not possible when we are checking the hash carefully and maintaining the set of pieces in an ds which ensure that there is no redundant piece in the final file

type PieceState int

const (
	NotRequested PieceState = iota
	Requested
	Downloaded
	Verified
)

type Piece struct {
	Index      int
	State      PieceState
	Rarity     int
	Hash       []byte
	IsVerified bool
}

// Initialize pieces and structures
// I think so the torrent file is uploaded and this func runs
// Now what have to make sure is to integrate diff learning methods in to that
func (tc *TorrentClient) InitPieces(pieceHashes [][]byte) {

	tc.TotalPieces = len(pieceHashes)
	tc.OwnBitfield = make([]bool, tc.TotalPieces)
	tc.Downloading = make(map[int]bool)
	tc.PieceHashMap = make(map[int][]byte)
	tc.Pieces = make([]*Piece, tc.TotalPieces)

	for i, hash := range pieceHashes {
		tc.PieceHashMap[i] = hash
		tc.Pieces[i] = &Piece{
			Index:  i,
			State:  NotRequested,
			Rarity: 0,
			Hash:   hash,
		}
	}
}

func (tc *TorrentClient) RarestPiece() {

}
