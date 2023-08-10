<script lang="ts">
	import { Badge, Button, Card, Group, Image, Text } from '@svelteuidev/core';

	import { type paths, type components } from '$lib/mangal';
	import createClient from 'openapi-fetch';

	const client = createClient<paths>({ baseUrl: '/api' });

	async function getProviders() {
		const { data, error } = await client.GET('/providers', {});
		if (error) {
			throw error;
		}

		return data;
	}
</script>

{#await getProviders()}
	<p>Loading providers...</p>
{:then providers}
	{#each providers as provider}
		<Card shadow="sm" padding="lg">
			<Group position="apart">
				<Text weight={500}>{provider.name}</Text>
				<Badge color="pink" variant="light">
					{provider.id}
				</Badge>
			</Group>
		</Card>
	{/each}
{:catch error}
	<p>{error.message}</p>
{/await}
