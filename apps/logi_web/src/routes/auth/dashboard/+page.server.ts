import { fail, redirect } from '@sveltejs/kit';
import type { PageServerLoad, Actions } from './$types';
import type { RequestEvent } from '@sveltejs/kit';
import { createUser, getAllUsers, logoutUser, deleteUser } from '@/server/users';
import type { User } from '@/components/users/columns';

export const load: PageServerLoad = async (event) => {
	const { locals } = event;
	if (!locals.user) {
		throw redirect(303, '/auth/login');
	}

	if (locals.user.role !== 'admin') {
		return {
			user: locals.user
		};
	}

	const users = (await getAllUsers(event)) as User[];
	if (!users) {
		return fail(500, { error: 'No se pudieron cargar los datos.' });
	}

	return {
		user: locals.user,
		users: users
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

export const actions: Actions = {
	logout,
	delete_user,
	create_user
};
