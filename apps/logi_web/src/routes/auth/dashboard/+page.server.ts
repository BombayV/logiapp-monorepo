import { fail, redirect } from '@sveltejs/kit';
import type { PageServerLoad, Actions } from './$types';
import type { RequestEvent } from '@sveltejs/kit';
import { createUser, getAllUsers, logoutUser, deleteUser } from '@/server/users';
import type { User } from '@/components/users/columns';
import { createOrder, getAllOrders, updateOrder } from '$lib/server/orders';
import type { Order } from '@/components/orders/columns';

export const load: PageServerLoad = async (event) => {
	const { locals } = event;
	if (!locals.user) {
		throw redirect(303, '/auth/login');
	}

	const orders = (await getAllOrders(event)) as Order[];

	if (locals.user.role !== 'admin') {
		return {
			user: locals.user,
			orders: orders
		};
	}

	const users = (await getAllUsers(event)) as User[];
	if (!users || !orders) {
		return fail(500, { error: 'No se pudieron cargar los datos.' });
	}

	return {
		user: locals.user,
		users: users,
		orders: orders
	};
};

const logout = async (event: RequestEvent) => {
	const result = await logoutUser(event);

	if (!result.success) {
		return fail(500, { error: result.error });
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

		// Server-side validation for order_number
		if (!orderDetails.order_number) {
			return fail(400, { error: 'El número de orden es requerido' });
		}

		if (orderDetails.order_number.length < 1 || orderDetails.order_number.length > 6) {
			return fail(400, { error: 'El número de orden debe tener entre 1 y 6 dígitos' });
		}

		if (!/^[0-9]+$/.test(orderDetails.order_number)) {
			return fail(400, { error: 'El número de orden solo puede contener dígitos' });
		}

		// Validate other required fields
		if (!orderDetails.address) {
			return fail(400, { error: 'La dirección es requerida' });
		}

		if (!orderDetails.email) {
			return fail(400, { error: 'El email es requerido' });
		}

		const result = await createOrder(event, orderDetails);

		// Handle validation errors from the orders service
		if (result && result.error) {
			return fail(400, { error: result.error });
		}

		// Handle case where no order was created
		if (!result) {
			return fail(500, { error: 'No se pudo crear el pedido.' });
		}

		// Return success message
		return { success: true, message: 'Orden creada exitosamente' };
	}
};

const create_user = async (event: RequestEvent) => {
	const formData = await event.request.formData();
	const userDetails = {
		email: formData.get('email') as string,
		password: formData.get('password') as string,
		first_name: formData.get('first_name') as string,
		last_name: formData.get('last_name') as string,
		phone_number: formData.get('phone_number') as string,
		role: formData.get('role') as 'sales' | 'driver'
	};

	// Validate required fields
	if (
		!userDetails.email ||
		!userDetails.password ||
		!userDetails.first_name ||
		!userDetails.last_name ||
		!userDetails.phone_number ||
		!userDetails.role
	) {
		return fail(400, { error: 'Todos los campos son requeridos.' });
	}

	const result = await createUser(event, userDetails);

	if (!result) {
		return fail(500, { error: 'Error al crear el usuario.' });
	}

	return { success: true, message: 'Usuario creado exitosamente' };
};

const delete_user = async (event: RequestEvent) => {
	const formData = await event.request.formData();
	const userId = formData.get('user_id') as string;
	if (!userId) {
		return fail(400, { error: 'ID de usuario no proporcionado.' });
	}

	const result = await deleteUser(event, userId);

	if (!result || !result.success) {
		return fail(500, { error: 'Error al eliminar el usuario.' });
	}

	return { success: true, message: 'Usuario eliminado exitosamente' };
};

const cancel_order = async (event: RequestEvent) => {
	const formData = await event.request.formData();
	const orderId = formData.get('order_id') as string;
	if (!orderId) {
		return fail(400, { error: 'ID de orden no proporcionado.' });
	}

	console.log('Cancelling order with ID:', orderId);

	// Update the order status to 'cancelled'
	const result = await updateOrder(event, orderId, { status: 'cancelled' });

	if (result && result.error) {
		return fail(400, { error: result.error });
	}

	if (!result) {
		return fail(500, { error: 'No se pudo cancelar la orden.' });
	}

	return { success: true, message: 'Orden cancelada exitosamente' };
};

export const actions: Actions = {
	logout,
	create_order,
	delete_user,
	create_user,
	cancel_order
};
