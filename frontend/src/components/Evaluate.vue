<script lang="ts" setup>
import { reactive } from 'vue'
import { useProjectStore } from '../project';
import { main } from '../../wailsjs/go/models';
import { EventsOn } from "../../wailsjs/runtime/runtime"

interface EvalLog {
    ID: number;
    Line: string;
}

const project = useProjectStore()
const data = reactive({
    caseID: 1,
    logID: 0,
    numCPUs: 1,
    logMap: new Map<number, string[]>()
})

// Setup listener for evaluation status updates
EventsOn("evalLog", (logEntry: EvalLog) => {
    console.log(logEntry)
    let log = data.logMap.get(logEntry.ID)
    if (!log) log = [] as string[]
    log.push(logEntry.Line)
    data.logMap.set(logEntry.ID, log)
})

function startEvaluate() {
    data.logMap.clear()
    project.startEvaluate(data.caseID, data.numCPUs)
}

function setLogID(id: number) {
    data.logID = id
}

function clearLogID() {
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
                    <label for="openfastExecutable" class="col-2 col-form-label">Executable</label>
                    <div class="col-10">
                        <div class="input-group">
                            <input type="text" :value="project.exec.Path" class="form-control" id="openfastExecutable"
                                aria-describedby="openfastExecutableHelp" readonly>
                            <button class="btn btn-outline-primary" type="button" id="openfastExecutable"
                                @click="project.selectExec">Browse</button>
                        </div>
                        <div id="openfastExecutableHelp" class="form-text">Path to OpenFAST executable</div>
                    </div>
                </div>
                <div class="row mt-3" v-if="project.exec.Version">
                    <label for="openfastVersion" class="col-2 col-form-label">Version</label>
                    <div class="col-10">
                        <textarea class="form-control" id="openfastVersion" readonly rows="8"
                            :value="project.exec.Version"></textarea>
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
                        <select class="form-select" id="selectedCase" v-model="data.caseID">
                            <option :value="c.ID" v-for="c in project.analysis.Cases">{{ c.ID }} - {{ c.Name }}</option>
                        </select>
                    </div>
                </div>
                <div class="row mb-3">
                    <label for="numCPUs" class="col-sm-2 col-form-label"># of CPUs</label>
                    <div class="col-sm-10">
                        <select class="form-select" id="numCPUs" v-model="data.numCPUs">
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
            <hr class="my-0" v-if="project.statusMap.size > 0" />
            <div class="card-body" v-if="project.statusMap.size > 0">
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
                        <tr v-for="[_, status] in project.statusMap">
                            <td class="text-center">{{ status.ID }}</td>
                            <td>{{ status.State }}</td>
                            <td>
                                <div v-if="status.Error" class="text-danger">Error: {{ status.Error }}</div>
                                <div v-else class="progress-stacked">
                                    <div class="progress" :style="{ width: (status.SimProgress * 0.5) + '%' }">
                                        <div class="progress-bar bg-primary">{{ status.SimProgress + "%" }} </div>
                                    </div>
                                    <div class="progress" :style="{ width: (status.LinProgress * 0.5) + '%' }">
                                        <div class="progress-bar bg-info">{{ status.LinProgress + "%" }} </div>
                                    </div>
                                    <!-- <div :class="{ 'progress-bar': true, 'bg-info': status.State == 'Linearization', 'bg-success': status.State == 'Complete', 'bg-warning': status.State == 'Canceled' }"
                                        :style="{ width: status.Progress + '%' }">{{ status.Progress }}%</div> -->
                                </div>
                            </td>
                            <td class="text-end">
                                <button class="btn btn-outline-primary btn-sm ms-3"
                                    @click="setLogID(status.ID)">Log</button>
                            </td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
        <div class="offcanvas offcanvas-bottom" :class="{ show: data.logID > 0 }" tabindex="-1" style="height: 50vh">
            <div class="offcanvas-header">
                <h5 class="offcanvas-title">Operating Point {{ data.logID }} Evaluation Log</h5>
                <button class="btn-close" @click="clearLogID"></button>
            </div>
            <div class="offcanvas-body">
                <pre><code v-for="line in data.logMap.get(data.logID)">{{ line + "\n" }}</code></pre>
            </div>
        </div>
    </main>
</template>

<style scoped></style>
