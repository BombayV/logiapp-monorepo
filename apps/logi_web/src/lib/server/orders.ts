import { fetchAuth } from '@/fetchAuth';
import type { RequestEvent } from '@sveltejs/kit';

export const getAllOrders = async (event: RequestEvent) => {
	try {
		const response = await fetchAuth(
			`/v1/orders`,
			{
				method: 'GET',
				headers: {
					'Content-Type': 'application/json'
				},
				credentials: 'include' // Ensure cookies are sent with the request
			},
			event
		);

		if (!response.ok) {
			const errorData = await response.json().catch(() => ({ message: 'Failed to fetch orders.' }));
			console.log('Error fetching orders:', errorData);
			return [];
		}

		const orders = await response.json();
		return orders.orders || [];
	} catch (error) {
		console.error('Error fetching orders:', error);
		return [];
	}
};

export const getOrderById = async (event: RequestEvent, orderId: string) => {
	try {
		const response = await fetchAuth(
			`/v1/orders/${orderId}`,
			{
				method: 'GET',
				headers: {
					'Content-Type': 'application/json'
				},
				credentials: 'include'
			},
			event
		);

		if (!response.ok) {
			const errorData = await response.json().catch(() => ({ message: 'Failed to fetch order.' }));
			console.log('Error fetching order:', errorData);
			return null;
		}

		const order = await response.json();
		return order;
	} catch (error) {
		console.error('Error fetching order:', error);
		return null;
	}
};

export const getOrderItems = async (event: RequestEvent, orderId: string) => {
	try {
		const response = await fetchAuth(
			`/v1/orders/${orderId}/items`,
			{
				method: 'GET',
				headers: {
					'Content-Type': 'application/json'
				},
				credentials: 'include'
			},
			event
		);

		if (!response.ok) {
			const errorData = await response
				.json()
				.catch(() => ({ message: 'Failed to fetch order items.' }));
			console.log('Error fetching order items:', errorData);
			return [];
		}

		const items = await response.json();
		return items.items || [];
	} catch (error) {
		console.error('Error fetching order items:', error);
		return [];
	}
};

export const createOrderItem = async (
	event: RequestEvent,
	orderId: string,
	itemData: { product_name: string; quantity: number }
) => {
	try {
		const response = await fetchAuth(
			`/v1/orders/${orderId}/items`,
			{
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(itemData),
				credentials: 'include'
			},
			event
		);

		if (!response.ok) {
			const errorData = await response
				.json()
				.catch(() => ({ error: 'Failed to create order item.' }));
			console.log('Error creating order item:', errorData);
			return { error: errorData.error || 'Error al crear el item' };
		}

		const newItem = await response.json();
		return newItem;
	} catch (error) {
		console.error('Error creating order item:', error);
		return { error: 'Error de conexión al crear el item' };
	}
};

export const createOrderItemsBulk = async (
	event: RequestEvent,
	orderId: string,
	items: { product_name: string; quantity: number }[]
) => {
	try {
		const response = await fetchAuth(
			`/v1/orders/${orderId}/items/bulk`,
			{
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ items }),
				credentials: 'include'
			},
			event
		);

		if (!response.ok) {
			const errorData = await response
				.json()
				.catch(() => ({ error: 'Failed to create order items.' }));
			console.log('Error creating order items:', errorData);
			return { error: errorData.error || 'Error al crear los items' };
		}

		const result = await response.json();
		return result;
	} catch (error) {
		console.error('Error creating order items:', error);
		return { error: 'Error de conexión al crear los items' };
	}
};

export const createOrder = async (
	event: RequestEvent,
	orderData: { email: string; address: string; order_number: string }
) => {
	console.log('Creating order with data:', orderData);

	// Client-side validation before sending to API
	if (!orderData.order_number) {
		console.log('Validation error: Order number is required');
		return { error: 'El número de orden es requerido' };
	}

	if (orderData.order_number.length < 1 || orderData.order_number.length > 6) {
		console.log('Validation error: Order number length invalid');
		return { error: 'El número de orden debe tener entre 1 y 6 dígitos' };
	}

	if (!/^[0-9]+$/.test(orderData.order_number)) {
		console.log('Validation error: Order number must be numeric');
		return { error: 'El número de orden solo puede contener dígitos' };
	}

	try {
		const response = await fetchAuth(
			`/v1/orders`,
			{
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(orderData),
				credentials: 'include'
			},
			event
		);

		if (!response.ok) {
			const errorData = await response.json().catch(() => ({ error: 'Failed to create order.' }));
			console.log('API Error creating order:', errorData);

			// Handle specific API validation errors
			if (errorData.error === 'order_number must contain only numbers') {
				return { error: 'El número de orden solo puede contener dígitos' };
			}

			// Handle duplicate order number error (if API returns it)
			if (
				(errorData.error && errorData.error.includes('duplicate')) ||
				errorData.error.includes('unique')
			) {
				return { error: 'Este número de orden ya existe' };
			}

			return { error: errorData.error || 'Error al crear la orden' };
		}

		const newOrder = await response.json();
		console.log('Order created successfully:', newOrder);
		return newOrder;
	} catch (error) {
		console.error('Error creating order:', error);
		return { error: 'Error de conexión al crear la orden' };
	}
};

export const deleteOrder = async (event: RequestEvent, orderId: string) => {
	try {
		const response = await fetchAuth(
			`/v1/orders/${orderId}`,
			{
				method: 'DELETE',
				credentials: 'include'
			},
			event
		);

		if (!response.ok) {
			const errorData = await response.json().catch(() => ({ error: 'Failed to delete order.' }));
			console.log('Error deleting order:', errorData);
			return { error: errorData.error || 'Error al eliminar la orden' };
		}

		return { success: true };
	} catch (error) {
		console.error('Error deleting order:', error);
		return { error: 'Error de conexión al eliminar la orden' };
	}
};

export const updateOrder = async (
	event: RequestEvent,
	orderId: string,
	updateData: { assigned_to?: string | null; address?: string; status?: string }
) => {
	try {
		console.log('Updating order with ID:', orderId, 'and data:', updateData);
		const response = await fetchAuth(
			`/v1/orders/${orderId}`,
			{
				method: 'PUT',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(updateData),
				credentials: 'include'
			},
			event
		);

		if (!response.ok) {
			const errorData = await response.json().catch(() => ({ error: 'Failed to update order.' }));
			console.log('Error updating order:', errorData);
			return { error: errorData.error || 'Error al actualizar la orden' };
		}

		const updatedOrder = await response.json();
		return updatedOrder;
	} catch (error) {
		console.error('Error updating order:', error);
		return { error: 'Error de conexión al actualizar la orden' };
	}
};

export const deleteOrderItem = async (event: RequestEvent, orderId: string, itemId: string) => {
	try {
		const response = await fetchAuth(
			`/v1/orders/${orderId}/items/${itemId}`,
			{
				method: 'DELETE',
				credentials: 'include'
			},
			event
		);

		if (!response.ok) {
			const errorData = await response
				.json()
				.catch(() => ({ error: 'Failed to delete order item.' }));
			console.log('Error deleting order item:', errorData);
			return { error: errorData.error || 'Error al eliminar el item' };
		}

		return { success: true };
	} catch (error) {
		console.error('Error deleting order item:', error);
		return { error: 'Error de conexión al eliminar el item' };
	}
};
