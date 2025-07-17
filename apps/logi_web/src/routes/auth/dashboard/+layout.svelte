<script lang="ts">
	import { enhance } from '$app/forms';
	import { page } from '$app/state';
	import { toast } from 'svelte-sonner';
	import { fly, fade } from 'svelte/transition';
	import { Menu, X, LayoutDashboard, MapPin, Package2 } from '@lucide/svelte';

	let { children } = $props();
	let mobileMenuOpen = $state(false);

	function isActive(path: string) {
		return page.url.pathname === path;
	}

	function toggleMobileMenu() {
		mobileMenuOpen = !mobileMenuOpen;
	}

	function closeMobileMenu() {
		mobileMenuOpen = false;
	}
</script>

<div class="min-h-screen bg-background">
	<!-- Navigation Header -->
	<nav class="bg-background shadow-sm border-b border-border">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<div class="flex justify-between items-center h-16">
				<div class="flex items-center space-x-8">
					<a href="/" class="flex items-center justify-center gap-x-2">
						<img src="/favicon/favicon-32x32.png" alt="LogiApp Logo" width={32} height={32} />
						<h1 class="text-xl font-bold text-foreground">LogiApp</h1>
					</a>
					<!-- Desktop Navigation -->
					<div class="hidden md:flex space-x-4">
						<a
							href="/auth/dashboard"
							class="px-3 py-2 rounded-md text-sm font-medium transition-colors flex items-center gap-2 {isActive(
								'/auth/dashboard'
							)
								? 'bg-primary/10 text-primary border border-primary/20'
								: 'text-muted-foreground hover:text-foreground hover:bg-muted'}"
						>
							<LayoutDashboard class="h-4 w-4" />
							Dashboard
						</a>
						<a
							href="/auth/dashboard/mapa"
							class="px-3 py-2 rounded-md text-sm font-medium transition-colors flex items-center gap-2 {isActive(
								'/auth/dashboard/mapa'
							)
								? 'bg-primary/10 text-primary border border-primary/20'
								: 'text-muted-foreground hover:text-foreground hover:bg-muted'}"
						>
							<MapPin class="h-4 w-4" />
							Mapa de Empleados
						</a>
						<a
							href="/auth/dashboard/orders"
							class="px-3 py-2 rounded-md text-sm font-medium transition-colors flex items-center gap-2 {isActive(
								'/auth/dashboard/orders'
							)
								? 'bg-primary/10 text-primary border border-primary/20'
								: 'text-muted-foreground hover:text-foreground hover:bg-muted'}"
						>
							<Package2 class="h-4 w-4" />
							Órdenes
						</a>
					</div>
				</div>
				<!-- Desktop Logout Button -->
				<div class="hidden md:flex items-center">
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
							class="px-4 py-2 text-sm font-medium text-red-700 bg-red-50 rounded-md hover:bg-red-100 hover:text-red-800 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-300 transition-all border border-red-200"
						>
							Cerrar sesión
						</button>
					</form>
				</div>
				<!-- Mobile Hamburger Button -->
				<div class="md:hidden">
					<button
						onclick={toggleMobileMenu}
						class="p-2 rounded-md text-muted-foreground hover:text-foreground hover:bg-muted focus:outline-none focus:ring-2 focus:ring-inset focus:ring-primary"
						aria-expanded={mobileMenuOpen}
						aria-label="Abrir menú de navegación"
					>
						<span class="sr-only">Abrir menú principal</span>
						<!-- Hamburger Icon -->
						{#if !mobileMenuOpen}
							<Menu class="w-6 h-6" />
						{/if}
						<!-- Close Icon -->
						{#if mobileMenuOpen}
							<X class="w-6 h-6" />
						{/if}
					</button>
				</div>
			</div>
		</div>
	</nav>

	<!-- Mobile Menu Overlay -->
	{#if mobileMenuOpen}
		<div class="fixed inset-0 z-50 md:hidden">
			<!-- Backdrop -->
			<button
				class="fixed bg-accent/20 inset-0 bg-opacity-50 transition-opacity cursor-default"
				onclick={closeMobileMenu}
				onkeydown={(e) => {
					if (e.key === 'Escape') {
						closeMobileMenu();
					}
				}}
				aria-label="Cerrar menú de navegación"
				type="button"
				in:fade={{ duration: 200 }}
				out:fade={{ duration: 150 }}
			></button>

			<!-- Sidebar -->
			<div
				class="fixed inset-y-0 right-0 flex max-w-xs w-full"
				in:fly={{ x: 320, duration: 300 }}
				out:fly={{ x: 320, duration: 250 }}
			>
				<div class="w-full bg-background shadow-xl flex flex-col h-full border-l border-border">
					<div
						class="flex items-center justify-between h-16 px-4 border-b border-border flex-shrink-0"
					>
						<div class="flex items-center gap-x-2">
							<img src="/favicon/favicon-32x32.png" alt="LogiApp Logo" width={24} height={24} />
							<h2 class="text-lg font-semibold text-foreground">LogiApp</h2>
						</div>
						<button
							onclick={closeMobileMenu}
							class="p-2 rounded-md text-muted-foreground hover:text-foreground hover:bg-muted"
							aria-label="Cerrar menú de navegación"
						>
							<X class="w-6 h-6" />
						</button>
					</div>

					<!-- Navigation Links -->
					<nav class="flex-1 px-4 py-6 space-y-2 overflow-y-auto">
						<a
							href="/auth/dashboard"
							onclick={closeMobileMenu}
							class="flex items-center px-4 py-3 text-sm font-medium rounded-lg transition-colors gap-3 {isActive(
								'/auth/dashboard'
							)
								? 'bg-primary/10 text-primary border border-primary/20'
								: 'text-muted-foreground hover:text-foreground hover:bg-muted'}"
						>
							<LayoutDashboard class="h-4 w-4" />
							Dashboard
						</a>
						<a
							href="/auth/dashboard/mapa"
							onclick={closeMobileMenu}
							class="flex items-center px-4 py-3 text-sm font-medium rounded-lg transition-colors gap-3 {isActive(
								'/auth/dashboard/mapa'
							)
								? 'bg-primary/10 text-primary border border-primary/20'
								: 'text-muted-foreground hover:text-foreground hover:bg-muted'}"
						>
							<MapPin class="h-4 w-4" />
							Mapa de Empleados
						</a>
						<a
							href="/auth/dashboard/orders"
							onclick={closeMobileMenu}
							class="flex items-center px-4 py-3 text-sm font-medium rounded-lg transition-colors gap-3 {isActive(
								'/auth/dashboard/orders'
							)
								? 'bg-primary/10 text-primary border border-primary/20'
								: 'text-muted-foreground hover:text-foreground hover:bg-muted'}"
						>
							<Package2 class="h-4 w-4" />
							Órdenes
						</a>
					</nav>

					<!-- Logout Button -->
					<div class="p-4 border-t border-border flex-shrink-0">
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
								class="w-full px-4 py-3 text-sm font-medium text-red-700 bg-red-50 rounded-lg hover:bg-red-100 hover:text-red-800 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-300 transition-all border border-red-200"
							>
								Cerrar sesión
							</button>
						</form>
					</div>
				</div>
			</div>
		</div>
	{/if}

	<!-- Main Content -->
	<main>
		{@render children()}
	</main>
</div>
