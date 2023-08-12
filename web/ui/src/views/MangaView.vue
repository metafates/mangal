<script setup lang="ts">
import { useRoute } from 'vue-router';

const route = useRoute()
</script>

<template>
    {{ $route.params }}
</template>