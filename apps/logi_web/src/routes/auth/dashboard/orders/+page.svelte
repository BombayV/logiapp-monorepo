<script lang="ts">
	import { enhance } from '$app/forms';
	import { DataTable } from '@/components/ui/data-table';
	import { orderColumns, type Order } from '@/components/orders/columns';
	import { toast } from 'svelte-sonner';
	import { buttonVariants, Button } from '@/components/ui/button';
	import * as Dialog from '@/components/ui/dialog';
	import { Label } from '@/components/ui/label';
	import { Input } from '@/components/ui/input';
	import { XCircle, CheckCircle, Package2, Plus } from '@lucide/svelte';
	import type { UserData } from '../../../../app';
	import type { ActionData } from '../$types';

	let { data, form }: { data: { orders: Order[]; user: UserData }; form: ActionData } = $props();

	$effect(() => {
		if (form?.success) {
			toast.success(form.message || 'Operación completada exitosamente');
		} else if (form?.error) {
			toast.error(form.error);
		}
	});

	let orderNumberError = $state('');
	let createOrderDialogOpen = $state(false);

	function validateOrderNumber(value: string): boolean {
		orderNumberError = '';

		if (!value) {
			orderNumberError = 'El número de orden es requerido';
			return false;
		}

		if (value.length < 1 || value.length > 6) {
			orderNumberError = 'El número de orden debe tener entre 1 y 6 dígitos';
			return false;
		}

		if (!/^[0-9]+$/.test(value)) {
			orderNumberError = 'El número de orden solo puede contener dígitos';
			return false;
		}

		return true;
	}

	function handleOrderNumberInput(event: Event) {
		const target = event.target as HTMLInputElement;
		const value = target.value;

		// Remove any non-numeric characters
		const numericValue = value.replace(/[^0-9]/g, '');

		// Limit to 6 characters
		const limitedValue = numericValue.slice(0, 6);

		// Update the input value
		target.value = limitedValue;

		// Validate
		validateOrderNumber(limitedValue);
	}
</script>

