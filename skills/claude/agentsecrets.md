---
name: agentsecrets
description: "Trigger when: .env files present, hardcoded API keys/tokens in code, curl with auth headers, process.env/os.environ/os.Getenv credential references, HTTP 401/403 errors, user mentions secrets/API keys/credentials/tokens/environment variables, or user says agentsecrets."
---

# AgentSecrets Router

AgentSecrets is a **zero-knowledge** secrets manager: the server never sees plaintext values. All encryption and decryption happen locally on the agent's machine.

## Diagnostics

Run these two commands first and read their output before dispatching:

```
agentsecrets --version 2>/dev/null
```

```
agentsecrets status 2>/dev/null
```

## Dispatch Table

Based on the diagnostic output and user intent, invoke exactly one sub-skill:

| Condition | Invoke skill |
|---|---|
| Binary not found OR status shows no user/workspace/project | `agentsecrets-setup` |
| User needs to make an authenticated API call | `agentsecrets-call` |
| User needs to list, sync, diff, push, pull, or delete secrets | `agentsecrets-ops` |
| Hardcoded credentials in code, .env generation needed, or secure coding guidance | `agentsecrets-code` |

**Default:** When the situation does not clearly match one sub-skill, invoke `agentsecrets-ops`.

## Zero-Knowledge Rules

1. **Never display secret values.** Do not print, log, or echo plaintext secrets.
2. **Never ask the user to paste keys.** Do not request raw API keys, tokens, or credentials in chat.
3. **Always use `agentsecrets call`** to inject secrets into HTTP requests instead of interpolating values.
4. **Suggest the user delete messages** that accidentally contain raw keys.
5. **When a key is missing,** tell the user to run `agentsecrets secrets set KEY=value` in their terminal — never accept the value yourself.
