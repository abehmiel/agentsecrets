# AgentSecrets Claude Code Skill — Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Create a router + sub-skills suite that teaches Claude Code how to detect credential needs, operate AgentSecrets CLI, and guide secure coding — without ever seeing secret values.

**Architecture:** A lightweight router skill (`agentsecrets.md`) auto-detects credential situations and dispatches to one of four focused sub-skills (setup, ops, call, code). Only the relevant sub-skill loads into context.

**Tech Stack:** Claude Code custom commands (markdown files with YAML frontmatter in `.claude/commands/`)

**Design Doc:** `docs/plans/2026-03-03-claude-code-skill-design.md`

---

### Task 1: Create directory structure

**Files:**
- Create: `skills/claude/` (directory)

**Step 1: Create the directory**

```bash
mkdir -p skills/claude
```

**Step 2: Verify**

```bash
ls -la skills/claude/
```
Expected: empty directory exists

**Step 3: Commit**

```bash
git add skills/claude/.gitkeep
git commit -m "chore: create skills/claude directory for Claude Code skill suite"
```

> Note: If git won't track an empty dir, create a `.gitkeep` file: `touch skills/claude/.gitkeep`

---

### Task 2: Write the router skill (`agentsecrets.md`)

**Files:**
- Create: `skills/claude/agentsecrets.md`

**Reference:** Design doc section "Router Skill" for decision tree.

**Step 1: Write the router skill**

The router must:
1. Have YAML frontmatter with `name: agentsecrets` and a `description` that covers all auto-detect triggers
2. State the zero-knowledge principle upfront
3. Run two diagnostic commands: `agentsecrets --version` and `agentsecrets status`
4. Based on results, invoke exactly one sub-skill via the Skill tool
5. Be ~50-60 lines max — it's a dispatcher, not a reference manual

```markdown
---
name: agentsecrets
description: "Zero-knowledge secrets management for AI agents. TRIGGER when: .env files present, hardcoded API keys/tokens in code, curl with auth headers, process.env/os.environ/os.Getenv credential references, HTTP 401/403 errors, user mentions secrets/API keys/credentials/tokens/environment variables, or user says 'agentsecrets'."
---

# AgentSecrets — Zero-Knowledge Credential Management

You manage credentials without ever seeing their values. AgentSecrets keeps secrets in the OS keychain and injects them at the HTTP transport layer.

**Core rule:** You never see, display, or handle secret values. Only key names.

## Dispatch

Run these two commands silently to understand the current state:

1. `agentsecrets --version 2>/dev/null` — is it installed?
2. `agentsecrets status 2>/dev/null` — is there an active user/workspace/project?

Then dispatch:

| Condition | Action |
|---|---|
| Binary not found OR status shows no user/workspace/project | Invoke skill `agentsecrets-setup` |
| User needs to make an authenticated API call | Invoke skill `agentsecrets-call` |
| User needs to list, sync, diff, push, pull, or delete secrets | Invoke skill `agentsecrets-ops` |
| Hardcoded credentials in code, .env generation needed, or secure coding guidance | Invoke skill `agentsecrets-code` |

If the situation doesn't clearly match one sub-skill, default to `agentsecrets-ops` (it starts with a status check that will catch setup issues).

## Zero-Knowledge Rules (apply to ALL sub-skills)

1. NEVER display, echo, print, or log an actual secret value
2. NEVER ask the user to paste a key value into chat
3. NEVER use curl or direct HTTP with credentials — always `agentsecrets call`
4. ALWAYS suggest the user delete any chat message where they accidentally shared a raw key value
5. When a secret is needed but missing, tell the user to run `agentsecrets secrets set KEY=value` in their terminal
```

**Step 2: Read back the file and verify**

- Confirm YAML frontmatter has `name` and `description`
- Confirm description mentions all auto-detect triggers
- Confirm dispatch table references all 4 sub-skills by exact name
- Confirm zero-knowledge rules section exists
- Confirm file is under 70 lines

**Step 3: Commit**

```bash
git add skills/claude/agentsecrets.md
git commit -m "feat: add agentsecrets router skill for Claude Code"
```

---

### Task 3: Write the setup sub-skill (`agentsecrets-setup.md`)

**Files:**
- Create: `skills/claude/agentsecrets-setup.md`

**Reference:** Design doc section "Sub-Skill: agentsecrets-setup" and OpenClaw SKILL.md Steps 1-4.

**Step 1: Write the setup sub-skill**

