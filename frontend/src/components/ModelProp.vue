<script lang="ts" setup>
import { reactive, onMounted, PropType, ref } from 'vue'
import { main as m } from "../../wailsjs/go/models"
import { Field } from '../types'

import { useProjectStore } from '../project';
const project = useProjectStore()

defineProps<{ field: Field }>()

function isInteger(field: Field): field is m.Integer { return field.Type === 'Integer' }
function isBool(field: Field): field is m.Bool { return field.Type === 'Bool' }
function isString(field: Field): field is m.String { return field.Type === 'String' }
function isPath(field: Field): field is m.Path { return field.Type === 'Path' }
function isPaths(field: Field): field is m.Paths { return field.Type === 'Paths' }
function isReal(field: Field): field is m.Real { return field.Type === 'Real' }
function isReals(field: Field): field is m.Reals { return field.Type === 'Reals' }

</script>

<template>
    <div class="row mb-2 g-3 align-items-center" v-if="field && field.Line > 0">
        <label :for="field.Name" class="col-sm-2 col-form-label">{{ field.Name }}</label>
        <div class="col-2">
            <div v-if="isPath(field)">
                <input type="text" class="form-control-plaintext" :id="field.Name" v-model="field.Value" readonly />
            </div>
            <div v-else-if="isPaths(field)">
                <input v-for="(_, i) in field.Value" type="text" class="form-control-plaintext"
                    :id="field.Name + (i == 0 ? '' : i)" v-model="field.Value[i]" readonly
                    :class="i == 0 ? '' : 'mt-1'" />
            </div>
            <div v-else-if="isBool(field)">
                <select class="form-select" :id="field.Name" v-model="field.Value" @change="project.updateModel">
                    <option :value="true">True</option>
                    <option :value="false">False</option>
                </select>
            </div>
            <div v-if="isString(field)">
                <input type="text" class="form-control" :id="field.Name" v-model="field.Value"
                    @change="project.updateModel" />
            </div>
            <div v-else-if="isInteger(field)">
                <input type="text" class="form-control" :id="field.Name" v-model.number="field.Value"
                    @change="project.updateModel" />
            </div>
            <div v-else-if="isReal(field)">
                <input type="text" class="form-control" :id="field.Name" v-model.number="field.Value"
                    @change="project.updateModel" />
            </div>
            <div v-if="isReals(field)">
                <input v-for="(_, i) in field.Value" type="text" class="form-control"
                    :id="field.Name + (i == 0 ? '' : i)" v-model.number="field.Value[i]" :class="i == 0 ? '' : 'mt-1'"
                    @change="project.updateModel" />
            </div>
        </div>
        <div class="col-8"><span class="form-text">{{ field.Desc }}</span></div>
    </div>
</template>

<style scoped></style>
