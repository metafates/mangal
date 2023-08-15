<script setup lang="ts">
import { useRoute } from 'vue-router';
import client from '@/api/client';
import { onMounted, ref } from 'vue';
import type { MangaPage } from '@/api/schemas';
import type { Manga } from '@/api/schemas';

const route = useRoute()
const error = ref("")
const loading = ref(true)
const data = ref(null as MangaPage | null)

onMounted(async () => {
    loading.value = true

    const query = route.params['query'] as string
    const provider = route.params['provider'] as string
    const manga = route.params['manga'] as string

    const res = await client.GET('/mangaPage', {
        params: {
            query: { provider, query, manga }
        }
    })
    loading.value = false

    if (res.error) {
        error.value = res.error.message
        return
    }

    data.value = res.data
})
</script>

<template>
    <div v-if="loading" class="d-flex justify-content-center">
        <div class="spinner-border" role="status">
            <span class="visually-hidden">Loading...</span>
        </div>
    </div>
    <p v-else-if="error">{{ error }}</p>
    <div v-else>
        <img :src="data?.anilistManga?.coverImage.extraLarge ?? data?.manga.cover" alt="">
        <h2>{{ data?.manga.title }}</h2>
        <div v-html="data?.anilistManga?.description"></div>
        <ul v-for="volume in data?.volumes">
            <li>Vol. {{ volume.volume.number }}</li>
            <ul v-for="chapter in volume.chapters">
                <li>Chapter {{ chapter.number }}</li>
            </ul>
        </ul>
    </div>
</template>