<script lang="ts">
	import type { PageData } from './$types';
	export let data: PageData;
</script>

<div class="min-h-screen bg-gray-50 p-8">
	<div class="max-w-4xl mx-auto bg-white rounded-lg shadow-md p-6">
		<div class="flex justify-between items-center mb-6">
			<h1 class="text-3xl font-bold text-gray-800">Dashboard</h1>
			<form method="POST" action="?/logout">
				<button
					type="submit"
					class="px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-md hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
				>
					Logout
				</button>
			</form>
		</div>

		<div class="bg-indigo-50 p-4 rounded-lg">
			<h2 class="text-xl font-semibold text-gray-700 mb-2">
				Bienvenue, {data.user.profile?.first_name}!
			</h2>
			<p class="text-gray-600">
				Your role is: <span class="font-mono bg-gray-200 text-gray-800 px-2 py-1 rounded-md"
					>{data.user.role}</span
				>
			</p>
		</div>

		<div class="mt-6">
			<h3 class="text-lg font-semibold text-gray-700">Session Information</h3>
			<pre class="bg-gray-100 p-4 rounded-md mt-2 overflow-x-auto text-sm"><code
					>{JSON.stringify(data, null, 2)}</code
				></pre>
		</div>
	</div>
</div>
