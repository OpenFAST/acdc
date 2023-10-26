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

function setDefaults() {

    // Get the files object
    let files = project.model.Files!

    // Set main file linearization defaults
    files.Main[0].Linearize.Value = true
    files.Main[0].CalcSteady.Value = true
    files.Main[0].TrimCase.Value = 2
    files.Main[0].TrimTol.Value = 0.001
    files.Main[0].TrimGain.Value = 100
    files.Main[0].Twr_Kdmp.Value = 100
    files.Main[0].Bld_Kdmp.Value = 100
    files.Main[0].NLinTimes.Value = 36
    files.Main[0].Gravity.Value = 0
    files.Main[0].OutFmt.Value = "ES16.9E2"
    files.Main[0].WrVTK.Value = 3
    files.Main[0].VTK_type.Value = 2

    // Set ElastoDyn file linearization defaults
    files.ElastoDyn[0].ShftTilt.Value = 0
    files.ElastoDyn[0].YawDOF.Value = false

    // Set AeroDyn file linearization defaults
    if (files.AeroDyn.length > 0) {
        files.AeroDyn[0].AFAeroMod.Value = 1
        files.AeroDyn[0].TwrPotent.Value = 0
        files.AeroDyn[0].TwrShadow.Value = 0
        files.AeroDyn[0].FrozenWake.Value = true
    }

    // Set ServoDyn file linearization defaults
    if (files.ServoDyn.length > 0) {
        files.ServoDyn[0].PCMode.Value = 0
        files.ServoDyn[0].VSContrl.Value = 1
        files.ServoDyn[0].HSSBrMode.Value = 0
        files.ServoDyn[0].YCMode.Value = 0
    }

    // Save changes to model
    project.updateModel()
}

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
                    <div class="fw-bold mb-2">Imported Files</div>
                    <div v-for="path in project.model.ImportedPaths">{{ path }}</div>
                </li>
                <li class="list-group-item" v-if="project.model.Notes.length > 0">
                    <div class="fw-bold mb-2">Notes</div>
                    <div v-for="note in project.model.Notes">{{ note }}</div>
                </li>
            </ul>
        </div>

        <div class="card mb-3" v-if="project.modelFileOptions.length > 0">
            <div class="card-header d-flex justify-content-between align-items-center">
                <span>Linearization Quick Setup</span>
                <a class="btn btn-primary btn-sm" @click="setDefaults">Defaults</a>
            </div>
            <ul class="list-group list-group-flush" v-if="project.model.Files">
                <li class="list-group-item">
                    <div class="fw-bold">Main</div>
                    <div>
                        <ModelProp :field="project.model.Files.Main[0].Linearize" />
                        <ModelProp :field="project.model.Files.Main[0].CalcSteady" />
                        <ModelProp :field="project.model.Files.Main[0].TrimCase" />
                        <ModelProp :field="project.model.Files.Main[0].TrimTol" />
                        <ModelProp :field="project.model.Files.Main[0].TrimGain" />
                        <ModelProp :field="project.model.Files.Main[0].Twr_Kdmp" />
                        <ModelProp :field="project.model.Files.Main[0].Bld_Kdmp" />
                        <ModelProp :field="project.model.Files.Main[0].NLinTimes" />
                        <ModelProp :field="project.model.Files.Main[0].Gravity" />
                    </div>
                </li>
                <li class="list-group-item">
                    <div class="fw-bold">ElastoDyn</div>
                    <div>
                        <ModelProp :field="project.model.Files.ElastoDyn[0].ShftTilt" />
                        <ModelProp :field="project.model.Files.ElastoDyn[0].YawDOF" />
                    </div>
                </li>
                <li class="list-group-item" v-if="project.model.Files.AeroDyn.length > 0">
                    <div class="fw-bold">AeroDyn</div>
                    <div>
                        <ModelProp :field="project.model.Files.AeroDyn[0].AFAeroMod" />
                        <ModelProp :field="project.model.Files.AeroDyn[0].TwrPotent" />
                        <ModelProp :field="project.model.Files.AeroDyn[0].TwrShadow" />
                        <ModelProp :field="project.model.Files.AeroDyn[0].FrozenWake" />
                        <ModelProp :field="project.model.Files.AeroDyn[0].SkewMod" />
                    </div>
                </li>
                <li class="list-group-item" v-if="project.model.Files.ServoDyn.length > 0">
                    <div class="fw-bold">ServoDyn</div>
                    <div>
                        <ModelProp :field="project.model.Files.ServoDyn[0].PCMode" />
                        <ModelProp :field="project.model.Files.ServoDyn[0].VSContrl" />
                        <ModelProp :field="project.model.Files.ServoDyn[0].HSSBrMode" />
                        <ModelProp :field="project.model.Files.ServoDyn[0].YCMode" />
                    </div>
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
