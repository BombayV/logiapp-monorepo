// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
declare global {
	namespace App {
		// interface Error {}
		interface Locals {
			user?: {
				email: string;
				role: 'admin' | 'sales' | 'driver';
				user_id: string;
				profile: {
					first_name: string;
					last_name: string;
					phone_number: string;
					last_connection: string;
				}
			} | null; // User object or null if not authenticated
		}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}
}

export {};
