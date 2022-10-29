package game

const (
	BoardSize int = 19
)

type Piece int

const NilPiece Piece = -1

const (
	BlackPiece Piece = iota
	WhitePiece
)

type Player int

const NilPlayer Player = -1

const (
	BlackPlayer Piece = iota
	WhitePlayer
)

type State struct {
	Board [][]Piece
	Turn  Player
}

func NewState(boardSize int /*local*/, starter Player) *State {
	newState := &State{} //local
	newState.Board = make([][]Piece, boardSize)
	for i := 0; i < boardSize; i++ {
		newState.Board[i] = make([]Piece, boardSize)
	}
	newState.Turn = starter
	return newState
}
