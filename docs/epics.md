# Epics & Stories — Lumen

> BMAD artifact · epics · owner: Kartikeya Chauhan · 2026-10-20
> This file is the source Syncline turns into a GitHub Projects board.
> Status legend: `todo` · `in-progress` · `done`

---

## Epic 1: Foundation & Infrastructure
**Goal:** Repo, schema, and hosting skeleton exist and build cleanly.
**Status:** done

- **Story 1.1** — Scaffold repo structure (web, importer, supabase, content, docs). `done`
- **Story 1.2** — Design Postgres schema: profiles, goals, entries, links, tags. `done`
- **Story 1.3** — Write RLS policies and admin SECURITY DEFINER functions. `done`
- **Story 1.4** — Auto-create profile on signup; first user becomes admin. `done`
- **Story 1.5** — Cloudflare Pages config: static adapter, `_headers`, `_redirects`. `done`

## Epic 2: Authentication
**Goal:** Users register and sign in with no email confirmation.
**Status:** done

- **Story 2.1** — Supabase client + auth store (session + profile, Svelte 5 runes). `done`
- **Story 2.2** — Register page (username/email/password), instant session. `done`
- **Story 2.3** — Login page + sign-out. `done`
- **Story 2.4** — Route guard: redirect unauthenticated, hide auth pages when signed in. `done`

## Epic 3: Goals
**Goal:** Users define and manage the goals they point learning at.
**Status:** done

- **Story 3.1** — Create goal (title, colour, target hours, status). `done`
- **Story 3.2** — List goals with per-goal hours and progress bar. `done`
- **Story 3.3** — Cycle goal status; delete goal. `done`

## Epic 4: Entries & Logging
**Goal:** Logging a learning session takes seconds.
**Status:** done

- **Story 4.1** — Log form: date, goal, title, minutes (quick-pick), notes, links. `done`
- **Story 4.2** — Timeline feed, newest-first, with delete. `done`
- **Story 4.3** — Link management (add/remove multiple links per entry). `done`

## Epic 5: Dashboard & Visualisation
**Goal:** Momentum is visible at a glance.
**Status:** done

- **Story 5.1** — Stat cards: total, last-7-days, streak, active days. `done`
- **Story 5.2** — 18-week contribution heatmap. `done`
- **Story 5.3** — Goals summary + recent feed on the dashboard. `done`

## Epic 6: Admin Panel
**Goal:** The admin moderates users and dashboards.
**Status:** done

- **Story 6.1** — Admin-only route guard + DB enforcement. `done`
- **Story 6.2** — List all users with aggregate stats. `done`
- **Story 6.3** — Delete a user (account + data); block self-delete. `done`
- **Story 6.4** — Clear a user's dashboard (keep account). `done`

## Epic 7: Markdown Import (Go CLI)
**Goal:** Existing SRE markdown logs populate the dashboard, idempotently.
**Status:** done

- **Story 7.1** — Define markdown schema (`goals.yaml` + daily/topic files). `done`
- **Story 7.2** — Go parser for goals, single-entry, multi-entry, `time:` durations. `done`
- **Story 7.3** — Supabase REST client (service_role) + email→UUID resolution. `done`
- **Story 7.4** — Idempotent upsert of goals, entries, links, tags; `--dry-run`. `done`

## Epic 8: Deployment & Launch
**Goal:** The app is live, documented, and tracked.
**Status:** in-progress

- **Story 8.1** — Create GitHub repo and push. `done`
- **Story 8.2** — Create GitHub Project board via Syncline. `done`
- **Story 8.3** — Provision Supabase project + run schema + disable email confirm. `todo`
- **Story 8.4** — Connect Cloudflare Pages and deploy. `todo`
- **Story 8.5** — SETUP.md end-to-end walkthrough. `done`

---

## Future (backlog, not in v1)

- CSV / Notion export.
- `--tracker` import mode for existing `crunch-tracker.md` tables.
- Weekly email/report of flux.
- Public read-only share link per goal.
- Password reset flow.
