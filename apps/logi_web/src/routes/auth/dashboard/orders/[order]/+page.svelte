<script lang="ts">
	import type { PageData } from './$types';
	import { Button } from '@/components/ui/button';
	import * as Card from '@/components/ui/card';
	import * as Dialog from '@/components/ui/dialog';
	import * as Select from '@/components/ui/select';
	import { Label } from '@/components/ui/label';
	import { Input } from '@/components/ui/input';
	import { Textarea } from '@/components/ui/textarea';
	import {
		ArrowLeft,
		Copy,
		Edit,
		Trash2,
		MapPin,
		User,
		Calendar,
		Package,
		Plus
	} from '@lucide/svelte';
	import { toast } from 'svelte-sonner';
	import { goto } from '$app/navigation';
	import { enhance } from '$app/forms';
	import { invalidateAll } from '$app/navigation';
	import StatusBadge from '$lib/components/orders/status-badge.svelte';

	let { data }: { data: PageData } = $props();

	const order = $derived(data.order);
	const user = $derived(data.user);
	const items = $derived(data.items);

	let addItemDialogOpen = $state(false);
	let bulkAddDialogOpen = $state(false);
	let deleteConfirmDialogOpen = $state(false);
	let editOrderDialogOpen = $state(false);
	let deleteItemDialogOpen = $state(false);
	let isSubmitting = $state(false);
	let isDeleting = $state(false);
	let isEditing = $state(false);
	let isDeletingItem = $state(false);
	let selectedItemId = $state('');
	let selectedItemName = $state('');

	// Bulk items state
	let bulkItems = $state([{ product_name: '', quantity: 1 }]);

	// Edit order state
	let editData = $state({
		assigned_to: '',
		delivery_address: '',
		status: 'pending'
	});

	// Update editData when order changes
	$effect(() => {
		editData.assigned_to = order.assigned_to || '';
		editData.delivery_address = order.delivery_address || '';
		editData.status = order.status || 'pending';
	});

	// Drivers list for assignment
	let drivers: Array<{
		user_id: string;
		email: string;
		first_name: string;
		last_name: string;
		phone_number: string;
		role: string;
		last_connection: string;
		created_at: string;
		updated_at: string;
	}> = $state([]);

	const copyOrderId = async () => {
		try {
			await navigator.clipboard.writeText(order.order_id);
			toast.success('ID de orden copiado al portapapeles');
		} catch (error) {
			toast.error('Error al copiar ID');
		}
	};

	const copyOrderNumber = async () => {
		try {
			await navigator.clipboard.writeText(order.order_number);
			toast.success('Número de orden copiado al portapapeles');
		} catch (error) {
			toast.error('Error al copiar número de orden');
		}
	};

	const goBack = () => {
		goto('/auth/dashboard/orders');
	};

	const editOrder = async () => {
		// Fetch drivers
		try {
			const response = await fetch('/api/drivers');
			if (response.ok) {
				const data = await response.json();
				drivers = data.drivers || [];
			}
		} catch (error) {
			console.error('Error fetching drivers:', error);
			toast.error('Error al cargar conductores');
		}

		editOrderDialogOpen = true;
	};

	const deleteOrder = () => {
		deleteConfirmDialogOpen = true;
	};

	const deleteItem = (itemId: string, itemName: string) => {
		selectedItemId = itemId;
		selectedItemName = itemName;
		deleteItemDialogOpen = true;
	};

	const addItem = () => {
		addItemDialogOpen = true;
	};

	const addBulkItem = () => {
		bulkItems.push({ product_name: '', quantity: 1 });
	};

	const removeBulkItem = (index: number) => {
		if (bulkItems.length > 1) {
			bulkItems.splice(index, 1);
		}
	};

	const resetBulkItems = () => {
		bulkItems = [{ product_name: '', quantity: 1 }];
	};

	const formatDate = (dateString: string) => {
		return new Date(dateString).toLocaleDateString('es-ES', {
			year: 'numeric',
			month: 'long',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	};
</script>

<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 mt-6">
	<!-- Header -->
	<Button variant="outline" size="sm" class="mb-4" onclick={goBack}>
		<ArrowLeft class="h-4 w-4 mr-2" />
		Volver a Órdenes
	</Button>
	<div class="flex flex-col sm:flex-row sm:items-center justify-between mb-6 gap-4">
		<div class="flex items-center gap-4">
			<div>
				<h1 class="text-2xl font-bold">Orden #{order.order_number}</h1>
				<p class="text-muted-foreground">ID: {order.order_id.substring(0, 8)}...</p>
			</div>
		</div>
		<div class="flex items-center gap-2 flex-wrap">
			<StatusBadge status={order.status} />
			<Button variant="outline" size="sm" onclick={editOrder}>
				<Edit class="h-4 w-4 mr-2" />
				<span class="hidden sm:inline">Editar</span>
				<span class="sm:hidden">Editar</span>
			</Button>
			<Button variant="destructive" size="sm" onclick={deleteOrder}>
				<Trash2 class="h-4 w-4 mr-2" />
				<span class="hidden sm:inline">Eliminar</span>
				<span class="sm:hidden">Eliminar</span>
			</Button>
		</div>
	</div>

	<!-- Order Details Grid -->
	<div class="grid grid-cols-1 lg:grid-cols-2 gap-4 sm:gap-6">
		<!-- Order Information -->
		<Card.Root>
			<Card.Header>
				<Card.Title class="flex items-center gap-2">
					<Package class="h-5 w-5" />
					Información de la Orden
				</Card.Title>
			</Card.Header>
			<Card.Content class="space-y-4">
				<div class="flex items-center justify-between">
					<span class="text-sm font-medium">Número de Orden:</span>
					<div class="flex items-center gap-2">
						<span class="font-mono">{order.order_number}</span>
						<Button variant="ghost" size="sm" onclick={copyOrderNumber}>
							<Copy class="h-4 w-4" />
						</Button>
					</div>
				</div>
				<div class="flex items-center justify-between">
					<span class="text-sm font-medium">ID de Orden:</span>
					<div class="flex items-center gap-2">
						<span class="font-mono text-xs">{order.order_id.substring(0, 8)}...</span>
						<Button variant="ghost" size="sm" onclick={copyOrderId}>
							<Copy class="h-4 w-4" />
						</Button>
					</div>
				</div>
				<div class="flex items-center justify-between">
					<span class="text-sm font-medium">Estado:</span>
					<StatusBadge status={order.status} />
				</div>
			</Card.Content>
		</Card.Root>

		<!-- Assignment Information -->
		<Card.Root>
			<Card.Header>
				<Card.Title class="flex items-center gap-2">
					<User class="h-5 w-5" />
					Asignación
				</Card.Title>
			</Card.Header>
			<Card.Content class="space-y-4">
				<div class="flex items-center justify-between">
					<span class="text-sm font-medium">Creado por:</span>
					<span>{order.created_by_username || 'No disponible'}</span>
				</div>
				<div class="flex items-center justify-between">
					<span class="text-sm font-medium">Asignado a:</span>
					<span>{order.assigned_to_username || 'Sin asignar'}</span>
				</div>
			</Card.Content>
		</Card.Root>

		<!-- Delivery Information -->
		<Card.Root>
			<Card.Header>
				<Card.Title class="flex items-center gap-2">
					<MapPin class="h-5 w-5" />
					Información de Entrega
				</Card.Title>
			</Card.Header>
			<Card.Content>
				<div class="space-y-2">
					<span class="text-sm font-medium">Dirección de Entrega:</span>
					<p class="text-sm bg-muted p-3 rounded-md">{order.delivery_address}</p>
				</div>
			</Card.Content>
		</Card.Root>

		<!-- Timeline -->
		<Card.Root>
			<Card.Header>
				<Card.Title class="flex items-center gap-2">
					<Calendar class="h-5 w-5" />
					Cronología
				</Card.Title>
			</Card.Header>
			<Card.Content class="space-y-4">
				<div class="flex items-center justify-between">
					<span class="text-sm font-medium">Creado el:</span>
					<span class="text-sm">{formatDate(order.created_at)}</span>
				</div>
				<div class="flex items-center justify-between">
					<span class="text-sm font-medium">Última actualización:</span>
					<span class="text-sm">{formatDate(order.updated_at)}</span>
				</div>
			</Card.Content>
		</Card.Root>
	</div>

	<!-- Order Items Section -->
	<Card.Root class="mt-4 sm:mt-6">
		<Card.Header>
			<div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
				<Card.Title class="flex items-center gap-2">
					<Package class="h-5 w-5" />
					Items de la Orden
					<span class="text-sm font-normal text-muted-foreground">
						({items.length}
						{items.length === 1 ? 'item' : 'items'})
					</span>
				</Card.Title>
				<div class="flex items-center gap-2">
					<Button size="sm" onclick={addItem}>
						<Plus class="h-4 w-4 mr-2" />
						<span class="hidden sm:inline">Agregar Item</span>
						<span class="sm:hidden">Agregar</span>
					</Button>
					<Button size="sm" variant="outline" onclick={() => (bulkAddDialogOpen = true)}>
						<Plus class="h-4 w-4 mr-2" />
						<span class="hidden sm:inline">Agregar Múltiples</span>
						<span class="sm:hidden">Múltiples</span>
					</Button>
				</div>
			</div>
		</Card.Header>
		<Card.Content>
			{#if items.length > 0}
				<div class="space-y-4">
					{#each items as item, index (item.item_id)}
						<div
							class="flex flex-col sm:flex-row sm:items-center justify-between p-4 border rounded-lg gap-4"
						>
							<div class="flex items-center gap-4">
								<div
									class="w-8 h-8 bg-primary/10 rounded-full flex items-center justify-center flex-shrink-0"
								>
									<span class="text-sm font-medium">{index + 1}</span>
								</div>
								<div class="min-w-0 flex-1">
									<h4 class="font-medium truncate">{item.product_name}</h4>
									<p class="text-sm text-muted-foreground">
										Item ID: {item.item_id.substring(0, 8)}...
									</p>
								</div>
							</div>
							<div class="flex items-center justify-between sm:justify-end gap-4">
								<div class="flex items-center gap-2">
									<span class="text-sm text-muted-foreground">Cantidad:</span>
									<span class="font-medium">{item.quantity}</span>
								</div>
								<div class="flex flex-col items-end gap-2">
									<p class="text-xs text-muted-foreground">
										Agregado: {formatDate(item.created_at)}
									</p>
									<Button
										variant="ghost"
										size="sm"
										onclick={() => deleteItem(item.item_id, item.product_name)}
										class="text-red-600 hover:text-red-700 hover:bg-red-50"
									>
										<Trash2 class="h-4 w-4" />
									</Button>
								</div>
							</div>
						</div>
					{/each}
				</div>
			{:else}
				<div class="text-center py-8">
					<Package class="h-12 w-12 text-muted-foreground mx-auto mb-4" />
					<p class="text-muted-foreground">Esta orden no tiene items asociados.</p>
				</div>
			{/if}
		</Card.Content>
	</Card.Root>
</div>

<!-- Add Item Dialog -->
<Dialog.Root bind:open={addItemDialogOpen}>
	<Dialog.Content class="sm:max-w-[425px]">
		<Dialog.Header>
			<Dialog.Title>Agregar Item a la Orden</Dialog.Title>
			<Dialog.Description>
				Agrega un nuevo item a la orden #{order.order_number}
			</Dialog.Description>
		</Dialog.Header>
		<form
			method="POST"
			action="?/add_item"
			use:enhance={({ formElement }) => {
				isSubmitting = true;
				return async ({ result, update }) => {
					isSubmitting = false;
					if (result.type === 'success') {
						toast.success('Item agregado exitosamente');
						addItemDialogOpen = false;
						formElement.reset();
						await update();
						await invalidateAll();
					} else if (result.type === 'failure') {
						const errorMessage = (result.data as any)?.error || 'Error al agregar el item';
						toast.error(String(errorMessage));
					}
				};
			}}
		>
			<div class="grid gap-4 py-4">
				<div class="grid grid-cols-4 items-center gap-4">
					<Label for="product_name" class="text-right">Producto</Label>
					<Input
						id="product_name"
						name="product_name"
						placeholder="Nombre del producto"
						class="col-span-3"
						required
						disabled={isSubmitting}
					/>
				</div>
				<div class="grid grid-cols-4 items-center gap-4">
					<Label for="quantity" class="text-right">Cantidad</Label>
					<Input
						id="quantity"
						name="quantity"
						type="number"
						placeholder="1"
						min="1"
						class="col-span-3"
						required
						disabled={isSubmitting}
					/>
				</div>
			</div>
			<Dialog.Footer>
				<Button
					type="button"
					variant="outline"
					onclick={() => (addItemDialogOpen = false)}
					disabled={isSubmitting}
				>
					Cancelar
				</Button>
				<Button type="submit" disabled={isSubmitting}>
					{isSubmitting ? 'Agregando...' : 'Agregar'}
				</Button>
			</Dialog.Footer>
		</form>
	</Dialog.Content>
</Dialog.Root>

<!-- Bulk Add Items Dialog -->
<Dialog.Root bind:open={bulkAddDialogOpen}>
	<Dialog.Content class="sm:max-w-[600px]">
		<Dialog.Header>
			<Dialog.Title>Agregar Múltiples Items</Dialog.Title>
			<Dialog.Description>Agrega varios items a la orden de una sola vez.</Dialog.Description>
		</Dialog.Header>
		<form
			method="POST"
			action="?/bulk_add_items"
			use:enhance={() => {
				isSubmitting = true;
				return async ({ result, update }) => {
					isSubmitting = false;
					if (result.type === 'success') {
						toast.success(String((result.data as any)?.message || 'Items agregados exitosamente'));
						bulkAddDialogOpen = false;
						resetBulkItems();
						await update();
						await invalidateAll();
					} else if (result.type === 'failure') {
						toast.error(String((result.data as any)?.error || 'Error al agregar items'));
					}
				};
			}}
		>
			<div class="space-y-4 max-h-[400px] overflow-y-auto">
				{#each bulkItems as item, index}
					<div class="flex items-center gap-2 p-3 border rounded-lg">
						<div class="flex-1 space-y-2">
							<Label for={`product_name_${index}`}>Producto</Label>
							<Input
								id={`product_name_${index}`}
								bind:value={item.product_name}
								placeholder="Nombre del producto"
								required
								disabled={isSubmitting}
							/>
						</div>
						<div class="w-24 space-y-2">
							<Label for={`quantity_${index}`}>Cantidad</Label>
							<Input
								id={`quantity_${index}`}
								bind:value={item.quantity}
								type="number"
								min="1"
								required
								disabled={isSubmitting}
							/>
						</div>
						<Button
							type="button"
							variant="outline"
							size="sm"
							onclick={() => removeBulkItem(index)}
							disabled={isSubmitting || bulkItems.length <= 1}
						>
							<Trash2 class="h-4 w-4" />
						</Button>
					</div>
				{/each}
			</div>

			<div class="flex justify-between items-center mt-4">
				<Button
					type="button"
					variant="outline"
					size="sm"
					onclick={addBulkItem}
					disabled={isSubmitting}
				>
					<Plus class="h-4 w-4 mr-2" />
					Agregar Más
				</Button>
				<div class="text-sm text-muted-foreground">
					{bulkItems.length}
					{bulkItems.length === 1 ? 'item' : 'items'}
				</div>
			</div>

			<input type="hidden" name="items" value={JSON.stringify(bulkItems)} />

			<Dialog.Footer class="mt-6">
				<Button
					type="button"
					variant="outline"
					onclick={() => {
						bulkAddDialogOpen = false;
						resetBulkItems();
					}}
					disabled={isSubmitting}
				>
					Cancelar
				</Button>
				<Button type="submit" disabled={isSubmitting}>
					{isSubmitting ? 'Agregando...' : `Agregar ${bulkItems.length}`}
				</Button>
			</Dialog.Footer>
		</form>
	</Dialog.Content>
</Dialog.Root>

<!-- Edit Order Dialog -->
<Dialog.Root bind:open={editOrderDialogOpen}>
	<Dialog.Content class="sm:max-w-[500px]">
		<Dialog.Header>
			<Dialog.Title>Editar Orden</Dialog.Title>
			<Dialog.Description>
				Modifica la información de la orden #{order.order_number}
			</Dialog.Description>
		</Dialog.Header>
		<form
			method="POST"
			action="?/update_order"
			use:enhance={() => {
				isEditing = true;
				return async ({ result, update }) => {
					isEditing = false;
					if (result.type === 'success') {
						toast.success(
							String((result.data as any)?.message || 'Orden actualizada exitosamente')
						);
						editOrderDialogOpen = false;
						await update();
						await invalidateAll();
					} else if (result.type === 'failure') {
						toast.error(String((result.data as any)?.error || 'Error al actualizar la orden'));
					}
				};
			}}
		>
			<div class="grid gap-4 py-4">
				<div class="grid grid-cols-4 items-center gap-4">
					<Label for="status" class="text-right">Estado</Label>
					<Select.Root type="single" bind:value={editData.status}>
						<Select.Trigger class="col-span-3">
							{editData.status === 'pending'
								? 'Pendiente'
								: editData.status === 'in_progress'
									? 'En Progreso'
									: editData.status === 'completed'
										? 'Completada'
										: editData.status === 'cancelled'
											? 'Cancelada'
											: 'Seleccionar estado'}
						</Select.Trigger>
						<Select.Content>
							<Select.Item value="pending">Pendiente</Select.Item>
							<Select.Item value="in_progress">En Progreso</Select.Item>
							<Select.Item value="completed">Completada</Select.Item>
							<Select.Item value="cancelled">Cancelada</Select.Item>
						</Select.Content>
					</Select.Root>
					<input type="hidden" name="status" value={editData.status} />
				</div>

				<div class="grid grid-cols-4 items-center gap-4">
					<Label for="assigned_to" class="text-right">Asignado a</Label>
					<Select.Root type="single" bind:value={editData.assigned_to}>
						<Select.Trigger class="col-span-3">
							{editData.assigned_to
								? (() => {
										const driver = drivers.find((d) => d.user_id === editData.assigned_to);
										return driver
											? `${driver.first_name} ${driver.last_name} - ${driver.email}`
											: 'Sin asignar';
									})()
								: 'Sin asignar'}
						</Select.Trigger>
						<Select.Content>
							<Select.Item value="">Sin asignar</Select.Item>
							{#each drivers as driver}
								<Select.Item value={driver.user_id}>
									{driver.first_name}
									{driver.last_name} - {driver.email}
								</Select.Item>
							{/each}
						</Select.Content>
					</Select.Root>
					<input type="hidden" name="assigned_to" value={editData.assigned_to} />
				</div>

				<div class="grid grid-cols-4 items-center gap-4">
					<Label for="address" class="text-right">Dirección</Label>
					<Textarea
						id="address"
						name="address"
						bind:value={editData.delivery_address}
						placeholder="Dirección de entrega"
						class="col-span-3"
						rows={3}
						disabled={isEditing}
					/>
				</div>
			</div>

			<Dialog.Footer>
				<Button
					type="button"
					variant="outline"
					onclick={() => (editOrderDialogOpen = false)}
					disabled={isEditing}
				>
					Cancelar
				</Button>
				<Button type="submit" disabled={isEditing}>
					{isEditing ? 'Actualizando...' : 'Actualizar Orden'}
				</Button>
			</Dialog.Footer>
		</form>
	</Dialog.Content>
</Dialog.Root>

<!-- Delete Confirmation Dialog -->
<Dialog.Root bind:open={deleteConfirmDialogOpen}>
	<Dialog.Content class="sm:max-w-[425px]">
		<Dialog.Header>
			<Dialog.Title>Confirmar Eliminación</Dialog.Title>
			<Dialog.Description>
				¿Estás seguro de que quieres eliminar la orden #{order.order_number}? Esta acción no se
				puede deshacer.
			</Dialog.Description>
		</Dialog.Header>
		<form
			method="POST"
			action="?/delete_order"
			use:enhance={() => {
				isDeleting = true;
				return async ({ result }) => {
					isDeleting = false;
					if (result.type === 'success') {
						toast.success(String((result.data as any)?.message || 'Orden eliminada exitosamente'));
						deleteConfirmDialogOpen = false;
						goto('/auth/dashboard');
					} else if (result.type === 'failure') {
						toast.error(String((result.data as any)?.error || 'Error al eliminar la orden'));
					}
				};
			}}
		>
			<Dialog.Footer>
				<Button
					type="button"
					variant="outline"
					onclick={() => (deleteConfirmDialogOpen = false)}
					disabled={isDeleting}
				>
					Cancelar
				</Button>
				<Button type="submit" variant="destructive" disabled={isDeleting}>
					{isDeleting ? 'Eliminando...' : 'Eliminar Orden'}
				</Button>
			</Dialog.Footer>
		</form>
	</Dialog.Content>
</Dialog.Root>

<!-- Delete Item Confirmation Dialog -->
<Dialog.Root bind:open={deleteItemDialogOpen}>
	<Dialog.Content class="sm:max-w-[425px]">
		<Dialog.Header>
			<Dialog.Title>Confirmar Eliminación de Item</Dialog.Title>
			<Dialog.Description>
				¿Estás seguro de que quieres eliminar el item "{selectedItemName}"? Esta acción no se puede
				deshacer.
			</Dialog.Description>
		</Dialog.Header>
		<form
			method="POST"
			action="?/delete_item"
			use:enhance={() => {
				isDeletingItem = true;
				return async ({ result, update }) => {
					isDeletingItem = false;
					if (result.type === 'success') {
						toast.success(String((result.data as any)?.message || 'Item eliminado exitosamente'));
						deleteItemDialogOpen = false;
						await update();
						await invalidateAll();
					} else if (result.type === 'failure') {
						toast.error(String((result.data as any)?.error || 'Error al eliminar el item'));
					}
				};
			}}
		>
			<input type="hidden" name="item_id" value={selectedItemId} />
			<Dialog.Footer>
				<Button
					type="button"
					variant="outline"
					onclick={() => (deleteItemDialogOpen = false)}
					disabled={isDeletingItem}
				>
					Cancelar
				</Button>
				<Button type="submit" variant="destructive" disabled={isDeletingItem}>
					{isDeletingItem ? 'Eliminando...' : 'Eliminar Item'}
				</Button>
			</Dialog.Footer>
		</form>
	</Dialog.Content>
</Dialog.Root>
