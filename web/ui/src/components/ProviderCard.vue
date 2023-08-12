<script setup lang="ts">
import type { components } from '@/api/mangal';

const props = defineProps<{
    provider: components['schemas']['Provider']
    btnText?: string
}>()

const emit = defineEmits<{
    (event: "click"): void
}>()

function handleClick() {
    emit('click')
}
</script>

<template>
    <div class="card" style="width: 18rem;">
        <div class="card-body">
            <h5 class="card-title">{{ provider.name }}</h5>
            <p class="card-text">{{ provider.description ?? "No description" }}</p>
            <button @click="handleClick" class="btn btn-primary">{{ btnText ?? "Use" }}</button>
        </div>
        <div class="card-footer">
            Version {{ provider.version }}
        </div>
    </div>
</template>