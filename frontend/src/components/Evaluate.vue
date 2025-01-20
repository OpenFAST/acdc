<script lang="ts" setup>
import { reactive, onMounted } from 'vue'
import { useProjectStore } from '../project';
import { main } from '../../wailsjs/go/models';
import { GetEvaluateLog } from "../../wailsjs/go/main/App"

const project = useProjectStore()
const data = reactive({
    logID: -1,
    logContents: "",
})

onMounted(() => {
    project.fetchAnalysis()
    project.fetchEvaluate()
})

function startEvaluate() {
    project.clearEvalStatus()
    project.startEvaluate(project.currentCaseID)
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
    data.logID = -1
}

</script>

<template>
    <main>
        <div class="card mb-3" v-if="project.evaluate != null">
            <div class="card-header">
                OpenFAST Executable
            </div>
            <div class="card-body">
                <div class="row">
                    <label for="executable" class="col-2 col-form-label">Path</label>
                    <div class="col-10">
                        <div class="input-group">
                            <input type="text" :value="project.evaluate.ExecPath" class="form-control" id="executable"
                                aria-describedby="executableHelp" readonly>
                            <a class="btn btn-outline-primary" id="executable" @click="project.selectExec">Browse</a>
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

        <div class="card mb-3" v-if="project.analysis != null && project.evaluate != null">
            <div class="card-header">Evaluate Case</div>
            <div class="card-body">
                <div class="hstack gap-3">
                    <label for="currentCaseID" class="col-form-label">Case</label>
                    <select class="form-select w-25" id="currentCaseID" v-model="project.currentCaseID">
                        <option :value="c.ID" v-for="c in project.analysis.Cases">{{ c.ID }} - {{ c.Name }}</option>
                    </select>
                    <label class="ms-3 col-form-label" for="numCPUs">CPUs:</label>
                    <input type="range" class="form-range w-25" min="1" :max="project.evaluate.MaxCPUs" id="numCPUs"
                        v-model.number="project.evaluate.NumCPUs" @change="project.updateEvaluate()">
                    <label class="col-form-label">{{ project.evaluate.NumCPUs }}</label>
                    <input class="form-check-input" type="checkbox" value="" id="fileOnly-checkbox"
                        v-model="project.evaluate.FilesOnly"  @change="project.updateEvaluate()">
                    <label class="form-check-label" for="fileOnly-checkbox">
                        Files Only
                    </label>
                    <a class="btn btn-success ms-auto" @click="startEvaluate">Start</a>
                    <a class="btn btn-danger" @click="project.cancelEvaluate()">Cancel</a>
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
                                <a class="btn btn-outline-primary btn-sm ms-3" @click="getLog(stat)"
                                    :disabled="stat.LogPath == ''">Log</a>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
        <div class="offcanvas offcanvas-bottom" :class="{ show: data.logID >= 0 }" tabindex="-1" style="height: 50vh">
            <div class="offcanvas-header">
                <h5 class="offcanvas-title">Operating Point {{ data.logID }} Evaluation Log</h5>
                <a class="btn-close" @click="closeLog"></a>
            </div>
            <div class="offcanvas-body">
                <pre><code>{{ data.logContents }}</code></pre>
            </div>
        </div>
    </main>
</template>

<style scoped></style>
