<script lang="ts">
	import EllipsisIcon from '@lucide/svelte/icons/ellipsis';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { toast } from 'svelte-sonner';
	import { enhance } from '$app/forms';

	let { id, orderNumber }: { id: string; orderNumber?: string } = $props();
	let showCancelDialog = $state(false);

	const copyOrderId = async () => {
		try {
			await navigator.clipboard.writeText(id);
			toast.success('ID de orden copiado al portapapeles');
		} catch (error) {
			toast.error('Error al copiar ID de orden');
		}
	};

	const copyOrderNumber = async () => {
		try {
			const numberToCopy = orderNumber || id;
			await navigator.clipboard.writeText(numberToCopy);
			toast.success('Número de orden copiado al portapapeles');
		} catch (error) {
			toast.error('Error al copiar número de orden');
		}
	};

	const handleCancelOrder = () => {
		showCancelDialog = true;
	};

	const submitCancelForm = () => {
		const form = document.getElementById(`cancel-order-form-${id}`) as HTMLFormElement;
		if (form) {
			form.submit();
		}
		showCancelDialog = false;
	};
</script>

<AlertDialog.Root bind:open={showCancelDialog}>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Confirmar cancelación</AlertDialog.Title>
			<AlertDialog.Description>
				¿Estás seguro de que deseas cancelar este pedido? Esta acción no se puede deshacer.
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel>Cancelar</AlertDialog.Cancel>
			<AlertDialog.Action onclick={submitCancelForm}>Eliminar</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>

<DropdownMenu.Root>
	<DropdownMenu.Trigger>
		{#snippet child({ props })}
			<Button {...props} variant="ghost" size="icon" class="relative size-8 p-0">
				<span class="sr-only">Open menu</span>
				<EllipsisIcon />
			</Button>
		{/snippet}
	</DropdownMenu.Trigger>
	<DropdownMenu.Content>
		<DropdownMenu.Group>
			<DropdownMenu.Label>Acciones</DropdownMenu.Label>
			<DropdownMenu.Item onclick={copyOrderId}>Copiar ID de orden</DropdownMenu.Item>
			<DropdownMenu.Item onclick={copyOrderNumber}>Copiar numero de orden</DropdownMenu.Item>
		</DropdownMenu.Group>
		<DropdownMenu.Separator />
		<a href={`/auth/dashboard/orders/${id}`}>
			<DropdownMenu.Item>Ver detalles</DropdownMenu.Item>
		</a>
		<DropdownMenu.Separator />
		<DropdownMenu.Item variant="destructive" onclick={handleCancelOrder}
			>Cancelar orden</DropdownMenu.Item
		>
	</DropdownMenu.Content>
</DropdownMenu.Root>

<form
	id={`cancel-order-form-${id}`}
	method="POST"
	action={`?/cancel_order`}
	use:enhance={({ formElement }) => {
		return async ({ result, update }) => {
			if (result.type === 'success') {
				toast.success('Orden cancelada correctamente');
				await update();
			} else if (result.type === 'failure') {
				toast.error('Error al cancelar la orden');
			}
		};
	}}
>
	<input type="hidden" name="order_id" value={id} />
</form>
