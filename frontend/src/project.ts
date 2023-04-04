import { defineStore } from 'pinia'
import { computed, ref, reactive } from 'vue'

import { Open, Create, Save } from '../wailsjs/go/main/Project'
import { main } from "../wailsjs/go/models"

export const useProjectStore = defineStore('project', () => {

    const loaded = ref(false)

    const info = reactive<main.Info>(new main.Info)

    // const name = ref('Eduardo')
    // const doubleCount = computed(() => count.value * 2)

    // function increment() {
    //     count.value++
    // }

    function $reset() {
        loaded.value = false
    }

    function create() {
        Create().then(result => {
            Object.assign(info, result)
            loaded.value = true
        }).catch(err => {
            console.log(err)
        })
    }

    function open() {
        Open().then(result => {
            Object.assign(info, result)
            loaded.value = true
        }).catch(err => {
            console.log(err)
        })
    }

    function save() {
        Save().then(result => {
            console.log(result)
            loaded.value = true
        }).catch(err => {
            console.log(err)
        })
    }

    // return { loaded, name, doubleCount, increment }
    return { loaded, info, $reset, open, create, save }
})
