import { error, redirect, fail } from '@sveltejs/kit';
import type { PageServerLoad, Actions } from './$types';
import {
	getOrderItems,
	getOrderById,
	createOrderItem,
	createOrderItemsBulk,
	deleteOrder,
	updateOrder,
	deleteOrderItem,
	getOrderForm,
	createOrderForm
} from '$lib/server/orders';
import type { Order, OrderItem } from '@/components/orders/columns';

export const load: PageServerLoad = async (event) => {
	const { locals, params } = event;

	if (!locals.user) {
		throw redirect(303, '/auth/login');
	}

	const orderId = params.order;

	if (!orderId) {
		throw error(404, 'Order ID not provided');
	}

	try {
		const order = (await getOrderById(event, orderId)) as Order;

		if (!order) {
			throw error(404, 'Order not found');
		}

		const items = (await getOrderItems(event, orderId)) as OrderItem[];
		if (!items) {
			throw error(404, 'Order items not found');
		}

		const form = await getOrderForm(event, orderId);

		return {
			user: locals.user,
			order: order,
			items: items,
			form: form
		};
	} catch (err) {
		console.error('Error loading order:', err);
		throw error(500, 'Failed to load order details');
	}
};

export const actions: Actions = {
	add_item: async (event) => {
		const { request, params, locals } = event;

		if (!locals.user) {
			throw redirect(303, '/auth/login');
		}

		const orderId = params.order;
		if (!orderId) {
			return fail(400, { error: 'Order ID not provided' });
		}

		const formData = await request.formData();
		const productName = formData.get('product_name') as string;
		const quantity = parseInt(formData.get('quantity') as string, 10);

		// Validation
		if (!productName || productName.trim() === '') {
			return fail(400, { error: 'El nombre del producto es requerido' });
		}

		if (!quantity || quantity <= 0) {
			return fail(400, { error: 'La cantidad debe ser mayor a 0' });
		}

		try {
			const result = await createOrderItem(event, orderId, {
				product_name: productName.trim(),
				quantity
			});

			if (result.error) {
				return fail(400, { error: result.error });
			}

			return { success: true, message: 'Item agregado exitosamente' };
		} catch (error) {
			console.error('Error adding item:', error);
			return fail(500, { error: 'Error interno del servidor' });
		}
	},

	bulk_add_items: async (event) => {
		const { request, params, locals } = event;

		if (!locals.user) {
			throw redirect(303, '/auth/login');
		}

		const orderId = params.order;
		if (!orderId) {
			return fail(400, { error: 'Order ID not provided' });
		}

		const formData = await request.formData();
		const itemsJson = formData.get('items') as string;

		if (!itemsJson) {
			return fail(400, { error: 'Items data is required' });
		}

		let items;
		try {
			items = JSON.parse(itemsJson);
		} catch (error) {
			return fail(400, { error: 'Invalid items data format' });
		}

		if (!Array.isArray(items) || items.length === 0) {
			return fail(400, { error: 'Items must be a non-empty array' });
		}

		// Validate each item
		for (const item of items) {
			if (!item.product_name || !item.quantity) {
				return fail(400, { error: 'Each item must have product_name and quantity' });
			}

			const quantityNum = parseInt(item.quantity);
			if (isNaN(quantityNum) || quantityNum <= 0) {
				return fail(400, { error: 'All quantities must be valid positive numbers' });
			}
		}

		try {
			const result = await createOrderItemsBulk(event, orderId, items);

			if (result.error) {
				return fail(400, { error: result.error });
			}

			return { success: true, message: `${items.length} items added successfully` };
		} catch (error) {
			console.error('Error adding bulk items:', error);
			return fail(500, { error: 'Failed to add items' });
		}
	},

	delete_order: async (event) => {
		const { params, locals } = event;

		if (!locals.user) {
			throw redirect(303, '/auth/login');
		}

		const orderId = params.order;
		if (!orderId) {
			return fail(400, { error: 'Order ID not provided' });
		}

		try {
			const result = await deleteOrder(event, orderId);

			if (result.error) {
				return fail(400, { error: result.error });
			}

			return { success: true, message: 'Orden eliminada exitosamente' };
		} catch (error) {
			console.error('Error deleting order:', error);
			return fail(500, { error: 'Error interno del servidor' });
		}
	},

	update_order: async (event) => {
		const { request, params, locals } = event;

		if (!locals.user) {
			throw redirect(303, '/auth/login');
		}

		const orderId = params.order;
		if (!orderId) {
			return fail(400, { error: 'Order ID not provided' });
		}

		const formData = await request.formData();
		const assigned_to = formData.get('assigned_to') as string;
		const address = formData.get('address') as string;
		const status = formData.get('status') as string;

		// Build update data object
		const updateData: { assigned_to?: string | null; address?: string; status?: string } = {};

		// Handle assigned_to field - include it if it's provided (even if empty string for unassignment)
		if (assigned_to !== null) {
			updateData.assigned_to = assigned_to === '' ? null : assigned_to;
		}

		if (address && address.trim() !== '') {
			updateData.address = address.trim();
		}
		if (status && status !== '') {
			updateData.status = status;
		}

		try {
			const result = await updateOrder(event, orderId, updateData);

			if (result.error) {
				return fail(400, { error: result.error });
			}

			return { success: true, message: 'Orden actualizada exitosamente' };
		} catch (error) {
			console.error('Error updating order:', error);
			return fail(500, { error: 'Error interno del servidor' });
		}
	},

	delete_item: async (event) => {
		const { request, params, locals } = event;

		if (!locals.user) {
			throw redirect(303, '/auth/login');
		}

		const orderId = params.order;
		if (!orderId) {
			return fail(400, { error: 'Order ID not provided' });
		}

		const formData = await request.formData();
		const itemId = formData.get('item_id') as string;

		if (!itemId) {
			return fail(400, { error: 'Item ID not provided' });
		}

		try {
			const result = await deleteOrderItem(event, orderId, itemId);

			if (result.error) {
				return fail(400, { error: result.error });
			}

			return { success: true, message: 'Item eliminado exitosamente' };
		} catch (error) {
			console.error('Error deleting item:', error);
			return fail(500, { error: 'Error interno del servidor' });
		}
	},

	create_survey: async (event) => {
		const { params, locals } = event;

		if (!locals.user) {
			throw redirect(303, '/auth/login');
		}

		const orderId = params.order;
		if (!orderId) {
			return fail(400, { error: 'Order ID not provided' });
		}

		try {
			const order = (await getOrderById(event, orderId)) as Order;
			if (!order) {
				return fail(404, { error: 'Order not found' });
			}

			const result = await createOrderForm(event, orderId, order.assigned_to);

			if (result.error) {
				return fail(400, { error: result.error });
			}

			return {
				success: true,
				message: 'Encuesta creada y enviada por correo exitosamente',
				form: result.form
			};
		} catch (error) {
			console.error('Error creating survey:', error);
			return fail(500, { error: 'Error interno del servidor' });
		}
	}
};
