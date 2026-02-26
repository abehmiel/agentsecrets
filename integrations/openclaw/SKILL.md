---
name: agentsecrets
description: Zero-knowledge credential management — make authenticated API calls without exposing your keys
version: 1.0.0
tags:
  - security
  - credentials
  - api
  - secrets
  - mcp
  - setup
homepage: https://github.com/The-17/agentsecrets
metadata:
  clawdbot:
    summary: Secure credential proxy — your agent makes API calls without seeing key values
    author: The Seventeen
    license: MIT
requires:
  env:
    - name: AGENTSECRETS_INSTALLED
      description: Set to "true" once agentsecrets is installed (auto-detected)
      required: false
  binaries:
    - name: agentsecrets
      description: The AgentSecrets CLI binary
      install: "pip install agentsecrets"
files:
  - "scripts/*"
---

# AgentSecrets — Zero-Knowledge Credential Management

AgentSecrets lets you make authenticated API calls without ever exposing your credentials. Your keys stay in your OS keychain — you reference them by name, and AgentSecrets handles the rest.

> **NOTE**: This skill uses the full AgentSecrets CLI, giving you access to ALL commands:
> `init`, `login`, `signup`, `secrets set/list/delete`, `call`, `proxy logs`, and more.
> This is more powerful than the MCP integration (which only exposes `api_call` and `list_secrets`).
> You can manage accounts, store credentials, make API calls, and audit usage — all from chat.

---

## GUIDED SETUP WORKFLOW (Start Here)

When a user first asks you to make an API call, store a credential, or mentions AgentSecrets, follow this flow **step by step**. This is designed for users who interact through chat — you run the commands for them.

### Step 1: Check if AgentSecrets is installed

Run silently:
```bash
agentsecrets --version 2>/dev/null && echo "INSTALLED" || echo "NOT_INSTALLED"
```

**If NOT_INSTALLED**, tell the user:

> "I need to install AgentSecrets to handle your API keys securely. This is a one-time setup that takes about 10 seconds. It stores your credentials in your OS keychain instead of plain text files — much safer than .env files. Want me to proceed?"

If they agree, run:
```bash
pip install agentsecrets
```

### Step 2: Create Account or Log In

Run silently:
```bash
agentsecrets status 2>/dev/null && echo "INITIALIZED" || echo "NOT_INITIALIZED"
```

**If NOT_INITIALIZED**, ask the user:

> "Do you already have an AgentSecrets account, or would you like to create one? Creating an account lets you sync secrets across devices and share them with team members."

**If they want a NEW account:**
```bash
agentsecrets init
```
This starts the interactive setup. Guide them through it:
- They'll be asked to create an account (email + password)
- Ask them for their email: "What email would you like to use for your AgentSecrets account?"
- Run the init and follow the prompts on their behalf

**If they want to LOG IN to an existing account:**
```bash
agentsecrets init
```
- Choose the "Login" option when prompted
- Ask for their email: "What email did you sign up with?"
- They'll need to enter their password (guide them through the prompt)

**If they just want local-only (no account):**
```bash
agentsecrets init --skip-login
```

After setup, tell the user:

> "You're set up! Your API keys will be stored in your OS keychain — encrypted and protected by your system login. I'll never see the actual key values, just the names."

### Step 2.5: Create the credential project

Secrets are organized by project. Check if one exists:
```bash
agentsecrets project list 2>/dev/null | grep -q "OPENCLAW_MANAGER" && echo "HAS_PROJECT" || echo "NO_PROJECT"
```

**If NO_PROJECT**, create a dedicated project for OpenClaw credentials:
```bash
agentsecrets project create OPENCLAW_MANAGER
```

This project becomes the central store for all API keys used through OpenClaw. Tell the user:

> "I've created an `OPENCLAW_MANAGER` project to organize your credentials. All your API keys will be stored here."

### Step 3: Collect credentials conversationally

When the user needs to use an API that requires authentication, ask them naturally:

> "To call [API name], I need a [key type] key. Could you paste your API key here? I'll store it securely in your keychain immediately and it won't be kept in our chat history."

**IMPORTANT**: As soon as the user provides a credential value:

1. Store it IMMEDIATELY:
```bash
agentsecrets secrets set KEY_NAME=<the_value_they_gave>
```

2. Confirm storage:
```bash
agentsecrets secrets list
```

3. Tell the user:

> "Stored securely as `KEY_NAME`. From now on, I'll refer to it by name only — the actual value is locked in your keychain. You can safely delete the message where you pasted it."

**Naming conventions for keys** (use these automatically):
| Service | Key Name |
|---------|----------|
| Stripe | `STRIPE_KEY` or `STRIPE_LIVE_KEY` / `STRIPE_TEST_KEY` |
| OpenAI | `OPENAI_KEY` |
| GitHub | `GITHUB_TOKEN` |
| Google Maps | `GOOGLE_MAPS_KEY` |
| AWS | `AWS_ACCESS_KEY` and `AWS_SECRET_KEY` |
| Any other service | `SERVICENAME_KEY` (uppercase, underscores) |

