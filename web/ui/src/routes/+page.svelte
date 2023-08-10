<script lang="ts">
	import { type paths, type components } from '$lib/mangal.d.ts';
	import createClient from 'openapi-fetch';

	const client = createClient<paths>({ baseUrl: '/api' });

	const providers: components['schemas']['Provider'][] = [];

	async function getProviders() {
		const { data, error } = await client.GET('/providers', {});
		if (error) {
			alert(error.message);
		}

		for (const d of data!) {
			providers.push(d);
		}
	}
</script>

<h1>Welcome to SvelteKit</h1>
<p>Visit <a href="https://kit.svelte.dev">kit.svelte.dev</a> to read the documentation</p>

<button on:click={getProviders}>Providers</button>

{#each providers as provider}
	<p>{provider.name}</p>
{/each}
