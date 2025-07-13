import { writable } from 'svelte/store';

// Define the shape of the user object
interface User {
	username: string;
	role: 'sales' | 'admin' | 'driver' | null;
}

// Define the shape of the session
interface Session {
	user: User | null;
}

// Create a writable store with an initial state
export const session = writable<Session>({
	user: null
});
