<script lang="ts" setup>
import { reactive, ref, onMounted, computed } from 'vue'
import Case from './Case.vue'
import { useProjectStore } from '../project';
import { main } from '../../wailsjs/go/models';
import { Scatter } from 'vue-chartjs'
import { Chart as ChartJS, Title, Tooltip, Legend, LineElement, PointElement, CategoryScale, LinearScale, ChartData } from 'chart.js'

ChartJS.register(Title, Tooltip, Legend, LineElement, PointElement, CategoryScale, LinearScale)

const project = useProjectStore()

onMounted(() => {
    project.fetchAnalysis()
})

function removeCase() {
    if (project.analysis == null || project.currentCaseID == null) return
    project.removeAnalysisCase(project.currentCaseID)
    project.currentCaseID = 0
}

</script>

<template>
    <main>
        <div class="card mb-3">
            <div class="card-header">Case</div>
            <div class="card-body">
                <div class="row">
                    <label for="currentCaseID" class="col-sm-2 col-form-label">Select</label>
                    <div class="col-sm-6">
                        <select class="form-select" id="currentCaseID" v-model="project.currentCaseID"
                            v-if="project.analysis != null">
                            <option v-for="c in project.analysis.Cases" :value="c.ID">{{ c.ID }} - {{ c.Name }}</option>
                        </select>
                    </div>
                    <div class="col-2 d-grid">
                        <button class="btn btn-primary" @click="project.addAnalysisCase()">Add</button>
                    </div>
                    <div class="col-2 d-grid">
                        <button class="btn btn-danger" @click="removeCase()"
                            :disabled="(project.analysis == null) || (project.analysis.Cases.length < 2)">Remove</button>
                    </div>
                </div>
            </div>
            <hr class="my-0" />
            <div class="card-body" v-if="project.analysis != null">
                <Case :Case="project.analysis.Cases[project.currentCaseID - 1]" />
            </div>
        </div>
    </main>
</template>

<style scoped></style>
