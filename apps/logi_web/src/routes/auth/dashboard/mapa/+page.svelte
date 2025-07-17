<script lang="ts">
	import { onMount } from 'svelte';
	import { browser } from '$app/environment';
	import { User, MapPin, Users } from '@lucide/svelte';
	import * as Card from '@/components/ui/card';
	import type { PageData } from './$types';
	import type { Map, LayerGroup, LeafletMouseEvent } from 'leaflet';

	let { data }: { data: PageData } = $props();

	// Leaflet will be loaded dynamically
	let L: typeof import('leaflet');
	let map: Map;
	let mapContainer: HTMLDivElement;
	let markersLayer: LayerGroup;

	// Create a hidden div to render the Lucide icon and get its HTML
	let iconContainer: HTMLDivElement;

	// Driver interface
	interface Driver {
		user_id: string;
		email: string;
		first_name: string;
		last_name: string;
		phone_number: string;
		latitude: number;
		longitude: number;
		last_connection: string;
		location_updated_at: string;
	}

	onMount(() => {
		let destroyed = false;
		(async () => {
			if (browser) {
				// Dynamically import Leaflet
				const leafletModule = await import('leaflet');
				L = leafletModule.default;

				// Initialize the map
				map = L.map(mapContainer, {
					center: [-0.1807, -78.4678], // Quito, Ecuador coordinates
					zoom: 10,
					zoomControl: true
				});

				// Add OpenStreetMap tiles
				L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
					attribution: '© OpenStreetMap contributors'
				}).addTo(map);

				// Create a layer group for markers
				markersLayer = L.layerGroup().addTo(map);

				// Add markers for drivers
				addDriverMarkers();
			}
		})();
		// Cleanup when component is destroyed
		return () => {
			destroyed = true;
			if (map) {
				map.remove();
			}
		};
	});

	// Helper function to get Lucide User icon HTML from the actual component
	function getUserIconHTML(): string {
		if (!iconContainer) return '';
		return iconContainer.innerHTML;
	}

	function addDriverMarkers() {
		if (!L || !markersLayer) return;

		// Clear existing markers
		markersLayer.clearLayers();

		// Add markers for each driver
		data.drivers.forEach((driver: Driver) => {
			if (driver.latitude && driver.longitude) {
				// Create custom icon for driver with Lucide User icon
				const driverIcon = L.divIcon({
					html: `
						<div class="driver-marker" data-driver-id="${driver.user_id}">
							<div class="driver-icon">
								${getUserIconHTML()}
							</div>
						</div>
					`,
					className: 'driver-marker-container',
					iconSize: [32, 32],
					iconAnchor: [16, 32],
					popupAnchor: [0, -32]
				});

				// Create marker
				const marker = L.marker([driver.latitude, driver.longitude], {
					icon: driverIcon
				});

				// Add popup with driver details
				const popupContent = `
					<div class="p-4 min-w-[280px]">
						<div class="flex items-center space-x-3 mb-4">
							<div class="w-10 h-10 bg-primary/10 rounded-full flex items-center justify-center">
								<svg class="w-5 h-5 text-primary" fill="none" stroke="currentColor" viewBox="0 0 24 24">
									<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
								</svg>
							</div>
							<div class="flex flex-col">
								<h4 class="font-semibold text-lg">${driver.first_name} ${driver.last_name}</h4>
								<p class="text-sm text-muted-foreground !m-0">${driver.email}</p>
							</div>
						</div>
						<div class="space-y-2 text-sm">
							<div class="flex items-center space-x-2">
								<span class="text-muted-foreground">Teléfono:</span>
								<span>${driver.phone_number}</span>
							</div>
							<div class="flex items-center space-x-2">
								<span class="text-muted-foreground">Última conexión:</span>
								<span>${formatDate(driver.last_connection)}</span>
							</div>
							${
								driver.latitude && driver.longitude
									? `
								<div class="flex items-center space-x-2">
									<span class="text-muted-foreground">Coordenadas:</span>
									<span class="font-mono text-xs">
										${driver.latitude.toFixed(6)}, ${driver.longitude.toFixed(6)}
									</span>
								</div>
							`
									: `
								<div class="flex items-center space-x-2">
									<span class="text-muted-foreground">Ubicación:</span>
									<span class="text-destructive">No disponible</span>
								</div>
							`
							}
						</div>
					</div>
				`;

				marker.bindPopup(popupContent, {
					maxWidth: 320,
					className: 'custom-popup'
				});

				markersLayer.addLayer(marker);
			}
		});

		// Fit map to show all markers if there are any
		if (data.drivers.length > 0) {
			const hasValidCoords = data.drivers.some((d: Driver) => d.latitude && d.longitude);
			if (hasValidCoords) {
				const group = L.featureGroup(Object.values(markersLayer.getLayers()));
				map.fitBounds(group.getBounds().pad(0.1));
			}
		}
	}

	function formatDate(dateString: string): string {
		const date = new Date(dateString);
		return date.toLocaleString('es-ES', {
			year: 'numeric',
			month: 'short',
			day: 'numeric',
			hour: '2-digit',
			minute: '2-digit'
		});
	}

	// Effect to update markers when data changes
	$effect(() => {
		if (browser && markersLayer && data.drivers) {
			addDriverMarkers();
		}
	});
</script>