The setup skill must:
1. Have YAML frontmatter with `name: agentsecrets-setup`
2. Check installation, guide install if missing (detect npx/brew/pip/go)
3. Guide `agentsecrets init` with storage-mode explanation
4. Guide workspace create/list/switch
5. Guide project create/list/use
6. Verify with `agentsecrets status` after each step
7. Never run install commands itself — user controls binary installation
8. Tell user what to expect from interactive prompts (passwords, etc.)

Content should cover:

**Section 1 — Installation Detection:**
```bash
agentsecrets --version 2>/dev/null && echo "INSTALLED" || echo "NOT_INSTALLED"
```
If not installed, detect package managers:
```bash
which npx && echo "npx available"
which brew && echo "brew available"
which pip && echo "pip available"
which go && echo "go available"
```
Recommend: `brew install The-17/tap/agentsecrets` (macOS), `npx @the-17/agentsecrets` (universal), `pip install agentsecrets` (Python), `go install github.com/The-17/agentsecrets/cmd/agentsecrets@latest` (Go).

Tell the user: "AgentSecrets keeps your API keys in your OS keychain. I'll manage credentials on your behalf — I'll never see the actual values, just the names."

**Section 2 — Account Setup:**
```bash
agentsecrets init --storage-mode 1
```
Explain storage modes:
- Mode 1 (keychain-only): Values only in OS keychain, `.env.example` with key names only — more secure
- Mode 2 (.env + keychain): Values in both `.env` file and keychain — compatible with tools expecting `.env`

The init command will prompt for email and password. Tell user what to expect.

**Section 3 — Workspace Setup:**
```bash
agentsecrets workspace list
agentsecrets workspace create "Workspace Name"
agentsecrets workspace switch "Workspace Name"
```

**Section 4 — Project Setup:**
```bash
agentsecrets project list
agentsecrets project create PROJECT_NAME
agentsecrets project use PROJECT_NAME
```

**Section 5 — Verification:**
```bash
agentsecrets status
```
Confirm: user logged in, workspace active, project active. If any are missing, loop back to the relevant section.

**Step 2: Read back the file and verify**

- Confirm frontmatter has `name: agentsecrets-setup`
- Confirm all 5 sections present
- Confirm no secret values appear anywhere
- Confirm install commands are presented as user instructions, not auto-run

**Step 3: Commit**

```bash
git add skills/claude/agentsecrets-setup.md
git commit -m "feat: add agentsecrets-setup sub-skill for installation and first-time config"
```

---

### Task 4: Write the ops sub-skill (`agentsecrets-ops.md`)

**Files:**
- Create: `skills/claude/agentsecrets-ops.md`

**Reference:** Design doc section "Sub-Skill: agentsecrets-ops" and OpenClaw SKILL.md Steps 5-6.

**Step 1: Write the ops sub-skill**

The ops skill must:
1. Have YAML frontmatter with `name: agentsecrets-ops`
2. Always start with `agentsecrets status` to confirm context
3. Cover all secrets commands: list, set (guidance only), diff, push, pull, delete
4. Include workspace/project switching for multi-environment workflows
5. Include standard key naming conventions table
6. Enforce zero-knowledge rules throughout

Content should cover:

**Section 1 — Status Check (always first):**
```bash
agentsecrets status
```

**Section 2 — List Secrets:**
```bash
agentsecrets secrets list
```
Returns key names only, never values.

**Section 3 — Add/Update Secrets:**
Never accept values in conversation. Tell the user:
> "I need `KEY_NAME` to proceed. Please run in your terminal:
> `agentsecrets secrets set KEY_NAME=your_value`
> Let me know when done."

Then verify: `agentsecrets secrets list`

**Section 4 — Drift Detection:**
```bash
agentsecrets secrets diff
```
Shows added/removed/changed/unchanged between local and cloud.

**Section 5 — Sync Operations:**
```bash
agentsecrets secrets pull   # cloud → local
agentsecrets secrets push   # local → cloud
```
Always run `diff` first to show user what will change.

**Section 6 — Delete:**
```bash
agentsecrets secrets delete KEY_NAME
```

**Section 7 — Environment Switching:**
```bash
agentsecrets workspace switch "production"
agentsecrets project use my-api
agentsecrets secrets pull
```

**Section 8 — Key Naming Conventions:**
Table of standard names: STRIPE_KEY, OPENAI_KEY, GITHUB_TOKEN, AWS_ACCESS_KEY, AWS_SECRET_KEY, etc.

**Step 2: Read back and verify**

- Confirm frontmatter correct
- Confirm status check is listed as mandatory first step
- Confirm no secret values appear
- Confirm set command is user-directed (not auto-run with values)

**Step 3: Commit**

```bash
git add skills/claude/agentsecrets-ops.md
git commit -m "feat: add agentsecrets-ops sub-skill for secrets lifecycle management"
```

---

