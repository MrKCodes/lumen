# Lumen content — markdown import schema

The Go importer (`../importer`) reads this directory and upserts everything into Supabase.
You write plain markdown; the dashboard renders it. Edit a file, re-run the importer, and the
dashboard reflects the change — nothing is duplicated, because every goal and entry has a
**stable key**.

```
content/
  sre-learnings/
    goals.yaml          # the goals this dataset tracks
    2026-10-15.md       # one file per day (recommended), OR
    linux-internals.md  # one file per topic — both work; see "Multiple entries per file"
```

The importer is run per dataset directory:

```bash
cd importer
go run . --dir ../content/sre-learnings --owner you@example.com
```

`--owner` is the email of the Lumen account the data belongs to. The entries land in that
user's dashboard.

---

## 1. `goals.yaml`

A list of goals. `slug` is the stable key — entries reference it. Change the title freely;
keep the slug.

```yaml
goals:
  - slug: linux-internals
    title: Linux Systems Internals
    description: The 9 surfaces the MongoDB loop exposed. Priority #1.
    color: "#36e2b4"
    target_hours: 40
    status: active          # active | paused | done | archived

  - slug: dsa-go
    title: DSA (Python → Go)
    description: NeetCode order, re-implemented in Go.
    color: "#7c5cff"
    target_hours: 60
    status: active

  - slug: system-design
    title: System Design
    color: "#ff6b9d"
    status: active
```

Only `slug` and `title` are required. Defaults: `color` random-from-palette, `status: active`.

---

## 2. Daily log file — `YYYY-MM-DD.md`

YAML frontmatter describes **one entry**; the markdown body is the entry's notes.

```markdown
---
goal: linux-internals          # must match a slug in goals.yaml
title: Log parsing with awk + combined one-liner
date: 2026-10-15               # optional — inferred from filename if omitted
minutes: 55
tags: [awk, pipes, drills]
links:
  - url: https://man7.org/linux/man-pages/man1/awk.1p.html
    label: awk POSIX man page
  - https://github.com/MrKCodes/sre-learnings/blob/main/linux-internals.md
---

Worked §1. Parsed a 2M-line access log for p99 latency with a single awk pass.
Got the combined pipe one-liner down cold. **Weak spot:** narrating the field-split
out loud — practice tomorrow.
```

### Field reference

| Field | Required | Notes |
|---|---|---|
| `goal` | ✅ | slug from `goals.yaml`. Unknown slug → importer warns and skips. |
| `title` | ✅ | short entry title. |
| `date` | — | `YYYY-MM-DD`. If omitted, taken from the filename. |
| `minutes` | — | integer, default `0`. Use `time: "1h30m"` instead if you prefer (parsed to minutes). |
| `tags` | — | list of strings; created per user on demand. |
| `links` | — | list; each item is either a bare URL string or `{url, label}`. |

The entry's **stable key** (`source_key`) is derived as `relative/path.md#<index>`. Re-importing
updates the matching entry instead of creating a new one.

---

## 3. Multiple entries per file

For a topic file that accumulates many sessions, separate entries with a `---` horizontal rule
and give each its own frontmatter-style block, OR use a top-level `entries:` list:

```markdown
---
goal: linux-internals
title: Page cache under memory pressure
entries:
  - date: 2026-10-18
    minutes: 50
    note: "8GB file on a 1GB box — traced the OOM path, reclaim vs swap."
    tags: [page-cache, oom]
  - date: 2026-10-19
    minutes: 30
    note: "Re-ran the drill, narrated the vmstat columns aloud."
---
```

Each list item becomes its own entry under the same goal/title, keyed by index.

---

## 4. Importing your existing SRE repo

Your `sre-learnings` repo uses tracker tables, not per-day frontmatter. Two paths:

- **Manual (now):** drop daily `YYYY-MM-DD.md` files here as you study — fastest, cleanest data.
- **Automated (later):** the importer has a `--tracker` mode stub that can parse the
  `crunch-tracker.md` table rows into entries. See `importer/README.md`.

Start manual. The schema above is the source of truth.
