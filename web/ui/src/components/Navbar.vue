<script setup lang="ts">
import client from '@/api/client';
import { onMounted, ref } from 'vue';
import { RouterLink } from 'vue-router'

const mangalVersion = ref("")
onMounted(async () => {
    const { data, error } = await client.GET('/mangalInfo', {})
    if (error) {
        return
    }

    mangalVersion.value = data.version
})
</script>

<template>
    <nav class="navbar shadow-sm sticky-top navbar-expand-lg bg-body-tertiary">
        <div class="container-fluid">
            <RouterLink class="navbar-brand" to="/">Mangal {{ mangalVersion }}</RouterLink>
            <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarSupportedContent"
                aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation">
                <span class="navbar-toggler-icon"></span>
            </button>
            <div class="collapse navbar-collapse" id="navarSupportedContent">
                <ul class="navbar-nav me-auto mb-2 mb-lg-0">
                    <li class="nav-item">
                        <RouterLink class="nav-link" to="/about">About</RouterLink>
                    </li>
                </ul>
            </div>
        </div>
    </nav>
</template>