import { fail, redirect } from '@sveltejs/kit';
import type { PageServerLoad, Actions } from './$types';
import { BACKEND_URL } from '$env/static/private';
import { fetchAuth } from '$lib/fetchAuth';
import type { RequestEvent } from '@sveltejs/kit';
import { getAllUsers } from '@/server/users';
import type { User } from '@/components/users/columns';
import { createOrder, getAllOrders } from '@/server/orders';
import type { Order } from '@/components/orders/columns';

export const load: PageServerLoad = async (event) => {
	const { locals } = event;
	if (!locals.user) {
		throw redirect(303, '/auth/login');
	}

	const users = (await getAllUsers(event)) as User[];
	const orders = (await getAllOrders(event)) as Order[];
	if (!users || !orders) {
		return fail(500, { error: 'No se pudieron cargar los datos.' });
	}

	console.log('Loaded users:', users);
	// Pass the user data to the page component
	return {
		user: locals.user,
		users: users,
		orders: orders
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
};

const create_order = async (event: RequestEvent) => {
	{
		const formData = await event.request.formData();
		const orderDetails = {
			address: formData.get('address') as string,
			email: formData.get('email') as string,
			order_number: formData.get('order_number') as string
		};

		const newOrder = await createOrder(event, orderDetails);
		if (!newOrder) {
			return fail(500, { error: 'No se pudo crear el pedido.' });
		}

		// Optionally, you can redirect or return a success message
		return { success: true };
	}
};

const delete_user = async (event: RequestEvent) => {
	const formData = await event.request.formData();
	const userId = formData.get('user_id') as string;
	console.log('Deleting user with ID:', userId);
	if (!userId) {
		return fail(400, { error: 'ID de usuario no proporcionado.' });
	}
}

export const actions: Actions = {
	logout,
	create_order
};
