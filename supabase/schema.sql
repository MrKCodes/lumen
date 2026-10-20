-- ╔══════════════════════════════════════════════════════════════════╗
-- ║  LUMEN — Supabase schema, RLS, triggers, admin functions          ║
-- ║  Run top-to-bottom in the Supabase SQL editor. Idempotent-ish:    ║
-- ║  safe to re-run (uses IF NOT EXISTS / CREATE OR REPLACE / DROP).   ║
-- ╚══════════════════════════════════════════════════════════════════╝

-- ── Extensions ──────────────────────────────────────────────────────
create extension if not exists "pgcrypto";   -- gen_random_uuid()

-- ════════════════════════════════════════════════════════════════════
-- 1. TABLES
-- ════════════════════════════════════════════════════════════════════

-- One profile per auth user. role drives admin powers.
create table if not exists public.profiles (
  id          uuid primary key references auth.users (id) on delete cascade,
  username    text unique,
  display_name text,
  role        text not null default 'user' check (role in ('user', 'admin')),
  created_at  timestamptz not null default now()
);

-- A goal the user is pointing learning at.
create table if not exists public.goals (
  id           uuid primary key default gen_random_uuid(),
  user_id      uuid not null references public.profiles (id) on delete cascade,
  slug         text not null,                 -- stable key for markdown import
  title        text not null,
  description  text,
  color        text not null default '#7c5cff',
  target_hours numeric(7,1),                  -- optional goal in hours
  status       text not null default 'active'
                 check (status in ('active', 'paused', 'done', 'archived')),
  created_at   timestamptz not null default now(),
  unique (user_id, slug)
);

-- A single learning entry: time spent on something, on a day.
create table if not exists public.entries (
  id         uuid primary key default gen_random_uuid(),
  user_id    uuid not null references public.profiles (id) on delete cascade,
  goal_id    uuid references public.goals (id) on delete set null,
  entry_date date not null,
  title      text not null,
  notes      text,
  minutes    integer not null default 0 check (minutes >= 0),
  -- stable key so re-importing the same markdown updates instead of duplicating
  source_key text,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  unique (user_id, source_key)
);

-- Links attached to an entry (docs, repos, articles).
create table if not exists public.entry_links (
  id       uuid primary key default gen_random_uuid(),
  entry_id uuid not null references public.entries (id) on delete cascade,
  user_id  uuid not null references public.profiles (id) on delete cascade,
  url      text not null,
  label    text
);

-- Tags, scoped per user, and the join to entries.
create table if not exists public.tags (
  id      uuid primary key default gen_random_uuid(),
  user_id uuid not null references public.profiles (id) on delete cascade,
  name    text not null,
  unique (user_id, name)
);

create table if not exists public.entry_tags (
  entry_id uuid not null references public.entries (id) on delete cascade,
  tag_id   uuid not null references public.tags (id) on delete cascade,
  primary key (entry_id, tag_id)
);

create index if not exists entries_user_date_idx on public.entries (user_id, entry_date);
create index if not exists entries_goal_idx      on public.entries (goal_id);
create index if not exists goals_user_idx         on public.goals (user_id);

-- ════════════════════════════════════════════════════════════════════
-- 2. HELPERS
-- ════════════════════════════════════════════════════════════════════

-- Is the CURRENT request made by an admin?
-- SECURITY DEFINER + search_path pinned: reads profiles WITHOUT triggering
-- RLS, so it is safe to call from inside a profiles policy (no recursion).
create or replace function public.is_admin()
returns boolean
language sql
security definer
set search_path = public
stable
as $$
  select exists (
    select 1 from public.profiles
    where id = auth.uid() and role = 'admin'
  );
$$;

-- keep entries.updated_at honest
create or replace function public.touch_updated_at()
returns trigger language plpgsql as $$
begin
  new.updated_at := now();
  return new;
end;
$$;

drop trigger if exists entries_touch on public.entries;
create trigger entries_touch
  before update on public.entries
  for each row execute function public.touch_updated_at();

-- ════════════════════════════════════════════════════════════════════
-- 3. AUTO-PROFILE ON SIGNUP
--    First user to ever register becomes admin; everyone else = user.
-- ════════════════════════════════════════════════════════════════════

create or replace function public.handle_new_user()
returns trigger
language plpgsql
security definer
set search_path = public
as $$
declare
  is_first boolean;
begin
  select count(*) = 0 into is_first from public.profiles;

  insert into public.profiles (id, username, display_name, role)
  values (
    new.id,
    coalesce(new.raw_user_meta_data ->> 'username', split_part(new.email, '@', 1)),
    coalesce(new.raw_user_meta_data ->> 'display_name', split_part(new.email, '@', 1)),
    case when is_first then 'admin' else 'user' end
  );
  return new;
end;
$$;

