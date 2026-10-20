# Lumen — setup & deploy

End-to-end, from zero to a live dashboard at $0. ~20 minutes.

---

## 0. Prerequisites

- Node 18+ and `pnpm` (`npm i -g pnpm`)
- Go 1.22+ (`brew install go`) — only for importing markdown
- A free [Supabase](https://supabase.com) account
- A free [Cloudflare](https://dash.cloudflare.com) account
- `gh` CLI authed (for repo + Syncline board)

---

## 1. Supabase

1. **Create a project** at app.supabase.com. Pick a region near you. Save the DB password.
2. **Run the schema.** Open **SQL Editor** → paste all of [`supabase/schema.sql`](../supabase/schema.sql) → **Run**.
   It creates tables, RLS, triggers, and admin functions. Safe to re-run.
3. **Turn off email confirmation.** Authentication → **Providers** → **Email** → disable
   *"Confirm email"* → Save. (This is what lets registration log users in instantly.)
4. **Grab your keys.** Project Settings → **API**:
   - `Project URL`  → `SUPABASE_URL` / `VITE_SUPABASE_URL`
   - `anon public`  → `VITE_SUPABASE_ANON_KEY`
   - `service_role` → `SUPABASE_SERVICE_ROLE_KEY` (**secret** — importer only)

> **First account = admin.** The very first user to register becomes the admin automatically
> (see `handle_new_user()` in the schema). Register yourself first.

---

## 2. Frontend (local dev)

```bash
cd web
cp .env.example .env          # fill VITE_SUPABASE_URL + VITE_SUPABASE_ANON_KEY
pnpm install
pnpm dev                      # http://localhost:5173
```

Register your account (becomes admin), create a goal, log a session.

---

## 3. Import your SRE markdown (optional)

```bash
cd importer
cp ../.env.example ../.env     # fill SUPABASE_URL + SUPABASE_SERVICE_ROLE_KEY
go mod tidy
go run . --dir ../content/sre-learnings --owner you@example.com --dry-run   # preview
go run . --dir ../content/sre-learnings --owner you@example.com             # write
```

`--owner` is the email you registered with. Re-running is safe — it updates, never duplicates.
Add your own `content/<dataset>/` dirs with their own `goals.yaml`; the schema is in
[`content/README.md`](../content/README.md).

---

## 4. Deploy frontend → Cloudflare Pages

1. Push this repo to GitHub (see §5).
2. Cloudflare dash → **Workers & Pages** → **Create** → **Pages** → **Connect to Git** → pick the repo.
3. Build settings:
   - **Framework preset:** SvelteKit (or none)
   - **Build command:** `pnpm build`
   - **Build output directory:** `build`
   - **Root directory:** `web`
4. **Environment variables:** add `VITE_SUPABASE_URL` and `VITE_SUPABASE_ANON_KEY`.
5. **Save and Deploy.** Every push to `main` redeploys; PRs get preview URLs.

> Add your Pages domain to Supabase **Authentication → URL Configuration → Site URL / Redirect URLs**.

CLI alternative:
```bash
cd web && pnpm build
npx wrangler pages deploy build --project-name lumen
```

---

## 5. GitHub repo + project board (Syncline)

```bash
# from the lumen/ root
gh repo create lumen --public --source=. --remote=origin --push

# Syncline board (see syncline/README.md for the detailed flow)
gh auth refresh -s project       # ensure 'project' scope
python3 syncline/bin/syncline_build_manifest.py --config syncline/configs/lumen.config.json
python3 syncline/bin/syncline_import.py --config syncline/configs/lumen.config.json --dry-run
python3 syncline/bin/syncline_import.py --config syncline/configs/lumen.config.json
```

---

## Troubleshooting

| Symptom | Fix |
|---|---|
| "Missing VITE_SUPABASE_*" in console | Copy `web/.env.example` → `web/.env`, restart `pnpm dev`. |
| Registration says "email not confirmed" | You didn't disable Confirm email in step 1.3. |
| Admin tab missing | You weren't the first user. Promote in SQL: `update profiles set role='admin' where username='you';` |
| Importer "no user with email" | Register that email in the app first. |
| 404 on refresh after deploy | Ensure `static/_redirects` shipped; Pages needs the SPA fallback. |
