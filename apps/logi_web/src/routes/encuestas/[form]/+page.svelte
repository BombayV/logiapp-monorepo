<script>
	import Button from '@/components/ui/button/button.svelte';
	import Input from '@/components/ui/input/input.svelte';
	import Label from '@/components/ui/label/label.svelte';
	import * as RadioGroup from '$lib/components/ui/radio-group/index.js';
	import Textarea from '@/components/ui/textarea/textarea.svelte';
	import { enhance } from '$app/forms';
	import { toast } from 'svelte-sonner';

	let { data } = $props();

	let driverRating = $state('');
	let cargoCondition = $state('');
	let isSubmitting = $state(false);
	let isSuccess = $state(false);

	let isFormFilled = $derived(driverRating !== '' && cargoCondition !== '');
</script>

<svelte:head>
	<title>LogiApp | Encuesta de Satisfaccion</title>
</svelte:head>

<div class="flex flex-col min-h-screen bg-linear-to-br from-background to-neutral-300 px-6 py-12">
	<div class="w-full max-w-xl mx-auto flex flex-col items-center relative">
		<h1 class="text-3xl font-semibold text-left w-full">Encuesta de Satisfaccion</h1>
		<p class="text-sm text-muted-foreground w-full text-left mt-1">ID: {data.form.public_id}</p>
		<p class="w-full text-left text-muted-foreground mt-2">
			Por favor conteste las preguntas acerca de su pedido. Todas las respuestas son privadas y solo
			se usaran para mejorar la experiencia con el client.
		</p>
		{#if isSuccess}
			<div class="w-full mt-8 p-8 bg-green-50 rounded-lg border border-green-200 text-center">
				<h2 class="text-2xl font-bold text-green-800 mb-2">¡Gracias por tu opinión!</h2>
				<p class="text-green-700">Hemos recibido tu encuesta correctamente.</p>
			</div>
		{:else if data.form?.is_finished}
			<div class="w-full mt-8 p-8 bg-yellow-50 rounded-lg border border-yellow-200 text-center">
				<h2 class="text-2xl font-bold text-yellow-800 mb-2">Encuesta ya respondida</h2>
				<p class="text-yellow-700">Esta encuesta ya ha sido completada anteriormente.</p>
			</div>
		{:else}
			<form
				class="w-full relative mt-8 flex flex-col gap-y-4"
				action="?/user_form"
				method="POST"
				use:enhance={() => {
					isSubmitting = true;
					return async ({ result }) => {
						isSubmitting = false;
						if (result.type === 'success') {
							isSuccess = true;
							toast.success('Encuesta enviada correctamente');
						} else if (result.type === 'failure') {
							toast.error(result.data?.error || 'Error al enviar la encuesta');
						}
					};
				}}
			>
				<input type="hidden" name="driver_rating" value={driverRating} />
				<input type="hidden" name="cargo_condition" value={cargoCondition} />
				<div class="flex flex-col w-full relative gap-y-3">
					<Label for="driver_email">Email de Conductor</Label>
					<Input class="grow" value={data.form.driver_email || 'No asignado'} disabled />
				</div>
				<div class="flex flex-col w-full relative gap-y-3">
					<Label for="driver_name">Nombre de Conductor</Label>
					<Input class="grow" value={data.form.driver_name || 'No asignado'} disabled />
				</div>
				<div class="flex flex-col w-full relative gap-y-3">
					<Label for="driver_rating">Calificacion de Conductor *</Label>
					<RadioGroup.Root bind:value={driverRating} class="grid grid-cols-5 mt-4">
						<div class="flex flex-col items-center justify-center gap-3">
							<RadioGroup.Item value="1" id="r1" class="size-6" />
							<Label for="r1" class="text-lg">1</Label>
						</div>
						<div class="flex flex-col items-center justify-center gap-3">
							<RadioGroup.Item value="2" id="r2" class="size-6" />
							<Label for="r2" class="text-lg">2</Label>
						</div>
						<div class="flex flex-col items-center justify-center gap-3">
							<RadioGroup.Item value="3" id="r3" class="size-6" />
							<Label for="r3" class="text-lg">3</Label>
						</div>
						<div class="flex flex-col items-center justify-center gap-3">
							<RadioGroup.Item value="4" id="r4" class="size-6" />
							<Label for="r4" class="text-lg">4</Label>
						</div>
						<div class="flex flex-col items-center justify-center gap-3">
							<RadioGroup.Item value="5" id="r5" class="size-6" />
							<Label for="r5" class="text-lg">5</Label>
						</div>
					</RadioGroup.Root>
				</div>
				<div class="flex flex-col w-full relative gap-y-3">
					<Label for="order_condition">Estado del Cargamento *</Label>
					<RadioGroup.Root bind:value={cargoCondition} class="flex flex-col mt-4">
						<div class="flex gap-3">
							<RadioGroup.Item value="good" id="r1" class="size-5" />
							<Label for="r1" class="text-md">Buen Estado</Label>
						</div>
						<div class="flex gap-3">
							<RadioGroup.Item value="regular" id="r2" class="size-5" />
							<Label for="r2" class="text-md">Regular</Label>
						</div>
						<div class="flex gap-3">
							<RadioGroup.Item value="bad" id="r3" class="size-5" />
							<Label for="r3" class="text-md">Malo</Label>
						</div>
						<div class="flex gap-3">
							<RadioGroup.Item value="very_bad" id="r4" class="size-5" />
							<Label for="r4" class="text-md">Muy Malo</Label>
						</div>
					</RadioGroup.Root>
				</div>
				<div class="flex flex-col w-full relative gap-y-3">
					<Label for="experience_comments">Experiencia con el Conductor</Label>
					<Textarea
						id="experience_comments"
						name="experience_comments"
						rows={4}
						placeholder="Escribe tus comentarios aqui..."
					/>
				</div>
				<Button
					disabled={!isFormFilled || isSubmitting}
					type="submit"
					size="lg"
					class="w-full mt-8"
				>
					{isSubmitting ? 'Enviando...' : 'Completar Encuesta'}
				</Button>
			</form>
		{/if}
	</div>
</div>
