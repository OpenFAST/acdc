<script lang="ts" setup>
import { ref, onMounted, computed } from 'vue'
import { useProjectStore } from '../project';
import { Scatter } from 'vue-chartjs'
import { Chart, ChartData, ChartOptions, LegendElement, Legend, ChartEvent, LegendItem, ActiveElement } from 'chart.js'

const freqChart = ref<typeof Scatter>()
const dampChart = ref<typeof Scatter>()

const xAxisWS = ref(true);

const project = useProjectStore()

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

const charts = computed(() => {

    const CD = project.results.CD
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
                <span>Linearization Files</span>
                <a class="btn btn-primary btn-sm" @click="project.openCaseDirectory">Import</a>
            </div>
        </div>

        <div class="card mb-3" v-if="project.results.CD != null">
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
        </div>

    </main>
</template>

<style scoped></style>
