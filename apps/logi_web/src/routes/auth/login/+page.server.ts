import { fail, redirect } from '@sveltejs/kit';
import type { Actions } from './$types';
import { BACKEND_URL } from '$env/static/private';
import { dev } from '$app/environment';

// The base URL for your Go backend. Use an environment variable in a real app.
export const actions: Actions = {
	login: async ({ cookies, request, fetch }) => {
		const data = await request.formData();
		const email = data.get('email');
		const password = data.get('password');

		if (!email || !password) {
			return fail(400, { error: 'Email y contraseña son obligatorios.' });
		}

		try {
			// Proxy the login request to the Go backend
			const response = await fetch(`${BACKEND_URL}/v1/users/login`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({ email, password }) // Adjust body to match Go backend expectations
			});
			console.log(BACKEND_URL);

			if (!response.ok) {
				const errorData = await response
					.json()
					.catch(() => ({ message: 'Invalid credentials or server error.' }));
				return fail(response.status, { error: errorData.error || 'Credenciales inválidas.' });
			}

			// Assuming the Go backend returns a JSON object with a "token" field
			const { token } = await response.json();

			if (!token) {
				return fail(500, { error: 'Inicio de sesión fallido, token no recibido.' });
			}

			// Set the token returned by Go in a secure, httpOnly cookie
			cookies.set('session', token, {
				path: '/',
				httpOnly: true,
				sameSite: 'strict',
				secure: !dev, // Use secure cookies in production
				maxAge: 60 * 60 * 24 // 24 hours (or match your JWT expiry)
			});
		} catch (error) {
			console.error('Error connecting to auth service:', error);
			return fail(500, { error: 'No se pudo conectar con el servicio de autenticación.' });
		}

		// Redirect to the dashboard on successful login
		throw redirect(303, '/auth/dashboard');
	}
};
