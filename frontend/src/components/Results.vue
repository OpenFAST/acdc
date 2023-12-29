<script lang="ts" setup>
import { ref, onMounted, computed, reactive } from 'vue'
import { useProjectStore } from '../project';
import { Scatter } from 'vue-chartjs'
import { Chart, ChartData, ChartOptions, LegendElement, Legend, ChartEvent, LegendItem, ActiveElement } from 'chart.js'
import { main } from '../../wailsjs/go/models';

const project = useProjectStore()

onMounted(() => {
    project.fetchResults()
})

const selectedOP = ref<main.OperatingPoint>()
const freqChart = ref<typeof Scatter>()
const dampChart = ref<typeof Scatter>()
const maxFreq = ref(0.0)
const doCluster = ref(false)
const xAxisWS = ref(true);
const rotorSpeedMods = [1, 3, 6, 9, 12, 15]

interface Graph {
    options: ChartOptions<'scatter'>
    data: ChartData<'scatter'>
}

function chartClick(event: ChartEvent, elements: ActiveElement[], chart: Chart<"scatter">) {
    console.log(elements)
    // chart.get
}

function toggleLine(index: number) {
    console.log("toggle", index)
    const charts = [freqChart.value!.chart, dampChart.value!.chart]
    for (const chart of charts) {
        chart.setDatasetVisibility(index, !chart.isDatasetVisible(index));
        chart.update();
    }
}

function importLinData() {
    project.openCaseDirDialog().then(results => {
        if (results.OPs.length > 0) {
            let maxRotSpeed = results.OPs[results.OPs.length - 1].RotSpeed;
            maxFreq.value = Number((maxRotSpeed / 60 * 15).toFixed(2));
        }
    })
}



const charts = computed(() => {

    const CD = project.diagram
    const xLabel = xAxisWS && CD.HasWind ? "Wind Speed (m/s)" : "Rotor Speed (RPM)"
    const xValues = xAxisWS && CD.HasWind ? CD.WindSpeeds : CD.RotSpeeds
    const freqMax = Math.max(...CD.Lines.filter(line => !line.Hide).map(line => Math.max(...line.Points.map(p => p.NaturalFreqHz))))
    const dampMax = Math.max(...CD.Lines.filter(line => !line.Hide).map(line => Math.max(...line.Points.map(p => p.DampingRatio))))

    let objs = new Array<Graph>;

    const configs = [
        { label: "Natural Frequency (Hz)", isNatFreq: true },
        { label: "Damping Ratio (-)", isNatfreq: false },
    ]

    for (const cfg of configs) {

        let data = { datasets: [] } as ChartData<'scatter'>

        // Loop through mode sets
        data.datasets = CD.Lines.map(line => ({
            label: line.ID + "",
            data: line.Points.map(p => ({
                x: xValues[p.OpPtID],
                y: cfg.isNatFreq ? p.NaturalFreqHz : p.DampingRatio,
            })),
            showLine: true,
        }))

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
                borderDash: [3, 3],
            })))
        }

        // Add data options to array of chart objects
        objs.push({
            data: data,
            options: {
                responsive: true,
                maintainAspectRatio: false,
                plugins: {
                    legend: { display: false, position: 'right', align: 'start' },
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
                        max: 1.05 * (cfg.isNatFreq ? freqMax : dampMax),
                    },
                },
                onClick: chartClick,
                interaction: {
                    mode: 'nearest'
                }
            } as ChartOptions<"scatter">,
        })
    }

    return objs
})



</script>

