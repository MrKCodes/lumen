import adapter from '@sveltejs/adapter-static';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	preprocess: vitePreprocess(),
	kit: {
		// Static SPA for Cloudflare Pages: everything renders client-side,
		// unknown routes fall back to index.html (true SPA routing).
		adapter: adapter({
			fallback: 'index.html',
			pages: 'build',
			assets: 'build'
		}),
		// We render entirely in the browser (auth + data are client-side).
		prerender: { entries: [] }
	}
};

export default config;
