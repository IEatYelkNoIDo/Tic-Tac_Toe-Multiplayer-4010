package main

import (
	"context"
	"encoding/json"
	"fmt"
	"image/color"
	_ "image/png"
	"log"
	"net"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
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
		winStartX, winStartY float64
		winEndX, winEndY     float64
		imageX, imageO       *ebiten.Image
		titleFont, smallFont font.Face
        state    GameState //Defines state as a GameState data type
        conn     net.Conn
}

type Input struct {
        Player int
        Row    int
        Col    int
}

type Message struct {
        Text string 
}

type Update struct {
        Player int
        Board  [3][3]int
        Turn   int
        Winner string
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


func (g *Game) Update() error {

        if !g.playing {
                return nil
        }

        x, y := ebiten.CursorPosition()

        switch g.state { //Switch is basically like a giant easier to use if/else statement

        case StateMenu: //equivalent to if g.state == "StateMenu"
                x, y := ebiten.CursorPosition()

                btnWidth := 240
				btnHeight := 80
				btnX := g.mX/2 - 120 // Centered button X
				btnY := g.mY/2 - 60 // Centered button Y
				btnY2 := g.mY/2 + 80

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
					
				// check player number before accepting any input
				expectedPlayer := 1 
				if g.turn % 2 == 0 {
					expectedPlayer = 2
				}

				// only process the click if the expected player is the current player
				if expectedPlayer == g.player {
					// send player input to server
						input := Input{Player: g.player, Row: row, Col: col}
						json.NewEncoder(g.conn).Encode(input) // Sends the player that made the move and the row and col that they made the move on
															// sends it something like {Player: int value, Col: int value, Col: int value} 
			
					// Click inside board
					if col >= 0 && col < 3 && row >= 0 && row < 3 {
						if g.board[row][col] == 0 {
							g.board[row][col] = g.player

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
				} else {					
					fmt.Println("Sorry not your turn")
				}
        	}
        

                if g.winner != "" { // Resets game if you press R
					if inpututil.IsKeyJustPressed(ebiten.KeyR) {
						// Keep assets that don't need reloading
						titleFont := g.titleFont
						smallFont := g.smallFont
						imageX := g.imageX
						imageO := g.imageO

						encoder := json.NewEncoder(g.conn)
						if err := encoder.Encode(1); err != nil {
							fmt.Println(err)
						}
						
						// Reset game logic
						*g = *NewGame()

						// Restore assets
						g.titleFont = titleFont
						g.smallFont = smallFont
						g.imageX = imageX
						g.imageO = imageO
					}
                }

        }
        return nil
}

func (g *Game) checkWin() {
    cell := float64(g.cellSize)
	offset := float64(g.offset)
	// Rows
	for i := 0; i < 3; i++ {
		if g.board[i][0] == g.board[i][1] && g.board[i][1] == g.board[i][2] && g.board[i][0] != 0 {
			g.winner = g.team
			g.winStartX = offset
			g.winStartY = offset + cell*float64(i) + cell/2
			g.winEndX = offset + 3*cell
			g.winEndY = g.winStartY
			return
		}
	}

	// Columns
	for i := 0; i < 3; i++ {
		if g.board[0][i] == g.board[1][i] && g.board[1][i] == g.board[2][i] && g.board[0][i] != 0 {
			g.winner = g.team
			g.winStartX = offset + cell*float64(i) + cell/2
			g.winStartY = offset
			g.winEndX = g.winStartX
			g.winEndY = offset + 3*cell
			return
		}
	}

	// Diagonal (top-left → bottom-right)
	if g.board[0][0] == g.board[1][1] && g.board[1][1] == g.board[2][2] && g.board[0][0] != 0 {
		g.winner = g.team
		g.winStartX = offset
		g.winStartY = offset
		g.winEndX = offset + 3*cell
		g.winEndY = offset + 3*cell
		return
	}

	// Diagonal (top-right → bottom-left)
	if g.board[0][2] == g.board[1][1] && g.board[1][1] == g.board[2][0] && g.board[0][2] != 0 {
		g.winner = g.team
		g.winStartX = offset + 3*cell
		g.winStartY = offset
		g.winEndX = offset
		g.winEndY = offset + 3*cell
		return
	}

	// Cat (Tie)
	if g.turn > 9 && g.winner == "" {
		g.winner = "CAT"
	}
}

func loadFont(filePath string, size float64) (font.Face, error) {
    data, err := os.ReadFile(filePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read font file: %v", err)
    }

    ttf, err := opentype.Parse(data)
    if err != nil {
        return nil, fmt.Errorf("failed to parse font: %v", err)
    }

    face, err := opentype.NewFace(ttf, &opentype.FaceOptions{
        Size:    size,
        DPI:     72,
        Hinting: font.HintingFull,
    })
    if err != nil {
        return nil, fmt.Errorf("failed to create font face: %v", err)
    }

    return face, nil
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
		text.Draw(screen, "Tic Tac Toe", g.titleFont, g.mX/4, g.mY/4, color.White)

		// Draw Play button rectangle
		ebitenutil.DrawRect(screen, float64(g.mX/2-120), float64(g.mY/2-60), 240, 80, pColor)
		text.Draw(screen, "Play", g.titleFont, g.mX/2-65, g.mY/2-5, color.White)

		// Draw Quit button rectangle
		ebitenutil.DrawRect(screen, float64(g.mX/2-120), float64(g.mY/2+80), 240, 80, qColor)
		text.Draw(screen, "Quit", g.titleFont, g.mX/2-55, g.mY/2+135, color.White)

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
				op := &ebiten.DrawImageOptions{}
			
				switch g.board[row][col] {
				case 1:
					scaleX := float64(g.cellSize) / float64(g.imageX.Bounds().Dx())
					scaleY := float64(g.cellSize) / float64(g.imageX.Bounds().Dy())
					op.GeoM.Scale(scaleX, scaleY)
					op.GeoM.Translate(x, y)
					screen.DrawImage(g.imageX, op)
				case 2:
					scaleX := float64(g.cellSize) / float64(g.imageO.Bounds().Dx())
					scaleY := float64(g.cellSize) / float64(g.imageO.Bounds().Dy())
					op.GeoM.Scale(scaleX, scaleY)
					op.GeoM.Translate(x, y)
					screen.DrawImage(g.imageO, op)
				}
			}
		}
	
		// Draws line through winner
		if g.winner != "" && g.winner != "CAT" {
			for i := 0; i <= 5; i++ {
				ebitenutil.DrawLine(
					screen,
					g.winStartX+float64(i),
					g.winStartY,
					g.winEndX+float64(i),
					g.winEndY,
					color.RGBA{0, 255, 0, 255}, //green
				)
			}
		}

		// Writes winner
		if g.winner != "" {
			text.Draw(screen, g.winner+" Wins! Press 'R' to play again!", g.smallFont, g.mX/20, g.mY/20, color.White)
			return
		}

		// Writes out turns
		if g.turn % 2 != 0 {
			text.Draw(screen, "Player 1's turn", g.smallFont, g.mX/20, g.mY/20, color.White)
		} else {
			text.Draw(screen, "Player 2's turn", g.smallFont, g.mX/20, g.mY/20, color.White)
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

	g.imageX, _, _ = ebitenutil.NewImageFromFile("X.png")
	g.imageO, _, _ = ebitenutil.NewImageFromFile("O.png")

	g.titleFont, err = loadFont("RasterForgeRegular-JpBgm.ttf", 48)
	g.smallFont, err = loadFont("RasterForgeRegular-JpBgm.ttf", 24)

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
		}
	}()

	ebiten.SetWindowSize(g.mX, g.mY)
	ebiten.SetWindowTitle("Tic Tac Toe - Go")
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
