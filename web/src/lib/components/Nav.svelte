<script>
	import { auth } from '$lib/stores/auth.svelte.js';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import Logo from './Logo.svelte';

	async function logout() {
		await auth.logout();
		goto('/login');
	}

	const links = [
		{ href: '/', label: 'Dashboard' },
		{ href: '/goals', label: 'Goals' },
		{ href: '/log', label: 'Log' }
	];
</script>

<nav>
	<div class="wrap inner">
		<a href="/" class="brand">
			<Logo size={26} />
			<span class="name">LUMEN</span>
			<span class="live-dot" title="live"></span>
		</a>

		{#if auth.isAuthed}
			<div class="tabs">
				{#each links as l}
					<a href={l.href} class="tab" class:active={$page.url.pathname === l.href}>{l.label}</a>
				{/each}
				{#if auth.isAdmin}
					<a href="/admin" class="tab admin-tab" class:active={$page.url.pathname === '/admin'}>
						Admin
					</a>
				{/if}
			</div>

			<div class="right">
				<span class="who mono">
					{auth.profile?.display_name ?? auth.profile?.username ?? 'you'}
					{#if auth.isAdmin}<span class="chip badge-admin">ADMIN</span>{/if}
				</span>
				<button class="btn btn--ghost" onclick={logout}>Sign out</button>
			</div>
		{/if}
	</div>
</nav>

<style>
	nav {
		position: sticky;
		top: 0;
		z-index: 50;
		backdrop-filter: blur(12px);
		background: color-mix(in oklab, var(--ink-900) 78%, transparent);
		border-bottom: 1px solid var(--line);
	}
	.inner {
		display: flex;
		align-items: center;
		gap: 24px;
		height: 62px;
	}
	.brand {
		display: flex;
		align-items: center;
		gap: 10px;
		color: var(--text);
	}
	.brand:hover {
		text-decoration: none;
	}
	.name {
		font-weight: 700;
		letter-spacing: 0.28em;
		font-size: 0.95rem;
		padding-left: 2px;
	}
	.tabs {
		display: flex;
		gap: 4px;
		margin-left: 12px;
	}
	.tab {
		font-family: var(--font-mono);
		font-size: 0.8rem;
		color: var(--text-dim);
		padding: 0.5em 0.85em;
		border-radius: var(--radius);
		border: 1px solid transparent;
	}
	.tab:hover {
		color: var(--text);
		text-decoration: none;
		background: var(--ink-700);
	}
	.tab.active {
		color: var(--flux);
		border-color: var(--line);
		background: var(--ink-700);
	}
	.admin-tab.active {
		color: var(--plasma);
	}
	.right {
		margin-left: auto;
		display: flex;
		align-items: center;
		gap: 14px;
	}
	.who {
		font-size: 0.8rem;
		color: var(--text-dim);
		display: flex;
		align-items: center;
		gap: 8px;
	}
	@media (max-width: 680px) {
		.who,
		.name {
			display: none;
		}
	}
</style>