drop trigger if exists on_auth_user_created on auth.users;
create trigger on_auth_user_created
  after insert on auth.users
  for each row execute function public.handle_new_user();

-- ════════════════════════════════════════════════════════════════════
-- 4. ROW-LEVEL SECURITY
--    Rule: a user touches only their own rows. Admins may also DELETE
--    any row (for moderation). Admins do NOT get blanket read/write of
--    private notes beyond what deletion needs.
-- ════════════════════════════════════════════════════════════════════

alter table public.profiles    enable row level security;
alter table public.goals       enable row level security;
alter table public.entries     enable row level security;
alter table public.entry_links enable row level security;
alter table public.tags        enable row level security;
alter table public.entry_tags  enable row level security;

-- ---- profiles ----
drop policy if exists profiles_select_self on public.profiles;
create policy profiles_select_self on public.profiles
  for select using (id = auth.uid() or public.is_admin());

drop policy if exists profiles_update_self on public.profiles;
create policy profiles_update_self on public.profiles
  for update using (id = auth.uid()) with check (id = auth.uid());

-- Deletes of users go through delete_user(); no direct delete policy needed.

-- ---- goals ----
drop policy if exists goals_rw_own on public.goals;
create policy goals_rw_own on public.goals
  for all using (user_id = auth.uid()) with check (user_id = auth.uid());

drop policy if exists goals_admin_delete on public.goals;
create policy goals_admin_delete on public.goals
  for delete using (public.is_admin());

-- ---- entries ----
drop policy if exists entries_rw_own on public.entries;
create policy entries_rw_own on public.entries
  for all using (user_id = auth.uid()) with check (user_id = auth.uid());

drop policy if exists entries_admin_delete on public.entries;
create policy entries_admin_delete on public.entries
  for delete using (public.is_admin());

-- ---- entry_links ----
drop policy if exists entry_links_rw_own on public.entry_links;
create policy entry_links_rw_own on public.entry_links
  for all using (user_id = auth.uid()) with check (user_id = auth.uid());

drop policy if exists entry_links_admin_delete on public.entry_links;
create policy entry_links_admin_delete on public.entry_links
  for delete using (public.is_admin());

-- ---- tags ----
drop policy if exists tags_rw_own on public.tags;
create policy tags_rw_own on public.tags
  for all using (user_id = auth.uid()) with check (user_id = auth.uid());

-- ---- entry_tags (owned transitively through the entry) ----
drop policy if exists entry_tags_rw_own on public.entry_tags;
create policy entry_tags_rw_own on public.entry_tags
  for all using (
    exists (select 1 from public.entries e
            where e.id = entry_tags.entry_id and e.user_id = auth.uid())
  ) with check (
    exists (select 1 from public.entries e
            where e.id = entry_tags.entry_id and e.user_id = auth.uid())
  );

-- ════════════════════════════════════════════════════════════════════
-- 5. ADMIN ACTIONS (SECURITY DEFINER, caller verified inside)
-- ════════════════════════════════════════════════════════════════════

-- List all users with aggregate stats — admin only.
create or replace function public.admin_list_users()
returns table (
  id uuid, username text, display_name text, role text,
  created_at timestamptz, goal_count bigint, entry_count bigint, total_minutes bigint
)
language sql
security definer
set search_path = public
as $$
  select p.id, p.username, p.display_name, p.role, p.created_at,
         (select count(*) from public.goals   g where g.user_id = p.id),
         (select count(*) from public.entries e where e.user_id = p.id),
         coalesce((select sum(e.minutes) from public.entries e where e.user_id = p.id), 0)
  from public.profiles p
  where public.is_admin()         -- empty result for non-admins
  order by p.created_at;
$$;

-- Delete a user's dashboard DATA but keep the account.
create or replace function public.admin_clear_dashboard(target uuid)
returns void
language plpgsql
security definer
set search_path = public
as $$
begin
  if not public.is_admin() then
    raise exception 'not authorized';
  end if;
  delete from public.entries where user_id = target;  -- links/tags cascade
  delete from public.goals   where user_id = target;
end;
$$;

-- Delete a user entirely (account + all data via FK cascade).
-- Guards: caller must be admin, and cannot delete themselves here.
create or replace function public.admin_delete_user(target uuid)
returns void
language plpgsql
security definer
set search_path = public
as $$
begin
  if not public.is_admin() then
    raise exception 'not authorized';
  end if;
  if target = auth.uid() then
    raise exception 'admins cannot delete their own account here';
  end if;
  delete from auth.users where id = target;  -- profile + data cascade
end;
$$;

-- Expose only to logged-in callers; the functions self-check admin.
grant execute on function public.admin_list_users()        to authenticated;
grant execute on function public.admin_clear_dashboard(uuid) to authenticated;
grant execute on function public.admin_delete_user(uuid)     to authenticated;
grant execute on function public.is_admin()                  to authenticated;