<template>
    <main>
        <div class="card mb-3">
            <div class="card-header d-flex justify-content-between align-items-center">
                <span>Linearization Data</span>
                <a class="btn btn-primary btn-sm" @click="importLinData">Import</a>
            </div>
            <div class="card-body" v-if="project.results.LinDir">
                <div class="mb-3 row">
                    <label for="case-dir" class="col-sm-2 col-form-label">Directory</label>
                    <div class="col-sm-10">
                        <input type="text" readonly class="form-control-plaintext" id="case-dir"
                            :value="project.results.LinDir">
                    </div>
                </div>
                <div class="mb-3 row">
                    <label for="inputPassword" class="col-sm-2 col-form-label">Operating Point</label>
                    <div class="col-sm-10">
                        <select class="form-select" v-model="selectedOP">
                            <option v-for="op in project.results.OPs" :value="op">
                                {{ op.ID + 1 }} -
                                {{ project.results.HasWind ? `${op.WindSpeed.toPrecision(3)} m/s` :
                                    `${op.RotSpeed.toPrecision(3)} RPM` }}
                            </option>
                        </select>
                    </div>
                </div>
                <table class="table table-bordered mb-0 text-center" v-if="selectedOP != null">
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

        <div class="card mb-3" v-if="project.results.OPs != null">
            <div class="card-header d-flex justify-content-between align-items-center">
                <span>Campbell Diagram</span>
            </div>
            <div class="card-body" v-if="project.results.LinDir">
                <div class="row g-3 align-items-center">
                    <div class="col-auto">
                        <label for="maxFreq" class="col-form-label">Max Frequency (Hz)</label>
                    </div>
                    <div class="col-auto">
                        <input type="text" class="form-control" id="maxFreq" v-model.number="maxFreq">
                    </div>
                    <div class="col-auto">
                        <label for="doCluster" class="col-form-label">Spectral Clustering</label>
                    </div>
                    <div class="col-auto">
                        <select class="form-select" v-model="doCluster">
                            <option :value="true" selected>Enable</option>
                            <option :value="false">Disable</option>
                        </select>
                    </div>
                    <div class="col-auto">
                        <a class="btn btn-primary btn-sm" @click="project.generateDiagram(maxFreq, doCluster)">Generate</a>
                    </div>
                </div>
            </div>

            <ul class="list-group list-group-flush" v-if="project.diagram.Lines != null">
                <li class="list-group-item d-flex justify-content-between">
                    <div style="position: relative; width: 50%; height: 65vh">
                        <Scatter ref="freqChart" :options="charts[0].options" :data="charts[0].data" />
                    </div>
                    <div style="position: relative; width: 50%; height: 65vh">
                        <Scatter ref="dampChart" :options="charts[1].options" :data="charts[1].data" />
                    </div>
                </li>
                <li class="list-group-item">
                    <table class="table table-sm table-borderless text-center mb-0">
                        <thead>
                            <tr>
                                <td>Show</td>
                                <td>Line</td>
                                <td>Label</td>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="line in project.diagram.Lines">
                                <td><input class="form-check-input" type="checkbox" v-model="line.Hide" :true-value="false"
                                        :false-value="true" @change="toggleLine(line.ID - 1)"></td>
                                <td>{{ line.ID }}</td>
                                <td><input class="form-control" v-model="line.Label" />
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </li>
            </ul>
        </div>

        <!-- <div class="card mb-3" v-if="project.results.CD != null">
            <div class="card-header d-flex justify-content-between align-items-center">
                <span>Campbell Diagram</span>
            </div>
            <ul class="list-group list-group-flush">
                <li class="list-group-item d-flex justify-content-between">
                    <div style="position: relative; width: 50%; height: 65vh">
                        <Scatter ref="freqChart" :options="charts[0].options" :data="charts[0].data" />
                    </div>
                    <div style="position: relative; width: 50%; height: 65vh">
                        <Scatter ref="dampChart" :options="charts[1].options" :data="charts[1].data" />
                    </div>
                </li>
                <li class="list-group-item">
                    <table class="table table-sm table-borderless text-center mb-0">
                        <thead>
                            <tr>
                                <td>Show</td>
                                <td>Line</td>
                                <td>Label</td>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="line in project.results.CD.Lines">
                                <td><input class="form-check-input" type="checkbox" v-model="line.Hide" :true-value="false"
                                        :false-value="true" @change="toggleLine(line.ID - 1)"></td>
                                <td>{{ line.ID }}</td>
                                <td><input class="form-control" v-model="line.Label" />
                                </td>
                            </tr>
                        </tbody>
                    </table>
                </li>
            </ul>
        </div> -->

    </main>
</template>

<style scoped></style>
