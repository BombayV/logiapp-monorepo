import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ locals, cookies }) => {
	const token = cookies.get('session');
	const user = locals.user ? { ...locals.user, token } : null;
	return {
		user
	};
};
