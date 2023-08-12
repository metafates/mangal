<script setup lang="ts">
import client from '@/api/client';
import type { components } from '@/api/mangal';
import type { Manga } from '@/api/schemas';

const props = defineProps<{
    manga: components['schemas']['Manga']
    btnText?: string
}>()

const emit = defineEmits<{
    (event: "click"): void
    (event: "load"): void
}>()

function handleClick() {
    emit('click')
}

function handleLoad() {
    emit('load')
}

function getImageURL(manga: Manga): string {
    if (!manga.cover) {
        return ""
    }

    const apiURL = new URL('http://localhost:6969/api/image')

    apiURL.searchParams.set('url', manga.cover)

    if (manga.url) {
        apiURL.searchParams.set('referer', manga.url)
    }

    return apiURL.toString()
}

</script>

<template>
    <div class="card" style="width: 18rem;">
        <img v-if="manga.cover" :src="getImageURL(manga)" class="card-img-top" alt="Cover image" @load="handleLoad">
        <div class="card-body">
            <h5 class="card-title">{{ manga.title }}</h5>
            <!-- <p class="card-text">{{ manga. }}</p> -->

            <div class="btn-group">
                <button @click="handleClick" class="btn btn-primary">{{ btnText ?? "Select" }}</button>
                <a v-if="manga.url" target="_blank" :href="manga.url" class="btn btn-secondary">Website</a>
            </div>
        </div>
    </div>
</template>