<script lang="ts" setup>
import { LOADED, LOADING, useProjectStore } from '../project';
import { reactive, onMounted } from 'vue'

const project = useProjectStore()

onMounted(() => { })

</script>

<template>
    <nav class="navbar navbar-expand-sm navbar-dark bg-primary fixed-top">
        <div class="container-fluid">
            <div class="collapse navbar-collapse" id="navbarNav">
                <ul class="navbar-nav me-auto">
                    <li class="nav-item">
                        <router-link class="nav-link" to="/">Project</router-link>
                    </li>
                    <li class="nav-item" v-if="project.status.project == LOADED">
                        <router-link class="nav-link" to="/model">Model</router-link>
                    </li>
                    <li class="nav-item" v-if="project.status.project == LOADED">
                        <router-link class="nav-link" to="/analysis">Analysis</router-link>
                    </li>
                    <li class="nav-item" v-if="project.status.project == LOADED">
                        <router-link class="nav-link" to="/evaluate">Evaluate</router-link>
                    </li>
                    <li class="nav-item" v-if="project.status.project == LOADED">
                        <router-link class="nav-link" to="/results">Results</router-link>
                    </li>
                </ul>
                <span class="navbar-text me-2" v-if="project.status.project == LOADING">
                    <div class="spinner-border spinner-border-sm text-light" role="status">
                        <span class="visually-hidden">Saving...</span>
                    </div>
                </span>
            </div>
        </div>
    </nav>
</template>

<style scoped></style>
