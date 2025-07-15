import type { ColumnDef } from '@tanstack/table-core';
import { renderComponent } from '../ui/data-table';
import StatusBadge from './status-badge.svelte';
import DataTableActions from './data-table-actions.svelte';

export type Status = 'pending' | 'in_progress' | 'completed' | 'cancelled';

export type Order = {
	order_id: string;
	created_by: string;
	assigned_to: string;
	address: string;
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
		accessorKey: 'created_by',
		header: 'Creado por'
	},
	{
		accessorKey: 'assigned_to',
		header: 'Asignado a',
		cell: ({ getValue }) => {
			const assignedTo = getValue<string>();
			return assignedTo || 'Sin asignar';
		}
	},
	{
		accessorKey: 'address',
		header: 'Dirección',
		cell: ({ getValue }) => {
			const address = getValue<string>();
			return address.length > 50 ? address.substring(0, 50) + '...' : address;
		}
	},
	{
		accessorKey: 'status',
		header: 'Estado',
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
			const { order_id } = row.original;
			return renderComponent(DataTableActions, { id: order_id });
		}
	}
];
