import { browser } from '$app/environment';
import { session } from '$lib/session';
import { get } from 'svelte/store';
import type { RequestEvent, ServerLoadEvent } from '@sveltejs/kit';
import { BACKEND_URL } from '$env/static/private';

// Define a type for the fetch function, which can accept a 'fetch' from a SvelteKit event or the global fetch
type Fetch = typeof fetch;

// The base URL for your Go backend. Use an environment variable in a real app.

export async function fetchAuth(
	input: RequestInfo | URL,
	init?: RequestInit,
	event?: RequestEvent | ServerLoadEvent
): Promise<Response> {
	let token: string | null | undefined = null;

	if (browser) {
		// Client-side: get the token from the session store
		const sessionData = get(session);
		// Assuming the token is stored in `user.token`
		// You might need to adjust this based on your actual user object structure
		token = sessionData.user?.token;
	} else if (event) {
		// Server-side: get the token from the cookies
		token = event.cookies.get('session');
	}

	const headers = new Headers(init?.headers);

	if (token) {
		headers.set('Authorization', `Bearer ${token}`);
	}

	const modifiedInit: RequestInit = {
		...init,
		headers
	};

	// Use the fetch from the event if available (server-side), otherwise use the global fetch
	const fetcher: Fetch = event?.fetch || fetch;

	// Construct the full URL if the input is a relative path
	const url = typeof input === 'string' && input.startsWith('/') ? `${BACKEND_URL}${input}` : input;

	return fetcher(url, modifiedInit);
}
