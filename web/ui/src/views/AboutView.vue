<script setup lang="ts">
import { type components, type paths } from '@/api/mangal';
import createClient from 'openapi-fetch';
import { onMounted, ref } from 'vue';

const client = createClient<paths>({ baseUrl: "/api" })
const info = ref({} as components['schemas']['MangalInfo'])

onMounted(async () => {
    const { data, error } = await client.GET('/mangalInfo', {})
    if (error) {
        throw error
    }

    info.value = data!
})
</script>

<template>
    <p>Mangal version {{ info.version }}</p>
</template>