<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
	<div class="bg-card rounded-lg shadow-md p-6 border border-border">
		<div class="flex justify-between items-center mb-6">
			<div class="flex items-center gap-3">
				<div class="p-2 bg-primary/10 rounded-lg">
					<Package2 class="h-6 w-6 text-primary" />
				</div>
				<h1 class="text-3xl font-bold text-foreground">Órdenes</h1>
			</div>
			<Dialog.Root bind:open={createOrderDialogOpen}>
				<Dialog.Trigger class={buttonVariants({ variant: 'outline' })}>
					<Plus class="h-4 w-4 mr-2" />
					Crear orden
				</Dialog.Trigger>
				<Dialog.Content class="sm:max-w-[425px]">
					<form
						method="POST"
						action="?/create_order"
						class="space-y-4"
						use:enhance={({ formElement }) => {
							return async ({ result, update }) => {
								if (result.type === 'success') {
									formElement.reset();
									orderNumberError = ''; // Clear validation error
									createOrderDialogOpen = false; // Close dialog
									await update({ reset: false });
								} else if (result.type === 'failure') {
									await update({ reset: false });
								}
							};
						}}
						onsubmit={(e) => {
							const formData = new FormData(e.currentTarget);
							const orderNumber = formData.get('order_number') as string;
							const address = formData.get('address') as string;
							const orderName = formData.get('order_name') as string;
							const orderPhoneNumber = formData.get('order_phone_number') as string;

							let isValid = true;

							if (!validateOrderNumber(orderNumber)) {
								isValid = false;
							}

							if (!address || address.trim().length === 0) {
								toast.error('La dirección es requerida');
								isValid = false;
							}

							if (!orderName || orderName.trim().length === 0) {
								toast.error('El nombre del cliente es requerido');
								isValid = false;
							}

							if (!orderPhoneNumber || orderPhoneNumber.trim().length === 0) {
								toast.error('El teléfono del cliente es requerido');
								isValid = false;
							}

							if (!isValid) {
								e.preventDefault();
							}
						}}
					>
						<Dialog.Header>
							<Dialog.Title>Crear Orden</Dialog.Title>
							<Dialog.Description>
								Complete los detalles de la orden a continuación.
							</Dialog.Description>
						</Dialog.Header>
						<div class="grid gap-4 py-4">
							<div class="grid grid-cols-4 items-center gap-4">
								<Label for="email" class="text-right">Email (Vendedor)</Label>
								<Input
									id="email"
									name="email"
									value={data.user?.email}
									class="col-span-3"
									readonly
								/>
							</div>

							<div class="grid grid-cols-4 items-center gap-4">
								<Label for="order_name" class="text-right">Nombre Cliente *</Label>
								<Input
									id="order_name"
									name="order_name"
									class="col-span-3"
									placeholder="Juan Pérez"
									required
								/>
							</div>

							<div class="grid grid-cols-4 items-center gap-4">
								<Label for="order_phone_number" class="text-right">Teléfono *</Label>
								<Input
									id="order_phone_number"
									name="order_phone_number"
									class="col-span-3"
									placeholder="0991234567"
									required
								/>
							</div>

							<div class="grid grid-cols-4 items-center gap-4">
								<Label for="order_email" class="text-right">Email Cliente</Label>
								<Input
									id="order_email"
									name="order_email"
									type="email"
									class="col-span-3"
									placeholder="cliente@ejemplo.com"
								/>
							</div>

							<div class="grid grid-cols-4 items-center gap-4">
								<Label for="order_cedula" class="text-right">Cédula/RUC</Label>
								<Input
									id="order_cedula"
									name="order_cedula"
									class="col-span-3"
									placeholder="1712345678"
									maxlength={13}
								/>
							</div>

							<div class="grid grid-cols-4 items-center gap-4">
								<Label for="address" class="text-right">Dirección *</Label>
								<Input
									id="address"
									name="address"
									class="col-span-3"
									placeholder="Guayacanes & Los Cipreses, Rumiñahui, Pichincha, 171101"
									required
								/>
							</div>
							<div class="grid grid-cols-4 items-center gap-4">
								<Label for="order_number" class="text-right">Número de orden *</Label>
								<div class="col-span-3">
									<Input
										id="order_number"
										name="order_number"
										placeholder="123456"
										type="text"
										minlength={1}
										maxlength={6}
										pattern={`[0-9]{1,6}`}
										title="El número de orden debe tener entre 1 y 6 dígitos"
										required
										class={orderNumberError ? 'border-red-500' : ''}
										oninput={handleOrderNumberInput}
									/>
									{#if orderNumberError}
										<p class="text-red-500 text-sm mt-1">{orderNumberError}</p>
									{/if}
								</div>
							</div>
						</div>
						<Dialog.Footer>
							<Button type="submit">Crear Orden</Button>
						</Dialog.Footer>
					</form>
				</Dialog.Content>
			</Dialog.Root>
		</div>

		{#if form?.error}
			<div class="mb-4 p-4 bg-red-50 border border-red-200 rounded-lg">
				<div class="flex">
					<XCircle class="w-5 h-5 text-red-400 mr-2" />
					<div>
						<p class="text-sm text-red-800">{form.error}</p>
					</div>
				</div>
			</div>
		{/if}

		{#if form?.success}
			<div class="mb-4 p-4 bg-green-50 border border-green-200 rounded-lg">
				<div class="flex">
					<CheckCircle class="w-5 h-5 text-green-400 mr-2" />
					<div>
						<p class="text-sm text-green-800">{form.message}</p>
					</div>
				</div>
			</div>
		{/if}

		<div class="bg-white shadow-sm rounded-lg border border-gray-200">
			<div class="px-4">
				<DataTable
					data={data.orders}
					columns={orderColumns}
					columnVisibility={{ order_id: false }}
					searchable={true}
					searchPlaceholder="Buscar orden por número/asignado/creado por..."
					customSearchColumns={['order_number', 'created_by_username', 'assigned_to_username']}
					paginated={true}
					sortable={true}
					pageSize={15}
					emptyMessage="No hay órdenes registradas."
				/>
			</div>
		</div>
	</div>
</div>
