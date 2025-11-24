# BattleSnake AI - Evil Bolt Project Context

## Project Overview
This is an elite-tier competitive BattleSnake AI named "Evil Bolt" - a crimson predator with advanced algorithms and strategic capabilities.

**Primary Goal**: Dominate the BattleSnake arena through intelligent decision-making, aggressive tactics, and space control.

## Technology Stack
- **Language**: Go 1.13+
- **Framework**: BattleSnake API v1
- **Algorithms**: Flood Fill (BFS), Manhattan Distance, Move Prediction
- **Deployment**: Local development, cloud deployment ready

## Architecture

### Core Files
- `main.go` - Game logic, move decisions, all AI algorithms (~620 lines)
- `models.go` - BattleSnake API data structures
- `server.go` - HTTP server and request handlers
- `BATTLESNAKE_SUMMARY.md` - Complete documentation

### Key Functions in main.go

**Helper Functions**:
- `manhattanDistance(a, b Coord) int` - O(1) distance calculation
- `scoreFoodSafety(food, head, w, h, opponents) int` - Food safety evaluation
- `getCoordFromMove(head, move) Coord` - Position prediction
- `floodFill(start, w, h, snakes) int` - BFS space calculation
- `predictOpponentMoves(snake, w, h) []Coord` - Enemy move prediction
- `scoreAggressiveMove(...) int` - Comprehensive aggression scoring

**Main Logic**:
- `info() BattlesnakeInfoResponse` - Snake appearance
- `start(state GameState)` - Game initialization
- `move(state GameState) BattlesnakeMoveResponse` - **CORE DECISION LOGIC**
- `end(state GameState)` - Game cleanup

## Development Phases (Completed)

### Phase 1: Foundation & Safety
- Boundary collision avoidance
- Self-collision detection
- Enemy body collision avoidance
- **Result**: 243 turns survival

### Phase 2: Food Strategy
- Manhattan distance pathfinding
- Smart food selection (safety scoring)
- 3-tier health management (urgent/moderate/healthy)
- **Result**: 163 turns, 9 food items

### Phase 3: Aggressive Tactics
- Head-to-head combat (size-based)
- Space control & opponent cutoff
- Wall-pushing strategy
- Smart tail chasing
- **Result**: 38 turns, 4-snake battle

### Phase 4: Advanced Aggression
- Flood fill algorithm
- Opponent move prediction
- Endgame 1v1 optimization
- **Result**: Elite competitive performance

## Decision System

### Move Scoring Formula
```
Total Score = Base Safety + Food Score + Aggression Score + Space Score + Endgame Bonus
```

### Scoring Values
**Aggression**:
- Head-to-head vs smaller snake: +50
- Head-to-head vs equal snake: +10
- Cut off opponent from food: +30
- Push enemy toward wall: +15 per wall
- Tail chasing (safe): +40
- Tail chasing (with threats): +10

**Space (Flood Fill)**:
- Plenty (2x body length): +100
- Adequate (body length + buffer): +50
- Tight space: -50
- Death trap: -200

**Food**:
- Corner food: -70
- Wall food: -20 per wall
- Contested food: -40
- Distance penalty: -(distance * 2)

**Endgame (1v1)**:
- Space dominance: +75
- Equal space: +25
- Losing space: -25
- Food denial: +50

### Health-Based Strategy
- **Critical (<30)**: Desperate food seeking, ignore safety
- **Moderate (<60)**: Balance food safety and distance
- **Healthy (>60)**: Full aggression mode
- **Endgame (1v1)**: Ultra-strategic space control

## Common Development Tasks

### Testing Locally
```bash
# Start server
go run .

# Run solo game
~/go/bin/battlesnake play -W 11 -H 11 --name "Evil Bolt" --url http://localhost:8000 -g solo

# Run multiplayer (4 snakes)
~/go/bin/battlesnake play -W 11 -H 11 \
  --name "Evil Bolt" --url http://localhost:8000 \
  --name "Opponent 1" --url http://localhost:8000 \
  --name "Opponent 2" --url http://localhost:8000 \
  --name "Opponent 3" --url http://localhost:8000
```

### Development Workflow
```bash
# Make changes to main.go
# Restart server (Ctrl+C, then go run .)
# Test with battlesnake CLI
# Commit changes
git add main.go
git commit -m "Description of changes"
git push
```

## Optimization Areas

### Performance
- Flood fill can be CPU-intensive on large boards
- Consider memoization for repeated calculations
- Parallel move evaluation for speed

### Strategic
- A* pathfinding for complex mazes
- Voronoi diagram for territory analysis
- Multi-step lookahead planning
- Opening book strategies

### Code Quality
- Extract scoring constants to configuration
- Create separate files for different strategies
- Add unit tests for core algorithms
- Benchmark flood fill performance

## Important Constraints

### BattleSnake API
- Must respond within 500ms
- Board sizes typically 11x11 or 19x19
- Health decreases by 1 each turn
- Health resets to 100 when eating food
- Snake dies at health 0 or collision

### Algorithm Complexity
- Flood fill: O(n) where n = board cells
- Manhattan distance: O(1)
- Move prediction: O(k) where k = number of opponents
- Total per-move: O(n + k) - very efficient

## Git Workflow

**Branches**:
- `main` - Production-ready code
- `dev` - Active development (current)

**Commit Format**:
```
Phase X: Feature Name - Brief description

- Detailed change 1
- Detailed change 2
- Test results

🤖 Generated with [Claude Code](https://claude.com/claude-code)

Co-Authored-By: Claude <noreply@anthropic.com>
```

## Deployment Checklist

- [ ] Test locally with multiple opponents
- [ ] Verify performance under 500ms response time
- [ ] Merge dev to main
- [ ] Deploy to cloud platform (Heroku, Railway, etc.)
- [ ] Register on play.battlesnake.com
- [ ] Enter deployed URL
- [ ] Test in real arena
- [ ] Monitor performance and iterate

## Future Enhancement Ideas

**High Priority**:
- [ ] A* pathfinding for guaranteed food routes
- [ ] Voronoi space analysis for territory control
- [ ] Multi-step lookahead (2-3 moves ahead)

**Medium Priority**:
- [ ] Opening book strategies (first 5-10 moves)
- [ ] Dynamic weight tuning based on opponent behavior
- [ ] Minimax for endgame scenarios

**Low Priority**:
- [ ] Machine learning for pattern recognition
- [ ] Monte Carlo tree search
- [ ] Neural network move evaluation

## Key Performance Indicators

**Success Metrics**:
- Survival time (turns)
- Eliminations per game
- Win rate in 1v1
- Food efficiency (food eaten / turns)
- Space control (flood fill average)

**Current Performance**:
- Solo survival: 163-243 turns
- Multiplayer survival: 36-38 turns
- Eliminations: 1 per 4-snake game
- Food eaten: 9 in 163 turns

## Debugging Tips

**Common Issues**:
1. **Self-trapping**: Check flood fill calculation
2. **Slow response**: Profile flood fill performance
3. **Bad food choices**: Review food safety scoring
4. **Lost head-to-heads**: Verify size comparison logic
5. **Endgame losses**: Check space control in 1v1

**Debug Logging**:
```go
log.Printf("MOVE %d: %s (Health: %d)\n", state.Turn, nextMove, myHealth)
```

## Contact & Resources

- **Repository**: https://github.com/AbdullaQeblawi/starter-snake-go
- **BattleSnake Docs**: https://docs.battlesnake.com
- **Community**: https://play.battlesnake.com/discord
- **Platform**: https://play.battlesnake.com

---

*This context file helps Claude understand the project for future development sessions.*
