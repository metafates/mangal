<script setup lang="ts">
import type { components } from '@/api/mangal';
import type { Manga } from '@/api/schemas';
import { useRoute } from 'vue-router';

const route = useRoute()

const props = defineProps<{
    manga: components['schemas']['Manga']
    btnText?: string
}>()

const emit = defineEmits<{
    (event: "load"): void
}>()

function handleLoad() {
    emit('load')
}

function getImageURL(manga: Manga): string {
    if (!manga.cover) {
        return ""
    }

    const apiURL = new URL('/api/image', `${location.protocol}//${location.host}`)

    apiURL.searchParams.set('url', manga.cover)

    if (manga.url) {
        apiURL.searchParams.set('referer', manga.url)
    }

    return apiURL.toString()
}

</script>

<template>
    <div class="card shadow" style="width: 18rem;">
        <img v-if="manga.cover" :src="getImageURL(manga)" class="card-img-top" alt="Cover image" @load="handleLoad">
        <div class="card-body">
            <h5 class="card-title">{{ manga.title }}</h5>
            <!-- <p class="card-text">{{ manga. }}</p> -->

            <a v-if="manga.url" target="_blank" :href="manga.url" class="link-primary icon-link icon-link-hover">
                Website
                <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor"
                    class="bi bi-arrow-up-left" viewBox="0 0 16 16">
                    <path fill-rule="evenodd"
                        d="M2 2.5a.5.5 0 0 1 .5-.5h6a.5.5 0 0 1 0 1H3.707l10.147 10.146a.5.5 0 0 1-.708.708L3 3.707V8.5a.5.5 0 0 1-1 0v-6z" />
                </svg>
            </a>
        </div>

        <div class="card-footer">
            ID {{ manga.id }}
        </div>
    </div>
</template>

<style scoped>
.card {
    cursor: pointer;
}
</style>