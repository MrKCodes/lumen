# lumen-import

A small Go CLI that reads markdown learning logs and upserts them into your Lumen Supabase
project. It is stdlib-first (only `yaml.v3` is external) on purpose, so it also works as a
readable Go learning artifact.

## Files

| File | Role |
|---|---|
| `main.go` | flags, wiring, dry-run |
| `config.go` | reads `SUPABASE_URL` + `SUPABASE_SERVICE_ROLE_KEY`, loads `.env` |
| `parse.go` | markdown + `goals.yaml` → `Goal` / `Entry` structs |
| `client.go` | thin Supabase REST + Auth-admin client (service_role) |
| `import.go` | idempotent upsert of goals, entries, links, tags |

## Prerequisites

- Go 1.22+ (`brew install go`)
- A Supabase project with `supabase/schema.sql` applied
- The account you import into must already exist — register it in the web app first
- Env set (or a repo-root `.env`, see `.env.example`):
  - `SUPABASE_URL`
  - `SUPABASE_SERVICE_ROLE_KEY` — **secret**, bypasses RLS, CLI-only

## Setup

```bash
cd importer
go mod tidy        # fetches yaml.v3, writes go.sum
```

## Use

```bash
# parse only, write nothing
go run . --dir ../content/sre-learnings --owner you@example.com --dry-run

# real import
go run . --dir ../content/sre-learnings --owner you@example.com

# build a binary
go build -o lumen-import .
./lumen-import --dir ../content/sre-learnings --owner you@example.com
```

## Idempotency

Re-running never duplicates:

- goals keyed by `(user_id, slug)`
- entries keyed by `(user_id, source_key)` where `source_key = relative/path.md#index`

Edit a log, re-run, and the matching rows update in place. Links and tags for an entry are
replaced on each import (delete-then-insert), so removing a link from the markdown removes it
from the dashboard.

## Security

The `service_role` key is a full-access secret. Keep it in your shell env or a gitignored
`.env`. Never pass it to the frontend and never commit it.

## Roadmap

- `--tracker` mode to parse `crunch-tracker.md`-style tables into entries automatically.
