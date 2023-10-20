<script lang="ts" setup>
import { reactive, ref, onMounted, computed } from 'vue'
import { main } from '../../wailsjs/go/models';
import { useProjectStore } from '../project';
import { Scatter } from 'vue-chartjs'
import { ChartData, ChartOptions } from 'chart.js'

const project = useProjectStore()


const props = defineProps<{
    Case: main.Case
}>()


function updateAnalysis() {
    project.updateAnalysis()
}

function changeCurveSize(event: Event) {
    const newLength: number = parseInt((event.target as HTMLInputElement).value)
    const diff = newLength - props.Case.Curve.length
    if (diff < 0) {
        props.Case.Curve = props.Case.Curve.slice(0, newLength)
    } else if (diff > 0) {
        for (let i = 0; i < diff; i++) {
            props.Case.Curve.push(new main.Condition({ WindSpeed: 0, RotorSpeed: 0, BladePitch: 0 }))
        }
    }
    updateAnalysis()
}


const myChart = computed(() => {
    let d = {
        data: { datasets: [] } as ChartData<'scatter'>,
        options: {
            maintainAspectRatio: false,
            responsive: true,
            plugins: {
                legend: {
                    position: 'left',
                    align: 'start'
                }
            },
            scales: {
                x: {
                    title: {
                        display: true,
                        text: props.Case.IncludeAero ? "Wind Speed (m/s)" : "Rotor Speed (RPM)",
                    }
                },
            },
        } as ChartOptions<"scatter">,
    }
    if (props.Case.IncludeAero) {
        d.data.datasets.push({
            label: "Rotor Speed OP (RPM)",
            data: props.Case.OperatingPoints.map(p => ({ x: p.WindSpeed, y: p.RotorSpeed })),
            showLine: true,
            borderColor: 'limegreen',
            backgroundColor: 'limegreen',
            pointBackgroundColor: 'limegreen'
        }, {
            label: "Rotor Speed Curve (RPM)",
            data: props.Case.Curve.map(p => ({ x: p.WindSpeed, y: p.RotorSpeed })),
            pointRadius: 5,
            pointBorderWidth: 2,
            pointBorderColor: 'darkgreen',
            backgroundColor: 'darkgreen',
            pointStyle: 'crossRot',
        }, {
            label: "Blade Pitch OP (deg)",
            data: props.Case.OperatingPoints.map(p => ({ x: p.WindSpeed, y: p.BladePitch })),
            showLine: true,
            borderColor: 'dodgerblue',
            backgroundColor: 'dodgerblue',
            pointBackgroundColor: 'dodgerblue'
        }, {
            label: "Blade Pitch Curve (deg)",
            data: props.Case.Curve.map(p => ({ x: p.WindSpeed, y: p.BladePitch })),
            pointRadius: 5,
            pointBorderWidth: 2,
            pointBorderColor: 'darkblue',
            backgroundColor: 'darkblue',
            pointStyle: 'crossRot',
        })
    } else {
        d.data.datasets.push({
            label: "Blade Pitch OP (deg)",
            data: props.Case.OperatingPoints.map(p => ({ x: p.RotorSpeed, y: p.BladePitch })),
            showLine: true,
            borderColor: 'dodgerblue',
            backgroundColor: 'dodgerblue',
            pointBackgroundColor: 'dodgerblue'
        }, {
            label: "Blade Pitch Curve (deg)",
            data: props.Case.Curve.map(p => ({ x: p.RotorSpeed, y: p.BladePitch })),
            pointRadius: 5,
            pointBorderWidth: 2,
            pointBorderColor: 'darkblue',
            backgroundColor: 'darkblue',
            pointStyle: 'crossRot',
        })
    }
    return d
})

</script>

<template>
    <div>
        <div class="row mb-3">
            <label for="caseName" class="col-2 col-form-label">Name</label>
            <div class="col-10">
                <input class="form-control" id="caseName" v-model="Case.Name" @change="updateAnalysis" />
            </div>
        </div>
        <div class="row mb-3">
            <label for="includeAero" class="col-2 col-form-label">Effects</label>
            <div class="col-10">
                <select class="form-select" id="includeAero" v-model="Case.IncludeAero" @change="updateAnalysis">
                    <option :value="false">Structure</option>
                    <option :value="true" v-if="project.model.HasAero">Structure
                        + Aerodynamics</option>
                </select>
            </div>
        </div>
        <div class="row mb-3" v-if="Case.IncludeAero">
            <label class="col-2 col-form-label">Wind Speed (m/s)</label>
            <div class="col-10">
                <div class="input-group">
                    <span class="input-group-text">Min</span>
                    <input class="form-control" v-model.number="Case.WindSpeedRange.Min" @change="updateAnalysis" />
                    <span class="input-group-text">Max</span>
                    <input class="form-control" v-model.number="Case.WindSpeedRange.Max" @change="updateAnalysis" />
                    <span class="input-group-text">Number</span>
                    <select class="form-select" id="StructureCurveSize" v-model="Case.WindSpeedRange.Num"
                        @change="updateAnalysis">
                        <option :value="n" v-for="n in 30">{{ n }}</option>
                    </select>
                </div>
            </div>
        </div>
        <div class="row mb-3" v-else>
            <label class="col-2 col-form-label">Rotor Speed (RPM)</label>
            <div class="col-10">
                <div class="input-group">
                    <span class="input-group-text">Min</span>
                    <input class="form-control" v-model.number="Case.RotorSpeedRange.Min" @change="updateAnalysis" />
                    <span class="input-group-text">Max</span>
                    <input class="form-control" v-model.number="Case.RotorSpeedRange.Max" @change="updateAnalysis" />
                    <span class="input-group-text">Number</span>
                    <select class="form-select" id="StructureCurveSize" v-model="Case.RotorSpeedRange.Num"
                        @change="updateAnalysis">
                        <option :value="n" v-for="n in 30">{{ n }}</option>
                    </select>
                </div>
            </div>
        </div>
        <div class="row mb-3">
            <div class="col-2">
                <label for="CurveTable" class="col-form-label">Condition Curve</label>
                <select class="form-select" id="StructureCurveSize" :value="Case.Curve.length" @change="changeCurveSize">
                    <option :value="n + 1" v-for="n in 29">{{ n + 1 }} Points</option>
                </select>
            </div>
            <div class="col-10">
                <table class="table table-small table-borderless align-middle mb-0" id="CurveTable">
                    <thead>
                        <tr>
                            <td class="text-center">Point #</td>
                            <td class="text-center" v-if="Case.IncludeAero">Wind Speed
                                (m/s)</td>
                            <td class="text-center">Rotor Speed (RPM)</td>
                            <td class="text-center">Blade Pitch (&deg;)</td>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="(row, i) in Case.Curve">
                            <td class="text-center">{{ i + 1 }}</td>
                            <td v-if="Case.IncludeAero">
                                <input v-model.number="row.WindSpeed" class="form-control" @change="updateAnalysis" />
                            </td>
                            <td><input v-model.number="row.RotorSpeed" class="form-control" @change="updateAnalysis" /></td>
                            <td><input v-model.number="row.BladePitch" class="form-control" @change="updateAnalysis" /></td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>
        <hr />
        <div style="height:350px; position: relative;">
            <Scatter :options="myChart.options" :data="myChart.data" />
        </div>
    </div>
</template>

<style scoped></style>
