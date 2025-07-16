<script lang="ts">
	import { enhance } from '$app/forms';
	import { page } from '$app/state';
	import { toast } from 'svelte-sonner';

	let { children } = $props();

	function isActive(path: string) {
		return page.url.pathname === path;
	}
</script>

<div class="min-h-screen bg-gray-50">
	<!-- Navigation Header -->
	<nav class="bg-white shadow-sm border-b">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="flex justify-between items-center h-16">
				<div class="flex items-center space-x-8">
					<a href="/" class="flex items-center justify-center gap-x-2">
						<img src="/favicon/favicon-32x32.png" alt="LogiApp Logo" width={32} height={32} />
						<h1 class="text-xl font-bold text-gray-900">LogiApp</h1>
					</a>
					<div class="hidden md:flex space-x-4">
						<a
							href="/auth/dashboard"
							class="px-3 py-2 rounded-md text-sm font-medium transition-colors {isActive(
								'/auth/dashboard'
							)
								? 'bg-blue-100 text-blue-700'
								: 'text-gray-500 hover:text-gray-700 hover:bg-gray-100'}"
						>
							Dashboard
						</a>
						<a
							href="/auth/dashboard/mapa"
							class="px-3 py-2 rounded-md text-sm font-medium transition-colors {isActive(
								'/auth/dashboard/mapa'
							)
								? 'bg-blue-100 text-blue-700'
								: 'text-gray-500 hover:text-gray-700 hover:bg-gray-100'}"
						>
							Mapa de Empleados
						</a>
						<a
							href="/auth/dashboard/orders"
							class="px-3 py-2 rounded-md text-sm font-medium transition-colors {isActive(
								'/auth/dashboard/orders'
							)
								? 'bg-blue-100 text-blue-700'
								: 'text-gray-500 hover:text-gray-700 hover:bg-gray-100'}"
						>
							Órdenes
						</a>
					</div>
				</div>
				<div class="flex items-center">
					<form
						method="POST"
						action="/auth/dashboard?/logout"
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
			</div>
		</div>
	</nav>

	<!-- Main Content -->
	<main>
		{@render children()}
	</main>
</div>
