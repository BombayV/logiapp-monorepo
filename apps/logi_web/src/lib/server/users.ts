import { fetchAuth } from "@/fetchAuth";
import { type RequestEvent } from "@sveltejs/kit";

export const getAllUsers = async (event: RequestEvent) => {
  try {
    const response = await fetchAuth(
      `/v1/users`,
      {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json'
        },
        credentials: 'include' // Ensure cookies are sent with the request
      },
      event
    );

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ message: 'Failed to fetch users.' }));
      console.log('Error fetching users:', errorData);
      return [];
    }

    const users = await response.json();
    return users.users || [];
  } catch (error) {
    console.error('Error fetching users:', error);
    return [];
  }
}

export const createUser = async (event: RequestEvent, userData: { email: string; password: string; role: 'sales' | 'driver', first_name: string; last_name: string, phone: string }) => {
  try {
    const response = await fetchAuth(
      `/v1/users`,
      {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(userData),
        credentials: 'include' // Ensure cookies are sent with the request
      },
      event
    );

    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ message: 'Failed to create user.' }));
      console.log('Error creating user:', errorData);
      return null;
    }

    const newUser = await response.json();
    return newUser;
  } catch (error) {
    console.error('Error creating user:', error);
    return null;
  }
}