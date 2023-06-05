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

        <div class="card mt-3" v-if="project.loaded">
            <div class="card-header">
                OpenFAST
            </div>
            <div class="card-body">
                <div class="mb-3">
                    <label for="openfastExecutable" class="form-label">Executable</label>
                    <div class="input-group">
                        <input type="text" :value="project.exec.Path" class="form-control" id="openfastExecutable"
                            aria-describedby="openfastExecutableHelp" readonly>
                        <button class="btn btn-outline-primary" type="button" id="openfastExecutable"
                            @click="project.selectExec">Browse</button>
                    </div>
                    <div id="openfastExecutableHelp" class="form-text">Path to OpenFAST executable</div>
                </div>

                <div class="mb-3">
                    <label for="openfastVersion" class="form-label">Version</label>
                    <textarea class="form-control" id="openfastVersion" readonly rows="18"
                        :value="project.exec.Version"></textarea>
                </div>

            </div>
        </div>
    </main>
</template>

<style scoped></style>
