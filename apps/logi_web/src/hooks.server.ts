import { BACKEND_URL } from '$env/static/private';
import { redirect } from '@sveltejs/kit';
import type { Handle } from '@sveltejs/kit';

// The base URL for your Go backend. Use an environment variable in a real app.
export const handle: Handle = async ({ event, resolve }) => {
	// Get the token from the cookies
	const token = event.cookies.get('session');

	if (!token) {
		// If there's no token, ensure the user is not set and continue.
		// Protected routes will handle redirection in their `load` functions or below.
		event.locals.user = null;
	} else {
		// We have a token, so we need to verify it with the Go backend
		try {
			// Send a request to your backend's verification endpoint
			const response = await event.fetch(`${BACKEND_URL}/v1/users/me`, {
				// Or /verify, /user, etc.
				headers: {
					// Send the token in the Authorization header
					Authorization: `Bearer ${token}`
				}
			});

			if (response.ok) {
				// Token is valid, the backend returns user info.
				const userData = await response.json();
				// Attach user info to event.locals. Adjust field names to match your Go backend's response.
				event.locals.user = {
					...userData,
					role: userData.role
				};
			} else {
				// Token is invalid (expired, tampered, etc.), so clear it.
				console.warn(`Backend token validation failed with status: ${response.status}`);
				event.locals.user = null;
				event.cookies.delete('session', { path: '/' });
			}
		} catch (error) {
			// This catches network errors if the backend is down.
			console.error('Could not connect to backend for token verification:', error);
			event.locals.user = null;
			// For security, log the user out if the verification service is unavailable.
			event.cookies.delete('session', { path: '/' });
		}
	}

	// --- Route Protection Logic ---

	const isUserLoggedIn = !!event.locals.user;
	const isTryingToAccessApp = event.url.pathname.startsWith('/auth/dashboard');

	// If user is not logged in and tries to access a protected route, redirect to login
	if (!isUserLoggedIn && isTryingToAccessApp) {
		throw redirect(303, '/auth/login');
	}

	// If user is logged in and tries to access the login page, redirect to the dashboard
	if (isUserLoggedIn && event.url.pathname === '/auth/login') {
		throw redirect(303, '/auth/dashboard');
	}

	// Resolve the request and continue to the page
	return resolve(event);
};
