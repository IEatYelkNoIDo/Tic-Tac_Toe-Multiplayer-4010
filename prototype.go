package main

import (
    "fmt"

    "image/color"

    "math"

    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Game struct {
    board    [3][3]int // 0=empty, 1=X, 2=O
    playing  bool
    player   int       // 1=X, 2=O
    turn     int
    cellSize int
    offset   int
    wins     int
    team     string
    winner   string
}

// Constructor
func NewGame() *Game {
    return &Game{
        playing:  true,
        player:   1,
        cellSize: 150,
        offset:   50,
        turn:     1,
        team:     "X",
    }
}

func (g *Game) Update() error {
    if !g.playing {
        return nil
    }

    if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
        x, y := ebiten.CursorPosition()
        col := (x - g.offset) / g.cellSize
        row := (y - g.offset) / g.cellSize

        // Click inside board
        if col >= 0 && col < 3 && row >= 0 && row < 3 {
            if g.board[row][col] == 0 {
                g.board[row][col] = g.player
                g.turn++

                if g.turn >= 5 {
                    g.checkWin()
                }
                // Alternate turns
                if g.player == 1 {
                    g.player = 2
                    g.team = "O"
                } else {
                    g.player = 1
                    g.team = "X"
                }
            }
        }
    }

    return nil
}

func (g *Game) checkWin() {
    // Check rows and columns
    for i := 0; i < 3; i++ {
        // Row check
        if g.board[i][0] == g.board[i][1] && g.board[i][1] == g.board[i][2] && g.board[i][0] != 0 {
            g.winner = g.team
            return
        }
        // Column check
        if g.board[0][i] == g.board[1][i] && g.board[1][i] == g.board[2][i] && g.board[0][i] != 0 {
            g.winner = g.team
            return
        }
    }

    // Diagonal checks
    if g.board[0][0] == g.board[1][1] && g.board[1][1] == g.board[2][2] && g.board[0][0] != 0 {
        g.winner = g.team
        return
    }
    if g.board[0][2] == g.board[1][1] && g.board[1][1] == g.board[2][0] && g.board[0][2] != 0 {
        g.winner = g.team
        return
    }

    // Cat (Tie)
    if g.turn >= 9 {
        g.winner = "CAT"
    }
}

func drawCircle(screen *ebiten.Image, cx, cy, r float64, clr color.Color) {
    segments := 50
    for i := 0; i < segments; i++ {
        theta1 := 2 * math.Pi * float64(i) / float64(segments)
        theta2 := 2 * math.Pi * float64(i+1) / float64(segments)
        x1 := cx + r*math.Cos(theta1)
        y1 := cy + r*math.Sin(theta1)
        x2 := cx + r*math.Cos(theta2)
        y2 := cy + r*math.Sin(theta2)
        ebitenutil.DrawLine(screen, x1, y1, x2, y2, clr)
    }
}

func (g *Game) Draw(screen *ebiten.Image) {
    // Draw grid lines
    ebitenutil.DrawLine(screen, float64(g.offset+g.cellSize), float64(g.offset), float64(g.offset+g.cellSize), float64(g.offset+3*g.cellSize), color.White)
    ebitenutil.DrawLine(screen, float64(g.offset+2*g.cellSize), float64(g.offset), float64(g.offset+2*g.cellSize), float64(g.offset+3*g.cellSize), color.White)
    ebitenutil.DrawLine(screen, float64(g.offset), float64(g.offset+g.cellSize), float64(g.offset+3*g.cellSize), float64(g.offset+g.cellSize), color.White)
    ebitenutil.DrawLine(screen, float64(g.offset), float64(g.offset+2*g.cellSize), float64(g.offset+3*g.cellSize), float64(g.offset+2*g.cellSize), color.White)

    // Draw X/O
    for row := 0; row < 3; row++ {
        for col := 0; col < 3; col++ {
            x := float64(g.offset + col*g.cellSize)
            y := float64(g.offset + row*g.cellSize)
            if g.board[row][col] == 1 {
                // Draw X
                ebitenutil.DrawLine(screen, x, y, x+float64(g.cellSize), y+float64(g.cellSize), color.RGBA{255, 0, 0, 255})
                ebitenutil.DrawLine(screen, x+float64(g.cellSize), y, x, y+float64(g.cellSize), color.RGBA{255, 0, 0, 255})
            } else if g.board[row][col] == 2 {
                // Draw O
                drawCircle(screen, x + float64(g.cellSize)/2, y + float64(g.cellSize)/2, float64(g.cellSize)/2, color.RGBA{0, 255, 0, 255})            
            }
        }
    }

    // Writes winner
    if g.winner != "" {
        ebitenutil.DebugPrint(screen, g.winner + " Wins!")
        return
    }

    // Writes out turns
    ebitenutil.DebugPrint(screen, fmt.Sprintf("Player %d's turn", g.player))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
    return 600, 600
}

func main() {
    ebiten.SetWindowSize(600, 600)
    ebiten.SetWindowTitle("Tic Tac Toe - Go")
    if err := ebiten.RunGame(NewGame()); err != nil {
        panic(err)
    }
}