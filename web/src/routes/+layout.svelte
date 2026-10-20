<script>
	import '$lib/styles/app.css';
	import { onMount } from 'svelte';
	import { auth } from '$lib/stores/auth.svelte.js';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import Nav from '$lib/components/Nav.svelte';
	import Footer from '$lib/components/Footer.svelte';
	import SideRail from '$lib/components/SideRail.svelte';

	let { children } = $props();
	let ready = $state(false);

	const PUBLIC = ['/login', '/register'];

	onMount(async () => {
		await auth.init();
		ready = true;
	});

	// Route guard: bounce unauthenticated users to /login, and signed-in users
	// away from the auth pages. Runs reactively once auth has initialised.
	$effect(() => {
		if (!ready || auth.loading) return;
		const path = $page.url.pathname;
		const isPublic = PUBLIC.includes(path);
		if (!auth.isAuthed && !isPublic) goto('/login');
		else if (auth.isAuthed && isPublic) goto('/');
		else if (path === '/admin' && !auth.isAdmin) goto('/');
	});
</script>

{#if ready && !auth.loading}
	<div class="shell">
		{#if auth.isAuthed}<Nav />{/if}
		<main>
			{@render children()}
		</main>
		<Footer />
	</div>
	<SideRail />
{:else}
	<div class="boot">
		<span class="live-dot"></span>
		<span class="mono faint">initialising lumen…</span>
	</div>
{/if}

<style>
	.shell {
		display: flex;
		flex-direction: column;
		min-height: 100vh;
	}
	main {
		flex: 1 0 auto;
		padding: 34px 0 60px;
	}
	.boot {
		min-height: 100vh;
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 12px;
	}
</style>
