# Claude Code Configuration

This directory contains Claude Code configuration files to help with ongoing BattleSnake development.

## Project Context

**`project-context.md`** - Comprehensive project overview including:
- Architecture and file structure
- Algorithm documentation
- Scoring system details
- Development phases
- Optimization areas
- Future enhancement ideas

This file provides Claude with full context about the project in future sessions.

## Slash Commands

Custom slash commands for common tasks:

### `/test-snake`
Runs a local 4-snake battle test and reports results.

**Usage**: Simply type `/test-snake` in chat

**What it does**:
- Starts the server if needed
- Runs a multiplayer test game
- Reports survival time and eliminations

### `/optimize-strategy`
Helps optimize scoring parameters based on behavior.

**Usage**: `/optimize-strategy`

**What it does**:
- Analyzes current scoring values
- Asks what behavior needs adjustment
- Suggests specific parameter changes
- Tests and validates improvements

### `/add-feature`
Guides you through adding a new strategic feature.

**Usage**: `/add-feature`

**What it does**:
- Reviews current architecture
- Designs new feature integration
- Implements incrementally
- Tests and commits changes

### `/analyze-performance`
Analyzes game logs and performance metrics.

**Usage**: `/analyze-performance`

**What it does**:
- Reviews recent game logs
- Calculates key metrics
- Identifies weaknesses
- Recommends improvements

## How to Use

1. **Open Claude Code in this directory**
   ```bash
   cd /Users/abdullaqeblawi/Dev/BattleSnake/starter-snake-go
   ```

2. **Claude automatically loads project-context.md** when you start a conversation

3. **Use slash commands** by typing `/command-name` in chat

4. **Context is preserved** across sessions for continuous development

## Adding New Commands

To add a custom command:

1. Create a new `.md` file in `.claude/commands/`
2. Name it `your-command.md`
3. Write instructions for Claude to follow
4. Use it with `/your-command`

**Example**:
```markdown
# .claude/commands/deploy.md
Deploy the BattleSnake to production:

1. Run tests locally
2. Merge dev to main
3. Push to cloud platform
4. Verify deployment
5. Register on play.battlesnake.com
```

## Tips for Effective Collaboration

**Be Specific**: When asking for changes, reference specific functions or line numbers

**Test Frequently**: Use `/test-snake` after each significant change

**Iterate**: Make small improvements and test rather than large rewrites

**Document**: Update project-context.md when adding major features

**Commit Often**: Keep git history clean with descriptive commits

## Future Enhancements

Potential additions to this configuration:

- [ ] Automated testing command with metrics tracking
- [ ] Performance benchmarking command
- [ ] Deployment automation command
- [ ] Strategy A/B testing command
- [ ] Game replay analysis command

---

*These configuration files help maintain context and streamline development across multiple Claude Code sessions.*
