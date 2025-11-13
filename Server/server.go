package main

import (
	"encoding/json"
	"fmt"
	"net"
)

type Input struct {
        Player int
        Row    int
        Col    int
}

type Update struct {
        Player int
        Board  [3][3]int
        Turn   int
        Winner string
}

// global variables
var (
        board [3][3]int
        turn  = 1
)

var conn1 net.Conn = nil
var conn2 net.Conn = nil

func main() {
        // listen
        // accept
        // handle connections
        playercount := 0

        fmt.Println("Server is up and running. Waiting for players to connect.")

        dstream, err := net.Listen("tcp", "100.67.88.56:8080")
        if err != nil {
                fmt.Println(err)
                return
        }
        defer dstream.Close()

        for {
                conn, err := dstream.Accept()
                if err != nil {
                        fmt.Println(err)
                        continue
                }
                playercount++
                playernum := playercount

                fmt.Println("Player connected:", playernum)

                // distinguish between the connections
                if playernum == 1 {
                        conn1 = conn
                } else if playernum == 2 {
                        conn2 = conn
                } else {
                        fmt.Println("Maximum players connected. Connection refused.")
                        conn.Close() // reject the connection
                }

                go func(conn net.Conn, playernum int) {

                        fmt.Println("Started goroutine for player ", playernum)

                        defer conn.Close() // close the connection after the go routine finishes
                        encoder := json.NewEncoder(conn)
                        decoder := json.NewDecoder(conn)

                        // assign a player id to the current connection.
                        if err := encoder.Encode(struct {Player int}{Player: playernum}); err != nil {
                                fmt.Println("error sending player number", err)
                                return
                        }

                        for {
                                var input Input

                                // decode the input from the player clicking on the board
                                if err := decoder.Decode(&input); err != nil {
                                        fmt.Println("Decode error:", err)
                                        return
                                }

                                fmt.Printf("Player %d move: row=%d, col=%d\n", input.Player, input.Row, input.Col)

                                // check if the player is allowed to make a move
                                expectedPlayer := 1
                                if turn % 2 == 0 {
                                        expectedPlayer = 2
                                }

                                if input.Player != expectedPlayer {
                                        fmt.Println("Player tried to move twice or more on one turn", input.Player, turn, expectedPlayer)
                                        continue
                                }


                                if board[input.Row][input.Col] == 0 {
                                board[input.Row][input.Col] = input.Player
                                turn++
                                }
                                winner := checkWin()

                                // currentlly connected player info
                                update := Update {
                                Player: input.Player,
                                Board: board,
                                Turn:  turn,
                                Winner: winner,
                                }


                                fmt.Println("currentUpdate values: ", update)

                                // update the current player
                                if err := encoder.Encode(update); err != nil {
                                        fmt.Println(err)
                                        return
                                }

                        if playernum == 1 {

                                otherPlayerUpdate := Update {
                                        Player: input.Player + 1,
                                        Board: board,
                                        Turn:  turn,
                                        Winner: winner,
                                }
                                // notify player 2
                                otherConn := conn2
                                otherEncoder := json.NewEncoder(otherConn)
                                if err := otherEncoder.Encode(otherPlayerUpdate); err != nil {
                                        fmt.Println(err)
                                        return
                                }
                        }
                        if playernum == 2 {
                                otherPlayerUpdate := Update {
                                        Player: input.Player - 1,
                                        Board: board,
                                        Turn:  turn,
                                        Winner: winner,
                                }
                                //notify player 1
                                otherConn := conn1
                                otherEncoder := json.NewEncoder(otherConn)
                                if err := otherEncoder.Encode(otherPlayerUpdate); err != nil {
                                        fmt.Println(err)
                                        return
                                }
                        }

                        }
               }(conn, playernum)
        }
}

func checkWin() string {

        for i := 0; i < 3; i++ {
                 if board[i][0] == board[i][1] && board[i][1] == board[i][2] && board[i][0] != 0 {
                        return fmt.Sprintf("Player %d", board[i][0])
                 }
                 if board[0][i] == board[1][i] && board[1][i] == board[2][i] && board[0][i] != 0 {
                        return fmt.Sprintf("Player %d", board[0][i])
                 }
        }

        if board [0][0] == board[1][1] && board[1][1] == board[2][2] && board[0][0] != 0 {
                return fmt.Sprintf("Player %d", board[0][0])
        }
        if board [0][2] == board[1][1] && board[1][1] == board[2][0] && board[0][2] != 0 {
                return fmt.Sprintf("Player %d", board[0][2])
        }

        full := true
        for i := 0; i < 3; i++ {
                for j := 0; j < 3; j++ {
                        if board[i][j] == 0 {
                                full = false
                        }
                }
        }
        if full {
                return "CAT"
        }
        return ""
}