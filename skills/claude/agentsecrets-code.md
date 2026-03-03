---
name: agentsecrets-code
description: "Secure coding patterns and credential hygiene: detect hardcoded secrets, generate .env.example files, write code using environment variables, and integrate with the AgentSecrets proxy. Invoked by the router when hardcoded credentials are found in source code, .env generation is needed, or secure coding guidance is required."
---

# AgentSecrets Code

## Step 1 — Detect Hardcoded Credentials

When reviewing or writing code, flag these patterns as potential hardcoded secrets:

- **API key prefixes:** long alphanumeric strings starting with `sk-`, `pk_`, `ghp_`, `xoxb-`, `AKIA`, `ya29.`
- **Suspicious variable assignments:** variables named `api_key`, `secret`, `token`, `password`, or `credential` assigned to string literals
- **Inline auth headers:** `Authorization: Bearer` with an inline string value
- **URL parameters:** URLs containing `?key=` or `?token=` with inline values
- **Committed `.env` files:** `.env` files containing actual values tracked by version control

When a hardcoded credential is found:

1. **Warn** the user about the specific credential and its location
2. **Suggest** replacing the value with an environment variable reference
3. **Invoke `agentsecrets-ops`** to help the user store the value securely via `agentsecrets secrets set KEY=value` in their terminal

## Step 2 — Generate `.env.example`

When a project needs its environment variables documented, first check what keys exist:

```bash
agentsecrets secrets list
```

Then generate a `.env.example` file with key names and descriptive comments only:

```
# Stripe API key for payment processing
STRIPE_KEY=
# OpenAI API key for LLM calls
OPENAI_KEY=
# Database connection string
DATABASE_URL=
```

Rules:

- **NEVER include actual values** in `.env.example`
- Add a descriptive comment above each key explaining its purpose
- Add `.env` to `.gitignore` if not already present
- The `.env.example` file is safe to commit to version control

## Step 3 — Write Code Using Environment Variables

When writing code that needs credentials, ALWAYS use environment variable patterns. Never hardcode.

**JavaScript / TypeScript:**

```javascript
const stripeKey = process.env.STRIPE_KEY;
if (!stripeKey) throw new Error("STRIPE_KEY not set");
```

**Python:**

```python
import os
api_key = os.environ["STRIPE_KEY"]
```

**Go:**

```go
apiKey := os.Getenv("STRIPE_KEY")
if apiKey == "" {
    log.Fatal("STRIPE_KEY not set")
}
```

**Rust:**

```rust
let api_key = std::env::var("STRIPE_KEY").expect("STRIPE_KEY not set");
```

**Ruby:**

```ruby
api_key = ENV.fetch("STRIPE_KEY")
```

## Step 4 — Code Review: Credential Hygiene

When reviewing code (PR reviews, refactoring, auditing), check for:

- **Hardcoded credential values** — replace with env var reference + store in AgentSecrets
- **`.env` files committed to git** — add `.env` to `.gitignore`, use `.env.example` instead
- **Secrets logged to console or files** — remove the log statements immediately
- **Credentials passed as function arguments** — read from env at point of use instead
- **Secrets in config files** — reference env vars instead of inline values

## Step 5 — Agent Framework Proxy Integration

When writing code for agent frameworks (LangChain, CrewAI, AutoGen, etc.), use the AgentSecrets proxy instead of passing credentials directly. The proxy must be running first:

```bash
agentsecrets proxy start
```

**Python (requests):**

```python
import requests

response = requests.post(
    "http://localhost:8765/proxy",
    headers={
        "X-AS-Target-URL": "https://api.openai.com/v1/chat/completions",
        "X-AS-Inject-Bearer": "OPENAI_KEY",
    },
    json={"model": "gpt-4", "messages": [{"role": "user", "content": "Hello"}]},
)
```

**JavaScript (fetch):**

```javascript
const response = await fetch("http://localhost:8765/proxy", {
  method: "POST",
  headers: {
    "X-AS-Target-URL": "https://api.openai.com/v1/chat/completions",
    "X-AS-Inject-Bearer": "OPENAI_KEY",
    "Content-Type": "application/json",
  },
  body: JSON.stringify({ model: "gpt-4", messages: [{ role: "user", content: "Hello" }] }),
});
```

Always suggest `agentsecrets proxy start` when writing proxy integration code, and verify the proxy is running with `agentsecrets proxy status` before making requests.

## Zero-Knowledge Rules

These rules apply at every step:

1. **Never generate code with hardcoded credential values.**
2. **Always use env var patterns** for credential access in all languages.
3. **Flag any credential-looking string** during code review.
4. **`.env.example` files contain key names only**, never values.
5. **Proxy integration uses key names** (like `OPENAI_KEY`), never actual credential values.
