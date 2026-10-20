<script>
	import { LINKS } from '$lib/links.js';
	// On-theme replacement for the BMC widget. The official widget renders its
	// popup in a cross-origin iframe we cannot style, so this is our own button
	// that opens the BMC page in a new tab. Message matches the old widget.
	const MESSAGE = 'Keep me off-call and on-coffee';
</script>

<a
	class="coffee"
	href={LINKS.coffee}
	target="_blank"
	rel="noreferrer"
	aria-label="Buy me a coffee"
	title={MESSAGE}
>
	<span class="cup" aria-hidden="true">
		<span class="steam s1"></span>
		<span class="steam s2"></span>
		<span class="steam s3"></span>
		<svg viewBox="0 0 24 24" width="20" height="20" fill="none" stroke="currentColor" stroke-width="1.8" stroke-linecap="round" stroke-linejoin="round">
			<path d="M4 9h13v4a5 5 0 0 1-5 5H9a5 5 0 0 1-5-5V9Z" />
			<path d="M17 10h2.2a2.3 2.3 0 0 1 0 4.6H17" />
			<path d="M4 21h14" />
		</svg>
	</span>
	<span class="msg">{MESSAGE}</span>
</a>

<style>
	.coffee {
		position: fixed;
		right: 20px;
		bottom: 20px;
		z-index: 60;
		display: inline-flex;
		align-items: center;
		gap: 0;
		height: 52px;
		padding: 0;
		width: 52px;
		border-radius: 999px;
		background: linear-gradient(180deg, var(--ink-700), var(--ink-800));
		border: 1px solid color-mix(in oklab, var(--flux) 45%, var(--line));
		color: var(--flux);
		box-shadow:
			0 0 0 1px rgba(198, 255, 74, 0.15),
			0 8px 30px -10px rgba(0, 0, 0, 0.8),
			0 0 26px -10px var(--flux);
		overflow: hidden;
		backdrop-filter: blur(8px);
		transition:
			width 0.28s cubic-bezier(0.2, 0.8, 0.2, 1),
			box-shadow 0.2s ease,
			transform 0.1s ease,
			border-color 0.2s ease;
	}
	.coffee:hover {
		/* wide enough for the full message: cup (52) + text + right padding.
		   bump this if you use a longer message. */
		width: 304px;
		border-color: var(--flux);
		text-decoration: none;
		box-shadow:
			0 0 0 1px var(--flux),
			0 8px 30px -8px rgba(0, 0, 0, 0.85),
			0 0 34px -8px var(--flux);
	}
	.coffee:active {
		transform: translateY(1px);
	}
	.coffee:focus-visible {
		outline: 2px solid var(--flux);
		outline-offset: 3px;
	}

	.cup {
		position: relative;
		flex: none;
		width: 52px;
		height: 52px;
		display: grid;
		place-items: center;
	}
	.msg {
		font-family: var(--font-mono);
		font-size: 0.76rem;
		white-space: nowrap;
		color: var(--text);
		opacity: 0;
		transform: translateX(6px);
		transition: opacity 0.18s ease, transform 0.18s ease;
		padding-right: 16px;
	}
	.coffee:hover .msg {
		opacity: 1;
		transform: translateX(0);
	}

	/* rising steam wisps */
	.steam {
		position: absolute;
		top: 12px;
		width: 2px;
		height: 7px;
		border-radius: 2px;
		background: var(--flux);
		opacity: 0;
		filter: blur(0.4px);
	}
	.s1 {
		left: 21px;
		animation: steam 2.4s ease-in-out infinite;
	}
	.s2 {
		left: 25px;
		animation: steam 2.4s ease-in-out 0.5s infinite;
	}
	.s3 {
		left: 29px;
		animation: steam 2.4s ease-in-out 1s infinite;
	}
	@keyframes steam {
		0% {
			opacity: 0;
			transform: translateY(2px) scaleY(0.6);
		}
		40% {
			opacity: 0.7;
		}
		100% {
			opacity: 0;
			transform: translateY(-7px) scaleY(1.1);
		}
	}

	@media (prefers-reduced-motion: reduce) {
		.coffee,
		.msg {
			transition: none;
		}
		.steam {
			animation: none;
		}
	}
</style>
