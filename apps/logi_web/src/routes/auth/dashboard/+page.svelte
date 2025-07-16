<script lang="ts">
	import type { PageData, ActionData } from './$types';
	import { buttonVariants, Button } from '@/components/ui/button';
	import * as Dialog from '@/components/ui/dialog';
	import { Label } from '@/components/ui/label';
	import { Input } from '@/components/ui/input';
	import { DataTable } from '@/components/ui/data-table';
	import * as Select from '@/components/ui/select';
	import { columns, type User } from '@/components/users/columns';
	import { orderColumns, type Order } from '@/components/orders/columns';
	import { toast } from 'svelte-sonner';
	import { enhance } from '$app/forms';
	import type { UserData } from '../../../app';

	let {
		data,
		form
	}: { data: { user: UserData; users: User[]; orders: Order[] }; form: ActionData } = $props();

	let role = $state<'sales' | 'driver'>('sales');

	// Show toast notifications for form results (except create_order which handles its own)
	$effect(() => {
		if (form?.success) {
			toast.success(form.message || 'Operación completada exitosamente');
		} else if (form?.error) {
			toast.error(form.error);
		}
	});

	// console.log('Backend users:', data.users);
	// console.log('Orders:', data.orders);

	// Client-side validation for order_number
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

<svelte:head>
	<title>LogiApp | Panel de control</title>
	<meta
		name="description"
		content="Bienvenido a tu panel de control en LogiApp. Aquí puedes gestionar tus órdenes y ver tu información de usuario."
	/>
</svelte:head>

