<script>
	let { size = 28 } = $props();
</script>

<!-- Lumen mark: a diamond "flux" core emitting rays. -->
<svg width={size} height={size} viewBox="0 0 32 32" fill="none" aria-hidden="true">
	<g stroke="var(--flux)" stroke-width="1.4" opacity="0.5">
		<line x1="16" y1="1" x2="16" y2="6" />
		<line x1="16" y1="26" x2="16" y2="31" />
		<line x1="1" y1="16" x2="6" y2="16" />
		<line x1="26" y1="16" x2="31" y2="16" />
		<line x1="6" y1="6" x2="9.5" y2="9.5" />
		<line x1="22.5" y1="22.5" x2="26" y2="26" />
		<line x1="26" y1="6" x2="22.5" y2="9.5" />
		<line x1="9.5" y1="22.5" x2="6" y2="26" />
	</g>
	<rect
		x="16"
		y="7"
		width="12.7"
		height="12.7"
		transform="rotate(45 16 16)"
		fill="var(--flux)"
		opacity="0.12"
	/>
	<rect
		x="16"
		y="10"
		width="8.5"
		height="8.5"
		transform="rotate(45 16 16)"
		fill="var(--flux)"
	/>
</svg>
