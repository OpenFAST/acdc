<script lang="ts" setup>
import { reactive, onMounted } from 'vue'
import { useProjectStore } from '../project';
import { Open, Create } from '../../wailsjs/go/main/Project'

const project = useProjectStore()

onMounted(() => {


})



</script>

<template>
    <main>
        <div class="card">
            <div class="card-header">Project</div>
            <div class="card-body">

                <form>

                    <div class="mb-3">
                        <label for="ProjectName" class="form-label">Name</label>
                        <input type="text" :value="project.info.Name" class="form-control" id="ProjectName">
                    </div>

                    <div class="mb-3" v-if="project.loaded">
                        {{ project.info.Date }}
                    </div>

                    <div class="mb-3">
                        <button type="button" class="btn btn-primary" @click="project.create">Create</button>
                        <button type="button" class="btn btn-primary ms-3" @click="project.open">Open</button>
                        <button type="button" class="btn btn-success float-end" @click="project.save">Save</button>
                    </div>
                </form>
            </div>
        </div>
    </main>
</template>

<style scoped></style>
