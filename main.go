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

// Helper function to get the coordinate for a move direction
func getCoordFromMove(head Coord, move string) Coord {
	next := head
	switch move {
	case "up":
		next.Y += 1
	case "down":
		next.Y -= 1
	case "left":
		next.X -= 1
	case "right":
		next.X += 1
	}
	return next
}

// Flood fill algorithm to calculate reachable space from a position
func floodFill(start Coord, boardWidth, boardHeight int, allSnakes []Battlesnake) int {
	// Create a visited map
	visited := make(map[Coord]bool)

	// Create occupied map (all snake bodies)
	occupied := make(map[Coord]bool)
	for _, snake := range allSnakes {
		for i, segment := range snake.Body {
			// Don't mark tails as occupied if snake just ate (tail won't move)
			// For simplicity, mark all segments except the very last tail segment
			if i < len(snake.Body)-1 {
				occupied[segment] = true
			}
		}
	}

	// BFS to count reachable cells
	queue := []Coord{start}
	visited[start] = true
	count := 0

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		count++

		// Check all 4 directions
		directions := []Coord{
			{X: current.X, Y: current.Y + 1},     // up
			{X: current.X, Y: current.Y - 1},     // down
			{X: current.X - 1, Y: current.Y},     // left
			{X: current.X + 1, Y: current.Y},     // right
		}

		for _, next := range directions {
			// Check bounds
			if next.X < 0 || next.X >= boardWidth || next.Y < 0 || next.Y >= boardHeight {
				continue
			}

			// Check if already visited or occupied
			if visited[next] || occupied[next] {
				continue
			}

			visited[next] = true
			queue = append(queue, next)
		}
	}

	return count
}

// Predict possible next positions for an opponent's head
func predictOpponentMoves(snake Battlesnake, boardWidth, boardHeight int) []Coord {
	possibleMoves := []Coord{}
	head := snake.Head

	// All possible moves
	candidates := []Coord{
		{X: head.X, Y: head.Y + 1},     // up
		{X: head.X, Y: head.Y - 1},     // down
		{X: head.X - 1, Y: head.Y},     // left
		{X: head.X + 1, Y: head.Y},     // right
	}

	for _, pos := range candidates {
		// Check bounds
		if pos.X < 0 || pos.X >= boardWidth || pos.Y < 0 || pos.Y >= boardHeight {
			continue
		}

		// Enemy probably won't move into their own neck (backwards)
		if len(snake.Body) > 1 {
			neck := snake.Body[1]
			if pos.X == neck.X && pos.Y == neck.Y {
				continue
			}
		}

		possibleMoves = append(possibleMoves, pos)
	}

	return possibleMoves
}

