---
name: agentsecrets-call
description: "Make authenticated API calls with zero-knowledge secret injection. AgentSecrets injects credentials at the HTTP transport layer so the agent never sees plaintext values."
---

# AgentSecrets Call

## Step 1 — Pre-flight Check (MANDATORY)

Run these two commands before every API call:

```bash
agentsecrets status
```

```bash
agentsecrets secrets list
```

Verify the required secret key exists in the list. If the key is missing, invoke `agentsecrets-ops` to help the user set it up before proceeding. Do not attempt any call until the key is confirmed present.

## Step 2 — One-Shot Calls

Base pattern:

```
agentsecrets call --url URL --method METHOD --AUTH_STYLE KEY_NAME
```

Default method is GET when `--method` is omitted.

### Auth Style Reference

| Style | Flag | Example |
|---|---|---|
| Bearer token | `--bearer KEY` | `agentsecrets call --url https://api.stripe.com/v1/balance --bearer STRIPE_KEY` |
| Custom header | `--header Name=KEY` | `agentsecrets call --url https://api.sendgrid.com/v3/mail/send --method POST --header X-Api-Key=SENDGRID_KEY --body '{"personalizations":[...]}'` |
| Query parameter | `--query param=KEY` | `agentsecrets call --url "https://maps.googleapis.com/maps/api/geocode/json?address=NYC" --query key=GOOGLE_MAPS_KEY` |
| Basic auth | `--basic KEY` | `agentsecrets call --url https://jira.example.com/rest/api/2/issue --basic JIRA_CREDS` |
| JSON body field | `--body-field path=KEY` | `agentsecrets call --url https://auth.example.com/token --method POST --body-field client_secret=OAUTH_SECRET` |
| Form field | `--form-field field=KEY` | `agentsecrets call --url https://api.example.com/auth --method POST --form-field api_key=API_KEY` |

### Multi-Method Examples

**POST with body:**

```bash
agentsecrets call --url https://api.stripe.com/v1/charges --method POST --bearer STRIPE_KEY --body '{"amount":1000}'
```

**PUT:**

```bash
agentsecrets call --url https://api.example.com/users/1 --method PUT --bearer AUTH_TOKEN --body '{"field":"value"}'
```

**DELETE:**

```bash
agentsecrets call --url https://api.example.com/users/1 --method DELETE --bearer AUTH_TOKEN
```

**Multiple auth styles combined:**

```bash
agentsecrets call --url https://api.example.com/data --bearer AUTH_TOKEN --header X-Org-ID=ORG_SECRET
```

## Step 3 — Proxy Mode

Use proxy mode for multi-call workflows, agent framework integration, or when making 3+ API calls in a session.

### Proxy Lifecycle

```bash
agentsecrets proxy start [--port 8765]
agentsecrets proxy status
agentsecrets proxy stop
```

### HTTP Proxy Pattern

Send requests to the local proxy with injection headers:

```
POST http://localhost:8765/proxy
X-AS-Target-URL: https://api.stripe.com/v1/balance
X-AS-Inject-Bearer: STRIPE_KEY
```

### Injection Headers

| Header | Purpose |
|---|---|
| `X-AS-Inject-Bearer: KEY` | Inject as Bearer token |
| `X-AS-Inject-Header: Name=KEY` | Inject as custom header |
| `X-AS-Inject-Query: param=KEY` | Inject as query parameter |
| `X-AS-Inject-Basic: KEY` | Inject as Basic auth |

Suggest proxy mode when:
- Making 3 or more API calls in a workflow
- Integrating with agent frameworks that issue their own HTTP requests
- Building multi-step workflows that chain API responses

## Step 4 — Audit Logs

Review API call history:

```bash
agentsecrets proxy logs
agentsecrets proxy logs --last 20
agentsecrets proxy logs --secret STRIPE_KEY
```

Logs show: timestamp, method, URL, key names used, status code, and response time. The log struct has no value field — it is structurally impossible for credential values to appear in logs.

## Zero-Knowledge Rules

These rules apply at every step:

1. **Always use `agentsecrets call` or proxy mode.** Never use curl, fetch, or direct HTTP with credentials.
2. **Never display secret values** in commands, output, or conversation.
3. **Check key existence before calling.** The pre-flight check in Step 1 is mandatory.
4. **Use proxy mode for multi-call workflows** to avoid repeating auth setup.
5. **All examples use key names only** (like `STRIPE_KEY`), never actual credential values.