<svelte:head>
	<title>Mapa de Empleados - LogiApp</title>
	<meta name="description" content="Mapa con la ubicación de los conductores activos" />
	<!-- Leaflet CSS -->
	<link rel="stylesheet" href="https://unpkg.com/leaflet@1.9.4/dist/leaflet.css" />
</svelte:head>

<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
	<!-- Hidden container to render the Lucide User icon -->
	<div bind:this={iconContainer} class="hidden">
		<User size={24} />
	</div>

	<div class="bg-card rounded-lg shadow-md p-6 border border-border">
		<div class="mb-6">
			<div class="flex items-center gap-3 mb-2">
				<div class="p-2 bg-primary/10 rounded-lg">
					<MapPin class="h-6 w-6 text-primary" />
				</div>
				<h1 class="text-3xl font-bold text-foreground">Mapa de Empleados</h1>
			</div>
			<p class="text-muted-foreground">
				Conductores activos: <span class="font-semibold text-primary">{data.count}</span>
			</p>
		</div>

		<!-- Map Container -->
		<div class="bg-card rounded-lg shadow-lg overflow-hidden relative border border-border z-0">
			<div bind:this={mapContainer} class="h-[600px] w-full" id="map"></div>
		</div>

		<!-- Drivers List -->
		<div class="mt-6 bg-card rounded-lg shadow-lg p-6 border border-border">
			<div class="flex items-center gap-3 mb-4">
				<div class="p-2 bg-secondary/20 rounded-lg">
					<Users class="h-5 w-5 text-secondary-foreground" />
				</div>
				<h2 class="text-xl font-semibold text-foreground">Conductores Activos</h2>
			</div>

			{#if data.drivers.length === 0}
				<p class="text-muted-foreground text-center py-8">
					No hay conductores activos en este momento
				</p>
			{:else}
				<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
					{#each data.drivers as driver}
						<Card.Root class="hover:shadow-md transition-shadow cursor-pointer">
							<Card.Header class="">
								<div class="flex items-center space-x-3">
									<div
										class="w-10 h-10 bg-primary/10 rounded-full flex items-center justify-center"
									>
										<User class="w-5 h-5 text-primary" />
									</div>
									<div>
										<Card.Title class="text-lg">{driver.first_name} {driver.last_name}</Card.Title>
										<Card.Description class="text-sm">{driver.email}</Card.Description>
									</div>
								</div>
							</Card.Header>
							<Card.Content class="pt-0">
								<div class="space-y-2 text-sm">
									<div class="flex items-center space-x-2">
										<span class="text-muted-foreground">Teléfono:</span>
										<span>{driver.phone_number}</span>
									</div>
									<div class="flex items-center space-x-2">
										<span class="text-muted-foreground">Última conexión:</span>
										<span>{formatDate(driver.last_connection)}</span>
									</div>
									{#if driver.latitude && driver.longitude}
										<div class="flex items-center space-x-2">
											<span class="text-muted-foreground">Coordenadas:</span>
											<span class="font-mono text-xs">
												{driver.latitude.toFixed(6)}, {driver.longitude.toFixed(6)}
											</span>
										</div>
									{:else}
										<div class="flex items-center space-x-2">
											<span class="text-muted-foreground">Ubicación:</span>
											<span class="text-destructive">No disponible</span>
										</div>
									{/if}
								</div>
							</Card.Content>
						</Card.Root>
					{/each}
				</div>
			{/if}
		</div>
	</div>
</div>

<style>
	/* Map container z-index fix */
	:global(.leaflet-container) {
		z-index: 1 !important;
	}

	:global(.leaflet-control-container) {
		z-index: 2 !important;
	}

	:global(.driver-marker-container) {
		background: none !important;
		border: none !important;
	}

	:global(.driver-marker) {
		background: white;
		border: 2px solid oklch(0.72 0.06 180);
		border-radius: 50%;
		width: 32px;
		height: 32px;
		display: flex;
		align-items: center;
		justify-content: center;
		color: oklch(0.72 0.06 180);
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
		cursor: pointer;
	}

	:global(.driver-marker:hover) {
		background: oklch(0.72 0.06 180);
		color: white;
		transform: scale(1.1);
		transition: all 0.2s ease;
	}

	/* Custom popup styling */
	:global(.custom-popup .leaflet-popup-content-wrapper) {
		border-radius: 8px;
		box-shadow:
			0 4px 6px -1px rgba(0, 0, 0, 0.1),
			0 2px 4px -1px rgba(0, 0, 0, 0.06);
		border: 1px solid #e5e7eb;
	}

	:global(.custom-popup .leaflet-popup-content) {
		margin: 0;
		line-height: 1.4;
		font-family: inherit;
	}

	:global(.custom-popup .leaflet-popup-tip) {
		background: white;
		border: 1px solid #e5e7eb;
		border-top: none;
		border-right: none;
	}

	:global(.custom-popup .leaflet-popup-close-button) {
		color: #9ca3af;
		font-size: 16px;
		padding: 8px;
		font-weight: bold;
		text-decoration: none;
		border-radius: 4px;
		transition: all 0.2s ease;
		width: 24px;
		height: 24px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: #f3f4f6;
		border: 1px solid #e5e7eb;
		margin-right: 8px;
		margin-top: 8px;
	}

	:global(.custom-popup .leaflet-popup-close-button:hover) {
		color: #374151;
		background: #e5e7eb;
		transform: scale(1.1);
	}

	:global(.custom-popup .leaflet-popup-close-button:active) {
		transform: scale(0.95);
	}
</style>
