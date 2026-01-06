import { fetchAuth } from '$lib/fetchAuth';
import type { RequestEvent } from '@sveltejs/kit';

export interface OrderForm {
	form_id: string;
	public_id: string;
	order_id: string;
	driver_id?: string;
	driver_rating?: number;
	cargo_condition?: string;
	comments?: string;
	is_finished: boolean;
	driver_name?: string;
	driver_email?: string;
	created_at: string;
	updated_at: string;
}

export async function getFormByPublicId(
	event: RequestEvent,
	publicId: string
): Promise<OrderForm | null> {
	try {
		const response = await fetchAuth(
			`/v1/public/forms/${publicId}`,
			{
				method: 'GET',
				headers: {
					'Content-Type': 'application/json'
				}
			},
			event
		);

		if (!response.ok) {
			if (response.status === 404) {
				return null;
			}
			throw new Error(`Failed to fetch form: ${response.statusText}`);
		}

		const form = await response.json();
		return form;
	} catch (error) {
		console.error('Error fetching form:', error);
		return null;
	}
}

export async function submitForm(
	event: RequestEvent,
	publicId: string,
	data: { driver_rating: number; cargo_condition: string; comments: string }
) {
	try {
		const response = await fetchAuth(
			`/v1/public/forms/${publicId}/submit`,
			{
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(data)
			},
			event
		);

		if (!response.ok) {
			const error = await response.json().catch(() => ({ error: 'Failed to submit form' }));
			return { error: error.error || 'Failed to submit form' };
		}

		return { success: true };
	} catch (error) {
		console.error('Error submitting form:', error);
		return { error: 'Network error' };
	}
}
