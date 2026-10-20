<script>
	let { label, value, sub = '', accent = 'var(--flux)' } = $props();
</script>

<div class="panel hud stat" style="--accent:{accent}">
	<div class="eyebrow">{label}</div>
	<div class="value mono">{value}</div>
	{#if sub}<div class="sub faint mono">{sub}</div>{/if}
</div>

<style>
	.stat {
		padding: 20px 20px 18px;
		overflow: hidden;
	}
	.stat::before {
		content: '';
		position: absolute;
		left: 0;
		top: 0;
		bottom: 0;
		width: 3px;
		background: var(--accent);
		box-shadow: 0 0 16px -2px var(--accent);
	}
	.value {
		font-size: 2.3rem;
		font-weight: 500;
		letter-spacing: -0.03em;
		line-height: 1.1;
		margin-top: 6px;
	}
	.sub {
		font-size: 0.72rem;
		margin-top: 4px;
	}
</style>
