<script lang="ts" setup>
import { reactive, onMounted } from 'vue'
import { useProjectStore } from '../project';
import { main } from '../../wailsjs/go/models';
import { GetEvaluateLog } from "../../wailsjs/go/main/App"

const project = useProjectStore()
const data = reactive({
    logID: 0,
    logContents: "",
})

onMounted(() => {
    project.fetchEvaluate()
    project.fetchAnalysis()
})

function startEvaluate() {
    project.startEvaluate(project.evalCaseID)
}

function getLog(status: main.EvalStatus) {
    GetEvaluateLog(status.LogPath).then((result) => {
        data.logContents = result
        data.logID = status.ID
    }).catch((err) => {
        console.log(err)
    })
}

function closeLog() {
    data.logID = 0
}

</script>

<template>
    <main>
        <div class="card mb-3">
            <div class="card-header">
                OpenFAST
            </div>
            <div class="card-body">
                <div class="row">
                    <label for="executable" class="col-2 col-form-label">Executable</label>
                    <div class="col-10">
                        <div class="input-group">
                            <input type="text" :value="project.evaluate.ExecPath" class="form-control" id="executable"
                                aria-describedby="executableHelp" readonly>
                            <button class="btn btn-outline-primary" type="button" id="executable"
                                @click="project.selectExec">Browse</button>
                        </div>
                        <div id="executableHelp" class="form-text">Path to OpenFAST executable</div>
                    </div>
                </div>
                <div class="row mt-3" v-if="project.evaluate.ExecVersion">
                    <label for="execVersion" class="col-2 col-form-label">Version</label>
                    <div class="col-10">
                        <textarea class="form-control" id="execVersion" readonly rows="8"
                            :value="project.evaluate.ExecVersion"></textarea>
                    </div>
                </div>
            </div>
        </div>

        <div class="card mb-3">
            <div class="card-header">Evaluate Case</div>
            <div class="card-body">
                <div class="row mb-3">
                    <label for="selectedCase" class="col-sm-2 col-form-label">Select Case</label>
                    <div class="col-sm-10">
                        <select class="form-select" id="selectedCase" v-model="project.evalCaseID">
                            <option :value="c.ID" v-for="c in project.analysis.Cases">{{ c.ID }} - {{ c.Name }}</option>
                        </select>
                    </div>
                </div>
                <div class="row mb-3">
                    <label for="numCPUs" class="col-sm-2 col-form-label"># of CPUs</label>
                    <div class="col-sm-10">
                        <select class="form-select" id="numCPUs" v-model="project.evaluate.NumCPUs"
                            @change="project.updateEvaluate()">
                            <option :value="n" v-for="n in [1, 2, 4, 6, 8, 12, 16, 24]">{{ n }}</option>
                        </select>
                    </div>
                </div>
                <div class="row">
                    <div class="col-2"></div>
                    <div class="col-10">
                        <button class="btn btn-success me-3" @click="startEvaluate">Start</button>
                        <button class="btn btn-danger" @click="project.cancelEvaluate()">Cancel</button>
                    </div>
                </div>

            </div>
            <hr class="my-0" v-if="project.evalStatus.length > 0" />
            <div class="card-body" v-if="project.evalStatus.length > 0">
                <table class="table table-sm table-borderless align-middle mb-0">
                    <thead>
                        <tr>
                            <th scope="col" class="text-center">OP</th>
                            <th scope="col" width="15%">State</th>
                            <th scope="col" width="80%">Progress</th>
                            <th scope="col"></th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="stat in project.evalStatus">
                            <td class="text-center">{{ stat.ID }}</td>
                            <td>{{ stat.State }}</td>
                            <td>
                                <div v-if="stat.Error" class="text-danger">Error: {{ stat.Error }}</div>
                                <div v-else class="progress-stacked">
                                    <div class="progress" :style="{ width: (stat.SimProgress * 0.5) + '%' }">
                                        <div class="progress-bar bg-primary">{{ stat.SimProgress + "%" }} </div>
                                    </div>
                                    <div class="progress" :style="{ width: (stat.LinProgress * 0.5) + '%' }">
                                        <div class="progress-bar bg-info">{{ stat.LinProgress + "%" }} </div>
                                    </div>
                                </div>
                            </td>
                            <td class="text-end">
                                <button class="btn btn-outline-primary btn-sm ms-3" @click="getLog(stat)"
                                    :disabled="stat.LogPath == ''">Log</button>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
        <div class="offcanvas offcanvas-bottom" :class="{ show: data.logID > 0 }" tabindex="-1" style="height: 50vh">
            <div class="offcanvas-header">
                <h5 class="offcanvas-title">Operating Point {{ data.logID }} Evaluation Log</h5>
                <button class="btn-close" @click="closeLog"></button>
            </div>
            <div class="offcanvas-body">
                <pre><code>{{ data.logContents }}</code></pre>
            </div>
        </div>
    </main>
</template>

<style scoped></style>
