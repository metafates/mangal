import { defineStore } from "pinia";
import { ref } from "vue";

export const useMangalStore = defineStore('mangal', () => {
    const providerID = ref("")
    const searchQuery = ref("")
    const mangaID = ref("")
    const volumeID = ref("")
    const chapterID = ref("")

    return {
        providerID,
        searchQuery,
        mangaID,
        volumeID,
        chapterID,
    }
})