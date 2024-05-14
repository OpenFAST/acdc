<script lang="ts" setup>
import { ref, onMounted, computed } from 'vue'
import { useProjectStore, LOADING, LOADED } from '../project';
import { Scatter } from 'vue-chartjs'
import { Chart, ChartData, ChartOptions, ChartEvent, ActiveElement } from 'chart.js'
import { ChartComponentRef } from "vue-chartjs"
import { main, diagram, viz } from "../../wailsjs/go/models"
import { ExportDiagramDataJSON } from "../../wailsjs/go/main/App"
import chroma from 'chroma-js'
import ModeViz from "./ModeViz.vue"

const project = useProjectStore()

onMounted(() => {
    project.fetchAnalysis()
})

const selectedOP = ref<main.OperatingPoint>()
const selectedLine = ref<diagram.Line | null>(null)
const targetLine = ref<diagram.Line | null>(null)
const selectedPoint = ref<diagram.Point | null>(null)
const freqChart = ref<ChartComponentRef<'scatter'> | null>(null)
const dampChart = ref<ChartComponentRef<'scatter'> | null>(null)
const showNodePaths = ref(true)
const xAxisWS = ref(true)
const rotorSpeedMods = [1, 3, 6, 9, 12, 15]
const vizScale = ref(20)
const vizScaleOptions = [0.5, 1, 2, 3, 5, 10, 20, 50, 75, 100, 150, 200, 300, 400, 500, 1000, 2000]

interface Graph {
    options: ChartOptions<'scatter'>
    data: ChartData<'scatter'>
}

function selectPoint(event: ChartEvent, elements: ActiveElement[], chart: Chart<"scatter">) {
    if (elements.length == 0 || project.diagram == null) return
    if (elements[0].datasetIndex >= project.diagram.Lines.length) return;
    selectedLine.value = project.diagram.Lines[elements[0].datasetIndex];
    selectedPoint.value = selectedLine.value.Points[elements[0].index];
}

function selectLine(line: diagram.Line) {
    selectedLine.value = line;
    selectedPoint.value = line.Points[0];
}

function toggleLineVisibility() {
    if (selectedLine.value == null) return
    selectedLine.value.Hidden = !selectedLine.value.Hidden;
    project.updateDiagram()
}

function swapModeLine() {

    // If the selected line, selected point, or swap line is null, return
    if (selectedLine.value == null || selectedPoint.value == null || targetLine.value == null) return

    // Get the operating point to swap
    const swapOP = selectedPoint.value.OpPtID

    // Partition points from each line to keep and swap
    const ptsSelKeep = selectedLine.value.Points.filter((p: diagram.Point) => p.OpPtID < swapOP)
    const ptsSelSwap = selectedLine.value.Points.filter((p: diagram.Point) => p.OpPtID >= swapOP)
    const ptsTarKeep = targetLine.value.Points.filter((p: diagram.Point) => p.OpPtID < swapOP)
    const ptsTarSwap = targetLine.value.Points.filter((p: diagram.Point) => p.OpPtID >= swapOP)

    // Concatenate points and assign to lines
    selectedLine.value.Points = ptsSelKeep.concat(ptsTarSwap)
    targetLine.value.Points = ptsTarKeep.concat(ptsSelSwap)

    // Update the line number on the moved points
    for (const p of selectedLine.value.Points) { p.Line = selectedLine.value.ID }
    for (const p of targetLine.value.Points) { p.Line = targetLine.value.ID }

    // Change selected line to the swap line
    selectedLine.value = targetLine.value

    // Clear swap line selection
    targetLine.value = null

    // Clear visualizations because the lines may have changed
    project.clearModeViz()

    // Save changes to diagram
    project.updateDiagram()
}

function exportDiagramDataJSON() {
    if (project.diagram == null) return
    ExportDiagramDataJSON(project.diagram).catch(err => {
        console.log(err)
    })
}

function getModeViz() {
    if (selectedPoint.value == null) return
    project.getModeViz(selectedPoint.value, vizScale.value)
}

function getLineViz() {
    if (selectedLine.value == null) return
    console.log(selectedLine.value)
    project.clearModeViz()
    console.log(selectedLine.value)
    let vizPromises = new Array<Promise<viz.ModeData>>
    for (const point of selectedLine.value.Points) {
        vizPromises.push(project.getModeViz(point, vizScale.value))
    }
    Promise.all(vizPromises).then(() => {
        project.modeViz.sort((a, b) => a.LineID != b.LineID ? a.LineID - b.LineID : a.OPID - b.OPID)
    }).catch((err) => { console.log(err) })
}

