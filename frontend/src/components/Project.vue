<script lang="ts" setup>
import { reactive, onMounted } from 'vue'
import { useProjectStore } from '../project';
const project = useProjectStore()
</script>

<template>
    <main>
        <div class="card mb-3">
            <div class="card-header">Project</div>
            <div class="card-body">
                <p class="card-text" v-if="project.loaded">Open: {{ project.info.Path }}</p>
                <p class="card-text" v-if="project.loaded">Updated: {{ project.info.Date }}</p>
                <a class="btn btn-success" @click="project.saveDialog">{{ project.loaded ? "Save As" : "New" }}</a>
                <a class="btn btn-primary ms-3" @click="project.openDialog">Open</a>
                <a class="btn btn-success float-end" @click="project.save">Save</a>
            </div>
        </div>
        <div class="card" v-if="project.config.RecentProjects.length > 0">
            <div class="card-header">Recent</div>
            <ul class="list-group list-group-flush">
                <a href="#" class="list-group-item list-group-item-action" v-for="item in project.config.RecentProjects"
                    @click="project.open(item)">{{ item }}</a>
            </ul>
        </div>
    </main>
</template>

<style scoped></style>
