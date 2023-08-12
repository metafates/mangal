<script setup lang="ts">
import client from '@/api/client'
import { type Manga } from '@/api/schemas'
import MangaCard from '@/components/MangaCard.vue';
import { ref } from 'vue';
import { useRoute } from 'vue-router';
import { inject } from 'vue';

const redrawVueMasonry = inject('redrawVueMasonry') as ((container: number) => void);

const route = useRoute()
const mangas = ref([] as Manga[])

const query = ref("")
const loading = ref(false)
const containerId = 42

async function searchMangas(query: string) {
    loading.value = true
    const { data, error } = await client.GET('/searchMangas', {
        params: {
            query: {
                provider: route.params['provider'] as string,
                query: query
            }
        }
    })
    loading.value = false

    if (error) {
        throw error
    }

    mangas.value = data!
}
</script>

<template>
    <h2 class="text-center">Search</h2>
    <div class="input-group mb-3">
        <input @keypress.enter="() => searchMangas(query)" v-model="query" class="form-control form-control-lg" type="text"
            :placeholder="`Search mangas`" aria-describedby="button-addon2">
        <button @click="() => searchMangas(query)" class="btn btn-primary" type="button" id="button-addon2">Search</button>
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
        <div v-for="manga in mangas" v-masonry-tile class="col-4 mb-4 manga">
            <MangaCard @load="redrawVueMasonry(containerId)" :manga="manga" />
        </div>
    </div>
</template>