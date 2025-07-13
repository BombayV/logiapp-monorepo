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
	import type { ActionData } from './$types';

	export let form: ActionData;
	let email = '';
	let password = '';
</script>

<div class="flex items-center justify-center min-h-screen">
	<Card class="w-full max-w-sm">
		<form action="?/login" method="post">
			<CardHeader>
				<CardTitle class="text-2xl">Iniciar sesión</CardTitle>
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
				{#if form?.error}
					{#if form.error === 'Invalid credentials.'}
						<p class="text-red-500 text-sm">Email o contraseña incorrectos.</p>
					{:else}
						<p class="text-red-500 text-sm">{form.error}</p>
					{/if}
				{/if}
			</CardContent>
			<CardFooter class="flex flex-col mt-4">
				<Button type="submit" class="w-full">Entrar</Button>
				<p class="mt-4 text-xs text-center text-muted-foreground">
					Toda la información es confidencial y no será compartida.
				</p>
				<a href="/" class="mt-2 text-sm text-blue-500 hover:underline"> Volver al inicio </a>
			</CardFooter>
		</form>
	</Card>
</div>
