<script lang="ts" setup>
import { storeToRefs } from 'pinia';
import { ref, watch, onMounted } from 'vue'
import { useProjectStore } from '../project';
import ModelProp from './ModelProp.vue'
import { main as m } from "../../wailsjs/go/models"
import { File, Field } from "../types"
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
            <div class="card-header">OpenFAST Model</div>
            <div class="card-body">
                <a class="btn btn-primary" @click="project.importModel">Import</a>
            </div>
        </div>

        <div class="card mb-3" v-if="project.modelFileOptions">
            <div class="card-header">Modify Fields</div>
            <div class="card-body">
                <div class="row">
                    <label for="fileSelect" class="col-sm-2 col-form-label">File</label>
                    <div class="col-sm-10">
                        <select class="form-select" v-model="selectedFile">
                            <option v-for="item in project.modelFileOptions" :value="item">
                                {{ item.Type }} - {{ item.Name }}
                            </option>
                        </select>
                    </div>
                </div>
            </div>

            <div class="card-body" v-if="selectedFile">
                <ModelProp :field="field" v-for="field in selectedFile.Fields" />
                <div class="text-center" v-if="selectedFile.Fields.length == 0">No fields in file can be modified</div>
            </div>

            <!-- <ul class="list-group list-group-flush" v-if="selectedFile && selectedFile.Fields">
                <li v-for="field in selectedFile.Fields" class="list-group-item">
                    <ModelProp :field="field" />
                </li>
            </ul> -->
            <!-- <div class="card-body">
                <small>
                    <pre><code>{{ selectedFile?.Text }}</code></pre>
                </small>
            </div> -->
        </div>
    </main>
</template>

<style scoped></style>