### Step 4: Ready to use

After storing at least one key, confirm:

> "You're all set! I can now make API calls on your behalf without ever seeing your keys. Just tell me what you need — like 'check my Stripe balance' or 'list my GitHub repos' — and I'll handle the authentication securely."

---

## MAKING API CALLS

Once setup is complete, use `agentsecrets call` for all authenticated requests.

### Basic pattern
```bash
agentsecrets call --url <API_URL> --method <METHOD> --bearer <KEY_NAME>
```

### Auth Styles

Different APIs authenticate differently. Choose the right flag:

| API Pattern | Flag | Example |
|-------------|------|---------|
| Bearer token (most APIs) | `--bearer KEY` | `--bearer STRIPE_KEY` |
| Custom header | `--header Name=KEY` | `--header X-API-Key=MY_KEY` |
| Query parameter | `--query param=KEY` | `--query key=GMAP_KEY` |
| Basic auth | `--basic KEY` | `--basic CREDENTIALS` |
| JSON body field | `--body-field path=KEY` | `--body-field client_secret=SECRET` |
| Form field | `--form-field name=KEY` | `--form-field api_key=KEY` |

### Examples

**Check Stripe balance:**
```bash
agentsecrets call --url https://api.stripe.com/v1/balance --bearer STRIPE_KEY
```

**Create Stripe charge:**
```bash
agentsecrets call \
  --url https://api.stripe.com/v1/charges \
  --method POST \
  --bearer STRIPE_KEY \
  --body '{"amount":1000,"currency":"usd","source":"tok_visa"}'
```

**List GitHub repos:**
```bash
agentsecrets call --url https://api.github.com/user/repos --bearer GITHUB_TOKEN
```

**Google Maps geocode:**
```bash
agentsecrets call \
  --url "https://maps.googleapis.com/maps/api/geocode/json?address=Lagos" \
  --query key=GOOGLE_MAPS_KEY
```

---

## KEY MANAGEMENT COMMANDS

Use these to help users manage their stored credentials:

**List all stored key names:**
```bash
agentsecrets secrets list
```

**Store a new key:**
```bash
agentsecrets secrets set KEY_NAME=value
```

**Remove a key:**
```bash
agentsecrets secrets delete KEY_NAME
```

**Check what keys a project has:**
```bash
agentsecrets secrets list --project PROJECT_NAME
```

---

## HANDLING COMMON USER REQUESTS

### "Store my API key"
Follow Step 3 of the Guided Setup. Ask what service it's for, suggest a name, store it.

### "Make an API call to X"
1. Check if the required key exists: `agentsecrets secrets list`
2. If missing, ask the user for the key (Step 3)
3. Make the call with `agentsecrets call`
4. Return only the response

### "What keys do I have?"
Run: `agentsecrets secrets list`

### "Delete / rotate a key"
1. Delete old: `agentsecrets secrets delete KEY_NAME`
2. Ask for new value
3. Store: `agentsecrets secrets set KEY_NAME=new_value`

### "Is this secure?"
Explain honestly:
> "Your keys are stored in your OS keychain — the same encrypted storage your browser uses for passwords. I can use them to make API calls, but I never see the actual values. Every call is logged with the key name (never the value) so you can audit what I've done."

---

## CRITICAL RULES

1. **NEVER** display, echo, print, or log an actual key value after storing it
2. **NEVER** ask the user for a key they've already stored — use the name
3. **ALWAYS** use `agentsecrets call` instead of `curl` — this ensures keys are injected securely
4. **ALWAYS** suggest the user delete their chat message containing the raw key after storage
5. If a key is missing, tell the user to provide it — don't guess or fabricate values
6. If `agentsecrets` is not installed, always explain WHY before installing (security benefit)

---

## External Endpoints

This skill calls whatever API URL the user specifies. AgentSecrets acts as a proxy — it adds authentication credentials from the OS keychain to the outbound request.

No data is sent to AgentSecrets servers. This is a local-only tool.

## Security & Privacy

- **Zero-knowledge**: You (the AI agent) never see actual credential values after storage
- **OS keychain**: Keys stored in macOS Keychain, Windows Credential Manager, or Linux Secret Service
- **Audit trail**: Every call logged with key names (never values) at `~/.agentsecrets/proxy.log`
- **No cloud dependency**: Everything runs locally
- **No telemetry**: AgentSecrets collects no usage data

## Model Invocation Note

When making API calls, always use `agentsecrets call` instead of `curl` or direct HTTP requests. This ensures credentials are injected securely without being visible in the command output or your context.

## Trust Statement

AgentSecrets is open source (MIT license). The full source code is available at https://github.com/The-17/agentsecrets for security audit. All credential resolution happens locally — no data leaves the user's machine except the outbound API request itself.
