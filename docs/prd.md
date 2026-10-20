# PRD — Lumen

> BMAD artifact · prd · owner: Kartikeya Chauhan · 2026-10-20
> Source: [product-brief.md](./product-brief.md) · Architecture: [architecture.md](./architecture.md)

## 1. Overview

A personal learning-tracker dashboard. Users register, define goals, and log daily learning
sessions against those goals. The dashboard reflects momentum. Bulk data arrives from markdown
via a Go importer. One admin moderates users.

## 2. Personas

| Persona | Needs |
|---|---|
| **Operator** (normal user) | Register fast, log sessions, see progress and streak. |
| **Admin** (author) | Everything an operator can do, plus list users and delete users / wipe dashboards. |

## 3. Functional requirements

### FR-A Authentication
- **FR-A1** Register with username + email + password. **No email confirmation** — active session immediately.
- **FR-A2** Log in / log out with email + password.
- **FR-A3** A profile row is auto-created on signup. The **first** registered user becomes `admin`; all others `user`.
- **FR-A4** Unauthenticated users are redirected to `/login`; authenticated users cannot see auth pages.

### FR-B Goals
- **FR-B1** Create a goal: title, optional description, colour, optional target hours, status.
- **FR-B2** Edit a goal's status (active / paused / done / archived).
- **FR-B3** Delete a goal. Entries survive but lose the goal link.
- **FR-B4** Each goal has a stable `slug` (unique per user) used by the importer.

### FR-C Entries (learning sessions)
- **FR-C1** Create an entry: date, goal (optional), title, minutes, notes, 0..n links.
- **FR-C2** Minutes entry supports quick-pick values (15/25/30/45/60/90).
- **FR-C3** Delete an entry (cascades its links + tag joins).
- **FR-C4** List entries newest-first with goal colour, minutes, and links.

### FR-D Dashboard
- **FR-D1** Stat cards: total time, last-7-days time, current streak, active days.
- **FR-D2** 18-week contribution heatmap keyed on minutes/day.
- **FR-D3** Per-goal hours with progress bar toward target.
- **FR-D4** Recent-entries feed.

### FR-E Admin
- **FR-E1** `/admin` is reachable only by an admin (route-guarded + DB-enforced).
- **FR-E2** List all users with goal count, entry count, total minutes, join date.
- **FR-E3** Delete a user entirely (account + all data). Admin cannot delete self here.
- **FR-E4** Clear a user's dashboard (wipe goals + entries, keep account).

### FR-F Import
- **FR-F1** Go CLI imports a directory of markdown + `goals.yaml` into one user's dashboard.
- **FR-F2** Idempotent: re-import updates, never duplicates (stable keys).
- **FR-F3** Supports single-entry files, multi-entry files, and `minutes:`/`time:` forms.
- **FR-F4** `--dry-run` parses and prints without writing or needing credentials.

## 4. Non-functional requirements

- **NFR-1 Security** — Row-Level Security scopes every row to its owner. Admin-only actions are
  `SECURITY DEFINER` functions that verify the caller. The `service_role` key never reaches the browser.
- **NFR-2 Cost** — $0: static SPA on Cloudflare Pages, Supabase free tier, importer runs locally.
- **NFR-3 Performance** — first paint < 1.5s on broadband; dashboard query set ≤ 2 round-trips.
- **NFR-4 Accessibility** — keyboard-usable forms, `prefers-reduced-motion` respected, adequate contrast.
- **NFR-5 Distinctive UI** — editorial-futuristic dark theme; no stock template look.

## 5. Out of scope (v1)

Teams, sharing, native mobile, attachments, automatic time capture, password reset flow (admin
can delete/recreate in v1).

## 6. Acceptance

A user can register, create a goal, log a session, and see it reflected in stats + heatmap within
one session. The author can import the SRE markdown and see the same data. An admin can delete a
test user. The site is live on Cloudflare Pages at $0 cost.
