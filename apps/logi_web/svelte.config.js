import adapter from '@sveltejs/adapter-cloudflare';
import { vitePreprocess } from '@sveltejs/vite-plugin-svelte';

/** @type {import('@sveltejs/kit').Config} */
const config = {
	// Consult https://svelte.dev/docs/kit/integrations
	// for more information about preprocessors
	preprocess: vitePreprocess(),

	kit: {
		// adapter-cloudflare optimizes your app for deployment to Cloudflare Pages
		// See https://svelte.dev/docs/kit/adapters for more information about adapters.
		adapter: adapter({
			// pages: true, // Use this if deploying to Cloudflare Pages
			// fallback: null, // Optional: specify fallback for SPA mode
		}),
		alias: {
			'@/*': './src/lib/*'
		}
	}
};

export default config;
