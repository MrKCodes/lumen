# Syncline — Lumen project board sync

This mirrors Lumen's BMAD plan ([`docs/epics.md`](../docs/epics.md) + [`plan/sprint-status.yaml`](./plan/sprint-status.yaml))
to a GitHub Projects v2 board: epics and stories become issues, with epic→story sub-issue links,
a **Status** field, and a **Work Type** (Epic/Story) field.

Upstream tooling: https://github.com/MrKCodes/Syncline (vendored under `bin/`, see `UPSTREAM_README.md`).

## Prerequisites

```bash
gh auth status            # must be logged in
gh auth refresh -s project  # ensure the 'project' scope
```

## Files

| Path | Role |
|---|---|
| `configs/lumen.config.json` | repo/project/status-map + plan paths (relative to this file) |
| `plan/sprint-status.yaml` | the plan the sync reads — one line per epic/story at 2-space indent |
| `bin/*.py` | vendored Syncline scripts (stdlib only) |
| `configs/state/` | manifest + import-state + sync-cache (gitignored, local only) |

## Workflow

```bash
cd syncline

# 1. Build the manifest from epics.md + sprint-status.yaml
python3 bin/syncline_build_manifest.py --config configs/lumen.config.json

# 2. One-time import to the board (dry-run first)
python3 bin/syncline_import.py --config configs/lumen.config.json --dry-run
python3 bin/syncline_import.py --config configs/lumen.config.json

# 3. Keep it in sync after editing sprint-status.yaml
python3 bin/syncline_sync.py --config configs/lumen.config.json
```

## Keeping the board current

When a story's status changes, edit `plan/sprint-status.yaml`, then run the sync (step 3).
`docs/epics.md` and `plan/sprint-status.yaml` are the two human-edited files; keep them aligned.
To automate, wire the Stop hook described in `docs/HOOK.md`.