function setCurrentVizID(id: number) {
    project.currentVizID = id
    const lineID = project.modeViz[id].LineID
    if (project.diagram == null || lineID >= project.diagram.Lines.length) return
    const line = project.diagram.Lines[project.modeViz[id].LineID]
    const opID = project.modeViz[id].OPID
    const point = line.Points.find(p => p.OpPtID == opID)
    if (point === undefined) return
    selectedLine.value = line
    selectedPoint.value = point
}

const charts = computed(() => {
    let objs = new Array<Graph>;
    if (project.diagram == null) return objs
    const CD = project.diagram
    const xLabel = (xAxisWS && CD.HasWind) ? "Wind Speed (m/s)" : "Rotor Speed (RPM)"
    const xValues = (xAxisWS && CD.HasWind) ? CD.WindSpeeds : CD.RotSpeeds
    const freqMin = Math.min(...CD.Lines.filter(line => !line.Hidden).map(line => Math.min(...line.Points.map(p => p.NaturalFreqHz))))
    const freqMax = Math.max(...CD.Lines.filter(line => !line.Hidden).map(line => Math.max(...line.Points.map(p => p.NaturalFreqHz))))
    const dampMin = Math.min(...CD.Lines.filter(line => !line.Hidden).map(line => Math.min(...line.Points.map(p => p.DampingRatio))))
    const dampMax = Math.max(...CD.Lines.filter(line => !line.Hidden).map(line => Math.max(...line.Points.map(p => p.DampingRatio))))

    const configs = [
        { label: "Natural Frequency (Hz)", isNatFreq: true },
        { label: "Damping Ratio (-)", isNatfreq: false },
    ]

    const lineColors = chroma.cubehelix().lightness([0.4, 0.75]).rotations(2).scale().colors(CD.Lines.length)

    for (const cfg of configs) {

        let data = { datasets: [] } as ChartData<'scatter'>

        // Set line color if not defined
        for (let i = 0; i < CD.Lines.length; i++) {
            if (CD.Lines[i].Color == "") {
                CD.Lines[i].Color = lineColors[i % lineColors.length]
            }
        }

        // Loop through mode sets
        data.datasets = data.datasets.concat(CD.Lines.map((line, i) => ({
            label: line.Label,
            data: line.Points.map(p => ({
                x: (xAxisWS && CD.HasWind) ? p.WindSpeed : p.RotSpeed,
                y: cfg.isNatFreq ? p.NaturalFreqHz : p.DampingRatio,
            })),
            borderColor: line.Color,
            showLine: true,
            hidden: line.Hidden,
        })))

        // Loop through rotor speed multipliers
        if (cfg.isNatFreq) {
            data.datasets = data.datasets.concat(rotorSpeedMods.map(rsm => ({
                label: rsm + 'P',
                data: CD.RotSpeeds.map((RotSpeed, i) => {
                    return {
                        x: xValues[i],
                        y: RotSpeed / 60 * rsm,
                    }
                }),
                pointStyle: false,
                showLine: true,
                borderDash: [4, 6],
                borderColor: "slategray",
                borderWidth: 2
            })))
        }

        // Plot selected point and visualization points if one is selected
        if (selectedPoint.value != null) {
            const p = selectedPoint.value;
            data.datasets.push({
                label: 'selectedPoint',
                data: [{
                    x: (xAxisWS && CD.HasWind) ? p.WindSpeed : p.RotSpeed,
                    y: cfg.isNatFreq ? p.NaturalFreqHz : p.DampingRatio,
                }],
                pointStyle: 'crossRot',
                borderColor: 'red',
                pointRadius: 12,
                pointHoverRadius: 12,
                pointBorderWidth: 5,
                pointHoverBorderWidth: 5,
            });
        }

        // Add data options to array of chart objects
        objs.push({
            data: data,
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: { display: false },
                    tooltip: {
                        filter: (e) => {
                            // Exclude selected point from tooltip
                            return e.dataset.label != "selectedPoint"
                        }
                    },
                },
                scales: {
                    x: {
                        title: { display: true, text: xLabel, font: { size: 18 } },
                        min: Math.min(...xValues),
                        max: Math.max(...xValues),
                        ticks: { font: { size: 16 } }
                    },
                    y: {
                        title: { display: true, text: cfg.label, font: { size: 18 } },
                        ticks: { font: { size: 16 } },
                        min: 0.95 * (cfg.isNatFreq ? freqMin : dampMin),
                        max: 1.05 * (cfg.isNatFreq ? freqMax : dampMax),
                    },
                },
                onClick: selectPoint,
                interaction: {
                    mode: 'nearest'
                },
                animation: { duration: 0 }
            } as ChartOptions<"scatter">,
        })
    }

    return objs
})


