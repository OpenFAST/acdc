<script lang="ts" setup>
import { reactive, onMounted, computed } from 'vue'
import { useProjectStore } from '../project';
import { Scatter } from 'vue-chartjs'
import { ChartData, ChartOptions } from 'chart.js'

const project = useProjectStore()
const freqChart = computed(() => {
    let d = {
        data: { datasets: [] } as ChartData<'scatter'>,
        options: {
            maintainAspectRatio: false,
            responsive: true,
            plugins: {
                legend: { position: 'right', align: 'start' }
            },
            scales: {
                x: {
                    title: {
                        display: true,
                        text: project.results.HasWind ? "Wind Speed (m/s)" : "Rotor Speed (RPM)",
                    }
                },
                y: {
                    title: {
                        display: true,
                        text: "Natural Frequency (Hz)",
                    }
                },
            },
        } as ChartOptions<"scatter">,
    }

    // Loop through mode sets
    const r = project.results
    for (const ms of r.ModeSets) {
        if (ms.FrequencyMean > 10) {
            continue
        }
        d.data.datasets.push({
            label: ms.Label,
            data: ms.Indices.map(ind => {
                const op = r.OPs[ind.OP]
                const mode = op.Modes[ind.Mode]
                return {
                    x: project.results.HasWind ? op.WindSpeed : op.RotSpeed,
                    y: mode.NaturalFreqHz,
                }
            }),
            showLine: true,
        })
    }

    return d
})

const dampingChart = computed(() => {
    let d = {
        data: { datasets: [] } as ChartData<'scatter'>,
        options: {
            maintainAspectRatio: false,
            responsive: true,
            plugins: {
                legend: { position: 'right', align: 'start' }
            },
            scales: {
                x: {
                    title: {
                        display: true,
                        text: project.results.HasWind ? "Wind Speed (m/s)" : "Rotor Speed (RPM)",
                    }
                },
                y: {
                    title: {
                        display: true,
                        text: "Damping Ratio (-)",
                    }
                },
            },
        } as ChartOptions<"scatter">,
    }

    // Loop through mode sets
    const r = project.results
    for (const ms of r.ModeSets) {
        if (ms.FrequencyMean > 10) {
            continue
        }
        d.data.datasets.push({
            label: ms.Label,
            data: ms.Indices.map(ind => {
                const op = r.OPs[ind.OP]
                const mode = op.Modes[ind.Mode]
                return {
                    x: project.results.HasWind ? op.WindSpeed : op.RotSpeed,
                    y: mode.DampingRatio,
                }
            }),
            showLine: true,
        })
    }

    return d
})


</script>

<template>
    <main>
        <div class="card mb-3">
            <div class="card-header d-flex justify-content-between align-items-center">
                <span>Case Directory</span>
                <a class="btn btn-primary btn-sm" @click="project.openCaseDirectory">Import</a>
            </div>
            <!-- <ul class="list-group list-group-flush" v-if="project.results.OPs != null">
                <li class="list-group-item" v-for="op in project.results.OPs">
                    <div v-for="path in op.Files">{{ path }}</div>
                </li>
            </ul> -->
        </div>

        <!-- {{ charts.data }} -->

        <div class="card mb-3" v-if="project.results.ModeSets != null">
            <div class="card-header">Campbell Diagram</div>
            <div class="card-body">
                <div style="height:350px; position: relative;">
                    <Scatter :options="freqChart.options" :data="freqChart.data" />
                </div>
                <div style="height:350px; position: relative;">
                    <Scatter :options="dampingChart.options" :data="dampingChart.data" />
                </div>
            </div>
        </div>

    </main>
</template>

<style scoped></style>
