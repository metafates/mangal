<script setup lang="ts">
import ProviderCard from '@/components/ProviderCard.vue';
import { onMounted, ref } from 'vue';
import { type components, type paths } from '@/api/mangal';
import createClient from 'openapi-fetch';
import router from '@/router';

const providers = ref([] as components['schemas']['Provider'][])

onMounted(async () => {
  const client = createClient<paths>({ baseUrl: "/api" })

  const { data, error } = await client.GET('/providers', {})
  if (error) {
    throw error
  }

  providers.value = data!
})

function handleClick(provider: components['schemas']['Provider']) {
  router.push({ name: 'search', params: { provider: provider.id } })
}
</script>

<template>
  <main>
    <h2 class="text-center">Providers</h2>
    <div class="row">
      <div v-for="provider in providers" class="col-12 col-sm-6 col-md-4 mb-4">
        <ProviderCard :provider="provider" @click="() => handleClick(provider)" />
      </div>
    </div>
  </main>
</template>
