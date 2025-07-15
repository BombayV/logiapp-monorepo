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

export const createOrder = async (
	event: RequestEvent,
	orderData: { email: string; address: string; order_number: string }
) => {
	try {
		const response = await fetchAuth(
			`/v1/orders`,
			{
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(orderData),
				credentials: 'include' // Ensure cookies are sent with the request
			},
			event
		);

		if (!response.ok) {
			const errorData = await response.json().catch(() => ({ message: 'Failed to create order.' }));
			console.log('Error creating order:', errorData);
			return null;
		}

		const newOrder = await response.json();
		return newOrder;
	} catch (error) {
		console.error('Error creating order:', error);
		return null;
	}
};
