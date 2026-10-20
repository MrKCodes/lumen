<script>
	import { onMount } from 'svelte';
	import { listGoals, listEntries, createEntry, deleteEntry } from '$lib/api.js';
	import { fmtMinutes, fmtDate, todayISO } from '$lib/util.js';

	let goals = $state([]);
	let entries = $state([]);
	let err = $state('');
	let busy = $state(false);

	let form = $state(blank());
	function blank() {
		return {
			entry_date: todayISO(),
			goal_id: '',
			title: '',
			minutes: 30,
			notes: '',
			links: [{ url: '', label: '' }]
		};
	}

	onMount(load);
	async function load() {
		try {
			[goals, entries] = await Promise.all([listGoals(), listEntries()]);
			if (!form.goal_id && goals[0]) form.goal_id = goals[0].id;
		} catch (ex) {
			err = ex.message;
		}
	}

	function addLink() {
		form.links = [...form.links, { url: '', label: '' }];
	}
	function removeLink(i) {
		form.links = form.links.filter((_, idx) => idx !== i);
	}

	async function submit(e) {
		e.preventDefault();
		err = '';
		busy = true;
		try {
			await createEntry({ ...form, links: form.links.filter((l) => l.url) });
			form = blank();
			if (goals[0]) form.goal_id = goals[0].id;
			await load();
		} catch (ex) {
			err = ex.message;
		} finally {
			busy = false;
		}
	}

	async function remove(entry) {
		if (!confirm(`Delete "${entry.title}"?`)) return;
		try {
			await deleteEntry(entry.id);
			await load();
		} catch (ex) {
			err = ex.message;
		}
	}

	const quick = [15, 25, 30, 45, 60, 90];
</script>

<div class="wrap">
	<header class="head">
		<div class="eyebrow">capture</div>
		<h1 class="display sm">Log a session</h1>
	</header>

	{#if err}<p class="error">{err}</p>{/if}

	<div class="layout">
		<form class="panel hud form" onsubmit={submit}>
			<div class="grid grid-2">
				<label class="field">
					<span>Date</span>
					<input type="date" bind:value={form.entry_date} required />
				</label>
				<label class="field">
					<span>Goal</span>
					<select bind:value={form.goal_id}>
						<option value="">— none —</option>
						{#each goals as g}<option value={g.id}>{g.title}</option>{/each}
					</select>
				</label>
			</div>

			<label class="field">
				<span>What did you learn?</span>
				<input bind:value={form.title} required placeholder="§4 page cache — traced the OOM path" />
			</label>

			<div class="field">
				<span>Minutes</span>
				<div class="mins">
					<input type="number" min="0" bind:value={form.minutes} class="minbox" />
					<div class="quick">
						{#each quick as q}
							<button type="button" class="chip qbtn" onclick={() => (form.minutes = q)}>{q}</button>
						{/each}
					</div>
				</div>
			</div>

			<label class="field">
				<span>Notes (markdown ok)</span>
				<textarea bind:value={form.notes} rows="4" placeholder="What clicked, what's still shaky…"></textarea>
			</label>

			<div class="field">
				<span>Links</span>
				{#each form.links as link, i}
					<div class="linkrow">
						<input bind:value={link.url} placeholder="https://…" class="lurl" />
						<input bind:value={link.label} placeholder="label" class="llabel" />
						{#if form.links.length > 1}
							<button type="button" class="btn btn--ghost xbtn" onclick={() => removeLink(i)}>✕</button>
						{/if}
					</div>
				{/each}
				<button type="button" class="chip addlink" onclick={addLink}>+ add link</button>
			</div>

			<button class="btn btn--flux submit" disabled={busy}>
				{busy ? 'saving…' : 'Log it ✦'}
			</button>
		</form>

		<div class="feedcol">
			<div class="eyebrow fh">timeline · {entries.length}</div>
			<div class="feed">
				{#each entries as e}
					<article class="panel row">
						<div class="when mono faint">{fmtDate(e.entry_date)}</div>
						<div class="body">
							<div class="r1">
								{#if e.goal}<span class="gdot" style="--c:{e.goal.color}"></span>{/if}
								<span class="et">{e.title}</span>
								<span class="min mono">{fmtMinutes(e.minutes)}</span>
							</div>
							{#if e.notes}<p class="notes muted">{e.notes}</p>{/if}
							{#if e.entry_links?.length}
								<div class="links">
									{#each e.entry_links as l}
										<a href={l.url} target="_blank" rel="noreferrer" class="chip">↗ {l.label || 'link'}</a>
									{/each}
								</div>
							{/if}
						</div>
						<button class="del" onclick={() => remove(e)} title="delete">✕</button>
					</article>
				{/each}
				{#if entries.length === 0}<p class="faint mono">no entries yet.</p>{/if}
			</div>
		</div>
	</div>
</div>

<style>
	.head {
		margin-bottom: 24px;
	}
	.display.sm {
		font-size: clamp(2rem, 4vw, 3rem);
		margin-top: 6px;
	}
	.layout {
		display: grid;
		gap: 28px;
	}
	@media (min-width: 940px) {
		.layout {
			grid-template-columns: 1.1fr 1fr;
			align-items: start;
		}
	}
	.form {
		padding: 24px;
		position: sticky;
		top: 80px;
	}
	.mins {
		display: flex;
		gap: 12px;
		align-items: center;
	}
	.minbox {
		width: 110px;
	}
	.quick {
		display: flex;
		gap: 6px;
		flex-wrap: wrap;
	}
	.qbtn {
		cursor: pointer;
	}
	.qbtn:hover {
		border-color: var(--flux);
		color: var(--flux);
	}
	.linkrow {
		display: flex;
		gap: 8px;
		margin-bottom: 8px;
	}
	.lurl {
		flex: 2;
	}
	.llabel {
		flex: 1;
	}
	.xbtn {
		padding: 0 0.7em;
	}
	.addlink {
		cursor: pointer;
	}
	.submit {
		width: 100%;
		justify-content: center;
		margin-top: 8px;
	}
	.fh {
		margin-bottom: 14px;
	}
	.feed {
		display: grid;
		gap: 10px;
	}
	.row {
		display: flex;
		gap: 14px;
		padding: 14px 16px;
		position: relative;
	}
	.when {
		font-size: 0.75rem;
		min-width: 46px;
		padding-top: 2px;
	}
	.body {
		flex: 1;
		min-width: 0;
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
		padding-right: 18px;
	}
	.notes {
		font-size: 0.86rem;
		margin-top: 6px;
		white-space: pre-wrap;
	}
	.links {
		display: flex;
		flex-wrap: wrap;
		gap: 6px;
		margin-top: 8px;
	}
	.del {
		position: absolute;
		top: 10px;
		right: 10px;
		background: none;
		border: none;
		color: var(--text-faint);
		cursor: pointer;
		font-size: 0.8rem;
	}
	.del:hover {
		color: var(--bad);
	}
</style>
