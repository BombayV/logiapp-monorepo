<script lang="ts">
	import { Button } from '$lib/components/ui/button';
	import { ArrowRight, Truck, Map, Download } from '@lucide/svelte';
	import type { PageData } from './$types';

	let { data }: { data: PageData } = $props();
</script>

<svelte:head>
	<title>LogiApp | Plataforma de Logística</title>
	<meta
		name="description"
		content="Bienvenido a LogiApp, la plataforma de logística que te ayuda a gestionar tus envíos de manera eficiente y efectiva.
			Desde la planificación hasta la entrega, LogiApp te ofrece herramientas para optimizar tu cadena de suministro y mejorar la experiencia del cliente."
	/>
</svelte:head>

<header class="flex items-center justify-between p-4 border-b bg-background">
	<a href="/" class="flex items-center space-x-2">
		<img src="/favicon/favicon-32x32.png" alt="LogiApp Logo" width={32} height={32} />
		<h1 class="text-2xl font-bold">LogiApp</h1>
	</a>
	<a href={data.loggedIn ? '/auth/dashboard' : '/auth/login'}>
		<Button>{data.loggedIn ? 'Dashboard' : 'Iniciar Sesión'}</Button>
	</a>
</header>
<div class="flex flex-col min-h-screen bg-background">
	<main class="flex-1">
		<section
			class="w-full py-12 md:py-24 lg:py-32 xl:py-48 bg-[url('data:image/svg+xml,%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20viewBox%3D%220%200%2032%2032%22%20width%3D%2232%22%20height%3D%2232%22%20fill%3D%22none%22%20stroke%3D%22rgb(228%20228%20231%2F.5)%22%3E%3Cpath%20d%3D%22M0%20.5H31.5V32%22%2F%3E%3C%2Fsvg%3E')]"
		>
			<div class="px-4 md:px-6">
				<div class="flex flex-col items-center justify-center space-y-4 text-center">
					<div class="space-y-2">
						<div class="inline-block rounded-lg bg-neutral-200 px-3 py-1 text-sm">
							Bienvenido a LogiApp
						</div>
						<h1 class="text-3xl font-bold tracking-tighter sm:text-5xl">
							Tu Plataforma de Logística Inteligente
						</h1>
						<p
							class="max-w-[900px] text-gray-500 md:text-xl/relaxed lg:text-base/relaxed xl:text-xl/relaxed"
						>
							Maneja pedidos, envíos y almacenes de manera eficiente con nuestra ayuda inteligente.
						</p>
					</div>
					<a href="/auth/login" class="inline-block">
						<Button size="lg" class="rounded-full px-8">Comienza Ahora</Button>
					</a>
				</div>
			</div>
		</section>
		<section class="w-full py-12 md:py-24 lg:py-32 bg-neutral-50">
			<div class="px-4 md:px-6">
				<div class="flex flex-col items-center justify-center space-y-4 text-center">
					<div class="space-y-2">
						<div class="inline-block rounded-lg bg-neutral-200 px-3 py-1 text-sm">
							Características Destacadas
						</div>
						<h2 class="text-3xl font-bold tracking-tighter sm:text-4xl">
							Todo lo que Necesitas para tu Logística
						</h2>
						<p
							class="max-w-[900px] text-gray-500 md:text-xl/relaxed lg:text-base/relaxed xl:text-xl/relaxed"
						>
							Nuestra plataforma ofrece todas las herramientas necesarias para gestionar tu
							logística de manera efectiva y eficiente.
						</p>
					</div>
				</div>
				<div class="mx-auto grid max-w-5xl items-center gap-6 py-12 lg:grid-cols-3 lg:gap-12">
					<div class="grid gap-1">
						<div class="flex justify-center">
							<Truck class="h-12 w-12 text-gray-500" />
						</div>
						<h3 class="text-xl font-bold text-center">Gestión a Tiempo Real</h3>
						<p class="text-sm text-gray-500 text-center">
							Monitorea tus envíos y pedidos en tiempo real para una mejor toma de decisiones.
						</p>
					</div>
					<div class="grid gap-1">
						<div class="flex justify-center">
							<ArrowRight class="h-12 w-12 text-gray-500" />
						</div>
						<h3 class="text-xl font-bold text-center">Integración Inteligente</h3>
						<p class="text-sm text-gray-500 text-center">
							Conecta con el personal más cercano y optimiza tus rutas de entrega con nuestra
							tecnologia.
						</p>
					</div>
					<div class="grid gap-1">
						<div class="flex justify-center">
							<Map class="h-12 w-12 text-gray-500" />
						</div>
						<h3 class="text-xl font-bold text-center">Mapa de Empleados</h3>
						<p class="text-sm text-gray-500 text-center">
							Manten un control total de todos los empleados disponibles y sus ubicaciones en tiempo
							real.
						</p>
					</div>
				</div>
			</div>
		</section>
		<section class="w-full py-12 md:py-24 lg:py-32 bg-white">
			<div class="px-4 md:px-6">
				<div class="flex flex-col items-center justify-center space-y-4 text-center">
					<div class="space-y-2">
						<div class="inline-block rounded-lg bg-neutral-200 px-3 py-1 text-sm">
							Disponible en Móvil
						</div>
						<h2 class="text-3xl font-bold tracking-tighter sm:text-4xl">Descarga la App LogiApp</h2>
						<p
							class="max-w-[900px] text-gray-500 md:text-xl/relaxed lg:text-base/relaxed xl:text-xl/relaxed"
						>
							Lleva tu logística contigo. Gestiona pedidos, monitorea envíos y mantén el control
							total desde tu dispositivo móvil.
						</p>
					</div>
					<div class="flex flex-col sm:flex-row gap-4 mt-8">
						<a href="https://www.mediafire.com/file/mli7h5kjqj92joy/LogiApp1.0.0.apk/file" class="inline-block">
							<Button size="lg" class="rounded-full px-8 flex items-center gap-2">
								<Download class="h-5 w-5" />
								Descargar para Android
							</Button>
						</a>
						<a href="/" class="inline-block">
							<Button size="lg" variant="outline" class="rounded-full px-8 flex items-center gap-2">
								<Download class="h-5 w-5" />
								Descargar para iOS
							</Button>
						</a>
					</div>
				</div>
			</div>
		</section>
	</main>
	<footer class="flex items-center justify-center w-full h-24 border-t">
		<p class="text-gray-500">© 2025 LogiApp. Todos los derechos reservados.</p>
	</footer>
</div>
