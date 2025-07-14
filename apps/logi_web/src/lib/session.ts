import { writable } from 'svelte/store';

// Define the shape of the user object
interface User {
	role?: string;
	token?: string;
}

// Define the shape of the session
interface Session {
	user: User | null;
}

// Create a writable store with an initial state
export const session = writable<Session>({
	user: null
});

// Function to update the session
export const updateSession = (userData: User | null) => {
	session.set({ user: userData });
};
