<script lang="ts" setup>
import { ref, watch, onMounted } from 'vue'
import { useProjectStore } from '../project';
import ModelProp from './ModelProp.vue'
import { File } from "../types"
const project = useProjectStore()

const selectedFile = ref<File>()

onMounted(() => {
    if (project.modelFileOptions.length > 0) {
        selectedFile.value = project.modelFileOptions[0]
    }
})

watch(
    () => project.modelFileOptions,
    () => {
        if (project.modelFileOptions.length > 0) {
            selectedFile.value = project.modelFileOptions[0]
        } else {
            selectedFile.value = undefined
        }
    })

</script>

<template>
    <main>
        <div class="card mb-3">
            <div class="card-header d-flex justify-content-between align-items-center">
                <span>OpenFAST Turbine</span>
                <a class="btn btn-primary btn-sm" @click="project.importModel">Import</a>
            </div>
            <ul class="list-group list-group-flush" v-if="project.model.Files">
                <li class="list-group-item" v-if="project.model.ImportedPaths.length > 0">
                    <h6 class="card-title">Imported Files</h6>
                    <div v-for="path in project.model.ImportedPaths">{{ path }}</div>
                </li>
                <li class="list-group-item" v-if="project.model.Notes.length > 0">
                    <h6 class="card-title">Notes</h6>
                    <div v-for="note in project.model.Notes">{{ note }}</div>
                </li>
            </ul>
        </div>

        <div class="card mb-3" v-if="project.modelFileOptions.length > 0">
            <div class="card-header">Modify File</div>
            <div class="card-body">
                <div class="row">
                    <label for="fileSelect" class="col-sm-2 col-form-label">File</label>
                    <div class="col-sm-10">
                        <select class="form-select" v-model="selectedFile">
                            <option v-for="(item, i) in project.modelFileOptions" :value="item">
                                {{ item.Type }} - {{ item.Name }}
                            </option>
                        </select>
                    </div>
                </div>
            </div>
            <hr class="my-0" v-if="selectedFile" />
            <div class="card-body" v-if="selectedFile">
                <ModelProp :field="field" v-for="field in selectedFile.Fields" />
                <div class="text-center" v-if="selectedFile.Fields.length == 0">No fields in file can be modified</div>
            </div>
        </div>
    </main>
</template>

<style scoped></style>
