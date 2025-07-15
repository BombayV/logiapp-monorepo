<script lang="ts">
	import type { PageData } from './$types';
	import { buttonVariants, Button } from '@/components/ui/button';
	import * as Dialog from '@/components/ui/dialog';
	import { Label } from '@/components/ui/label';
	import { Input } from '@/components/ui/input';
	import { DataTable } from '@/components/ui/data-table';
	import { columns, type User } from '@/components/users/columns';
	import { orderColumns, type Order } from '@/components/orders/columns';
	import type { UserData } from '../../../app';
	export let data: {
		user: UserData;
		users: User[];
		orders: Order[];
	};

	// console.log('Backend users:', data.users);
	// console.log('Orders:', data.orders);
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
			<form method="POST" action="?/logout">
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
									<Input id="email" class="col-span-3" placeholder="me@bombayv.com" />
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
							</div>
							<Dialog.Footer>
								<Button type="submit">Crear Usuario</Button>
							</Dialog.Footer>
						</form>
					</Dialog.Content>
				</Dialog.Root>
			</div>
			<div class="mt-4">
				<DataTable
					data={data.users}
					{columns}
					columnVisibility={{ user_id: false }}
					searchable={true}
					searchPlaceholder="Buscar usuarios..."
					paginated={true}
					sortable={true}
					pageSize={100}
					emptyMessage="No hay usuarios registrados."
				/>
			</div>
		</div>

		<div class="mt-6 flex flex-col">
			<div class="flex items-center justify-between">
				<h3 class="text-lg font-semibold text-gray-700">Ordenes recientes</h3>
				<Dialog.Root>
					<Dialog.Trigger class={buttonVariants({ variant: 'outline' })}>Crear orden</Dialog.Trigger
					>
					<Dialog.Content class="sm:max-w-[425px]">
						<form method="POST" action="?/create_order" class="space-y-4">
							<Dialog.Header>
								<Dialog.Title>Crear Orden</Dialog.Title>
								<Dialog.Description>
									Complete los detalles de la orden a continuación.
								</Dialog.Description>
							</Dialog.Header>
							<div class="grid gap-4 py-4">
								<div class="grid grid-cols-4 items-center gap-4">
									<Label for="email" class="text-right">Email</Label>
									<Input id="email" value={data.user?.email} class="col-span-3" disabled />
								</div>
								<div class="grid grid-cols-4 items-center gap-4">
									<Label for="address" class="text-right">Dirección</Label>
									<Input
										id="address"
										name="address"
										class="col-span-3"
										placeholder="Guayacanes & Los Cipreses, Rumiñahui, Pichincha, 171101"
									/>
								</div>
							</div>
							<Dialog.Footer>
								<Button type="submit">Crear Orden</Button>
							</Dialog.Footer>
						</form>
					</Dialog.Content>
				</Dialog.Root>
			</div>
			<div class="mt-4">
				<DataTable
					data={data.orders}
					columns={orderColumns}
					columnVisibility={{ order_id: false }}
					searchable={true}
					searchPlaceholder="Buscar órdenes..."
					paginated={true}
					sortable={true}
					pageSize={100}
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