// Helper function to score aggressive moves (higher is better)
func scoreAggressiveMove(nextPos Coord, myLength int, opponents []Battlesnake, myTail Coord, food []Coord, boardWidth, boardHeight int) int {
	score := 0

	// For each opponent, check if we can win a head-to-head
	for _, snake := range opponents {
		enemyHead := snake.Head

		// Check if enemy could move to same position (head-to-head)
		distance := manhattanDistance(nextPos, enemyHead)

		if distance == 1 {
			// We're adjacent to enemy head - potential head-to-head next turn
			if myLength > snake.Length {
				// We're bigger - be aggressive!
				score += 50
			} else if myLength == snake.Length {
				// Equal size - be cautious but slightly aggressive
				score += 10
			} else {
				// We're smaller - this is already handled by collision avoidance
				score -= 30
			}
		} else if distance == 2 {
			// Enemy is 2 moves away - position for potential confrontation
			if myLength > snake.Length {
				score += 20
			}
		}

		// Space control: Cut off opponent from food
		if len(food) > 0 && myLength >= snake.Length {
			// Find nearest food to enemy
			nearestFoodToEnemy := food[0]
			minDist := manhattanDistance(enemyHead, food[0])
			for _, f := range food {
				dist := manhattanDistance(enemyHead, f)
				if dist < minDist {
					minDist = dist
					nearestFoodToEnemy = f
				}
			}

			// Check if we're moving between enemy and their food
			distEnemyToFood := manhattanDistance(enemyHead, nearestFoodToEnemy)
			distNextPosToFood := manhattanDistance(nextPos, nearestFoodToEnemy)
			distEnemyToNextPos := manhattanDistance(enemyHead, nextPos)

			// If we're getting between enemy and food, bonus!
			if distNextPosToFood < distEnemyToFood && distEnemyToNextPos < distEnemyToFood {
				score += 30 // Cut them off!
			}
		}

		// Space control: Push enemy toward walls
		enemyDistFromWalls := 0
		if enemyHead.X <= 1 || enemyHead.X >= boardWidth-2 {
			enemyDistFromWalls++
		}
		if enemyHead.Y <= 1 || enemyHead.Y >= boardHeight-2 {
			enemyDistFromWalls++
		}

		// If enemy is near wall and we're moving closer, bonus
		if enemyDistFromWalls > 0 && distance <= 3 && myLength >= snake.Length {
			score += 15 * enemyDistFromWalls // Push them into danger!
		}
	}

	// Tail chasing for safe space control
	tailDistance := manhattanDistance(nextPos, myTail)

	// Check if there are any close threats
	hasCloseThreats := false
	for _, snake := range opponents {
		if manhattanDistance(nextPos, snake.Head) < 4 {
			hasCloseThreats = true
			break
		}
	}

	// When no close threats, strongly prefer following tail to avoid self-trap
	if !hasCloseThreats {
		if tailDistance == 1 {
			score += 40 // Very close to tail - great for territory control
		} else if tailDistance == 2 {
			score += 25
		} else if tailDistance == 3 {
			score += 15
		}
	} else {
		// With threats nearby, still follow tail but less strongly
		if tailDistance < 3 {
			score += 10
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

		// Step 3.5 - Avoid potential head-to-head collisions (unless we're larger)
		possibleEnemyMoves := predictOpponentMoves(snake, boardWidth, boardHeight)
		for _, enemyNextPos := range possibleEnemyMoves {
			// Only avoid if we're smaller or equal size (unless we want to be aggressive)
			if state.You.Length <= snake.Length {
				// Enemy might move here, so avoid it
				if myHead.X-1 == enemyNextPos.X && myHead.Y == enemyNextPos.Y {
					isMoveSafe["left"] = false
				}
				if myHead.X+1 == enemyNextPos.X && myHead.Y == enemyNextPos.Y {
					isMoveSafe["right"] = false
				}
				if myHead.Y-1 == enemyNextPos.Y && myHead.X == enemyNextPos.X {
					isMoveSafe["down"] = false
				}
				if myHead.Y+1 == enemyNextPos.Y && myHead.X == enemyNextPos.X {
					isMoveSafe["up"] = false
				}
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

		// Score each safe move based on distance to nearest food AND available space
		type FoodMove struct {
			move  string
			score int
		}
		foodMoves := []FoodMove{}
		allSnakes := state.Board.Snakes

		for _, move := range safeMoves {
			nextHead := getCoordFromMove(myHead, move)

			// Calculate distance from new position to food
			distance := manhattanDistance(nextHead, nearestFood)
			foodScore := -distance // Negative distance (closer is better)

			// Calculate available space - critical even when seeking food!
			spaceAvailable := floodFill(nextHead, boardWidth, boardHeight, allSnakes)
			minSpaceNeeded := state.You.Length + 3

			spaceScore := 0
			if needsFoodUrgently {
				// When desperate, take more risk but still avoid death traps
				if spaceAvailable < state.You.Length {
					spaceScore = -100 // Definite death trap
				} else {
					spaceScore = 0 // Take the risk
				}
			} else {
				// When not desperate, be more careful about space
				if spaceAvailable >= minSpaceNeeded*2 {
					spaceScore = 50
				} else if spaceAvailable >= minSpaceNeeded {
					spaceScore = 20
				} else if spaceAvailable >= state.You.Length {
					spaceScore = -30
				} else {
					spaceScore = -150 // Avoid!
				}
			}

			totalScore := foodScore + spaceScore
			foodMoves = append(foodMoves, FoodMove{move: move, score: totalScore})
		}

		// Find best food moves considering both distance and space
		bestMoves := []string{}
		bestFoodScore := -999999

		for _, fm := range foodMoves {
			if fm.score > bestFoodScore {
				bestFoodScore = fm.score
				bestMoves = []string{fm.move}
			} else if fm.score == bestFoodScore {
				bestMoves = append(bestMoves, fm.move)
			}
		}

		// Choose randomly among best moves
		if len(bestMoves) > 0 {
			nextMove = bestMoves[rand.Intn(len(bestMoves))]
		}
	} else {
		// Health is high or no food available - BE AGGRESSIVE!
		myLength := state.You.Length
		myTail := state.You.Body[len(state.You.Body)-1]

		// Check if we're in endgame (1v1)
		isEndgame := len(opponents) == 1
		var endgameOpponent Battlesnake
		if isEndgame && len(opponents) > 0 {
			endgameOpponent = opponents[0]
		}

		// Score each safe move based on aggression
		type ScoredMove struct {
			move  string
			score int
		}
		scoredMoves := []ScoredMove{}

		// Get all snakes for flood fill
		allSnakes := state.Board.Snakes

		for _, move := range safeMoves {
			nextPos := getCoordFromMove(myHead, move)
			aggressionScore := scoreAggressiveMove(nextPos, myLength, opponents, myTail, food, boardWidth, boardHeight)

			// Endgame bonus: In 1v1, heavily favor controlling more space
			endgameBonus := 0
			if isEndgame {
				mySpace := floodFill(nextPos, boardWidth, boardHeight, allSnakes)
				enemySpace := floodFill(endgameOpponent.Head, boardWidth, boardHeight, allSnakes)

				// If we control more space, big bonus!
				if mySpace > enemySpace {
					endgameBonus = 75 // Dominate space control
				} else if mySpace == enemySpace {
					endgameBonus = 25 // Equal footing
				} else {
					endgameBonus = -25 // Losing space battle
				}

				// If we're longer and enemy is low on health, be ultra aggressive
				if myLength > endgameOpponent.Length && endgameOpponent.Health < 40 {
					// Cut them off from food!
					if len(food) > 0 {
						nearestFoodToEnemy := food[0]
						minDist := manhattanDistance(endgameOpponent.Head, food[0])
						for _, f := range food {
							dist := manhattanDistance(endgameOpponent.Head, f)
							if dist < minDist {
								minDist = dist
								nearestFoodToEnemy = f
							}
						}

						distToFood := manhattanDistance(nextPos, nearestFoodToEnemy)
						enemyDistToFood := manhattanDistance(endgameOpponent.Head, nearestFoodToEnemy)

						if distToFood < enemyDistToFood {
							endgameBonus += 50 // Cut them off from food!
						}
					}
				}
			}

			aggressionScore += endgameBonus

			// Calculate available space using flood fill
			spaceAvailable := floodFill(nextPos, boardWidth, boardHeight, allSnakes)

			// Space is critical - heavily weight it
			// Need at least our body length + buffer in available space
			minSpaceNeeded := myLength + 5
			spaceScore := 0

			if spaceAvailable >= minSpaceNeeded*2 {
				spaceScore = 100 // Plenty of space
			} else if spaceAvailable >= minSpaceNeeded {
				spaceScore = 50 // Adequate space
			} else if spaceAvailable >= myLength {
				spaceScore = -50 // Tight space - risky!
			} else {
				spaceScore = -200 // Death trap - avoid!
			}

			totalScore := aggressionScore + spaceScore
			scoredMoves = append(scoredMoves, ScoredMove{move: move, score: totalScore})
		}

		// Find moves with highest total score (aggression + space)
		bestMoves := []string{}
		bestTotalScore := -999999

		for _, sm := range scoredMoves {
			if sm.score > bestTotalScore {
				bestTotalScore = sm.score
				bestMoves = []string{sm.move}
			} else if sm.score == bestTotalScore {
				bestMoves = append(bestMoves, sm.move)
			}
		}

		// Choose randomly among most aggressive moves
		if len(bestMoves) > 0 {
			nextMove = bestMoves[rand.Intn(len(bestMoves))]
		} else {
			nextMove = safeMoves[rand.Intn(len(safeMoves))]
		}
	}

	log.Printf("MOVE %d: %s (Health: %d)\n", state.Turn, nextMove, myHealth)
	return BattlesnakeMoveResponse{Move: nextMove}
}

func main() {
	RunServer()
}
