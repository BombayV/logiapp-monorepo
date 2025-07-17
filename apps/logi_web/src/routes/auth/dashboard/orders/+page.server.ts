import { fail, redirect } from '@sveltejs/kit';
import type { PageServerLoad, Actions } from './$types';
import type { RequestEvent } from '@sveltejs/kit';
import { createOrder, getAllOrders, updateOrder } from '$lib/server/orders';
import type { Order } from '@/components/orders/columns';

export const load: PageServerLoad = async (event) => {
	const { locals } = event;
	if (!locals.user) {
		throw redirect(303, '/auth/login');
	}

	try {
		const orders = (await getAllOrders(event)) as Order[];
		if (!orders) {
			return fail(500, { error: 'No se pudieron cargar los datos.' });
		}

		return {
			user: locals.user,
			orders: orders || []
		};
	} catch (err) {
		console.error('Error loading orders page data:', err);
		return {
			user: locals.user,
			orders: [],
			drivers: []
		};
	}
};

const create_order = async (event: RequestEvent) => {
	const formData = await event.request.formData();
	const orderDetails = {
		email: formData.get('email') as string,
		address: formData.get('address') as string,
		order_number: formData.get('order_number') as string
	};

	// Validate required fields
	if (!orderDetails.email || !orderDetails.address || !orderDetails.order_number) {
		return fail(400, { error: 'Todos los campos son requeridos' });
	}

	// Validate email format
	const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
	if (!emailRegex.test(orderDetails.email)) {
		return fail(400, { error: 'El formato del email no es válido' });
	}

	// Validate order number format
	const orderNumberRegex = /^[0-9]{1,6}$/;
	if (!orderNumberRegex.test(orderDetails.order_number)) {
		return fail(400, { error: 'El número de orden debe tener entre 1 y 6 dígitos' });
	}

	const result = await createOrder(event, orderDetails);

	if (!result.success) {
		return fail(500, { error: result.error });
	}

	return { success: true, message: 'Orden creada exitosamente' };
};

const update_order = async (event: RequestEvent) => {
	const formData = await event.request.formData();
	const orderDetails = {
		order_id: formData.get('order_id') as string,
		assigned_to: formData.get('assigned_to') as string,
		status: formData.get('status') as string
	};

	const updateData: { assigned_to?: string | null; address?: string; status?: string } = {};

	if (orderDetails.assigned_to) {
		updateData.assigned_to = orderDetails.assigned_to;
	}
	if (orderDetails.status) {
		updateData.status = orderDetails.status;
	}

	const result = await updateOrder(event, orderDetails.order_id, updateData);

	if (!result.success) {
		return fail(500, { error: result.error });
	}

	return { success: true, message: 'Orden actualizada exitosamente' };
};

export const actions: Actions = {
	create_order,
	update_order
};
