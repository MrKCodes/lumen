<script>
	import { onMount } from 'svelte';
	import { listGoals, listEntries } from '$lib/api.js';
	import { auth } from '$lib/stores/auth.svelte.js';
	import { fmtMinutes, fmtDate, currentStreak, minutesByDay } from '$lib/util.js';
	import StatCard from '$lib/components/StatCard.svelte';
	import Heatmap from '$lib/components/Heatmap.svelte';
	import GoalBar from '$lib/components/GoalBar.svelte';

	let goals = $state([]);
	let entries = $state([]);
	let loading = $state(true);
	let err = $state('');

	onMount(load);
	async function load() {
		loading = true;
		try {
			[goals, entries] = await Promise.all([listGoals(), listEntries()]);
		} catch (ex) {
			err = ex.message;
		} finally {
			loading = false;
		}
	}

	const totalMin = $derived(entries.reduce((s, e) => s + (e.minutes || 0), 0));
	const streak = $derived(currentStreak(entries));
	const activeDays = $derived(Object.keys(minutesByDay(entries)).length);
	const last7 = $derived(() => {
		const cut = new Date();
		cut.setDate(cut.getDate() - 7);
		const iso = cut.toISOString().slice(0, 10);
		return entries.filter((e) => e.entry_date >= iso).reduce((s, e) => s + (e.minutes || 0), 0);
	});

	// minutes per goal
	function minutesForGoal(goalId) {
		return entries.filter((e) => e.goal_id === goalId).reduce((s, e) => s + (e.minutes || 0), 0);
	}
	const recent = $derived(entries.slice(0, 8));
</script>

<div class="wrap">
	<header class="hero">
		<div>
			<div class="eyebrow">operator · {auth.profile?.username ?? ''}</div>
			<h1 class="display">
				Keep the <span class="glow">light</span> on.
			</h1>
			<p class="muted lead">
				Every logged hour adds flux. Point it at a goal and watch it climb.
			</p>
		</div>
		<a href="/log" class="btn btn--flux big-cta">+ Log a session</a>
	</header>

	{#if err}<p class="error">{err}</p>{/if}

	{#if loading}
		<p class="faint mono">loading telemetry…</p>
	{:else}
		<section class="grid grid-4 stats">
			<StatCard label="Total flux" value={fmtMinutes(totalMin)} sub="all time" />
			<StatCard label="Last 7 days" value={fmtMinutes(last7())} sub="rolling" accent="var(--ice)" />
			<StatCard label="Streak" value={`${streak}d`} sub="consecutive" accent="var(--plasma)" />
			<StatCard label="Active days" value={activeDays} sub="logged ≥1 session" accent="var(--aurora)" />
		</section>

		<section class="panel hud block">
			<div class="block-head">
				<div class="eyebrow">activity · last 18 weeks</div>
				<span class="faint mono small">min/day → brightness</span>
			</div>
			<Heatmap {entries} />
		</section>

		<section class="split">
			<div>
				<div class="sec-title">
					<h2>Goals</h2>
					<a href="/goals" class="faint mono small">manage →</a>
				</div>
				{#if goals.length === 0}
					<div class="panel empty">
						<p class="muted">No goals yet.</p>
						<a href="/goals" class="btn btn--ghost">Create your first goal</a>
					</div>
				{:else}
					<div class="grid" style="gap:12px">
						{#each goals as g}
							<GoalBar goal={g} minutes={minutesForGoal(g.id)} />
						{/each}
					</div>
				{/if}
			</div>

			<div>
				<div class="sec-title">
					<h2>Recent</h2>
					<a href="/log" class="faint mono small">log →</a>
				</div>
				{#if recent.length === 0}
					<div class="panel empty"><p class="muted">Nothing logged yet.</p></div>
				{:else}
					<div class="feed">
						{#each recent as e}
							<article class="row panel">
								<div class="when mono faint">{fmtDate(e.entry_date)}</div>
								<div class="body">
									<div class="r1">
										{#if e.goal}<span class="gdot" style="--c:{e.goal.color}"></span>{/if}
										<span class="et">{e.title}</span>
										<span class="min mono">{fmtMinutes(e.minutes)}</span>
									</div>
									{#if e.entry_links?.length}
										<div class="links">
											{#each e.entry_links as l}
												<a href={l.url} target="_blank" rel="noreferrer" class="chip">
													↗ {l.label || new URL(l.url).hostname}
												</a>
											{/each}
										</div>
									{/if}
								</div>
							</article>
						{/each}
					</div>
				{/if}
			</div>
		</section>
	{/if}
</div>

<style>
	.hero {
		display: flex;
		align-items: flex-end;
		justify-content: space-between;
		gap: 24px;
		margin-bottom: 34px;
		flex-wrap: wrap;
	}
	.display {
		margin-top: 8px;
	}
	.glow {
		color: var(--flux);
		text-shadow: 0 0 30px rgba(198, 255, 74, 0.5);
	}
	.lead {
		margin-top: 12px;
		max-width: 46ch;
	}
	.big-cta {
		font-size: 0.95rem;
		padding: 0.9em 1.4em;
	}
	.stats {
		margin-bottom: 18px;
	}
	.block {
		padding: 22px 24px;
		margin-bottom: 30px;
	}
	.block-head {
		display: flex;
		justify-content: space-between;
		align-items: baseline;
		margin-bottom: 16px;
	}
	.small {
		font-size: 0.7rem;
	}
	.split {
		display: grid;
		gap: 28px;
	}
	@media (min-width: 900px) {
		.split {
			grid-template-columns: 1fr 1fr;
		}
	}
	.sec-title {
		display: flex;
		align-items: baseline;
		justify-content: space-between;
		margin-bottom: 14px;
	}
	.sec-title h2 {
		font-size: 1.3rem;
	}
	.empty {
		padding: 28px;
		display: flex;
		flex-direction: column;
		gap: 14px;
		align-items: flex-start;
	}
	.feed {
		display: grid;
		gap: 10px;
	}
	.row {
		display: flex;
		gap: 16px;
		padding: 14px 16px;
	}
	.when {
		font-size: 0.75rem;
		min-width: 48px;
		padding-top: 2px;
	}
	.body {
		flex: 1;
	}
	.r1 {
		display: flex;
		align-items: center;
		gap: 9px;
	}
	.gdot {
		width: 8px;
		height: 8px;
		border-radius: 50%;
		background: var(--c);
		box-shadow: 0 0 8px -1px var(--c);
		flex: none;
	}
	.et {
		font-weight: 500;
	}
	.min {
		margin-left: auto;
		color: var(--flux);
		font-size: 0.8rem;
	}
	.links {
		display: flex;
		flex-wrap: wrap;
		gap: 6px;
		margin-top: 8px;
	}
</style>
