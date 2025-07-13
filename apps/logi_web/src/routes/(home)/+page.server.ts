import type { PageServerLoad, Actions } from './$types';

export const load: PageServerLoad = async ({ locals }) => {
	// Pass the user data to the page component
	return {
		loggedIn: !!locals.user
	};
};
