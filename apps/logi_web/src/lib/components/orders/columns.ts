import type { ColumnDef } from '@tanstack/table-core';
import { renderComponent } from '../ui/data-table';
import StatusBadge from './status-badge.svelte';
import DataTableActions from './data-table-actions.svelte';
import DataTableSortableHeader from '../users/data-table-sortable-header.svelte';

export type Status = 'pending' | 'in_progress' | 'completed' | 'cancelled';

export type Order = {
	order_id: string;
	order_number: string;
	created_by: string;
	created_by_username?: string;
	assigned_to: string;
	assigned_to_username?: string;
	order_name: string;
	order_phone_number: string;
	order_email?: string;
	order_cedula?: string;
	delivery_address: string;
	status: Status;
	created_at: string;
	updated_at: string;
};

export type OrderItem = {
	item_id: string;
	order_id: string;
	product_name: string;
	quantity: number;
	created_at: string;
	updated_at: string;
};

export const orderColumns: ColumnDef<Order>[] = [
	{
		accessorKey: 'order_id',
		header: 'ID',
		cell: ({ getValue }) => {
			const id = getValue<string>();
			return id.substring(0, 8) + '...';
		}
	},
	{
		accessorKey: 'order_number',
		header: ({ column }) =>
			renderComponent(DataTableSortableHeader, {
				onclick: column.getToggleSortingHandler(),
				sortDirection: column.getIsSorted(),
				text: 'Número de Orden'
			}),
		enableSorting: true,
		sortingFn: 'alphanumeric'
	},
	{
		accessorKey: 'created_by_username',
		header: 'Creado por'
	},
	{
		accessorKey: 'assigned_to_username',
		header: 'Asignado a',
		cell: ({ getValue }) => {
			const assignedTo = getValue<string>();
			return assignedTo || 'Sin asignar';
		}
	},
	{
		accessorKey: 'delivery_address',
		header: 'Dirección de Entrega',
		cell: ({ getValue }) => {
			const address = getValue<string>();
			return address.length > 50 ? address.substring(0, 50) + '...' : address;
		}
	},
	{
		accessorKey: 'status',
		header: ({ column }) =>
			renderComponent(DataTableSortableHeader, {
				onclick: column.getToggleSortingHandler(),
				sortDirection: column.getIsSorted(),
				text: 'Estado'
			}),
		enableSorting: true,
		sortingFn: 'alphanumeric',
		cell: ({ getValue }) => {
			const status = getValue<Status>();
			return renderComponent(StatusBadge, { status });
		}
	},
	{
		accessorKey: 'created_at',
		header: 'Creado el',
		cell: ({ getValue }) => {
			const date = new Date(getValue<string>());
			return date.toLocaleDateString('es-ES', {
				year: 'numeric',
				month: 'short',
				day: 'numeric',
				hour: '2-digit',
				minute: '2-digit'
			});
		}
	},
	{
		id: 'actions',
		header: 'Acciones',
		cell: ({ row }) => {
			const { order_id, order_number } = row.original;
			return renderComponent(DataTableActions, {
				id: order_id,
				orderNumber: order_number
			});
		}
	}
];
