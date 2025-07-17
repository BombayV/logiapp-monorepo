<script lang="ts">
	import {
		Card,
		CardContent,
		CardDescription,
		CardFooter,
		CardHeader,
		CardTitle
	} from '$lib/components/ui/card';
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';
	import { Button } from '$lib/components/ui/button';
	import type { ActionData, PageData } from './$types';
	import { LogIn } from '@lucide/svelte';

	export let form: ActionData;
	export let data: PageData;
	let email = '';
	let password = '';

	// Type guard to check if data has error property
	function hasError(data: PageData): data is PageData & { error: string } {
		return 'error' in data && typeof data.error === 'string';
	}
</script>

<svelte:head>
	<title>LogiApp | Iniciar sesión</title>
	<meta name="description" content="Inicia sesión para acceder a tu cuenta." />
</svelte:head>

<div class="flex items-center justify-center min-h-screen">
	<Card class="w-full max-w-sm">
		<form action="?/login" method="post">
			<CardHeader>
				<div class="flex items-center gap-3 mb-2">
					<div class="p-2 bg-primary/10 rounded-lg">
						<LogIn class="h-6 w-6 text-primary" />
					</div>
					<CardTitle class="text-2xl">Iniciar sesión</CardTitle>
				</div>
				<CardDescription>Ingresa su email y contraseña para acceder a su cuenta.</CardDescription>
			</CardHeader>
			<CardContent class="grid gap-4 mt-2">
				<div class="grid gap-2">
					<Label for="email">Email</Label>
					<Input
						id="email"
						type="email"
						name="email"
						placeholder="me@gmail.com"
						required
						bind:value={email}
					/>
				</div>
				<div class="grid gap-2">
					<Label for="password">Contraseña</Label>
					<Input
						id="password"
						type="password"
						name="password"
						required
						placeholder="********"
						bind:value={password}
					/>
				</div>
				{#if form?.error || hasError(data)}
					{#if form?.error === 'Invalid credentials.'}
						<p class="text-destructive text-sm">Email o contraseña incorrectos.</p>
					{:else}
						<p class="text-destructive text-sm">
							{form?.error || (hasError(data) ? data.error : '')}
						</p>
					{/if}
				{/if}
			</CardContent>
			<CardFooter class="flex flex-col mt-4">
				<Button type="submit" class="w-full">Entrar</Button>
				<p class="mt-4 text-xs text-center text-muted-foreground">
					Toda la información es confidencial y no será compartida.
				</p>
				<a href="/" class="mt-2 text-sm text-primary hover:underline"> Volver al inicio </a>
			</CardFooter>
		</form>
	</Card>
</div>
