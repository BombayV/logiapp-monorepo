import type { ColumnDef } from "@tanstack/table-core";
import { renderComponent } from "../ui/data-table";
import { Mail } from '@lucide/svelte'
import DataTableEmailButton from "./data-table-email-button.svelte";
import DataTableActions from "./data-table-actions.svelte";

export type User = {
  user_id: string;
  email: string;
  role: 'admin' | 'sales' | 'driver';
  created_at: string;
  updated_at: string;
}

export const columns: ColumnDef<User>[] = [
  {
    accessorKey: 'user_id',
    header: 'ID',
  },
  {
    accessorKey: 'email',
    header: ({ column }) => (
      renderComponent(DataTableEmailButton, {
        onclick: column.getToggleSortingHandler(),
      })
    )
  },
  {
    accessorKey: 'role',
    header: 'Rol',
  },
  {
    accessorKey: 'created_at',
    header: 'Creado el',
    cell: ({ getValue }) => {
      const date = new Date(getValue<string>());
      return `${date.toLocaleString()}`;
    },
  },
  {
    accessorKey: 'updated_at',
    header: 'Actualizado el',
    cell: ({ getValue }) => {
      const date = new Date(getValue<string>());
      return `${date.toLocaleString()}`;
    },
  },
  {
    id: 'actions',
    header: 'Acciones',
    cell: ({ row }) => {
      const { user_id } = row.original;
      return renderComponent(DataTableActions, { id: user_id });
    },
  }
]