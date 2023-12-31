package base

import "github.com/ComputerSaysYeah/RookMills/speed"
import . "github.com/ComputerSaysYeah/RookMills/api"

type gameSt struct {
	speed.Recyclable

	board          Board
	moveNo         int
	halfMoveNo     int
	nextPlayer     Piece
	enPassant      Square
	WK, WQ, bk, bq bool // castling

	boardPool       speed.Pool[Board]
	movesIterPool   speed.Pool[MovesIterator]
	squaresIterPool speed.Pool[SquaresIterator]

	returner func(any)
}

func NewGame(boardPool speed.Pool[Board],
	movesIterPool speed.Pool[MovesIterator],
	squaresIterPool speed.Pool[SquaresIterator]) Game {
	board := boardPool.Lease()
	board.SetStartingPieces()
	return &gameSt{
		board:           board,
		moveNo:          1,
		halfMoveNo:      0,
		nextPlayer:      White,
		enPassant:       None,
		WK:              true,
		WQ:              true,
		bk:              true,
		bq:              true,
		boardPool:       boardPool,
		movesIterPool:   movesIterPool,
		squaresIterPool: squaresIterPool,
		returner:        nil}
}

func (g *gameSt) Reset() {
	g.board.SetStartingPieces() //XXX see if we can avoid this on each node exploration
	g.moveNo = 1
	g.halfMoveNo = 0
	g.WK, g.WQ, g.bk, g.bq = true, true, true, true
	g.nextPlayer = White
}

func (g *gameSt) MoveNo() int {
	return g.moveNo
}

func (g *gameSt) HalfMoveNo() int {
	return g.halfMoveNo
}

func (g *gameSt) MoveNext() Piece {
	return g.nextPlayer
}

func (g *gameSt) EnPassant() Square {
	return g.enPassant
}

// Move applies the move, it does not verify the Move is valid, it applies it, for a Game it needs verifying the move is within the ValidMoves
func (g *gameSt) Move(move Move) {

	if move.Promote().IsEmpty() {
		piece := g.Board().Get(move.From())
		g.Board().Set(move.To(), piece)
		if piece.IsPawn() {
			if g.enPassant == move.To() {
				if piece.Colour() == Black {
					g.Board().Set(move.To().N(), Empty)
				} else {
					g.Board().Set(move.To().S(), Empty)
				}
			} else if move.Manhattan() == 2 && move.From().Col() == move.To().Col() {
				if piece.Colour() == Black && move.From().Row() == Row7 {
					g.enPassant = move.To().N()
				}
				if piece.Colour() == White && move.From().Row() == Row2 {
					g.enPassant = move.To().S()
				}
			} else {
				g.enPassant = None
			}
		}

	} else {
		g.Board().Set(move.To(), move.Promote())
	}
	g.Board().Set(move.From(), Empty)

	if g.nextPlayer.Colour() == Black {
		g.moveNo++
	}
	g.halfMoveNo++
	g.nextPlayer = g.nextPlayer.Opponent()
}

func (g *gameSt) Castling() (WK, WQ, bk, bq bool) {
	return g.WK, g.WQ, g.bk, g.bq
}

func (g *gameSt) SetReturnerFn(returner func(any)) {
	g.returner = returner
}

func (g *gameSt) Return() {
	g.returner(g)
}

func (g *gameSt) Board() Board {
	return g.board
}

func (g *gameSt) SetMoveNo(moveNo int) {
	g.moveNo = moveNo
}

func (g *gameSt) SetHalfMoveNo(halfMoveNo int) {
	g.halfMoveNo = halfMoveNo
}

func (g *gameSt) SetMoveNext(piece Piece) {
	g.nextPlayer = piece
}

func (g *gameSt) SetEnPassant(square Square) {
	g.enPassant = square
}

func (g *gameSt) SetCastling(WK, WQ, bk, bq bool) {
	g.WK = WK
	g.WQ = WQ
	g.bk = bk
	g.bq = bq
}

func (g *gameSt) CopyFrom(o Game) {
	g.board.CopyFrom(o.Board())
	g.moveNo = o.MoveNo()
	g.halfMoveNo = o.HalfMoveNo()
	g.nextPlayer = o.MoveNext()
	g.enPassant = o.EnPassant()
	g.SetCastling(o.Castling())
}
