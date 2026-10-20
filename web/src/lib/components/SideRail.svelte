<script>
	import { LINKS } from '$lib/links.js';
</script>

<!-- Fixed right-edge rail, vertically centred. Hidden on small screens
     (links fold into the Footer there). Clears the bottom-right BMC widget. -->
<aside class="rail" aria-label="Social links">
	<span class="tick" aria-hidden="true"></span>

	<a href={LINKS.repo} target="_blank" rel="noreferrer" data-label="project">
		<svg viewBox="0 0 24 24" width="19" height="19" fill="currentColor" aria-hidden="true"><path d="M12 .5A11.5 11.5 0 0 0 .5 12 11.5 11.5 0 0 0 8.4 23c.6.1.8-.3.8-.6v-2c-3.2.7-3.9-1.4-3.9-1.4-.5-1.3-1.3-1.7-1.3-1.7-1-.7.1-.7.1-.7 1.2.1 1.8 1.2 1.8 1.2 1 .9 2.6.6 3.2.4.1-.7.4-1.2.7-1.5-2.6-.3-5.3-1.3-5.3-5.7 0-1.3.5-2.3 1.2-3.1-.1-.3-.5-1.5.1-3.1 0 0 1-.3 3.3 1.2a11.5 11.5 0 0 1 6 0C17.3 4.7 18.3 5 18.3 5c.6 1.6.2 2.8.1 3.1.8.8 1.2 1.8 1.2 3.1 0 4.4-2.7 5.4-5.3 5.7.4.4.8 1.1.8 2.2v3.3c0 .3.2.7.8.6A11.5 11.5 0 0 0 23.5 12 11.5 11.5 0 0 0 12 .5Z"/></svg>
		<span class="label">project repo</span>
	</a>

	<a href={LINKS.profile} target="_blank" rel="noreferrer" data-label="profile">
		<svg viewBox="0 0 24 24" width="19" height="19" fill="none" stroke="currentColor" stroke-width="1.8" aria-hidden="true"><circle cx="12" cy="8" r="4"/><path d="M4 21a8 8 0 0 1 16 0"/></svg>
		<span class="label">@MrKCodes</span>
	</a>

	<a href={LINKS.org} target="_blank" rel="noreferrer" data-label="org">
		<svg viewBox="0 0 24 24" width="19" height="19" fill="none" stroke="currentColor" stroke-width="1.8" aria-hidden="true"><path d="M3 21h18M5 21V5a2 2 0 0 1 2-2h6a2 2 0 0 1 2 2v16M15 21V9h2a2 2 0 0 1 2 2v10M8 7h2M8 11h2M8 15h2"/></svg>
		<span class="label">OpsAlchemist</span>
	</a>

	<a href={LINKS.website} target="_blank" rel="noreferrer" data-label="site">
		<svg viewBox="0 0 24 24" width="19" height="19" fill="none" stroke="currentColor" stroke-width="1.8" aria-hidden="true"><circle cx="12" cy="12" r="9"/><path d="M3 12h18M12 3c3 3 3 15 0 18M12 3c-3 3-3 15 0 18"/></svg>
		<span class="label">website</span>
	</a>

	<a href={LINKS.x} target="_blank" rel="noreferrer" data-label="x">
		<svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor" aria-hidden="true"><path d="M18.9 1.6h3.5l-7.6 8.7 9 11.9h-7l-5.5-7.2-6.3 7.2H1.5l8.2-9.3L1 1.6h7.2l5 6.6 5.7-6.6Zm-1.2 18.4h1.9L6.4 3.6H4.3l13.4 16.4Z"/></svg>
		<span class="label">@KubeMonk</span>
	</a>

	<span class="tick" aria-hidden="true"></span>
</aside>

<style>
	.rail {
		position: fixed;
		right: 0;
		top: 50%;
		transform: translateY(-50%);
		z-index: 40;
		display: flex;
		flex-direction: column;
		align-items: flex-end;
		gap: 6px;
		padding: 10px 0;
	}
	/* little HUD ticks top & bottom of the rail */
	.tick {
		width: 22px;
		height: 1px;
		background: var(--line);
		margin-right: 14px;
	}
	.rail a {
		position: relative;
		display: flex;
		align-items: center;
		justify-content: center;
		width: 44px;
		height: 44px;
		color: var(--text-dim);
		background: color-mix(in oklab, var(--ink-700) 85%, transparent);
		border: 1px solid var(--line);
		border-right: 0;
		border-radius: 10px 0 0 10px;
		backdrop-filter: blur(8px);
		transition:
			color 0.15s ease,
			background 0.15s ease,
			width 0.22s cubic-bezier(0.2, 0.8, 0.2, 1),
			box-shadow 0.2s ease;
		overflow: hidden;
	}
	/* label hidden until hover; expands leftward */
	.rail a .label {
		position: absolute;
		right: 46px;
		white-space: nowrap;
		font-family: var(--font-mono);
		font-size: 0.74rem;
		color: var(--flux);
		opacity: 0;
		transform: translateX(6px);
		transition:
			opacity 0.18s ease,
			transform 0.18s ease;
		pointer-events: none;
	}
	.rail a:hover {
		color: var(--flux);
		background: var(--ink-600);
		border-color: color-mix(in oklab, var(--flux) 45%, var(--line));
		box-shadow:
			-2px 0 0 0 var(--flux),
			0 0 22px -8px var(--flux);
		text-decoration: none;
	}
	.rail a:hover .label {
		opacity: 1;
		transform: translateX(0);
	}
	.rail a:focus-visible {
		outline: 2px solid var(--flux);
		outline-offset: 2px;
	}

	/* Fold away on small screens — Footer carries the links there. */
	@media (max-width: 900px) {
		.rail {
			display: none;
		}
	}
	@media (prefers-reduced-motion: reduce) {
		.rail a,
		.rail a .label {
			transition: none;
		}
	}
</style>
