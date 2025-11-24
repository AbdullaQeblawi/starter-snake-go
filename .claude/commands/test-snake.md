Run a local BattleSnake test game:

1. Ensure the server is running (check if there's a background process)
2. If not running, start it with: `go run .` in background
3. Wait 2 seconds for server to start
4. Run a 4-snake battle test using:
   ```
   ~/go/bin/battlesnake play -W 11 -H 11 \
     --name "Evil Bolt" --url http://localhost:8000 \
     --name "Opponent 1" --url http://localhost:8000 \
     --name "Opponent 2" --url http://localhost:8000 \
     --name "Opponent 3" --url http://localhost:8000
   ```
5. Report the results: survival time, eliminations, and any interesting patterns
