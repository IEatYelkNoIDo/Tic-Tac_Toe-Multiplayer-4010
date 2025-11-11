/*
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"image/color"
	"log"
	"math"
	"net"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Creates a new data type named GameState and assigns it as a baseline int
type GameState int

const ( // Creates constant values for GameState
        StateMenu    GameState = iota //Automatically assigns numbers starting from 0 under this constant
        StatePlaying                  // GameState = 0, State Playing = 1
)

// Defines types that will be shared accross multiple funcitions by using a pointer
type Game struct {
        board    [3][3]int // 0=empty, 1=X, 2=O
        playing  bool
        h_play   bool // hover for playing
        h_quit   bool
        player   int // 1=X, 2=O
        turn     int
        cellSize int
        offset   int
        wins     int
        mX       int // Max X border
        mY       int // Max Y border
        team     string
        winner   string
        state    GameState //Defines state as a GameState data type
        conn     net.Conn
        restart bool
}

type Input struct {
        Player int
        Row    int
        Col    int
}

type Message struct {
        Text string 
}

//var player Input

type Update struct {
        Player int
        Board  [3][3]int
        Turn   int
        Winner string
        Restart bool
}

// Constructor
func NewGame() *Game {
        return &Game{
                playing:  true,
                h_play:   false,
                h_quit:   false,
                player:   1,
                cellSize: 150,
                offset:   50,
                turn:     1,
                mX:       600,
                mY:       600,
                team:     "X",
                state:    StateMenu,
        }
}

//var restart = false

func (g *Game) Update() error {
        if !g.playing {
                return nil
        }

        x, y := ebiten.CursorPosition()

        switch g.state { //Switch is basically like a giant easier to use if/else statement

        case StateMenu: //equivalent to if g.state == "StateMenu"
                x, y := ebiten.CursorPosition()

                btnWidth := 160
                btnHeight := 40
                btnX := g.mX/2 - 80 // Centered button X
                btnY := g.mY/2 - 20 // Centered button Y
                btnY2 := g.mY/2 + 40

                // Hover Play Check
                if x >= btnX && x <= btnX+btnWidth {
                        if y >= btnY && y <= btnY+btnHeight {
                                g.h_play = true
                        } else if y >= btnY2 && y <= btnY2+btnHeight {
                                g.h_quit = true
                        } else {
                                g.h_play = false
                                g.h_quit = false
                        }
                }

                // Click Check
                if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
                        if x >= btnX && x <= btnX+btnWidth {
                                if y >= btnY && y <= btnY+btnHeight {
                                        g.state = StatePlaying
                                } else if y >= btnY2 && y <= btnY2+btnHeight {
                                        os.Exit(0)
                                }
                        }
                }

        case StatePlaying: //else if g.state == "StatePlaying"
                if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && g.winner == "" {

                        col := (x - g.offset) / g.cellSize
                        row := (y - g.offset) / g.cellSize

                        // Click inside board
                        if col >= 0 && col < 3 && row >= 0 && row < 3 {
                                if g.board[row][col] == 0 {
                                g.board[row][col] = g.player

                                        input := Input{Player: g.player, Row: row, Col: col}
                                        json.NewEncoder(g.conn).Encode(input)

                                        // Alternate turns
                                        if g.player == 1 {
                                                //g.player = 2
                                                g.team = "X"
                                        } else {
                                                //g.player = 1
                                                g.team = "O"
                                        }

                                         if g.turn >= 5 {
                                                g.checkWin()
                                        }
                                        

                                }
                        }
                }
        

                if g.winner != "" { // Resets game if you press R
                        if inpututil.IsKeyJustPressed(ebiten.KeyR) {
                                msg := struct{ Restart bool } {Restart: true}
                                encoder := json.NewEncoder(g.conn)

                                if err := encoder.Encode(msg); err != nil {
                                        fmt.Println("restart error ", err)
                                }
                                *g = *NewGame()
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
        if g.turn > 9 && g.winner == "" {
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

        var pColor, qColor color.Color

        if g.h_play == false {
                pColor = color.RGBA{10, 10, 255, 255}
        } else {
                pColor = color.RGBA{100, 100, 200, 255}
        }

        if g.h_quit == false {
                qColor = color.RGBA{10, 10, 255, 255}
        } else {
                qColor = color.RGBA{100, 100, 200, 255}
        }

        switch g.state {

        case StateMenu:
                // Draw background
                screen.Fill(color.RGBA{30, 30, 30, 255})

                // Draw title
                ebitenutil.DebugPrintAt(screen, "Tic Tac Toe", g.mX/2-60, g.mY/2-80)

                // Draw Play button rectangle
                ebitenutil.DrawRect(screen, float64(g.mX/2-80), float64(g.mY/2-20), 160, 40, pColor)
                ebitenutil.DebugPrintAt(screen, "Play", g.mX/2-20, g.mY/2-10)

                // Draw Quit button rectangle
                ebitenutil.DrawRect(screen, float64(g.mX/2-80), float64(g.mY/2+40), 160, 40, qColor)
                ebitenutil.DebugPrintAt(screen, "Quit", g.mX/2-20, g.mY/2+50)

        case StatePlaying:
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
                                        drawCircle(screen, x+float64(g.cellSize)/2, y+float64(g.cellSize)/2, float64(g.cellSize)/2, color.RGBA{0, 255, 0, 255})
                                }
                        }
                }

                // Writes winner
                if g.winner != "" {
                        ebitenutil.DebugPrint(screen, g.winner+" Wins!"+" Press 'R' to play again!")
                        return
                }

                // Writes out turns
                if g.turn % 2 != 0 {
                ebitenutil.DebugPrint(screen, "Player l's turn")
                } else {
                        ebitenutil.DebugPrint(screen, "Player 2's turn")
                }

        }
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
        return g.mX, g.mY
}

func main() {

        // establish connection to server
        var d net.Dialer
        ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
        defer cancel()

        var err error

        // chose whichever server you like. Every server is connected via a tailscale ip address
        conn, err := d.DialContext(ctx, "tcp", "100.67.88.56:8080")
        //conn, err := d.DialContext(ctx, "tcp", "100.118.145.55:8080")
        //conn, err := d.DialContext(ctx, "tcp", "100.108.153.55:8080")
        
        if err != nil {
                log.Println("Dial error:", err)
        }
        defer conn.Close()

        decoder := json.NewDecoder(conn)
        var init struct{ Player int }

        // player number is the same as saying player id
        if err := decoder.Decode(&init); err != nil {
                fmt.Println("Failed to read player number", err)
                return
        }

        fmt.Println("You are Player: ", init.Player)

        g := NewGame()
        g.player = init.Player
        g.conn = conn

        // go routine that constantly updates the client
        go func() {
                decoder := json.NewDecoder(conn)

                for {
                        var update Update
                        if err := decoder.Decode(&update); err != nil {
                                fmt.Println("Decode error:", err)
                                return
                        }
                        // update the game state
                        g.board = update.Board
                        g.turn = update.Turn
                        g.winner = update.Winner
                        g.player = update.Player
                        //g.restart = update.Restart
                }
        }()

        ebiten.SetWindowSize(g.mX, g.mY)
        ebiten.SetWindowTitle("Tic Tac Toe - Go")
        if err := ebiten.RunGame(g); err != nil {
                panic(err)
        }
}
*/