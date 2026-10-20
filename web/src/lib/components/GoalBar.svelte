<script>
	import { fmtHours } from '$lib/util.js';
	let { goal, minutes = 0 } = $props();
	const hours = $derived(minutes / 60);
	const target = $derived(goal.target_hours || 0);
	const pct = $derived(target ? Math.min(100, (hours / target) * 100) : 0);
</script>

<div class="goal panel" style="--c:{goal.color}">
	<div class="top">
		<span class="dot"></span>
		<span class="title">{goal.title}</span>
		<span class="status chip">{goal.status}</span>
	</div>
	<div class="nums mono">
		<span class="big">{fmtHours(minutes)}h</span>
		{#if target}<span class="faint">/ {target}h</span>{/if}
	</div>
	{#if target}
		<div class="track"><div class="fill" style="width:{pct}%"></div></div>
		<div class="pct faint mono">{Math.round(pct)}%</div>
	{/if}
</div>

<style>
	.goal {
		padding: 16px 18px;
		border-left: 3px solid var(--c);
	}
	.top {
		display: flex;
		align-items: center;
		gap: 9px;
	}
	.dot {
		width: 9px;
		height: 9px;
		border-radius: 50%;
		background: var(--c);
		box-shadow: 0 0 10px -1px var(--c);
	}
	.title {
		font-weight: 500;
	}
	.status {
		margin-left: auto;
		font-size: 0.66rem;
		text-transform: uppercase;
	}
	.nums {
		margin-top: 10px;
		display: flex;
		align-items: baseline;
		gap: 6px;
	}
	.big {
		font-size: 1.6rem;
		color: var(--c);
	}
	.track {
		margin-top: 10px;
		height: 6px;
		border-radius: 999px;
		background: var(--ink-900);
		border: 1px solid var(--line);
		overflow: hidden;
	}
	.fill {
		height: 100%;
		background: var(--c);
		box-shadow: 0 0 12px -1px var(--c);
		transition: width 0.6s cubic-bezier(0.2, 0.8, 0.2, 1);
	}
	.pct {
		font-size: 0.7rem;
		margin-top: 4px;
		text-align: right;
	}
</style>
