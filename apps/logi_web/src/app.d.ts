// See https://svelte.dev/docs/kit/types#app.d.ts
// for information about these interfaces
declare global {
	namespace App {
		// interface Error {}
		interface Locals {
			user?: {
				profile?: {
					first_name?: string;
					last_name?: string;
					phone_number?: string;
				};
				role?: string;
				user?: string;
			} | null; // User object or null if not authenticated
		}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}
}

export {};
