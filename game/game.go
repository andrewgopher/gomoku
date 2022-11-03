package game

type Piece int

const (
	NilPiece Piece = iota
	BlackPiece
	WhitePiece
)

type Player int

const (
	NilPlayer Player = iota
	BlackPlayer
	WhitePlayer
)

var PlayerToPiece map[Player]Piece = map[Player]Piece{
	BlackPlayer: BlackPiece,
	WhitePlayer: WhitePiece,
}

var PieceToPlayer map[Piece]Player = map[Piece]Player{
	BlackPiece: BlackPlayer,
	WhitePiece: WhitePlayer,
}

var PlayerToString map[Player]string = map[Player]string{
	BlackPlayer: "Black",
	WhitePlayer: "White",
}

var Directions []Point[int] = []Point[int]{{1, 0}, {-1, 0}, {0, 1}, {0, -1}, {1, 1}, {-1, -1}, {-1, 1}, {1, -1}}

type State struct {
	Board     [][]Piece
	Turn      Player
	Winner    Player
	NumConsec int
}

type Point[T int | float32 | int32 | int64 | float64] struct {
	X T
	Y T
}

func NewState(boardSize int /*local*/, numConsec int, starter Player) *State {
	newState := &State{} //local
	newState.Board = make([][]Piece, boardSize)
	for i := 0; i < boardSize; i++ {
		newState.Board[i] = make([]Piece, boardSize)
	}
	newState.Turn = starter
	newState.Winner = NilPlayer
	newState.NumConsec = numConsec
	return newState
}

func (state *State) OnBoard(pos *Point[int]) bool {
	return pos.X >= 0 && pos.Y >= 0 && pos.X < len(state.Board) && pos.Y < len(state.Board)
}

func (state *State) CheckWin(pos *Point[int]) Player { //check wins thru pos; assumes that Directions pairs the opposites
	target := state.Board[pos.X][pos.Y]
	consecDirs := make([]int, 8)
	for i, dir := range Directions {
		currPos := &Point[int]{pos.X + dir.X, pos.Y + dir.Y} //TODO: add method
		for state.OnBoard(currPos) && state.Board[currPos.X][currPos.Y] == target {
			currPos.X += dir.X
			currPos.Y += dir.Y
			consecDirs[i]++
		}
	}
	for i := 0; i < 4; i++ {
		if consecDirs[i*2]+consecDirs[i*2+1] >= state.NumConsec-1 {
			return PieceToPlayer[target]
		}
	}
	return NilPlayer
}

func OtherPlayer(p Player) Player {
	return p%2 + 1
}

func (state *State) MakeMove(pos *Point[int]) {
	state.Board[pos.X][pos.Y] = PlayerToPiece[state.Turn]
	state.Turn = OtherPlayer(state.Turn)
	state.Winner = state.CheckWin(pos)
}
