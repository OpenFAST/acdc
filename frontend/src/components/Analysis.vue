<script lang="ts" setup>
import { reactive, ref, onMounted, computed } from 'vue'
import Case from './Case.vue'
import { useProjectStore } from '../project';
import { main } from '../../wailsjs/go/models';
import { Scatter } from 'vue-chartjs'
import { Chart as ChartJS, Title, Tooltip, Legend, LineElement, PointElement, CategoryScale, LinearScale, ChartData } from 'chart.js'

ChartJS.register(Title, Tooltip, Legend, LineElement, PointElement, CategoryScale, LinearScale)

const project = useProjectStore()
const selectedCaseID = ref(1)

function addCase() {
    project.addAnalysisCase().then((c) => {
        selectedCaseID.value = c.ID
    })
}

function removeCase(i: number) {
    if (project.analysis.Cases.length > 1) {
        project.removeAnalysisCase(selectedCaseID.value)
        selectedCaseID.value = 1
    }
}

</script>

<template>
    <main>
        <div class="card mb-3">
            <div class="card-header">Case</div>
            <div class="card-body">
                <div class="row">
                    <label for="selectedCase" class="col-sm-2 col-form-label">Select</label>
                    <div class="col-sm-6">
                        <select class="form-select" id="selectedCase" v-model="selectedCaseID">
                            <option :value="c.ID" v-for="c in project.analysis.Cases">{{ c.ID }} - {{ c.Name }}</option>
                        </select>
                    </div>
                    <div class="col-2 d-grid">
                        <button class="btn btn-primary" @click="addCase">Add</button>
                    </div>
                    <div class="col-2 d-grid">
                        <button class="btn btn-danger" @click="removeCase(selectedCaseID)"
                            :disabled="(project.analysis.Cases == null) || (project.analysis.Cases.length < 2)">Remove</button>
                    </div>
                </div>
            </div>
            <hr class="my-0" />
            <div class="card-body" v-if="project.analysis.Cases != null">
                <Case :Case="project.analysis.Cases[selectedCaseID - 1]" />
            </div>
        </div>
    </main>
</template>

<style scoped></style>