### Task 5: Write the call sub-skill (`agentsecrets-call.md`)

**Files:**
- Create: `skills/claude/agentsecrets-call.md`

**Reference:** Design doc section "Sub-Skill: agentsecrets-call" and OpenClaw SKILL.md Steps 7-9.

**Step 1: Write the call sub-skill**

The call skill must:
1. Have YAML frontmatter with `name: agentsecrets-call`
2. Cover all 6 auth styles with examples
3. Cover proxy mode for multi-call workflows
4. Cover audit log inspection
5. Enforce: always `agentsecrets call`, never raw curl/fetch
6. Check that required keys exist before attempting calls

Content should cover:

**Section 1 — Pre-flight Check:**
```bash
agentsecrets status
agentsecrets secrets list
```
Verify the key you need exists before calling.

**Section 2 — One-Shot Calls:**
Pattern: `agentsecrets call --url URL --method METHOD --AUTH_STYLE KEY_NAME`

Auth style reference table with examples:

| Style | Flag | Example |
|---|---|---|
| Bearer | `--bearer KEY` | `agentsecrets call --url https://api.stripe.com/v1/balance --bearer STRIPE_KEY` |
| Custom header | `--header Name=KEY` | `agentsecrets call --url https://api.sendgrid.com/v3/mail/send --method POST --header X-Api-Key=SENDGRID_KEY --body '{...}'` |
| Query param | `--query param=KEY` | `agentsecrets call --url "https://maps.googleapis.com/maps/api/geocode/json?address=NYC" --query key=GOOGLE_MAPS_KEY` |
| Basic auth | `--basic KEY` | `agentsecrets call --url https://jira.example.com/rest/api/2/issue --basic JIRA_CREDS` |
| JSON body | `--body-field path=KEY` | `agentsecrets call --url https://auth.example.com/token --method POST --body-field client_secret=OAUTH_SECRET` |
| Form field | `--form-field field=KEY` | `agentsecrets call --url https://api.example.com/auth --method POST --form-field api_key=API_KEY` |

Include multi-method examples (GET, POST, PUT, DELETE) and combining multiple auth flags.

**Section 3 — Proxy Mode:**
```bash
agentsecrets proxy start [--port 8765]
agentsecrets proxy status
agentsecrets proxy stop
```
Suggest proxy when: making 3+ calls, integrating with agent frameworks, or building multi-step workflows.

HTTP proxy pattern:
```
POST http://localhost:8765/proxy
X-AS-Target-URL: https://api.stripe.com/v1/balance
X-AS-Inject-Bearer: STRIPE_KEY
```

**Section 4 — Audit:**
```bash
agentsecrets proxy logs
agentsecrets proxy logs --last 20
agentsecrets proxy logs --secret STRIPE_KEY
```
Shows: timestamp, method, URL, key names, status code, duration. Never values.

**Step 2: Read back and verify**

- Confirm all 6 auth styles documented with examples
- Confirm proxy mode section present
- Confirm no secret values in any example
- Confirm pre-flight check is mandatory

**Step 3: Commit**

```bash
git add skills/claude/agentsecrets-call.md
git commit -m "feat: add agentsecrets-call sub-skill for authenticated API calls"
```

---

### Task 6: Write the code sub-skill (`agentsecrets-code.md`)

**Files:**
- Create: `skills/claude/agentsecrets-code.md`

**Reference:** Design doc section "Sub-Skill: agentsecrets-code".

**Step 1: Write the code sub-skill**

The code skill must:
1. Have YAML frontmatter with `name: agentsecrets-code`
2. Cover detecting hardcoded credentials in source code
3. Cover generating `.env.example` files
4. Cover writing code that reads from env vars (multi-language)
5. Cover flagging credential leaks in code review
6. Cover proxy integration for agent frameworks

Content should cover:

**Section 1 — Detect Hardcoded Credentials:**
When reviewing or writing code, flag any of these patterns:
- String literals that look like API keys (long alphanumeric strings, `sk-`, `pk_`, `ghp_`, `xoxb-`, etc.)
- Variables named `api_key`, `secret`, `token`, `password` assigned string literals
- `Authorization: Bearer` headers with inline values
- URLs containing `?key=` or `?token=` with inline values

When found: warn the user, suggest moving to AgentSecrets, invoke `agentsecrets-ops` to set the key.

**Section 2 — Generate `.env.example`:**
When a project needs env vars documented:
```bash
agentsecrets secrets list
```
Then generate `.env.example` with key names and descriptive comments:
```
# Stripe API key for payment processing
STRIPE_KEY=
# OpenAI API key for LLM calls
OPENAI_KEY=
# Database connection string
DATABASE_URL=
```
Never include values. Only key names and descriptions.

