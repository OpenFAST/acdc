<script lang="ts" setup>
import { reactive, onMounted } from 'vue'
import { useProjectStore } from '../project';
import { LOADED } from '../project';
const project = useProjectStore()
</script>

<template>
    <main>
        <div class="card mb-3">
            <div class="card-header">Project</div>
            <div class="card-body">
                <div class="hstack mb-3" v-if="project.info != null">
                    <div>File: {{ project.info.Path }}</div>
                    <div class="ms-auto">Updated: {{ project.info.Date }}</div>
                </div>
                <a class="btn btn-success" @click="project.saveDialog">{{ project.status.project == LOADED ? "Save As" :
                    "New" }}</a>
                <a class="btn btn-primary ms-3" @click="project.openDialog">Open</a>
            </div>
        </div>
        <div class="card" v-if="project.config != null">
            <div class="card-header">Recent</div>
            <ul class="list-group list-group-flush">
                <a href="#" class="list-group-item list-group-item-action" v-for="item in project.config.RecentProjects"
                    @click="project.open(item)">{{ item }}</a>
            </ul>
        </div>
    </main>
</template>

<style scoped></style>
