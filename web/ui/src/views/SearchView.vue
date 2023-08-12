<script setup lang="ts">
import client from '@/api/client'
import { type Manga } from '@/api/schemas'
import MangaCard from '@/components/MangaCard.vue';
import { onMounted, ref, watch } from 'vue';
import { useRoute } from 'vue-router';
import { inject } from 'vue';
import router from '@/router';

const redrawVueMasonry = inject('redrawVueMasonry') as ((container: number) => void);

const route = useRoute()
const mangas = ref([] as Manga[])

const query = ref(route.params['query'] as string)
const loading = ref(false)
const provider = route.params['provider'] as string
const containerId = 42

async function updateRouteWithQuery(query: string) {
    await router.push({
        name: 'search', params: {
            ...route.params,
            query
        }
    })
}

async function handleInput() {
    await updateRouteWithQuery(query.value)
    await searchMangas(query.value)
}

async function searchMangas(query: string) {
    loading.value = true
    const { data, error } = await client.GET('/searchMangas', {
        params: {
            query: { provider, query }
        }
    })
    loading.value = false

    if (error) {
        throw error
    }

    mangas.value = data!
}

function openMangaView(manga: Manga) {
    router.push({
        name: 'manga', params: {
            provider,
            query: query.value,
            manga: manga.id
        }
    })
}

onMounted(async () => {
    if (!query.value) {
        return
    }

    await searchMangas(query.value)
})
</script>

<template>
    <h2 class="text-center">Search</h2>
    <div class="input-group mb-3">
        <input @keypress.enter="handleInput" v-model="query" class="form-control form-control-lg" type="text"
            :placeholder="`Search mangas`" aria-describedby="button-addon2">
        <button @click="handleInput" class="btn btn-primary" type="button" id="button-addon2">Search</button>
    </div>
    <p class="lead text-center">
        {{ mangas.length }} results
    </p>

    <div v-if="loading" class="d-flex justify-content-center">
        <div class="spinner-border" role="status">
            <span class="visually-hidden">Loading...</span>
        </div>
    </div>
    <div v-else class="row" v-masonry="containerId" transition-duration="0.3s" item-selector=".manga">
        <div v-for="manga in mangas" v-masonry-tile class="col-12 col-sm-6 col-md-4 mb-4 manga">
            <MangaCard @click="openMangaView(manga)" @load="redrawVueMasonry(containerId)" :manga="manga" />
        </div>
    </div>
</template>