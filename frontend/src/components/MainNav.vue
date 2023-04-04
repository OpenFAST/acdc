<script lang="ts" setup>
import { useProjectStore } from '../project';
import { reactive, onMounted } from 'vue'

const project = useProjectStore()

onMounted(() => {


})

</script>

<template>
    <nav class="navbar navbar-expand-sm navbar-dark bg-primary fixed-top">
        <div class="container-fluid">
            <div class="collapse navbar-collapse" id="navbarNav">
                <ul class="navbar-nav">
                    <li class="nav-item">
                        <router-link class="nav-link" to="/">Home</router-link>
                    </li>
                    <li class="nav-item" v-if="project.loaded">
                        <router-link class="nav-link" to="/turbine">Turbine</router-link>
                    </li>
                    <li class="nav-item" v-if="project.loaded">
                        <router-link class="nav-link" to="/analyze">Analyze</router-link>
                    </li>
                    <li class="nav-item" v-if="project.loaded">
                        <router-link class="nav-link" to="/Results">Results</router-link>
                    </li>
                </ul>
            </div>
        </div>
    </nav>
</template>

<style scoped></style>
