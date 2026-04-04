---
name: using-git-worktrees
description: Use when feature work must stay isolated from the current checkout, parallel branches are needed without git switch or checkout, an implementation plan should run in a separate working tree, the working tree is dirty and blocks branch changes, or the user asks for a git worktree, isolated workspace, or second checkout of the same repo.
---

# Using Git Worktrees

## Overview

Git worktrees add another working directory that shares the same `.git` object database, so multiple branches can be edited at once without switching the main tree.

**Core principle:** Pick the worktree root in a fixed priority order, then verify project-local roots are ignored before `git worktree add`.

**Announce at start:** State that the using-git-worktrees skill is being used to set up an isolated workspace.

## Triggers (symptoms)

Prefer this skill when any of the following appear:

- Need to keep the current directory on its branch while coding another branch
- `git switch` / `git checkout` refused because of local changes
- Plan or subagent work should not touch the open workspace’s uncommitted files
- User names “worktree”, “second clone”, or “parallel branch folder”

## Directory selection

Apply this order; do not invent a new convention when a higher-priority signal exists.

### 1. Existing project directories

```bash
ls -d .worktrees 2>/dev/null     # preferred (hidden)
ls -d worktrees 2>/dev/null      # alternative
```

**If found:** Use that directory. If both exist, `.worktrees` wins.

### 2. Project docs (e.g. CLAUDE.md)

```bash
grep -i "worktree.*director" CLAUDE.md 2>/dev/null
```

**If a path is documented:** Use it without asking.

### 3. Ask the user

If no directory exists and no documented preference:

```
No worktree directory found. Where should worktrees live?

1. .worktrees/ (project-local, hidden)
2. ~/.config/superpowers/worktrees/<project-name>/ (outside the repo)

Which should be used?
```

## Safety verification

### Project-local roots (`.worktrees` or `worktrees`)

**Before `git worktree add`, confirm the directory is ignored** (covers repo, global, and system exclude rules):

```bash
git check-ignore -q .worktrees 2>/dev/null || git check-ignore -q worktrees 2>/dev/null
```

**If the chosen root is not ignored:**

1. Add the directory (or a parent pattern) to `.gitignore` (or the appropriate exclude file).
2. Commit that change when the project expects commits for config fixes.
3. Then create the worktree.

**Why:** Unignored worktree paths can show up as thousands of tracked or untracked files and break reviews.

### Global path (`~/.config/superpowers/worktrees/...`)

No ignore check inside the repo; the tree lives outside the project.

## Creation steps

### 1. Resolve project name and paths

```bash
project=$(basename "$(git rev-parse --show-toplevel)")
```

Build the worktree path from the chosen location:

- **Project-local:** `<chosen-root>/<branch-slug>` (e.g. `.worktrees/feature-auth`).
- **Global:** `$HOME/.config/superpowers/worktrees/$project/<branch-slug>` (expand `$HOME`; do not rely on `~` inside a variable for `cd`).

### 2. Create worktree and branch

```bash
git worktree add "<full-path-to-worktree>" -b "<branch-name>"
cd "<full-path-to-worktree>" || exit 1
```

Use a branch name that matches team conventions (often `feature/...` or ticket id).

### 3. Run project setup (auto-detect)

```bash
if [ -f package.json ]; then npm install; fi
if [ -f Cargo.toml ]; then cargo build; fi
if [ -f requirements.txt ]; then pip install -r requirements.txt; fi
if [ -f pyproject.toml ]; then poetry install; fi
if [ -f go.mod ]; then go mod download; fi
```

Skip steps when the marker file is absent. Prefer the project’s documented install command when it differs (Makefile, pnpm, uv, etc.).

### 4. Verify a clean baseline

Run the project’s usual test or check command (examples: `npm test`, `cargo test`, `pytest`, `go test ./...`).

**If tests fail:** Report output, distinguish repo baseline vs. environment, and ask whether to proceed or fix first.

**If tests pass:** State readiness.

### 5. Report to the user

```
Worktree ready at <absolute-path>
Tests: <command> — <summary>
Ready to implement <short feature or plan label>
```

## Quick reference

| Situation | Action |
|-----------|--------|
| `.worktrees/` exists | Use it; verify ignored |
| `worktrees/` exists | Use it; verify ignored |
| Both exist | Prefer `.worktrees/` |
| Neither exists | Read project docs → ask user |
| Project-local root not ignored | Fix ignore (and commit if required) → then add worktree |
| Baseline tests fail | Report and get explicit go/no-go |
| No `package.json` / `Cargo.toml` / etc. | Skip matching install step |

## Common mistakes

| Mistake | Why it hurts | Fix |
|---------|----------------|-----|
| Skipping `git check-ignore` | Accidental tracking or giant `git status` | Always verify project-local roots |
| Guessing the worktree root | Breaks team convention | Follow priority: existing → docs → ask |
| Continuing with failing baseline tests | New defects vs. old defects blur | Stop and confirm with the user |
| Hardcoding one stack’s commands | Wrong for Ruby, Java, monorepos, etc. | Detect markers; follow project docs |

## Example narrative (reference)

```
I'm using the using-git-worktrees skill to set up an isolated workspace.

[.worktrees exists → git check-ignore confirms it is ignored]
[git worktree add .worktrees/feature-auth -b feature/auth]
[npm install && npm test → 47 passing]

Worktree ready at <repo>/.worktrees/feature-auth
Tests: npm test — 47 passed, 0 failed
Ready to implement auth feature
```

## Red flags — stop and correct

- Creating a project-local worktree under a path that is not ignored
- Omitting baseline tests or checks without explicit user approval to skip
- Picking a worktree root when two conventions exist and priority was skipped
- Using `~` inside a quoted variable for `cd` instead of `$HOME` or an absolute path

## Integration

**Often chained from:**

- **brainstorming** (when design is approved and implementation follows)
- **subagent-driven-development** — before tasks that need a clean tree
- **executing-plans** — before executing plan steps in isolation
- Any workflow that needs a second checkout without disturbing the first

**Pairs with:**

- **finishing-a-development-branch** — for merge/PR/cleanup after the worktree branch is done
