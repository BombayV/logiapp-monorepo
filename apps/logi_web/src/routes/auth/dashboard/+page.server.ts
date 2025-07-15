import { fail, redirect } from '@sveltejs/kit';
import type { PageServerLoad, Actions } from './$types';
import { BACKEND_URL } from '$env/static/private';
import { fetchAuth } from '$lib/fetchAuth';
import type { RequestEvent } from '@sveltejs/kit';
import { getAllUsers } from '@/server/users';
import type { User } from '@/components/users/columns';

export const load: PageServerLoad = async (event) => {
	const { locals } = event;
	if (!locals.user) {
		throw redirect(303, '/auth/login');
	}

	const users = await getAllUsers(event) as User[];

	// Pass the user data to the page component
	return {
		user: locals.user,
		users: users,
	};
};

const logout = async (event: RequestEvent) => {
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
			return fail(response.status, { error: errorData.error || 'Cerrar sesión falló.' });
		}
		// Clear the session cookie
		event.cookies.delete('session', { path: '/' });
	} catch (error) {
		console.error('Error during logout:', error);
		return fail(500, { error: 'No se pudo conectar con el servicio de autenticación.' });
	}

	// Redirect to the login page after successful logout
	throw redirect(303, '/auth/login');
}

const create_order = async (event: RequestEvent) => {
	{
		const formData = await event.request.formData();
		const orderDetails = {
			address: formData.get('address'),
			email: formData.get('email')
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
				const errorData = await response
					.json()
					.catch(() => ({ message: 'Order creation failed.' }));
				return fail(response.status, { error: errorData.error || 'Error al crear el pedido.' });
			}

			// Optionally, you can redirect or return a success message
			return { success: true };
		} catch (error) {
			console.error('Error during order creation:', error);
			return fail(500, { error: 'No se pudo conectar con el servicio de pedidos.' });
		}
	}
};

export const actions: Actions = {
	logout,
	create_order
};
