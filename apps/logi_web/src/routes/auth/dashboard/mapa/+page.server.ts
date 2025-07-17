import { error } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch }) => {
	try {
		// Fetch active drivers with location data from our API route
		const response = await fetch('/api/drivers');

		if (!response.ok) {
			throw error(response.status, 'Failed to fetch drivers data');
		}

		const data = await response.json();

		return {
			drivers: data.drivers || [],
			count: data.count || 0
		};
	} catch (err) {
		console.error('Error fetching drivers:', err);
		throw error(500, 'Failed to load drivers data');
	}
};
