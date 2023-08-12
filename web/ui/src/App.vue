<script setup lang="ts">
import { RouterView } from 'vue-router'
import Navbar from './components/Navbar.vue';
</script>

<template>
  <Navbar />

  <main class="container p-3">
    <RouterView />
  </main>
</template>