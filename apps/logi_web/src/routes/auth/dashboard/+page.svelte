<script lang="ts">
	import type { PageData, ActionData } from './$types';
	import { buttonVariants, Button } from '@/components/ui/button';
	import * as Dialog from '@/components/ui/dialog';
	import { Label } from '@/components/ui/label';
	import { Input } from '@/components/ui/input';
	import { DataTable } from '@/components/ui/data-table';
	import * as Select from '@/components/ui/select';
	import { columns, type User } from '@/components/users/columns';
	import { toast } from 'svelte-sonner';
	import { enhance } from '$app/forms';
	import type { UserData } from '../../../app';

	let { data, form }: { data: { user: UserData; users: User[] }; form: ActionData } = $props();

	let role = $state<'sales' | 'driver'>('sales');

	// Show toast notifications for form results
	$effect(() => {
		if (form?.success) {
			toast.success(form.message || 'Operación completada exitosamente');
		} else if (form?.error) {
			toast.error(form.error);
		}
	});

	// console.log('Backend users:', data.users);
</script>

<svelte:head>
	<title>LogiApp | Panel de control</title>
	<meta
		name="description"
		content="Bienvenido a tu panel de control en LogiApp. Aquí puedes gestionar tus órdenes y ver tu información de usuario."
	/>
</svelte:head>

<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 mt-6">
	<div class="bg-white rounded-lg shadow-md p-6">
		<div class="flex justify-between items-center mb-6">
			<h1 class="text-3xl font-bold text-gray-800">Dashboard</h1>
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
		{/if}
	</div>
</div>
