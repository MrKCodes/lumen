<div align="center">

# ◆ LUMEN

**A learning-tracker dashboard. Point your effort at a goal and watch it light up.**

Log what you learn each day — time spent, links, notes — against the goals you set.
First dataset: an SRE interview-prep crunch sprint.

`SvelteKit` · `Supabase` · `Go importer` · `Cloudflare Pages`

</div>

---

## Why "Lumen"

A lumen is a unit of luminous flux — how much light a source throws. Here it is how much
learning you throw at a goal. Every logged hour adds light.

## Architecture

```
┌──────────────────┐      ┌──────────────────────┐
│  SvelteKit SPA   │ ───▶ │  Supabase            │
│  (Cloudflare     │ auth │  Postgres + Auth +   │
│   Pages, static) │ RLS  │  RLS + REST          │
└──────────────────┘      └──────────▲───────────┘
                                     │ service_role (upsert)
                          ┌──────────┴───────────┐
                          │  Go CLI importer     │
                          │  reads content/*.md  │
                          └──────────────────────┘
```

- **Frontend** — SvelteKit compiled to a static SPA. No server. Hosted free on Cloudflare Pages.
- **Backend** — Supabase. Postgres holds the data; Supabase Auth handles login (email + password,
  **no email confirmation**); Row-Level Security scopes every row to its owner; an admin role can
  delete users and their dashboards.
- **Importer** — a Go CLI that parses markdown learning logs (see [`content/`](./content)) and
  upserts them into Supabase using the secret `service_role` key. Runs on your machine, not hosted.

## Layout

| Path | What |
|---|---|
| `web/` | SvelteKit SPA (the dashboard + admin panel) |
| `importer/` | Go CLI that imports markdown learning logs → Supabase |
| `supabase/` | `schema.sql` — tables, RLS policies, triggers, admin functions |
| `content/` | Markdown learning logs + `goals.yaml` (the import source) |
| `docs/` | Product brief, PRD, architecture, epics & stories (BMAD) |
| `syncline/` | Config that syncs the BMAD plan to the GitHub Projects board |

## Quick start

See [`docs/SETUP.md`](./docs/SETUP.md) for the full walkthrough. Short version:

1. **Supabase** — create a free project, run `supabase/schema.sql` in the SQL editor,
   turn **off** "Confirm email" under Authentication → Providers → Email.
2. **Frontend** — `cd web && cp .env.example .env` (fill in URL + anon key), `pnpm i`, `pnpm dev`.
3. **Import your data** — `cd importer && go run . --dir ../content/sre-learnings`.
4. **Deploy** — push to GitHub; connect the repo to Cloudflare Pages (build dir `web`).

## Security note

The `service_role` key bypasses RLS and must never reach the browser. It lives only in the Go
importer's environment. The frontend uses the `anon` key, which is safe to ship because RLS
enforces per-user access in the database.
