<script lang="ts">
	import { onMount } from 'svelte';
	import { browser } from '$app/environment';
	import type { PageData } from './$types';

	export let data: PageData;

	// Leaflet will be loaded dynamically
	let L: any;
	let map: any;
	let mapContainer: HTMLDivElement;
	let markersLayer: any;

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
					center: [36.1699, -115.1398],
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

	function addDriverMarkers() {
		if (!L || !markersLayer) return;

		// Clear existing markers
		markersLayer.clearLayers();

		// Add markers for each driver
		data.drivers.forEach((driver: Driver) => {
			if (driver.latitude && driver.longitude) {
				// Create custom icon for driver
				const driverIcon = L.divIcon({
					html: `
						<div class="driver-marker">
							<div class="driver-icon">
								<svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
									<path d="M12 2C13.1 2 14 2.9 14 4C14 5.1 13.1 6 12 6C10.9 6 10 5.1 10 4C10 2.9 10.9 2 12 2ZM21 9V7L15 5.5V4C15 3.45 14.55 3 14 3H10C9.45 3 9 3.45 9 4V5.5L3 7V9H21ZM12 8C12.55 8 13 8.45 13 9V10H11V9C11 8.45 11.45 8 12 8Z" fill="currentColor"/>
								</svg>
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

				// Create popup content
				const popupContent = `
					<div class="driver-popup">
						<h3 class="font-semibold text-lg">${driver.first_name} ${driver.last_name}</h3>
						<div class="space-y-2 mt-2">
							<p><strong>Email:</strong> ${driver.email}</p>
							<p><strong>Teléfono:</strong> ${driver.phone_number}</p>
							<p><strong>Última conexión:</strong> ${formatDate(driver.last_connection)}</p>
							<p><strong>Ubicación actualizada:</strong> ${formatDate(driver.location_updated_at)}</p>
						</div>
					</div>
				`;

				marker.bindPopup(popupContent);
				markersLayer.addLayer(marker);
			}
		});

		// Fit map to show all markers if there are any
		if (data.drivers.length > 0) {
			const hasValidCoords = data.drivers.some((d: Driver) => d.latitude && d.longitude);
			if (hasValidCoords) {
				const group = new L.featureGroup(Object.values(markersLayer._layers));
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

	// Reactive statement to update markers when data changes
	$: if (browser && markersLayer && data.drivers) {
		addDriverMarkers();
	}
</script>

<svelte:head>
	<title>Mapa de Empleados - LogiApp</title>
	<meta name="description" content="Mapa con la ubicación de los conductores activos" />
	<!-- Leaflet CSS -->
	<link rel="stylesheet" href="https://unpkg.com/leaflet@1.9.4/dist/leaflet.css" />
</svelte:head>

<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 mt-6">
	<div class="bg-white rounded-lg shadow-md p-6">
		<div class="mb-6">
			<h1 class="text-3xl font-bold text-gray-800 mb-2">Mapa de Empleados</h1>
			<p class="text-gray-600">
				Conductores activos: <span class="font-semibold text-blue-600">{data.count}</span>
			</p>
		</div>

		<!-- Map Container -->
		<div class="bg-white rounded-lg shadow-lg overflow-hidden">
			<div bind:this={mapContainer} class="h-[600px] w-full" id="map"></div>
		</div>

		<!-- Drivers List -->
		<div class="mt-6 bg-white rounded-lg shadow-lg p-6">
			<h2 class="text-xl font-semibold mb-4">Conductores Activos</h2>

			{#if data.drivers.length === 0}
				<p class="text-gray-500 text-center py-8">No hay conductores activos en este momento</p>
			{:else}
				<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
					{#each data.drivers as driver}
						<div class="border rounded-lg p-4 hover:shadow-md transition-shadow">
							<h3 class="font-semibold text-lg">{driver.first_name} {driver.last_name}</h3>
							<div class="space-y-1 mt-2 text-sm text-gray-600">
								<p><strong>Email:</strong> {driver.email}</p>
								<p><strong>Teléfono:</strong> {driver.phone_number}</p>
								<p><strong>Última conexión:</strong> {formatDate(driver.last_connection)}</p>
								{#if driver.latitude && driver.longitude}
									<p>
										<strong>Ubicación:</strong>
										{driver.latitude.toFixed(6)}, {driver.longitude.toFixed(6)}
									</p>
								{:else}
									<p class="text-orange-600"><strong>Ubicación:</strong> No disponible</p>
								{/if}
							</div>
						</div>
					{/each}
				</div>
			{/if}
		</div>
	</div>
</div>

<style>
	:global(.driver-marker-container) {
		background: none !important;
		border: none !important;
	}

	:global(.driver-marker) {
		background: white;
		border: 2px solid #3b82f6;
		border-radius: 50%;
		width: 32px;
		height: 32px;
		display: flex;
		align-items: center;
		justify-content: center;
		color: #3b82f6;
		box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
	}

	:global(.driver-marker:hover) {
		background: #3b82f6;
		color: white;
		transform: scale(1.1);
		transition: all 0.2s ease;
	}

	:global(.driver-popup) {
		min-width: 250px;
	}

	:global(.driver-popup h3) {
		margin: 0 0 8px 0;
		color: #1f2937;
	}

	:global(.driver-popup p) {
		margin: 4px 0;
		font-size: 14px;
		color: #4b5563;
	}

	:global(.leaflet-popup-content-wrapper) {
		border-radius: 8px;
		box-shadow:
			0 4px 6px -1px rgba(0, 0, 0, 0.1),
			0 2px 4px -1px rgba(0, 0, 0, 0.06);
	}

	:global(.leaflet-popup-tip) {
		background: white;
	}
</style>
