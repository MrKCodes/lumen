<script>
	import { auth } from '$lib/stores/auth.svelte.js';
	import { goto } from '$app/navigation';
	import Logo from '$lib/components/Logo.svelte';

	let email = $state('');
	let username = $state('');
	let password = $state('');
	let err = $state('');
	let busy = $state(false);

	async function submit(e) {
		e.preventDefault();
		err = '';
		busy = true;
		try {
			await auth.register({ email, password, username, displayName: username });
			goto('/'); // no email confirmation — straight in
		} catch (ex) {
			err = ex.message ?? 'Registration failed';
		} finally {
			busy = false;
		}
	}
</script>

<div class="auth">
	<div class="panel panel--glow hud card">
		<div class="head">
			<Logo size={34} />
			<div>
				<div class="eyebrow">new operator</div>
				<h1 class="ttl">Create account</h1>
			</div>
		</div>

		<form onsubmit={submit}>
			<label class="field">
				<span>Username</span>
				<input bind:value={username} required autocomplete="username" placeholder="kartikeya" />
			</label>
			<label class="field">
				<span>Email</span>
				<input type="email" bind:value={email} required autocomplete="email" />
			</label>
			<label class="field">
				<span>Password</span>
				<input
					type="password"
					bind:value={password}
					required
					minlength="6"
					autocomplete="new-password"
				/>
			</label>

			{#if err}<p class="error">{err}</p>{/if}

			<button class="btn btn--flux full" disabled={busy}>
				{busy ? 'creating…' : 'Ignite ✦'}
			</button>
			<p class="note faint mono">No email confirmation — you're in instantly.</p>
		</form>

		<p class="alt faint mono">
			Have an account? <a href="/login">Sign in →</a>
		</p>
	</div>
</div>

<style>
	.auth {
		min-height: 80vh;
		display: grid;
		place-items: center;
		padding: 20px;
	}
	.card {
		width: 100%;
		max-width: 420px;
		padding: 34px;
	}
	.head {
		display: flex;
		align-items: center;
		gap: 14px;
		margin-bottom: 26px;
	}
	.ttl {
		font-size: 1.7rem;
		margin-top: 2px;
	}
	.full {
		width: 100%;
		justify-content: center;
		margin-top: 6px;
	}
	.note {
		font-size: 0.72rem;
		text-align: center;
		margin-top: 10px;
	}
	.alt {
		margin-top: 18px;
		text-align: center;
		font-size: 0.8rem;
	}
</style>
