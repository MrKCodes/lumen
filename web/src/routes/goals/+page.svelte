<script>
	import { onMount } from 'svelte';
	import { listGoals, createGoal, updateGoal, deleteGoal, listEntries } from '$lib/api.js';
	import { PALETTE, fmtHours } from '$lib/util.js';

	let goals = $state([]);
	let entries = $state([]);
	let err = $state('');
	let showForm = $state(false);

	let draft = $state(newDraft());
	function newDraft() {
		return { title: '', description: '', color: PALETTE[0], target_hours: '', status: 'active' };
	}

	onMount(load);
	async function load() {
		try {
			[goals, entries] = await Promise.all([listGoals(), listEntries()]);
		} catch (ex) {
			err = ex.message;
		}
	}

	function minutesForGoal(id) {
		return entries.filter((e) => e.goal_id === id).reduce((s, e) => s + (e.minutes || 0), 0);
	}

	async function save(e) {
		e.preventDefault();
		err = '';
		try {
			await createGoal(draft);
			draft = newDraft();
			showForm = false;
			await load();
		} catch (ex) {
			err = ex.message;
		}
	}

	async function remove(g) {
		if (!confirm(`Delete goal "${g.title}"? Entries stay but lose their goal link.`)) return;
		try {
			await deleteGoal(g.id);
			await load();
		} catch (ex) {
			err = ex.message;
		}
	}

	async function cycleStatus(g) {
		const order = ['active', 'paused', 'done', 'archived'];
		const next = order[(order.indexOf(g.status) + 1) % order.length];
		await updateGoal(g.id, { status: next });
		await load();
	}
</script>

<div class="wrap">
	<header class="head">
		<div>
			<div class="eyebrow">targets</div>
			<h1 class="display sm">Goals</h1>
		</div>
		<button class="btn btn--flux" onclick={() => (showForm = !showForm)}>
			{showForm ? 'Close' : '+ New goal'}
		</button>
	</header>

	{#if err}<p class="error">{err}</p>{/if}

	{#if showForm}
		<form class="panel hud form" onsubmit={save}>
			<div class="grid grid-2">
				<label class="field">
					<span>Title</span>
					<input bind:value={draft.title} required placeholder="Linux Systems Internals" />
				</label>
				<label class="field">
					<span>Target hours (optional)</span>
					<input type="number" min="0" step="0.5" bind:value={draft.target_hours} />
				</label>
			</div>
			<label class="field">
				<span>Description</span>
				<input bind:value={draft.description} placeholder="What this goal is about" />
			</label>
			<div class="field">
				<span>Colour</span>
				<div class="swatches">
					{#each PALETTE as c}
						<button
							type="button"
							class="sw"
							class:on={draft.color === c}
							style="--c:{c}"
							onclick={() => (draft.color = c)}
							aria-label={c}
						></button>
					{/each}
				</div>
			</div>
			<button class="btn btn--flux">Create goal</button>
		</form>
	{/if}

	<div class="grid grid-2 list">
		{#each goals as g}
			<div class="panel gcard" style="--c:{g.color}">
				<div class="gtop">
					<span class="dot"></span>
					<span class="gt">{g.title}</span>
					<button class="chip stat-chip" onclick={() => cycleStatus(g)} title="click to change">
						{g.status}
					</button>
				</div>
				{#if g.description}<p class="muted desc">{g.description}</p>{/if}
				<div class="gfoot mono">
					<span class="hrs">{fmtHours(minutesForGoal(g.id))}h{g.target_hours ? ` / ${g.target_hours}h` : ''}</span>
					<button class="btn btn--danger sm-btn" onclick={() => remove(g)}>Delete</button>
				</div>
			</div>
		{/each}
	</div>

	{#if goals.length === 0}
		<p class="faint mono empty">No goals yet. Create one to start tracking.</p>
	{/if}
</div>

<style>
	.head {
		display: flex;
		justify-content: space-between;
		align-items: flex-end;
		margin-bottom: 26px;
	}
	.display.sm {
		font-size: clamp(2rem, 4vw, 3rem);
	}
	.form {
		padding: 24px;
		margin-bottom: 26px;
	}
	.swatches {
		display: flex;
		gap: 8px;
	}
	.sw {
		width: 30px;
		height: 30px;
		border-radius: 6px;
		background: var(--c);
		border: 2px solid transparent;
		cursor: pointer;
	}
	.sw.on {
		border-color: var(--text);
		box-shadow: 0 0 12px -2px var(--c);
	}
	.list {
		gap: 16px;
	}
	.gcard {
		padding: 18px 20px;
		border-left: 3px solid var(--c);
	}
	.gtop {
		display: flex;
		align-items: center;
		gap: 10px;
	}
	.dot {
		width: 10px;
		height: 10px;
		border-radius: 50%;
		background: var(--c);
		box-shadow: 0 0 10px -1px var(--c);
	}
	.gt {
		font-weight: 500;
		font-size: 1.1rem;
	}
	.stat-chip {
		margin-left: auto;
		text-transform: uppercase;
		font-size: 0.64rem;
		cursor: pointer;
	}
	.desc {
		margin-top: 8px;
		font-size: 0.9rem;
	}
	.gfoot {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-top: 16px;
		font-size: 0.85rem;
	}
	.hrs {
		color: var(--c);
	}
	.sm-btn {
		padding: 0.4em 0.8em;
		font-size: 0.72rem;
	}
	.empty {
		margin-top: 20px;
	}
</style>
