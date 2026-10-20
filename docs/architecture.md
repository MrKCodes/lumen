# Architecture — Lumen

> BMAD artifact · architecture · owner: Kartikeya Chauhan · 2026-10-20

## 1. Shape

A static single-page app talks directly to Supabase. No application server exists or is hosted.
A local Go CLI handles bulk import. Plan tracking is pushed to GitHub Projects by Syncline.

```
            ┌──────────────────────────────┐
browser ──▶ │ SvelteKit SPA (static)       │
            │ Cloudflare Pages             │
            └───────┬──────────────────────┘
                    │ anon key + user JWT
                    ▼
            ┌──────────────────────────────┐
            │ Supabase                     │
            │  • Postgres (RLS)            │
            │  • Auth (GoTrue)             │
            │  • PostgREST auto API        │
            │  • RPC (SECURITY DEFINER)    │
            └───────▲──────────────────────┘
                    │ service_role (secret)
            ┌───────┴──────────────────────┐
            │ Go CLI importer (local)      │
            │  content/*.md → upsert       │
            └──────────────────────────────┘
```

## 2. Why this stack

| Decision | Reason | Alternative rejected |
|---|---|---|
| **No app server** | Supabase RLS + PostgREST make a server redundant; nothing to host → truly free, no cold starts. | A Go API on Fly.io — adds a host, cold starts, ops, for no v1 benefit. |
| **SvelteKit static SPA** | Smallest bundles, compiles the framework away, distinctive UI easy to craft. | React (heavier, more generic). |
| **Cloudflare Pages** | Free, unlimited bandwidth, true SPA fallback, preview deploys, header/redirect control. | GitHub Pages — SPA routing needs a 404 hack, no previews. |
| **Supabase** | Postgres + Auth + RLS + REST on a free tier; user already has an account. | Firebase (NoSQL, less SQL control). |
| **Go importer** | User wants to learn Go; bulk import is naturally a CLI; keeps the secret key off the browser. | A JS script — misses the learning goal. |

## 3. Data model

See [`supabase/schema.sql`](../supabase/schema.sql). Tables: `profiles`, `goals`, `entries`,
`entry_links`, `tags`, `entry_tags`.

```
auth.users 1─1 profiles 1─* goals 1─* entries 1─* entry_links
                                         └─* entry_tags *─1 tags
```

Stable keys for idempotent import:
- `goals (user_id, slug)` unique
- `entries (user_id, source_key)` unique — `source_key = relative/path.md#index`

## 4. Security model

- **RLS everywhere.** Every table has RLS on. Base policy: `user_id = auth.uid()` for all CRUD.
- **Admin reads/deletes** go through `public.is_admin()`, a `SECURITY DEFINER` function that reads
  `profiles.role` without re-triggering RLS (avoids policy recursion).
- **Admin mutations** (`admin_list_users`, `admin_delete_user`, `admin_clear_dashboard`) are
  `SECURITY DEFINER` RPCs that re-check `is_admin()` on entry and refuse otherwise. Admin cannot
  delete their own account via the RPC.
- **Key hygiene.** Browser uses only the public `anon` key (safe — RLS enforces access). The
  `service_role` key (bypasses RLS) lives only in the Go importer's env, never shipped.
- **Transport headers** set via Cloudflare `_headers` (XFO, nosniff, referrer, permissions policy).

## 5. Frontend structure

```
web/src/
  lib/
    supabase.js            # client (anon key)
    api.js                 # typed data calls; RLS does the scoping
    util.js                # time/heatmap/streak helpers
    stores/auth.svelte.js  # session + profile (Svelte 5 runes)
    components/            # Nav, StatCard, Heatmap, GoalBar, Logo
    styles/app.css         # design system
  routes/
    +layout.svelte         # auth init + route guard
    +page.svelte           # dashboard
    login, register        # auth (no-confirm)
    goals, log             # CRUD
    admin                  # moderation (guarded)
```

Rendering is fully client-side (`ssr=false`) because auth + data are per-user and dynamic.

## 6. Deploy

- **Frontend:** Cloudflare Pages, build command `pnpm build`, output dir `web/build`, root `web/`.
  Env: `VITE_SUPABASE_URL`, `VITE_SUPABASE_ANON_KEY`. SPA fallback via `static/_redirects`.
- **Database:** run `supabase/schema.sql` once in the Supabase SQL editor; disable email confirm.
- **Importer:** run locally with `SUPABASE_URL` + `SUPABASE_SERVICE_ROLE_KEY` in env.

## 7. Risks / tradeoffs

- **Client-trust surface.** All business logic that matters for security must live in RLS/RPC, not
  the SPA — the browser is untrusted. Mitigated by schema design above.
- **No email confirm** means anyone can register with any email. Acceptable: personal app, admin
  can delete abusive accounts. Revisit if opened publicly.
- **Supabase free tier** pauses after inactivity; fine for personal use.
