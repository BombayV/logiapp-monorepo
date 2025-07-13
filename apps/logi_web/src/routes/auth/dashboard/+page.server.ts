import { fail, redirect } from '@sveltejs/kit';
import type { PageServerLoad, Actions } from './$types';
import { BACKEND_URL } from '$env/static/private';

export const load: PageServerLoad = async ({ locals }) => {
	// The hook already protects this route, so if we reach here, `locals.user` should exist.
	// If for some reason it doesn't, we redirect.
	if (!locals.user) {
		throw redirect(303, '/auth/login');
	}

	// Pass the user data to the page component
	return {
		user: locals.user
	};
};

export const actions: Actions = {
	logout: async ({ cookies }) => {
		try {
			const token = cookies.get('session');
			if (!token) {
				console.warn('No session token found during logout.');
				return fail(400, { error: 'No session token found.' });
			}
			const response = await fetch(`${BACKEND_URL}/v1/users/logout`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					Authorization: `Bearer ${token}` // Include the token in the request
				},
				credentials: 'include' // Ensure cookies are sent with the request
			});
			if (!response.ok) {
				const errorData = await response.json().catch(() => ({ message: 'Logout failed.' }));
				return fail(response.status, { error: errorData.error || 'Logout failed.' });
			}
			// Clear the session cookie
			cookies.delete('session', { path: '/' });
		} catch (error) {
			console.error('Error during logout:', error);
			return fail(500, { error: 'Could not connect to the authentication service.' });
		}

		// Redirect to the login page after successful logout
		throw redirect(303, '/auth/login');
	}
};
