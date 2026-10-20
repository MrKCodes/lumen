<script>
	import { onMount } from 'svelte';
	import { adminListUsers, adminDeleteUser, adminClearDashboard } from '$lib/api.js';
	import { auth } from '$lib/stores/auth.svelte.js';
	import { fmtMinutes } from '$lib/util.js';

	let users = $state([]);
	let err = $state('');
	let msg = $state('');
	let loading = $state(true);

	onMount(load);
	async function load() {
		loading = true;
		err = '';
		try {
			users = await adminListUsers();
		} catch (ex) {
			err = ex.message;
		} finally {
			loading = false;
		}
	}

	async function clearDash(u) {
		if (!confirm(`Wipe ALL goals + entries for "${u.username}"? The account stays. This cannot be undone.`))
			return;
		try {
			await adminClearDashboard(u.id);
			msg = `Cleared ${u.username}'s dashboard.`;
			await load();
		} catch (ex) {
			err = ex.message;
		}
	}

	async function deleteUser(u) {
		if (u.id === auth.user?.id) {
			alert("You can't delete your own admin account here.");
			return;
		}
		if (!confirm(`Permanently DELETE user "${u.username}" and all their data? This cannot be undone.`))
			return;
		if (!confirm(`Final check — type-less confirm. Really delete ${u.username}?`)) return;
		try {
			await adminDeleteUser(u.id);
			msg = `Deleted ${u.username}.`;
			await load();
		} catch (ex) {
			err = ex.message;
		}
	}
</script>

<div class="wrap">
	<header class="head">
		<div>
			<div class="eyebrow" style="color:var(--plasma)">control room</div>
			<h1 class="display sm">Admin</h1>
			<p class="muted">Manage operators and their dashboards. Destructive actions are irreversible.</p>
		</div>
		<button class="btn btn--ghost" onclick={load}>↻ Refresh</button>
	</header>

	{#if err}<p class="error">{err}</p>{/if}
	{#if msg}<p class="ok">{msg}</p>{/if}

	{#if loading}
		<p class="faint mono">loading operators…</p>
	{:else}
		<div class="panel hud tablewrap">
			<table>
				<thead>
					<tr>
						<th>Operator</th>
						<th>Role</th>
						<th>Goals</th>
						<th>Entries</th>
						<th>Total</th>
						<th>Joined</th>
						<th class="ar">Actions</th>
					</tr>
				</thead>
				<tbody>
					{#each users as u}
						<tr>
							<td>
								<div class="uname">{u.display_name || u.username}</div>
								<div class="faint mono small">@{u.username}</div>
							</td>
							<td>
								{#if u.role === 'admin'}
									<span class="chip badge-admin">ADMIN</span>
								{:else}<span class="chip">user</span>{/if}
							</td>
							<td class="mono">{u.goal_count}</td>
							<td class="mono">{u.entry_count}</td>
							<td class="mono flux">{fmtMinutes(u.total_minutes)}</td>
							<td class="mono small faint">{u.created_at?.slice(0, 10)}</td>
							<td class="ar">
								<div class="acts">
									<button class="btn btn--ghost sm-btn" onclick={() => clearDash(u)}>
										Clear data
									</button>
									<button
										class="btn btn--danger sm-btn"
										onclick={() => deleteUser(u)}
										disabled={u.id === auth.user?.id}
									>
										Delete
									</button>
								</div>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
			{#if users.length === 0}<p class="faint mono pad">no operators found.</p>{/if}
		</div>
	{/if}
</div>

<style>
	.head {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		gap: 20px;
		margin-bottom: 26px;
	}
	.display.sm {
		font-size: clamp(2rem, 4vw, 3rem);
		margin: 6px 0 10px;
	}
	.tablewrap {
		overflow-x: auto;
		padding: 6px;
	}
	table {
		width: 100%;
		border-collapse: collapse;
		font-size: 0.9rem;
	}
	th {
		text-align: left;
		font-family: var(--font-mono);
		font-size: 0.68rem;
		letter-spacing: 0.14em;
		text-transform: uppercase;
		color: var(--text-faint);
		padding: 14px 14px;
		border-bottom: 1px solid var(--line);
	}
	td {
		padding: 14px 14px;
		border-bottom: 1px solid var(--ink-600);
		vertical-align: middle;
	}
	tr:last-child td {
		border-bottom: none;
	}
	.uname {
		font-weight: 500;
	}
	.small {
		font-size: 0.72rem;
	}
	.flux {
		color: var(--flux);
	}
	.ar {
		text-align: right;
	}
	.acts {
		display: inline-flex;
		gap: 8px;
	}
	.sm-btn {
		padding: 0.45em 0.8em;
		font-size: 0.72rem;
	}
	.pad {
		padding: 20px;
	}
</style>
