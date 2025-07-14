import { fail, redirect } from '@sveltejs/kit';
import type { PageServerLoad, Actions } from './$types';
import { BACKEND_URL } from '$env/static/private';
import { fetchAuth } from '$lib/fetchAuth';

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
	logout: async (event) => {
		try {
			const response = await fetchAuth(
				`${BACKEND_URL}/v1/users/logout`,
				{
					method: 'POST',
					headers: {
						'Content-Type': 'application/json'
					},
					credentials: 'include' // Ensure cookies are sent with the request
				},
				event
			);
			if (!response.ok) {
				const errorData = await response.json().catch(() => ({ message: 'Logout failed.' }));
				return fail(response.status, { error: errorData.error || 'Logout failed.' });
			}
			// Clear the session cookie
			event.cookies.delete('session', { path: '/' });
		} catch (error) {
			console.error('Error during logout:', error);
			return fail(500, { error: 'Could not connect to the authentication service.' });
		}

		// Redirect to the login page after successful logout
		throw redirect(303, '/auth/login');
	},
	create_order: async (event) => {
		const formData = await event.request.formData();
		const orderDetails = {
			address: formData.get('address'),
			email: formData.get('email'),
		};

		try {
			const response = await fetchAuth(
				`${BACKEND_URL}/v1/orders`,
				{
					method: 'POST',
					headers: {
						'Content-Type': 'application/json'
					},
					body: JSON.stringify(orderDetails),
					credentials: 'include' // Ensure cookies are sent with the request
				},
				event
			);

			if (!response.ok) {
				const errorData = await response.json().catch(() => ({ message: 'Order creation failed.' }));
				return fail(response.status, { error: errorData.error || 'Order creation failed.' });
			}

			// Optionally, you can redirect or return a success message
			return { success: true };
		} catch (error) {
			console.error('Error during order creation:', error);
			return fail(500, { error: 'Could not connect to the order service.' });
		}
	}
};