</script>

<template>
    <main>


        <div class="card mb-3">
            <div class="card-header hstack">
                <span>Linearization Folder</span>
            </div>
            <div class="card-body">
                <div class="row row-cols-auto g-3">
                    <div class="col" v-for="c in project.analysis?.Cases">
                        <a class="btn btn-outline-primary w-100" @click="project.selectCaseLinDir(c.ID)">Case {{ c.ID
                            }}: {{ c.Name }} </a>
                    </div>
                    <div class="col">
                        <a class="btn btn-outline-primary w-100" @click="project.selectCustomLinDir()">Custom </a>
                    </div>
                </div>
                <div class="mt-3">
                    Selected folder: {{ project.linDir }}
                </div>
            </div>
        </div>

        <div class="card mb-3" v-if="project.linDir">
            <div class="card-header hstack">
                <span>Linearization Files</span>
                <a class="btn btn-primary ms-auto" @click="project.processLinDir()">Process</a>
            </div>

            <div v-if="project.status.results == LOADING" class="spinner-border text-primary my-3 mx-auto"
                role="status">
                <span class="visually-hidden">Loading...</span>
            </div>
            <div class="card-body" v-if="project.results != null && project.status.results != LOADING">
                <div class="row">
                    <label for="inputPassword" class="col-sm-3 col-form-label">Operating Point Data</label>
                    <div class="col-sm-9">
                        <select class="form-select" v-model="selectedOP">
                            <option :value="null">None</option>
                            <option v-for="op in project.results.OPs" :value="op">
                                {{ op.ID }} -
                                {{ project.results.HasWind ? `${op.WindSpeed.toPrecision(3)} m/s` :
                        `${op.RotSpeed.toPrecision(3)} RPM` }}
                            </option>
                        </select>
                    </div>
                </div>
                <table class="table table-bordered mt-4 mb-0 text-center table-sm" v-if="selectedOP != null">
                    <thead>
                        <tr>
                            <th scope="col">Mode</th>
                            <th scope="col">Natural Frequency (Hz)</th>
                            <th scope="col">Damped Frequency (Hz)</th>
                            <th scope="col">Damping Ratio (-)</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="mode in selectedOP?.Modes">
                            <td>{{ mode.ID + 1 }}</td>
                            <td>{{ mode.NaturalFreqHz.toPrecision(5) }}</td>
                            <td>{{ mode.DampedFreqHz.toPrecision(5) }}</td>
                            <td>{{ mode.DampingRatio.toExponential(3) }}</td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>

        <div class="card mb-3" v-if="project.results != null">
            <div class="card-header hstack">
                <span>Campbell Diagram</span>
                <a class="btn btn-primary ms-auto" v-if="project.diagram != null"
                    @click="exportDiagramDataJSON()">Export Data (.json)</a>
            </div>
            <div class="card-body">
                <form class="row row-cols-auto g-3 align-items-center">
                    <div class="col">
                        <div class="input-group">
                            <span class="input-group-text">Freq. Min/Max (Hz)</span>
                            <input type="text" class="form-control" id="minFreq"
                                v-model.number="project.diagramOptions.MinFreq">
                            <input type="text" class="form-control" id="maxFreq"
                                v-model.number="project.diagramOptions.MaxFreq">
                        </div>
                    </div>
                    <div class="col">
                        <div class="form-check">
                            <input class="form-check-input" type="checkbox" id="doCluster"
                                v-model="project.diagramOptions.Cluster">
                            <label class="form-check-label" for="doCluster">
                                Spectral Clustering
                            </label>
                        </div>
                    </div>
                    <div class="col">
                        <div class="form-check">
                            <input class="form-check-input" type="checkbox" id="filterStructural"
                                v-model="project.diagramOptions.FilterStruct">
                            <label class="form-check-label" for="filterStructural">
                                Filter Non-structural Modes
                            </label>
                        </div>
                    </div>

                    <div class="col-12">
                        <a class="btn btn-primary" @click="project.generateDiagram()">Generate</a>
                        <span class="ms-2" v-if="project.diagram != null">Warning: replaces current diagram</span>
                    </div>
                </form>
            </div>

            <div class="card-body border-top" v-if="project.status.diagram == LOADING">
                <div style="width:100%;height: 65vh">
                    <div class="spinner-border text-primary align-middle mx-auto" role="status">
                        <span class="visually-hidden">Loading...</span>
                    </div>
                </div>
            </div>

            <div class="card-body border-top" v-if="project.diagram != null">
                <div class="row">
                    <div class="col-sm-12 col-lg-6">
                        <div style="position: relative; height: 65vh">
                            <Scatter ref="freqChart" :options="charts[0].options" :data="charts[0].data" />
                        </div>
                    </div>
                    <div class="col-sm-12 col-lg-6">
                        <div style="position: relative; height: 65vh">
                            <Scatter ref="dampChart" :options="charts[1].options" :data="charts[1].data" />
                        </div>
                    </div>
                    <div class="col-12">
                        <div class="row row-cols-auto g-2 mt-2">
                            <div class="col" v-for="(line, i) in project.diagram.Lines.filter((line) => !line.Hidden)">
                                <a class="btn btn-outline-dark" @click="selectLine(line)">
                                    <div style="position: relative">
                                        <div style="display: block; height: 4px; width: 24px; position: absolute; top: 49%; left: 0"
                                            :style="{ 'background-color': line.Color }">
                                        </div>
                                        <div style="margin-left: 32px;">{{ line.Label }}</div>
                                    </div>
                                </a>
                            </div>
                        </div>
                    </div>
                    <div class="col-12">
                        <div class="row row-cols-auto g-2 mt-3">
                            <div class="col">
                                <label class="col-form-label">Hidden:</label>
                            </div>
                            <div class="col" v-for="(line, i) in project.diagram.Lines.filter((line) => line.Hidden)">
                                <a class="btn btn-outline-dark" @click="selectLine(line)">
                                    <div style="position: relative">
                                        <div style="display: block; height: 4px; width: 24px; position: absolute; top: 49%; left: 0"
                                            :style="{ 'background-color': line.Color }">
                                        </div>
                                        <div style="margin-left: 32px;">{{ line.Label }}</div>
                                    </div>
                                </a>
                            </div>
                        </div>
                    </div>

                </div>
            </div>
        </div>

        <div class="row row-cols-1 row-cols-xl-2 mb-3 g-3">
            <div class="col-12 col-xl-5">
                <div class="card h-100" v-if="selectedLine != null">
                    <div class="card-header hstack">
                        <span>Line</span>
                        <a class="btn btn-primary ms-auto" @click="toggleLineVisibility">
                            {{ selectedLine.Hidden ? "Show" : "Hide" }}
                        </a>
                    </div>
                    <div class="card-body">
                        <form class="row g-3">
                            <div class="col-4">
                                <label for="lineLabel" class="col-form-label">Label</label>
                                <input type="text" class="form-control" id="lineLabel" v-model="selectedLine.Label"
                                    @change="project.updateDiagram()">
                            </div>
                            <div class="col-4">
                                <label for="lineColor" class="col-form-label">Color</Label>
                                <input type="color" class="form-control form-control-color w-100" id="lineColor"
                                    v-model="selectedLine.Color" @change="project.updateDiagram()">
                            </div>
                            <div class="col-12" v-if="!selectedLine.Hidden">
                                <label for="linePoints" class="col-form-label">Select Point</Label>
                                <select class="form-select" id="linePoints" v-model="selectedPoint">
                                    <option v-for="point in selectedLine.Points" :value="point">OP: {{ point.OpPtID }},
                                        Rotor Speed: {{ point.RotSpeed.toFixed(2) }}, Wind Speed: {{
                        point.WindSpeed.toFixed(2) }}, Natural Frequency: {{
                        point.NaturalFreqHz.toFixed(3) }}
                                    </option>
                                </select>
                            </div>
                            <div class="col-12" v-if="!selectedLine.Hidden">
                                <a class="btn btn-primary ms-auto" @click="getLineViz()">
                                    Visualize Line Modes
                                </a>
                            </div>
                        </form>
                    </div>
                </div>
            </div>

            <div class="col-12 col-xl-7">
                <div class="card h-100" v-if="selectedPoint != null">
                    <div class="card-header hstack">
                        <span>Mode</span>
                        <a class="btn btn-primary ms-auto" @click="getModeViz()">
                            Visualize
                        </a>
                    </div>
                    <div class="card-body">
                        <form class="row row-cols-auto g-3">
                            <div class="col-3">
                                <label for="modeOP" class="col-form-label">Operating Point</label>
                                <input type="text" readonly class="form-control-plaintext" id="modeOP"
                                    :value="selectedPoint.OpPtID">
                            </div>
                            <div class="col-3">
                                <label for="modeID" class="col-form-label">Mode ID</Label>
                                <input type="text" class="form-control-plaintext" id="modeID"
                                    :value="selectedPoint.ModeID">
                            </div>
                            <div class="col-3">
                                <label for="modeRotSpeed" class="col-form-label">Rotor Speed (RPM)</Label>
                                <input type="text" class="form-control-plaintext" id="modeRotSpeed"
                                    :value="selectedPoint.RotSpeed.toFixed(3)">
                            </div>
                            <div class="col-3">
                                <label for="modeWindSpeed" class="col-form-label">Wind Speed (m/s)</Label>
                                <input type="text" class="form-control-plaintext" id="modeWindSpeed"
                                    :value="selectedPoint.WindSpeed.toFixed(3)">
                            </div>
                            <div class="col-3">
                                <label for="modeOP" class="col-form-label">Natural Freq. (Hz)</label>
                                <input type="text" class="form-control-plaintext" id="modeOP"
                                    :value="selectedPoint.NaturalFreqHz.toFixed(3)">
                            </div>
                            <div class="col-3">
                                <label for="modeID" class="col-form-label">Damped Freq. (Hz)</Label>
                                <input type="text" class="form-control-plaintext" id="modeID"
                                    :value="selectedPoint.DampedFreqHz.toFixed(3)">
                            </div>
                            <div class="col-3">
                                <label for="modeRotSpeed" class="col-form-label">Damping Ratio (%)</Label>
                                <input type="text" class="form-control-plaintext" id="modeRotSpeed"
                                    :value="selectedPoint.DampingRatio.toFixed(3)">
                            </div>
                            <div class="col-3">
                                <label for="vizScale" class="col-form-label">Visualization Scale</Label>
                                <select class="form-select" v-model.number="vizScale">
                                    <option :value="v" v-for="v in vizScaleOptions">{{ v }}</option>
                                </select>
                            </div>
                            <div class="col-12 hstack" v-if="project.diagram != null && selectedLine != null">
                                <label for="vizScale" class="col-form-label">Select line to swap mode</Label>
                                <select class="form-select ms-3 w-auto" v-model="targetLine">
                                    <option :value="line"
                                        v-for="line in project.diagram.Lines.filter(l => l.ID != selectedLine?.ID && l.Points.find(p => p.OpPtID == selectedPoint?.OpPtID))">
                                        {{ line.Label }}</option>
                                </select>
                                <a class="btn btn-primary ms-3" @click="swapModeLine()"
                                    :disabled="targetLine == null">Swap</a>
                            </div>
                        </form>
                    </div>
                </div>
            </div>
        </div>

        <div class="card mb-3"
            v-if="selectedPoint != null && project.modeViz.length > 0 && project.currentVizID >= 0 && project.diagram != null">
            <div class="card-header hstack">
                <span>Mode Visualization</span>

            </div>
            <div class="card-body">
                <div class="row">
                    <div class="col-10">
                        <div style="width:100%; height: 80vh">
                            <ModeViz :ModeData="project.modeViz[project.currentVizID]" :showNodePaths="showNodePaths">
                            </ModeViz>
                        </div>
                    </div>
                    <div class="col-2">
                        <div class="d-grid gap-2 mb-3">
                            <a class="btn btn-primary" @click="project.clearModeViz">Clear</a>
                            <a class="btn btn-primary" @click="showNodePaths = !showNodePaths">{{
                        showNodePaths ? 'Hide' : 'Show' }} Node Paths</a>
                        </div>
                        <div class="list-group">
                            <a class="list-group-item list-group-item-action" v-for="(mv, i) in project.modeViz"
                                :class="{ active: (i == project.currentVizID) }" @click="setCurrentVizID(i)">
                                {{ project.diagram.Lines[mv.LineID].Label }}, OP {{ mv.OPID }}
                            </a>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </main>
</template>

<style scoped></style>
