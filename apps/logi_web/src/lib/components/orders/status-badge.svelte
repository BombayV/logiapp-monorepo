<script lang="ts">
	import type { Status } from '../orders/columns';
	
	let { status }: { status: Status } = $props();
	
	const statusConfig = {
		pending: {
			label: 'Pendiente',
			class: 'bg-yellow-100 text-yellow-800 border-yellow-200'
		},
		in_progress: {
			label: 'En progreso',
			class: 'bg-blue-100 text-blue-800 border-blue-200'
		},
		completed: {
			label: 'Completado',
			class: 'bg-green-100 text-green-800 border-green-200'
		},
		cancelled: {
			label: 'Cancelado',
			class: 'bg-red-100 text-red-800 border-red-200'
		}
	};
	
	const config = statusConfig[status];
</script>

<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border {config.class}">
	{config.label}
</span>
