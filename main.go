package main

// Welcome to
// __________         __    __  .__                               __
// \______   \_____ _/  |__/  |_|  |   ____   ______ ____ _____  |  | __ ____
//  |    |  _/\__  \\   __\   __\  | _/ __ \ /  ___//    \\__  \ |  |/ // __ \
//  |    |   \ / __ \|  |  |  | |  |_\  ___/ \___ \|   |  \/ __ \|    <\  ___/
//  |________/(______/__|  |__| |____/\_____>______>___|__(______/__|__\\_____>
//
// This file can be a nice home for your Battlesnake logic and helper functions.
//
// To get you started we've included code to prevent your Battlesnake from moving backwards.
// For more info see docs.battlesnake.com

import (
	"log"
	"math/rand"
)

// Helper function to calculate Manhattan distance between two points
func manhattanDistance(a, b Coord) int {
	dx := a.X - b.X
	if dx < 0 {
		dx = -dx
	}
	dy := a.Y - b.Y
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}

// Helper function to score food safety (higher is better)
func scoreFoodSafety(food Coord, myHead Coord, boardWidth, boardHeight int, opponents []Battlesnake) int {
	score := 100

	// Penalize food near walls (corners are especially dangerous)
	distanceFromWalls := 0
	if food.X == 0 || food.X == boardWidth-1 {
		score -= 20
		distanceFromWalls++
	}
	if food.Y == 0 || food.Y == boardHeight-1 {
		score -= 20
		distanceFromWalls++
	}
	// Extra penalty for corners
	if distanceFromWalls == 2 {
		score -= 30
	}

	// Penalize if larger snakes are closer to this food
	myDistance := manhattanDistance(myHead, food)
	for _, snake := range opponents {
		enemyDistance := manhattanDistance(snake.Head, food)
		if enemyDistance < myDistance && snake.Length >= len(opponents) {
			score -= 40 // Larger snake is closer, avoid this food
		}
	}

	return score
}

// info is called when you create your Battlesnake on play.battlesnake.com
// and controls your Battlesnake's appearance
// TIP: If you open your Battlesnake URL in a browser you should see this data
func info() BattlesnakeInfoResponse {
	log.Println("INFO")

	return BattlesnakeInfoResponse{
		APIVersion: "1",
		Author:     "AbdullaQeblawi",
		Color:      "#DC143C", // Crimson red - aggressive
		Head:       "evil",    // Evil eyes for intimidation
		Tail:       "bolt",    // Lightning bolt tail
	}
}

// start is called when your Battlesnake begins a game
func start(state GameState) {
	log.Println("GAME START")
}

// end is called when your Battlesnake finishes a game
func end(state GameState) {
	log.Printf("GAME OVER\n\n")
}

