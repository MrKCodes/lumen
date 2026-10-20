<script>
	import { LINKS } from '$lib/links.js';
</script>

<footer>
	<div class="wrap inner">
		<span class="made faint mono">
			built by <a href={LINKS.profile} target="_blank" rel="noreferrer">@KubeMonk</a>
			· <a href={LINKS.website} target="_blank" rel="noreferrer">opsalchemistlabs</a>
			· powered by caffeine
		</span>

		<!-- Mobile-only social row; desktop uses the SideRail instead. -->
		<nav class="social" aria-label="Social and project links">
			<a href={LINKS.repo} target="_blank" rel="noreferrer" aria-label="Lumen repository">
				<svg viewBox="0 0 24 24" width="18" height="18" fill="currentColor" aria-hidden="true"><path d="M12 .5A11.5 11.5 0 0 0 .5 12 11.5 11.5 0 0 0 8.4 23c.6.1.8-.3.8-.6v-2c-3.2.7-3.9-1.4-3.9-1.4-.5-1.3-1.3-1.7-1.3-1.7-1-.7.1-.7.1-.7 1.2.1 1.8 1.2 1.8 1.2 1 .9 2.6.6 3.2.4.1-.7.4-1.2.7-1.5-2.6-.3-5.3-1.3-5.3-5.7 0-1.3.5-2.3 1.2-3.1-.1-.3-.5-1.5.1-3.1 0 0 1-.3 3.3 1.2a11.5 11.5 0 0 1 6 0C17.3 4.7 18.3 5 18.3 5c.6 1.6.2 2.8.1 3.1.8.8 1.2 1.8 1.2 3.1 0 4.4-2.7 5.4-5.3 5.7.4.4.8 1.1.8 2.2v3.3c0 .3.2.7.8.6A11.5 11.5 0 0 0 23.5 12 11.5 11.5 0 0 0 12 .5Z"/></svg>
			</a>
			<a href={LINKS.profile} target="_blank" rel="noreferrer" aria-label="GitHub profile">
				<svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="1.8" aria-hidden="true"><circle cx="12" cy="8" r="4"/><path d="M4 21a8 8 0 0 1 16 0"/></svg>
			</a>
			<a href={LINKS.org} target="_blank" rel="noreferrer" aria-label="Organization">
				<svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="1.8" aria-hidden="true"><path d="M3 21h18M5 21V5a2 2 0 0 1 2-2h6a2 2 0 0 1 2 2v16M15 21V9h2a2 2 0 0 1 2 2v10M8 7h2M8 11h2M8 15h2"/></svg>
			</a>
			<a href={LINKS.website} target="_blank" rel="noreferrer" aria-label="Website">
				<svg viewBox="0 0 24 24" width="18" height="18" fill="none" stroke="currentColor" stroke-width="1.8" aria-hidden="true"><circle cx="12" cy="12" r="9"/><path d="M3 12h18M12 3c3 3 3 15 0 18M12 3c-3 3-3 15 0 18"/></svg>
			</a>
			<a href={LINKS.x} target="_blank" rel="noreferrer" aria-label="X / Twitter">
				<svg viewBox="0 0 24 24" width="16" height="16" fill="currentColor" aria-hidden="true"><path d="M18.9 1.6h3.5l-7.6 8.7 9 11.9h-7l-5.5-7.2-6.3 7.2H1.5l8.2-9.3L1 1.6h7.2l5 6.6 5.7-6.6Zm-1.2 18.4h1.9L6.4 3.6H4.3l13.4 16.4Z"/></svg>
			</a>
		</nav>
	</div>
</footer>

<style>
	footer {
		border-top: 1px solid var(--line);
		background: color-mix(in oklab, var(--ink-900) 70%, transparent);
		backdrop-filter: blur(8px);
		padding: 18px 0;
		margin-top: 40px;
	}
	.inner {
		display: flex;
		align-items: center;
		justify-content: space-between;
		gap: 14px;
		flex-wrap: wrap;
	}
	.made {
		font-size: 0.76rem;
	}
	.made a {
		color: var(--text-dim);
	}
	.made a:hover {
		color: var(--flux);
	}
	.social {
		display: flex;
		align-items: center;
		gap: 4px;
	}
	.social a {
		display: inline-flex;
		padding: 0.5em;
		color: var(--text-dim);
		border-radius: var(--radius);
		transition: color 0.15s ease, background 0.15s ease;
	}
	.social a:hover {
		color: var(--flux);
		background: var(--ink-700);
	}
	/* Desktop: the SideRail shows the socials, so hide the footer row. */
	@media (min-width: 901px) {
		.social {
			display: none;
		}
	}
	@media (max-width: 620px) {
		.inner {
			flex-direction: column;
			text-align: center;
		}
	}
</style>