**Section 3 — Write Code Using Env Vars:**
Language-specific patterns:

JavaScript/TypeScript:
```javascript
const apiKey = process.env.STRIPE_KEY;
if (!apiKey) throw new Error("STRIPE_KEY not set");
```

Python:
```python
import os
api_key = os.environ["STRIPE_KEY"]
```

Go:
```go
apiKey := os.Getenv("STRIPE_KEY")
if apiKey == "" {
    log.Fatal("STRIPE_KEY not set")
}
```

Rust:
```rust
let api_key = std::env::var("STRIPE_KEY").expect("STRIPE_KEY not set");
```

**Section 4 — Code Review: Credential Hygiene:**
When reviewing code (PR review, refactoring, etc.), check for:
- Hardcoded credential values → suggest env vars + AgentSecrets
- `.env` files committed to git → suggest `.gitignore` entry + `.env.example` instead
- Secrets logged to console/files → suggest removing log statements
- Credentials passed as function arguments → suggest reading from env at point of use

**Section 5 — Agent Framework Proxy Integration:**
For code using LangChain, CrewAI, AutoGen, or similar:

```python
import requests

# Instead of: headers = {"Authorization": f"Bearer {api_key}"}
# Use AgentSecrets proxy:
response = requests.post(
    "http://localhost:8765/proxy",
    headers={
        "X-AS-Target-URL": "https://api.openai.com/v1/chat/completions",
        "X-AS-Inject-Bearer": "OPENAI_KEY",
    },
    json={"model": "gpt-4", "messages": [...]},
)
```

**Step 2: Read back and verify**

- Confirm all 5 sections present
- Confirm code examples never contain actual credential values
- Confirm multi-language env var patterns included
- Confirm proxy integration example present

**Step 3: Commit**

```bash
git add skills/claude/agentsecrets-code.md
git commit -m "feat: add agentsecrets-code sub-skill for secure coding patterns"
```

---

### Task 7: Update skills/README.md

**Files:**
- Modify: `skills/README.md`

**Step 1: Update the README to document the Claude Code skill suite**

Add a section under "Available Skills" for Claude Code:

```markdown
### Claude Code

**Directory**: `claude/`

A router + sub-skills suite for Claude Code custom commands.

**Files:**
- `agentsecrets.md` — Router that auto-detects credential needs and dispatches
- `agentsecrets-setup.md` — Installation & first-time setup
- `agentsecrets-ops.md` — Secrets lifecycle (list, set, diff, push, pull, delete)
- `agentsecrets-call.md` — Authenticated API calls (6 auth styles + proxy)
- `agentsecrets-code.md` — Secure coding patterns & credential hygiene

**Install globally (all projects):**
```bash
cp skills/claude/*.md ~/.claude/commands/
```

**Install per-project:**
```bash
cp skills/claude/*.md .claude/commands/
```
```

Also update the "Using These Skills" section with Claude Code instructions and remove the "coming soon" language for platforms that still don't have skills.

**Step 2: Read back and verify**

- Claude Code section present with correct file list
- Both installation methods documented
- No broken markdown

**Step 3: Commit**

```bash
git add skills/README.md
git commit -m "docs: update skills README with Claude Code skill suite"
```

---

### Task 8: Verify the full suite

**Files:**
- Read: all 5 files in `skills/claude/`

**Step 1: Verify router references**

Read `skills/claude/agentsecrets.md` and confirm:
- It references `agentsecrets-setup`, `agentsecrets-ops`, `agentsecrets-call`, `agentsecrets-code` by exact name
- Each sub-skill file has a matching `name:` in frontmatter

**Step 2: Verify zero-knowledge compliance**

Grep all skill files for potential credential leaks:
```bash
grep -rn "sk-\|sk_live\|Bearer [A-Za-z0-9]\|password123\|secret_value" skills/claude/
```
Expected: zero matches

**Step 3: Verify no file exceeds reasonable length**

```bash
wc -l skills/claude/*.md
```
Expected: router ~50-60 lines, sub-skills ~80-150 lines each

**Step 4: Verify frontmatter format**

Each file must have:
```yaml
---
name: agentsecrets[-subskill]
description: "..."
---
```

Check all 5 files have valid frontmatter.

---

### Task 9: Remove `.gitkeep` and final commit

**Step 1: Clean up**

If `.gitkeep` exists in `skills/claude/`, remove it (the directory now has real files).

```bash
rm -f skills/claude/.gitkeep
```

**Step 2: Final commit with all files**

```bash
git add skills/claude/ skills/README.md
git status
git commit -m "feat: complete AgentSecrets Claude Code skill suite (router + 4 sub-skills)"
```
