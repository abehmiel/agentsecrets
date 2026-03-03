---
name: agentsecrets-setup
description: "Install AgentSecrets and complete first-time setup: account creation, workspace configuration, and project initialization."
---

# AgentSecrets Setup

## Step 1 — Installation Detection

Check whether the binary is already installed:

```bash
agentsecrets --version 2>/dev/null && echo "INSTALLED" || echo "NOT_INSTALLED"
```

If **INSTALLED**, skip to Step 2.

If **NOT_INSTALLED**, detect which package managers are available:

```bash
which npx && echo "npx available"
which brew && echo "brew available"
which pip && echo "pip available"
which go && echo "go available"
```

Present the matching installation options to the user. **Do NOT run any install command yourself** — the user decides what goes on their machine.

- **macOS/Linux (Homebrew):** `brew install The-17/tap/agentsecrets`
- **Universal (Node):** `npx @the-17/agentsecrets`
- **Python:** `pip install agentsecrets`
- **Go:** `go install github.com/The-17/agentsecrets/cmd/agentsecrets@latest`

Tell the user:

> AgentSecrets keeps your API keys in your OS keychain. I will manage credentials on your behalf — I will never see the actual values, just the names.

Wait for the user to confirm installation, then re-run the version check to verify.

## Step 2 — Account Setup

Initialize AgentSecrets with keychain-only storage (recommended):

```bash
agentsecrets init --storage-mode 1
```

Explain the two storage modes before running init:

- **Mode 1 — Keychain-only (recommended):** Secret values live exclusively in the OS keychain. A `.env.example` file is generated with key names only, containing no values. This is the more secure option.
- **Mode 2 — .env + keychain:** Values are stored in both a `.env` file and the keychain. Use this only if other tools in the project read `.env` directly.

Tell the user what to expect: the init command will interactively prompt for an email address and password. They should enter these directly in the terminal. Do not request or handle these values yourself.

## Step 3 — Workspace Setup

Check for existing workspaces:

```bash
agentsecrets workspace list
```

- If **no workspaces exist**, ask the user for a workspace name, then create and switch to it:

```bash
agentsecrets workspace create "Workspace Name"
agentsecrets workspace switch "Workspace Name"
```

- If **workspaces already exist**, list them and ask which one to use, then switch:

```bash
agentsecrets workspace switch "Workspace Name"
```

## Step 4 — Project Setup

Check for existing projects in the active workspace:

```bash
agentsecrets project list
```

- If **no projects exist**, derive a project name from the current repository or directory name. Convention: use SCREAMING_SNAKE_CASE or the repo name as-is. Confirm the name with the user, then create and activate:

```bash
agentsecrets project create PROJECT_NAME
agentsecrets project use PROJECT_NAME
```

- If **projects already exist**, list them and ask which one to use:

```bash
agentsecrets project use PROJECT_NAME
```

## Step 5 — Verification

Run the status check to confirm everything is configured:

```bash
agentsecrets status
```

Verify all three conditions are met:

1. **User is logged in** — if not, return to Step 2.
2. **Workspace is active** — if not, return to Step 3.
3. **Project is active** — if not, return to Step 4.

Once all three conditions pass, tell the user:

> Setup complete. Secrets management is ready. You can now use `/agentsecrets-ops` to manage secrets (list, set, sync, diff, push, pull) or `/agentsecrets-call` to make authenticated API calls with secret injection.

## Zero-Knowledge Rules

These rules apply at every step:

- **Never display, request, or handle actual secret values.**
- **Never run install commands yourself.** Present options and let the user choose.
- **Tell the user what to expect** from any interactive prompt before running it.
- When a credential is needed, direct the user to run `agentsecrets secrets set KEY=value` in their own terminal.
