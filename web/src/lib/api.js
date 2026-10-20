import { supabase } from '$lib/supabase.js';

/**
 * Data layer. Every call is scoped to the signed-in user by RLS, so these
 * queries never filter by user_id explicitly — the database does it.
 */

// ── goals ────────────────────────────────────────────────────────────
export async function listGoals() {
	const { data, error } = await supabase
		.from('goals')
		.select('*')
		.order('created_at', { ascending: true });
	if (error) throw error;
	return data;
}

export async function createGoal(g) {
	const { data: userData } = await supabase.auth.getUser();
	const payload = {
		user_id: userData.user.id,
		slug: g.slug || slugify(g.title),
		title: g.title,
		description: g.description ?? null,
		color: g.color,
		target_hours: g.target_hours ? Number(g.target_hours) : null,
		status: g.status || 'active'
	};
	const { data, error } = await supabase.from('goals').insert(payload).select().single();
	if (error) throw error;
	return data;
}

export async function updateGoal(id, patch) {
	const { data, error } = await supabase.from('goals').update(patch).eq('id', id).select().single();
	if (error) throw error;
	return data;
}

export async function deleteGoal(id) {
	const { error } = await supabase.from('goals').delete().eq('id', id);
	if (error) throw error;
}

// ── entries ──────────────────────────────────────────────────────────
export async function listEntries({ limit = 500 } = {}) {
	const { data, error } = await supabase
		.from('entries')
		.select('*, goal:goals(id, title, color, slug), entry_links(id, url, label)')
		.order('entry_date', { ascending: false })
		.limit(limit);
	if (error) throw error;
	return data;
}

export async function createEntry(e) {
	const { data: userData } = await supabase.auth.getUser();
	const { data, error } = await supabase
		.from('entries')
		.insert({
			user_id: userData.user.id,
			goal_id: e.goal_id || null,
			entry_date: e.entry_date,
			title: e.title,
			notes: e.notes ?? null,
			minutes: Number(e.minutes) || 0
		})
		.select()
		.single();
	if (error) throw error;

	if (e.links?.length) {
		const rows = e.links
			.filter((l) => l.url)
			.map((l) => ({ entry_id: data.id, user_id: userData.user.id, url: l.url, label: l.label || null }));
		if (rows.length) {
			const { error: linkErr } = await supabase.from('entry_links').insert(rows);
			if (linkErr) throw linkErr;
		}
	}
	return data;
}

export async function deleteEntry(id) {
	const { error } = await supabase.from('entries').delete().eq('id', id);
	if (error) throw error;
}

// ── admin (RPC — the functions self-check admin role) ─────────────────
export async function adminListUsers() {
	const { data, error } = await supabase.rpc('admin_list_users');
	if (error) throw error;
	return data;
}

export async function adminDeleteUser(targetId) {
	const { error } = await supabase.rpc('admin_delete_user', { target: targetId });
	if (error) throw error;
}

export async function adminClearDashboard(targetId) {
	const { error } = await supabase.rpc('admin_clear_dashboard', { target: targetId });
	if (error) throw error;
}

// ── helpers ──────────────────────────────────────────────────────────
export function slugify(s) {
	return s
		.toLowerCase()
		.trim()
		.replace(/[^a-z0-9]+/g, '-')
		.replace(/^-+|-+$/g, '');
}