// move is called on every turn and returns your next move
// Valid moves are "up", "down", "left", or "right"
// See https://docs.battlesnake.com/api/example-move for available data
func move(state GameState) BattlesnakeMoveResponse {

	isMoveSafe := map[string]bool{
		"up":    true,
		"down":  true,
		"left":  true,
		"right": true,
	}

	// We've included code to prevent your Battlesnake from moving backwards
	myHead := state.You.Body[0] // Coordinates of your head
	myNeck := state.You.Body[1] // Coordinates of your "neck"

	if myNeck.X < myHead.X { // Neck is left of head, don't move left
		isMoveSafe["left"] = false

	} else if myNeck.X > myHead.X { // Neck is right of head, don't move right
		isMoveSafe["right"] = false

	} else if myNeck.Y < myHead.Y { // Neck is below head, don't move down
		isMoveSafe["down"] = false

	} else if myNeck.Y > myHead.Y { // Neck is above head, don't move up
		isMoveSafe["up"] = false
	}

	// Step 1 - Prevent your Battlesnake from moving out of bounds
	boardWidth := state.Board.Width
	boardHeight := state.Board.Height

	if myHead.X == 0 {
		isMoveSafe["left"] = false
	}
	if myHead.X == boardWidth-1 {
		isMoveSafe["right"] = false
	}
	if myHead.Y == 0 {
		isMoveSafe["down"] = false
	}
	if myHead.Y == boardHeight-1 {
		isMoveSafe["up"] = false
	}

	// Step 2 - Prevent your Battlesnake from colliding with itself
	myBody := state.You.Body
	for _, segment := range myBody {
		if myHead.X-1 == segment.X && myHead.Y == segment.Y {
			isMoveSafe["left"] = false
		}
		if myHead.X+1 == segment.X && myHead.Y == segment.Y {
			isMoveSafe["right"] = false
		}
		if myHead.Y-1 == segment.Y && myHead.X == segment.X {
			isMoveSafe["down"] = false
		}
		if myHead.Y+1 == segment.Y && myHead.X == segment.X {
			isMoveSafe["up"] = false
		}
	}

	// Step 3 - Prevent your Battlesnake from colliding with other Battlesnakes
	opponents := state.Board.Snakes
	for _, snake := range opponents {
		if snake.ID == state.You.ID {
			continue // Skip our own snake, already handled above
		}
		for _, segment := range snake.Body {
			if myHead.X-1 == segment.X && myHead.Y == segment.Y {
				isMoveSafe["left"] = false
			}
			if myHead.X+1 == segment.X && myHead.Y == segment.Y {
				isMoveSafe["right"] = false
			}
			if myHead.Y-1 == segment.Y && myHead.X == segment.X {
				isMoveSafe["down"] = false
			}
			if myHead.Y+1 == segment.Y && myHead.X == segment.X {
				isMoveSafe["up"] = false
			}
		}
	}

	// Are there any safe moves left?
	safeMoves := []string{}
	for move, isSafe := range isMoveSafe {
		if isSafe {
			safeMoves = append(safeMoves, move)
		}
	}

	if len(safeMoves) == 0 {
		log.Printf("MOVE %d: No safe moves detected! Moving down\n", state.Turn)
		return BattlesnakeMoveResponse{Move: "down"}
	}

	// Step 4 - Move towards food to regain health and survive longer
	food := state.Board.Food
	myHealth := state.You.Health
	nextMove := safeMoves[0] // Default to first safe move

	// Health management: determine if we need food urgently
	needsFoodUrgently := myHealth < 30
	needsFoodSoon := myHealth < 60

	if len(food) > 0 && (needsFoodUrgently || needsFoodSoon) {
		// Find the best food (combination of nearest + safest)
		bestFood := food[0]
		var bestScore int

		if needsFoodUrgently {
			// When health is critical, only consider distance - ignore safety
			bestScore = -manhattanDistance(myHead, bestFood)
		} else {
			// When health is moderate, balance safety and distance
			bestScore = scoreFoodSafety(bestFood, myHead, boardWidth, boardHeight, opponents) - manhattanDistance(myHead, bestFood)*2
		}

		for _, f := range food {
			distance := manhattanDistance(myHead, f)
			var score int

			if needsFoodUrgently {
				// Critical health: only distance matters
				score = -distance
			} else {
				// Moderate health: balance safety and distance
				safety := scoreFoodSafety(f, myHead, boardWidth, boardHeight, opponents)
				score = safety - distance*2
			}

			if score > bestScore {
				bestScore = score
				bestFood = f
			}
		}

		nearestFood := bestFood

		// Score each safe move based on distance to nearest food
		bestMoves := []string{}
		bestDistance := 999999

		for _, move := range safeMoves {
			// Calculate where this move would take us
			nextHead := myHead
			switch move {
			case "up":
				nextHead.Y += 1
			case "down":
				nextHead.Y -= 1
			case "left":
				nextHead.X -= 1
			case "right":
				nextHead.X += 1
			}

			// Calculate distance from new position to food
			distance := manhattanDistance(nextHead, nearestFood)

			if distance < bestDistance {
				bestDistance = distance
				bestMoves = []string{move}
			} else if distance == bestDistance {
				bestMoves = append(bestMoves, move)
			}
		}

		// Choose randomly among best moves
		if len(bestMoves) > 0 {
			nextMove = bestMoves[rand.Intn(len(bestMoves))]
		}
	} else {
		// Health is high or no food available - choose safe move strategically
		// For now, just pick a random safe move (Phase 3 will add territory control)
		nextMove = safeMoves[rand.Intn(len(safeMoves))]
	}

	log.Printf("MOVE %d: %s (Health: %d)\n", state.Turn, nextMove, myHealth)
	return BattlesnakeMoveResponse{Move: nextMove}
}

func main() {
	RunServer()
}
