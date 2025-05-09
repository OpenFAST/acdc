<script lang="ts" setup>
import { ref, watch, onMounted, computed } from 'vue'
import { useProjectStore } from '../project';
import ModelProp from './ModelProp.vue'
import { File, instanceOfField } from "../types"
const project = useProjectStore()

const selectedFile = ref<File>()
const selectedFileID = ref<number>()

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
                ID: options.length,
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
    for (let ed of files.ElastoDyn) {
        ed.YawDOF.Value = false
    }

    // BeamDyn file linearization defaults
    for (let bd of files.BeamDyn) {
        bd.RotStates.Value = true
    }

    // Set AeroDyn file linearization defaults
    for (let ad of files.AeroDyn) {
        ad.AFAeroMod.Value = 0
        ad.UAMod.Value = 0
        ad.TwrPotent.Value = 1
        ad.TwrShadow.Value = 1
        ad.FrozenWake.Value = false
        ad.SkewMod.Value = 1
    }

    // Set HydroDyn file linearization defaults
    for (let hd of files.HydroDyn) {
        hd.WaveMod.Value = 0
        hd.WvDiffQTF.Value = false
        hd.WvSumQTF.Value = false
        hd.ExctnMod.Value = 0
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
                <li class="list-group-item" v-for="(main, i) in project.model.Files.Main">
                    <div class="fw-bold">Main</div>
                    <div>
                        <ModelProp :field="main.TMax" />
                        <ModelProp :field="main.DT" />
                        <ModelProp :field="main.Gravity" />
                        <ModelProp :field="main.OutFmt" />
                        <ModelProp :field="main.Linearize" />
                        <ModelProp :field="main.CalcSteady" />
                        <ModelProp :field="main.TrimTol" />
                        <ModelProp :field="main.Twr_Kdmp" />
                        <ModelProp :field="main.Bld_Kdmp" />
                        <ModelProp :field="main.NLinTimes" />
                    </div>
                </li>
                <li class="list-group-item" v-for="(ed, i) in project.model.Files.ElastoDyn">
                    <div class="fw-bold">ElastoDyn {{ i + 1 }}</div>
                    <div>
                        <ModelProp :field="ed.ShftTilt" />
                        <ModelProp :field="ed.DrTrDOF" />
                        <ModelProp :field="ed.GenDOF" />
                        <ModelProp :field="ed.YawDOF" />
                        <ModelProp :field="ed.TwFADOF1" />
                        <ModelProp :field="ed.TwFADOF2" />
                        <ModelProp :field="ed.TwSSDOF1" />
                        <ModelProp :field="ed.TwSSDOF2" />
                        <ModelProp :field="ed.ShftTilt" />
                    </div>
                </li>
                <li class="list-group-item" v-for="(bd, i) in project.model.Files.BeamDyn">
                    <div class="fw-bold">BeamDyn {{ i + 1 }}</div>
                    <div>
                        <ModelProp :field="bd.RotStates" />
                    </div>
                </li>
                <li class="list-group-item" v-for="(ad, i) in project.model.Files.AeroDyn">
                    <div class="fw-bold">AeroDyn {{ i + 1 }}</div>
                    <div>
                        <ModelProp :field="ad.AFAeroMod" />
                        <ModelProp :field="ad.UAMod" />
                        <ModelProp :field="ad.TwrPotent" />
                        <ModelProp :field="ad.TwrShadow" />
                        <ModelProp :field="ad.FrozenWake" />
                        <ModelProp :field="ad.SkewMod" />
                    </div>
                </li>
                <li class="list-group-item" v-for="(hd, i) in project.model.Files.HydroDyn">
                    <div class="fw-bold">HydroDyn {{ i + 1 }}</div>
                    <div>
                        <ModelProp :field="hd.WaveMod" />
                        <ModelProp :field="hd.WvDiffQTF" />
                        <ModelProp :field="hd.WvSumQTF" />
                        <ModelProp :field="hd.ExctnMod" />
                    </div>
                </li>
                <li class="list-group-item" v-for="(srvd, i) in project.model.Files.ServoDyn">
                    <div class="fw-bold">ServoDyn {{ i + 1 }}</div>
                    <div>
                        <ModelProp :field="srvd.VS_Rgn2K" />
                        <ModelProp :field="srvd.VS_RtGnSp" />
                        <ModelProp :field="srvd.VS_RtTq" />
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
                        <select class="form-select" v-model="selectedFileID">
                            <option v-for="(item, i) in modelFileOptions" :value="item.ID">
                                {{ item.Type }} - {{ item.Name }}
                            </option>
                        </select>
                    </div>
                </div>
            </div>
            <hr class="my-0" v-if="selectedFileID" />
            <div class="card-body" v-if="selectedFileID">
                <ModelProp :field="field" v-for="field in modelFileOptions[selectedFileID].Fields" />
                <div class="text-center" v-if="modelFileOptions[selectedFileID].Fields.length == 0">No fields in file
                    can be modified</div>
            </div>
        </div>
    </main>
</template>

<style scoped></style>
