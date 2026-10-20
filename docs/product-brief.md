# Product Brief — Lumen

> BMAD artifact · product-brief · owner: Kartikeya Chauhan · 2026-10-20

## Problem

Learning toward a hard goal (a Senior/Staff SRE loop, a language, a cert) happens in scattered
sessions across many days. The effort is real but invisible: there is no single place that shows
*how much* time went where, *whether* it is pointed at the right goal, and *whether* momentum is
holding. Spreadsheets and markdown trackers capture the data but do not reflect it back as signal.

## Vision

Lumen turns daily learning into visible flux. You set goals, log sessions (time + what +
links + notes), and the dashboard reflects momentum: totals per goal, a contribution heatmap, a
streak, progress toward target hours. The first real dataset is an SRE interview-prep crunch
sprint, imported from markdown.

## Target user

- **Primary:** the author — an SRE prepping for Staff-level loops, already keeping markdown logs.
- **Secondary:** anyone self-studying toward goals who wants a cool, personal, free dashboard.

## Goals / success signals

1. Logging a session takes < 20 seconds.
2. The dashboard answers "am I keeping momentum?" at a glance (streak + heatmap + per-goal hours).
3. Existing markdown logs import with zero data loss and no duplication on re-import.
4. Hosting cost: **$0**. Frontend static on Cloudflare Pages, backend on Supabase free tier.
5. A distinctive UI that does not look like generic generated output.

## Non-goals (v1)

- Teams / sharing / social features.
- Mobile native app (responsive web is enough).
- Rich text / file attachments on entries (links only).
- Automated time tracking (manual entry; markdown import for bulk).

## Constraints

- Free tier only. No always-on server (hence Supabase + static SPA, Go importer run locally).
- Auth without email confirmation (fast personal onboarding).
- One admin (the author) who can moderate users and dashboards.

## Stack decision (see architecture.md)

SvelteKit static SPA · Supabase (Postgres + Auth + RLS) · Go CLI importer · Cloudflare Pages ·
Syncline → GitHub Projects for plan tracking.
