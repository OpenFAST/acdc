<script lang="ts" setup>
import { reactive, onMounted } from 'vue'
import { SelectExec } from '../../wailsjs/go/main/Project'


onMounted(() => {


})

const data = reactive({ execPath: '' })

function selectExec() {
    SelectExec().then(result => {
        data.execPath = result
        console.log(result)
    }).catch(err => {
        console.log(err)
    })
}

</script>

<template>
    <main>
        <div class="card">
            <div class="card-header">
                Settings
            </div>
            <div class="card-body">
                <div class="mb-3">
                    <label for="openfastExecutable" class="form-label">OpenFAST Executable</label>
                    <div class="input-group">
                        <input type="text" :value="data.execPath" class="form-control" id="openfastExecutable"
                            aria-describedby="openfastExecutableHelp" readonly>
                        <button class="btn btn-outline-primary" type="button" id="openfastExecutable"
                            @click="selectExec">Browse</button>
                    </div>
                    <div id="openfastExecutableHelp" class="form-text">Path to OpenFAST executable</div>
                </div>
            </div>
        </div>
    </main>
</template>

<style scoped></style>
