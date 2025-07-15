import type { ColumnDef } from '@tanstack/table-core';
import { renderComponent } from '../ui/data-table';
import { Mail } from '@lucide/svelte';
import DataTableEmailButton from './data-table-email-button.svelte';
import DataTableActions from './data-table-actions.svelte';

export type User = {
	user_id: string;
	email: string;
	first_name?: string;
	last_name?: string;
	role: 'admin' | 'sales' | 'driver';
	phone_number?: string;
	last_connection?: string;
	created_at: string;
	updated_at: string;
};

export const columns: ColumnDef<User>[] = [
	{
		accessorKey: 'user_id',
		header: 'ID'
	},
	{
		accessorKey: 'email',
		header: ({ column }) =>
			renderComponent(DataTableEmailButton, {
				onclick: column.getToggleSortingHandler()
			})
	},
	{
		accessorKey: 'first_name',
		header: 'Nombre'
	},
	{
		accessorKey: 'last_name',
		header: 'Apellido'
	},
	{
		accessorKey: 'role',
		header: 'Rol'
	},
	{
		accessorKey: 'phone_number',
		header: 'Teléfono',
		cell: ({ getValue }) => {
			const phone = getValue<string>();
			return phone || 'No disponible';
		}
	},
	{
		accessorKey: 'last_connection',
		header: 'Última conexión',
		cell: ({ getValue }) => {
			const date = new Date(getValue<string>());
			return `${date.toLocaleString()}`;
		}
	},
	{
		id: 'actions',
		header: 'Acciones',
		cell: ({ row }) => {
			const user = row.original;
			return renderComponent(DataTableActions, { user });
		}
	}
];
