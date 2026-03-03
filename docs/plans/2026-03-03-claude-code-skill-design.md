# AgentSecrets Claude Code Skill — Design Document

**Date**: 2026-03-03
**Status**: Approved

## Problem

AgentSecrets has AI integrations for MCP and OpenClaw, but no native Claude Code skill. Developers using Claude Code in projects that need credentials have no structured way to use AgentSecrets. They either hardcode secrets, paste values into chat, or manually run CLI commands outside Claude Code.

## Solution

A Claude Code skill suite that teaches Claude Code how to detect credential needs, operate AgentSecrets CLI, and guide secure coding patterns — without ever seeing secret values.

## Architecture: Router + Sub-Skills

### Design Decision

A single monolithic skill would load 300-400 lines into context even when only one section is relevant. Instead, a lightweight **router skill** dispatches to focused **sub-skills**, keeping context usage minimal.

### File Layout

```
skills/claude/
├── agentsecrets.md          # Router (~50 lines) — detects situation, dispatches
├── agentsecrets-setup.md    # Installation & first-time setup
├── agentsecrets-ops.md      # Secrets operations (list, set, diff, push, pull, delete)
├── agentsecrets-call.md     # Authenticated API calls (6 auth styles)
└── agentsecrets-code.md     # Secure coding patterns & credential hygiene
```

### Router Skill (`agentsecrets.md`)

**Triggers** (broad auto-detect + explicit invocation):

Auto-detect:
- `.env`, `.env.example`, `.env.local` files in the project
- Hardcoded API keys or tokens in source code
- `curl` commands with auth headers
- Code referencing `process.env.`, `os.environ`, `os.Getenv` for credential-like keys
- HTTP 401/403 errors during development
- User mentions: secrets, API keys, credentials, tokens, environment variables

Explicit:
- User says "agentsecrets", "manage secrets", or invokes `/agentsecrets`

**Decision Tree**:

```
User message or auto-detect trigger
  │
  ├─ Is agentsecrets installed? (agentsecrets --version)
  │   └─ No → invoke agentsecrets-setup
  │
  ├─ Is there an active project? (agentsecrets status)
  │   └─ No → invoke agentsecrets-setup
  │
  ├─ User needs authenticated API call?
  │   └─ Yes → invoke agentsecrets-call
  │
  ├─ User needs secrets management? (list, sync, diff, etc.)
  │   └─ Yes → invoke agentsecrets-ops
  │
  └─ Detected credential issue in code? (.env generation, hardcoded keys, etc.)
      └─ Yes → invoke agentsecrets-code
```

The router runs `agentsecrets --version` and `agentsecrets status` before dispatching. If multiple sub-skills apply, it picks the most relevant one first and chains to others after.

### Sub-Skill: `agentsecrets-setup`

Handles:
- Detecting if binary is installed, guiding installation if not
- Running `agentsecrets init` (with storage-mode selection)
- Workspace creation/listing/switching
- Project creation/listing/use
- Verification via `agentsecrets status`

Key behaviors:
- Detects available package managers (npx, brew, pip, go) to recommend install method
- Guides user through interactive prompts (passwords, choices) — tells them what to expect
- Never runs install commands itself — user controls what binaries go on their machine
- Confirms each step completed before moving to next

### Sub-Skill: `agentsecrets-ops`

Handles:
- `agentsecrets secrets list` — show key names (never values)
- `agentsecrets secrets set` — guide user to run in their terminal (values never in chat)
- `agentsecrets secrets diff` — detect drift between local and cloud
- `agentsecrets secrets push` — upload local to cloud
- `agentsecrets secrets pull` — download cloud to local
- `agentsecrets secrets delete` — remove secrets
- Workspace/project switching when operating across environments

Key behaviors:
- Always runs `agentsecrets status` first to confirm context
- Runs `agentsecrets secrets diff` before push/pull operations
- Never displays, echoes, or logs secret values
- When a secret is missing, tells user exactly what to run: `agentsecrets secrets set KEY=value`
- Uses standard key naming conventions (STRIPE_KEY, OPENAI_KEY, etc.)

### Sub-Skill: `agentsecrets-call`

Handles:
- One-shot authenticated API calls via `agentsecrets call`
- All 6 auth styles: bearer, header, query, basic, body-field, form-field
- Proxy mode for multi-call workflows (start, status, stop, logs)
- Audit log inspection

Key behaviors:
- Always uses `agentsecrets call` — never raw curl/fetch with credentials
- Checks that required secret keys exist before attempting call
- Provides full command examples for each auth style
- Reports API responses without exposing credential values
- Suggests proxy mode when multiple calls are needed

### Sub-Skill: `agentsecrets-code`

Handles:
- Detecting hardcoded credentials in source code
- Generating `.env.example` files (key names only, no values)
- Writing code that reads from environment variables
- Flagging credential leaks during code review
- Guiding proxy integration for agent frameworks (LangChain, CrewAI, etc.)

Key behaviors:
- When writing code that needs credentials, uses `process.env.KEY_NAME` / `os.environ["KEY_NAME"]` / `os.Getenv("KEY_NAME")`
- Never generates code with hardcoded credential values
- When reviewing code, flags any string that looks like a key/token/secret
- Generates `.env.example` with comments describing each key's purpose
- For agent framework code, shows proxy integration pattern (`localhost:8765`)

## Installation

**Global (all projects):**
```bash
cp skills/claude/*.md ~/.claude/commands/
```

**Per-project:**
```bash
cp skills/claude/*.md .claude/commands/
```

One skill file, document both methods. Users choose based on whether they want AgentSecrets available everywhere or only in specific repos.

## Zero-Knowledge Constraints

All sub-skills enforce these rules:

1. NEVER display, echo, print, or log an actual secret value
2. NEVER ask the user to paste a key value into chat
3. NEVER use curl or direct HTTP with credentials — always `agentsecrets call`
4. ALWAYS run `agentsecrets status` before any secrets operation
5. ALWAYS run `agentsecrets secrets diff` before deployment workflows
6. ALWAYS suggest the user delete any chat message where they accidentally shared a raw key value
7. When a secret is needed but missing, tell the user to run `agentsecrets secrets set KEY=value` in their terminal — never accept the value in conversation

## Success Criteria

- Claude Code automatically detects credential needs in any project with the skill installed
- Setup is guided end-to-end for new users without them reading docs
- No credential value ever appears in Claude Code's conversation or tool output
- Code written by Claude Code uses env var patterns, not hardcoded values
- API calls go through `agentsecrets call`, not raw HTTP
- Context cost is minimal — only the relevant sub-skill loads
