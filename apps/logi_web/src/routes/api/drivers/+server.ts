import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { fetchAuth } from '@/fetchAuth';

export const GET: RequestHandler = async (event) => {
	try {
		const response = await fetchAuth(
			'/v1/users/drivers',
			{
				method: 'GET',
				headers: {
					'Content-Type': 'application/json'
				}
			},
			event
		);

		if (!response.ok) {
			return json({ error: 'Failed to fetch drivers' }, { status: response.status });
		}

		const data = await response.json();
		return json(data);
	} catch (error) {
		console.error('Error fetching drivers:', error);
		return json({ error: 'Internal server error' }, { status: 500 });
	}
};
