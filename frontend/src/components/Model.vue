<script lang="ts" setup>
import { ref, watch, onMounted, computed } from 'vue'
import { useProjectStore } from '../project';
import ModelProp from './ModelProp.vue'
import { File, instanceOfField } from "../types"
const project = useProjectStore()

const selectedFile = ref<File>()

onMounted(() => {
    project.fetchModel()
})

const modelFileOptions = computed<File[]>(() => {
    const options: File[] = []
    if (project.model == null || project.model.Files == null) return options
    for (const files of Object.values(project.model.Files) as File[][]) {
        if (files == null) continue
        for (const file of files) {
            options.push({
                Name: file.Name,
                Type: file.Type,
                Fields: Object.values(file).filter(instanceOfField),
            } as File)
        }
    }
    return options
})

function setDefaults() {

    if (project.model == null || project.model.Files == null) return

    // Get the files object
    let files = project.model.Files

    // Set main file linearization defaults
    if (files.Main.length > 0) {
        files.Main[0].Linearize.Value = true
        files.Main[0].CalcSteady.Value = true
        files.Main[0].TrimTol.Value = 0.001
        files.Main[0].Twr_Kdmp.Value = 100
        files.Main[0].Bld_Kdmp.Value = 100
        files.Main[0].NLinTimes.Value = 36
        files.Main[0].OutFmt.Value = "ES18.9E3"
    }

    // Set ElastoDyn file linearization defaults
    if (files.ElastoDyn.length > 0) {
        files.ElastoDyn[0].YawDOF.Value = false
    }

    // Set AeroDyn file linearization defaults
    if (files.AeroDyn.length > 0) {
        files.AeroDyn[0].AFAeroMod.Value = 1
        files.AeroDyn[0].TwrPotent.Value = 0
        files.AeroDyn[0].TwrShadow.Value = 0
        files.AeroDyn[0].FrozenWake.Value = true
        files.AeroDyn[0].SkewMod.Value = 0
    }

    // Set HydroDyn file linearization defaults
    if (files.HydroDyn.length > 0) {
        files.HydroDyn[0].WaveMod.Value = 0
        files.HydroDyn[0].WvDiffQTF.Value = false
        files.HydroDyn[0].WvSumQTF.Value = false
        files.HydroDyn[0].ExctnMod.Value = 0
    }

    // Save changes to model
    project.updateModel()
}

</script>

<template>
    <main>
        <div class="card mb-3">
            <div class="card-header hstack">
                <span>OpenFAST Model Files</span>
                <a class="btn btn-primary btn-sm ms-auto" @click="project.importModelDialog">Import</a>
            </div>
            <ul class="list-group list-group-flush" v-if="project.model != null">
                <li class="list-group-item" v-if="project.model.ImportedPaths.length > 0">
                    <div class="fw-bold mb-2">Imported Files</div>
                    <div class="row">
                        <div class="col-3 col-md-6" v-for="path in project.model.ImportedPaths">{{ path }}
                        </div>
                    </div>
                </li>
                <li class="list-group-item" v-if="project.model.Notes.length > 0">
                    <div class="fw-bold mb-2">Notes</div>
                    <div v-for="note in project.model.Notes">{{ note }}</div>
                </li>
            </ul>
        </div>

        <div class="card mb-3" v-if="project.model != null && project.model.Files != null">
            <div class="card-header hstack">
                <span>Linearization Quick Setup</span>
                <a class="btn btn-primary btn-sm ms-auto" @click="setDefaults">Defaults</a>
            </div>
            <ul class="list-group list-group-flush" v-if="project.model.Files != null">
                <li class="list-group-item">
                    <div class="fw-bold">Main</div>
                    <div>
                        <ModelProp :field="project.model.Files.Main[0].TMax" />
                        <ModelProp :field="project.model.Files.Main[0].DT" />
                        <ModelProp :field="project.model.Files.Main[0].Gravity" />
                        <ModelProp :field="project.model.Files.Main[0].OutFmt" />
                        <ModelProp :field="project.model.Files.Main[0].Linearize" />
                        <ModelProp :field="project.model.Files.Main[0].CalcSteady" />
                        <ModelProp :field="project.model.Files.Main[0].TrimTol" />
                        <ModelProp :field="project.model.Files.Main[0].Twr_Kdmp" />
                        <ModelProp :field="project.model.Files.Main[0].Bld_Kdmp" />
                        <ModelProp :field="project.model.Files.Main[0].NLinTimes" />
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
                        <ModelProp :field="project.model.Files.ServoDyn[0].VS_Rgn2K" />
                        <ModelProp :field="project.model.Files.ServoDyn[0].VS_RtGnSp" />
                        <ModelProp :field="project.model.Files.ServoDyn[0].VS_RtTq" />
                    </div>
                </li>
            </ul>
        </div>

        <div class="card mb-3" v-if="modelFileOptions.length > 0">
            <div class="card-header">Modify File</div>
            <div class="card-body">
                <div class="row">
                    <label for="fileSelect" class="col-sm-2 col-form-label">File</label>
                    <div class="col-sm-10">
                        <select class="form-select" v-model="selectedFile">
                            <option v-for="(item, i) in modelFileOptions" :value="item">
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