<div class="min-h-screen bg-gray-50 p-8">
	<div class="max-w-4xl mx-auto bg-white rounded-lg shadow-md p-6">
		<div class="flex justify-between items-center mb-6">
			<h1 class="text-3xl font-bold text-gray-800">Panel de control</h1>
			<form
				method="POST"
				action="?/logout"
				use:enhance={({ formElement }) => {
					return async ({ result, update }) => {
						if (result.type === 'redirect') {
							toast.success('Sesión cerrada exitosamente');
						}
						await update();
					};
				}}
			>
				<button
					type="submit"
					class="px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-md hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
				>
					Cerrar sesión
				</button>
			</form>
		</div>

		<div class="bg-indigo-50 p-4 rounded-lg">
			<h2 class="text-xl font-semibold text-gray-700">
				Bienvenido, {data.user?.profile?.first_name}
				{data.user?.profile?.last_name}!
			</h2>
			<p class="text-gray-600">
				Su rol es: <span class="font-mono bg-gray-200 text-gray-800 px-2 py-1 rounded-md"
					>{data.user?.role}</span
				>
			</p>
		</div>

		{#if data.user?.role === 'admin'}
			<div class="mt-6 flex flex-col">
				<div class="flex items-center justify-between">
					<h3 class="text-lg font-semibold text-gray-700">Usuarios</h3>
					<Dialog.Root>
						<Dialog.Trigger class={buttonVariants({ variant: 'outline' })}
							>Crear usuario</Dialog.Trigger
						>
						<Dialog.Content class="sm:max-w-[425px]">
							<form method="POST" action="?/create_user" class="space-y-4">
								<Dialog.Header>
									<Dialog.Title>Crear Usuario</Dialog.Title>
									<Dialog.Description>
										Complete los detalles del usuario a continuación.
									</Dialog.Description>
								</Dialog.Header>
								<div class="grid gap-4 py-4">
									<div class="grid grid-cols-4 items-center gap-4">
										<Label for="email" class="text-right">Email</Label>
										<Input
											id="email"
											name="email"
											class="col-span-3"
											placeholder="me@bombayv.com"
										/>
									</div>
									<div class="grid grid-cols-4 items-center gap-4">
										<Label for="password" class="text-right">Contraseña</Label>
										<Input
											id="password"
											name="password"
											class="col-span-3"
											placeholder="********"
											type="password"
										/>
									</div>
									<div class="grid grid-cols-4 items-center gap-4">
										<Label for="first_name" class="text-right">Nombre</Label>
										<Input
											id="first_name"
											name="first_name"
											class="col-span-3"
											placeholder="John"
											required
										/>
										<Label for="last_name" class="text-right">Apellido</Label>
										<Input
											id="last_name"
											name="last_name"
											class="col-span-3"
											placeholder="Doe"
											required
										/>
									</div>
									<div class="grid grid-cols-4 items-center gap-4">
										<Label for="phone_number" class="text-right">Teléfono</Label>
										<Input
											id="phone_number"
											name="phone_number"
											class="col-span-3"
											placeholder="+1234567890"
											type="tel"
											required
										/>
									</div>
									<div class="grid grid-cols-4 items-center gap-4">
										<Label for="role" class="text-right">Rol</Label>
										<Select.Root name="role" type="single" bind:value={role}>
											<Select.Trigger class="h-8 w-full col-span-3">
												{role === 'sales' ? 'Ventas' : 'Conductor'}
											</Select.Trigger>
											<Select.Content>
												<Select.Label>Seleccione un rol</Select.Label>
												<Select.Item value="sales">Ventas</Select.Item>
												<Select.Item value="driver">Conductor</Select.Item>
											</Select.Content>
										</Select.Root>
									</div>
								</div>
								<Dialog.Footer>
									<Button type="submit">Crear Usuario</Button>
								</Dialog.Footer>
							</form>
						</Dialog.Content>
					</Dialog.Root>
				</div>
				<div>
					<DataTable
						data={data.users}
						{columns}
						columnVisibility={{ user_id: false }}
						searchable={true}
						searchPlaceholder="Buscar usuarios por nombre/apellido/telefono/email..."
						customSearchColumns={['first_name', 'last_name', 'phone_number', 'email']}
						paginated={true}
						sortable={true}
						pageSize={10}
						emptyMessage="No hay usuarios registrados."
					/>
				</div>
			</div>
			<hr class="my-6 border-gray-200" />
		{/if}
		<div class="mt-6 flex flex-col">
			<div class="flex items-center justify-between">
				<h3 class="text-lg font-semibold text-gray-700">Ordenes recientes</h3>
				<Dialog.Root bind:open={createOrderDialogOpen}>
					<Dialog.Trigger class={buttonVariants({ variant: 'outline' })}>Crear orden</Dialog.Trigger
					>
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
										await update();
									} else if (result.type === 'failure') {
										await update();
									}
								};
							}}
							onsubmit={(e) => {
								const formData = new FormData(e.currentTarget);
								const orderNumber = formData.get('order_number') as string;
								const address = formData.get('address') as string;

								let isValid = true;

								if (!validateOrderNumber(orderNumber)) {
									isValid = false;
								}

								if (!address || address.trim().length === 0) {
									toast.error('La dirección es requerida');
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
									<Label for="email" class="text-right">Email</Label>
									<Input
										id="email"
										name="email"
										value={data.user?.email}
										class="col-span-3"
										readonly
									/>
								</div>
								<div class="grid grid-cols-4 items-center gap-4">
									<Label for="address" class="text-right">Dirección</Label>
									<Input
										id="address"
										name="address"
										class="col-span-3"
										placeholder="Guayacanes & Los Cipreses, Rumiñahui, Pichincha, 171101"
										required
									/>
								</div>
								<div class="grid grid-cols-4 items-center gap-4">
									<Label for="order_number" class="text-right">Número de orden</Label>
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
			<div>
				<DataTable
					data={data.orders}
					columns={orderColumns}
					columnVisibility={{ order_id: false }}
					searchable={true}
					searchPlaceholder="Buscar order por numero de orden/asignado/creado por..."
					customSearchColumns={['order_number', 'created_by_username', 'assigned_to_username']}
					paginated={true}
					sortable={true}
					pageSize={15}
					emptyMessage="No hay órdenes recientes."
				/>
			</div>
		</div>

		<div class="mt-6 flex flex-col">
			<div class="flex items-center justify-between">
				<h3 class="text-lg font-semibold text-gray-700">Mapa de empleados</h3>
			</div>
		</div>
	</div>
</div>
