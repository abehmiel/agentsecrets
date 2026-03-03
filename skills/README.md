# AgentSecrets AI Skills

This directory contains skills/instructions for different AI platforms to use AgentSecrets effectively.

## What Are Skills?

Skills teach AI assistants how to:
- Use AgentSecrets commands properly
- Never display secret values
- Help users with secrets management workflows
- Write code that uses secrets securely

## Available Skills

### Claude Code

**Directory**: `claude/`

A router + sub-skills suite for Claude Code custom commands. Auto-detects credential needs (`.env` files, hardcoded keys, auth errors) and dispatches to focused sub-skills.

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

### OpenClaw

**File**: `../integrations/openclaw/SKILL.md`

Native skill + exec provider for OpenClaw agents. Full credential lifecycle management with step-by-step operational guide.

### ChatGPT (Custom Instructions)

**Status**: Planned

Custom instructions for ChatGPT to use AgentSecrets in Code Interpreter or when helping with deployments.

### GitHub Copilot

**Status**: Planned

Instructions for integrating AgentSecrets into GitHub Copilot workflows.

## Using These Skills

### For Claude Code

1. Copy the skill files to your Claude Code commands directory:
   - **Global** (all projects): `cp skills/claude/*.md ~/.claude/commands/`
   - **Per-project**: `cp skills/claude/*.md .claude/commands/`
2. The router skill auto-detects credential needs and dispatches to sub-skills
3. You can also invoke explicitly with `/agentsecrets` in Claude Code

### For OpenClaw

1. Install via ClawHub: `claw install SteppaCodes/agentsecrets`
2. Or copy `integrations/openclaw/SKILL.md` to your OpenClaw skills directory

### For Other AI Platforms

Each platform has different integration methods. Check the specific skill file for instructions.

## Creating New Skills

Want to create a skill for a different AI platform? 

1. Create a new directory: `skills/[platform-name]/`
2. Add your skill instructions
3. Submit a PR

We welcome skills for:
- Cursor
- Tabnine
- Amazon Q
- Codeium
- Any other AI coding assistant

## Principles

All skills should follow these principles:

1. **Zero-Knowledge**: AI never sees secret values
2. **Security First**: Always reference secrets by key
3. **Helpful**: Guide users through workflows
4. **Clear**: Explain what commands do
5. **Safe**: Never log or display secrets

## Testing Skills

Test your skill by:
1. Using it with the target AI platform
2. Asking for help with secrets management
3. Verifying the AI never displays secret values
4. Checking it uses AgentSecrets commands correctly

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md) for contribution guidelines.

For skill-specific questions, open an issue tagged with `ai-integration`.

---

**Goal**: Make AgentSecrets the default way ALL AI assistants handle secrets.