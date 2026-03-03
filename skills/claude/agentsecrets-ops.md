---
name: agentsecrets-ops
description: "Secrets lifecycle management: list, set, diff, push, pull, delete, and environment switching. Invoked by the router when the user needs to manage secrets."
---

# AgentSecrets Ops

## Step 1 — Status Check (MANDATORY)

Run this before every operation:

```bash
agentsecrets status
```

Verify all three conditions:

1. **User is logged in**
2. **Workspace is active**
3. **Project is active**

If any condition fails, invoke `agentsecrets-setup` to resolve it before continuing.

## Step 2 — List Secrets

```bash
agentsecrets secrets list
```

This returns **key names only**, never values. Use it to check what exists before performing any other operation.

## Step 3 — Add or Update Secrets

**CRITICAL: Never accept credential values in conversation.**

When a secret is missing or needs to be added, tell the user:

> I need `KEY_NAME` to proceed. Please run in your terminal:
> `agentsecrets secrets set KEY_NAME=your_value`
> Let me know when done.

After the user confirms, verify the key was stored:

```bash
agentsecrets secrets list
```

Confirm the key now appears in the list before proceeding.

## Step 4 — Drift Detection

```bash
agentsecrets secrets diff
```

This shows what differs between local (keychain) and cloud:

- **Added** — keys only in local
- **Removed** — keys only in cloud
- **Changed** — keys with different values (names only, not values)
- **Unchanged** — keys that match

Use diff output to inform the user before any sync operation.

## Step 5 — Sync Operations

```bash
agentsecrets secrets pull   # cloud -> local
agentsecrets secrets push   # local -> cloud
```

**ALWAYS run `agentsecrets secrets diff` first** and show the user what will change before pushing or pulling. Do not push or pull without the user seeing and approving the diff.

## Step 6 — Delete Secrets

```bash
agentsecrets secrets delete KEY_NAME
```

Confirm with the user before running the delete command. Deletion removes the key from both cloud and local storage.

## Step 7 — Environment Switching

To work across environments (dev, staging, production), switch workspace and project, then pull:

```bash
agentsecrets workspace switch "production"
agentsecrets project use my-api
agentsecrets secrets pull
```

Always run `agentsecrets secrets list` after switching to confirm the correct secrets are loaded.

## Step 8 — Key Naming Conventions

Use these standard names when suggesting or identifying keys:

| Service | Key Name |
|---|---|
| Stripe (live) | `STRIPE_KEY` or `STRIPE_LIVE_KEY` |
| Stripe (test) | `STRIPE_TEST_KEY` |
| OpenAI | `OPENAI_KEY` |
| GitHub | `GITHUB_TOKEN` |
| Google Maps | `GOOGLE_MAPS_KEY` |
| AWS | `AWS_ACCESS_KEY` and `AWS_SECRET_KEY` |
| SendGrid | `SENDGRID_KEY` |
| Twilio | `TWILIO_SID` and `TWILIO_TOKEN` |
| Generic | `SERVICENAME_KEY` (uppercase, underscores) |

Convention: uppercase letters, underscores as separators, suffix with `_KEY`, `_TOKEN`, or `_SECRET` as appropriate.

## Zero-Knowledge Rules

These rules apply at every step:

1. **Never display secret values.** Do not print, log, or echo plaintext secrets.
2. **Never ask the user to paste keys into chat.** Do not request raw API keys, tokens, or credentials in conversation.
3. **Direct the user to their terminal** to run `agentsecrets secrets set KEY=value` for any credential input.
4. **Suggest the user delete messages** that accidentally contain raw keys.
