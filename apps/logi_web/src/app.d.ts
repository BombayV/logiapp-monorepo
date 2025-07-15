// See https://svelte.dev/docs/kit/types#app.d.ts
export type UserData = {
	email: string;
	role: 'admin' | 'sales' | 'driver';
	user_id: string;
	profile: {
		first_name: string;
		last_name: string;
		phone_number: string;
		last_connection: string;
	};
} | null;

// for information about these interfaces
declare global {
	namespace App {
		// interface Error {}
		interface Locals {
			user: UserData;
		}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}
}

export {};
