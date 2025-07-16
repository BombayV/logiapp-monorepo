<script lang="ts">
	import EllipsisIcon from '@lucide/svelte/icons/ellipsis';
	import { Button } from '$lib/components/ui/button/index.js';
	import * as DropdownMenu from '$lib/components/ui/dropdown-menu/index.js';
	import * as AlertDialog from '$lib/components/ui/alert-dialog/index.js';
	import { toast } from 'svelte-sonner';
	import { enhance } from '$app/forms';
	import type { User } from './columns.js';

	let { user }: { user: User } = $props();
	let showDeleteDialog = $state(false);

	const copyEmail = async () => {
		try {
			await navigator.clipboard.writeText(user.email);
			toast.success('Email copiado al portapapeles');
		} catch (error) {
			toast.error('Error al copiar el email');
		}
	};

	const copyPhone = async () => {
		if (user.phone_number) {
			try {
				await navigator.clipboard.writeText(user.phone_number);
				toast.success('Teléfono copiado al portapapeles');
			} catch (error) {
				toast.error('Error al copiar el teléfono');
			}
		} else {
			toast.error('No hay teléfono disponible para este usuario');
		}
	};

	const handleDeleteUser = () => {
		showDeleteDialog = true;
	};

	const submitDeleteForm = () => {
		const form = document.getElementById(`delete-user-form-${user.user_id}`) as HTMLFormElement;
		if (form) {
			form.submit();
		}
		showDeleteDialog = false;
	};
</script>

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
			<DropdownMenu.Item onclick={copyEmail}>Copiar email</DropdownMenu.Item>
			<DropdownMenu.Item onclick={copyPhone}>Copiar teléfono</DropdownMenu.Item>
		</DropdownMenu.Group>
		{#if user.role !== 'admin'}
			<DropdownMenu.Separator />
			<DropdownMenu.Item variant="destructive" onclick={handleDeleteUser}>
				Eliminar usuario
			</DropdownMenu.Item>
		{/if}
	</DropdownMenu.Content>
</DropdownMenu.Root>

<AlertDialog.Root bind:open={showDeleteDialog}>
	<AlertDialog.Content>
		<AlertDialog.Header>
			<AlertDialog.Title>Confirmar eliminación</AlertDialog.Title>
			<AlertDialog.Description>
				¿Estás seguro de que deseas eliminar este usuario? Esta acción no se puede deshacer.
			</AlertDialog.Description>
		</AlertDialog.Header>
		<AlertDialog.Footer>
			<AlertDialog.Cancel>Cancelar</AlertDialog.Cancel>
			<AlertDialog.Action onclick={submitDeleteForm}>Eliminar</AlertDialog.Action>
		</AlertDialog.Footer>
	</AlertDialog.Content>
</AlertDialog.Root>

<!-- Hidden form for user deletion -->
<form
	id="delete-user-form-{user.user_id}"
	method="POST"
	action="?/delete_user"
	style="display: none;"
	use:enhance={({ formElement }) => {
		return async ({ result, update }) => {
			if (result.type === 'success') {
				toast.success('Usuario eliminado correctamente');
				await update();
			} else if (result.type === 'failure') {
				toast.error('Error al eliminar el usuario');
			}
		};
	}}
>
	<input type="hidden" name="user_id" value={user.user_id} />
</form>
