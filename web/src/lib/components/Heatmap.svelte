<script>
	import { heatmap } from '$lib/util.js';
	let { entries = [], weeks = 18 } = $props();
	const cells = $derived(heatmap(entries, weeks));
</script>

<div class="heat">
	{#each cells as c}
		<span
			class="cell lvl{c.level}"
			title={`${c.date}: ${c.minutes}m`}
		></span>
	{/each}
</div>

<style>
	.heat {
		display: grid;
		grid-auto-flow: column;
		grid-template-rows: repeat(7, 1fr);
		gap: 3px;
		overflow-x: auto;
		padding-bottom: 4px;
	}
	.cell {
		width: 13px;
		height: 13px;
		border-radius: 3px;
		background: var(--ink-600);
		border: 1px solid var(--line);
	}
	.lvl1 {
		background: color-mix(in oklab, var(--flux) 25%, var(--ink-600));
	}
	.lvl2 {
		background: color-mix(in oklab, var(--flux) 50%, var(--ink-600));
	}
	.lvl3 {
		background: color-mix(in oklab, var(--flux) 75%, var(--ink-600));
	}
	.lvl4 {
		background: var(--flux);
		box-shadow: 0 0 8px -2px var(--flux);
	}
</style>
