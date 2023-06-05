<script lang="ts" setup>
import { reactive, onMounted } from 'vue'
import { useProjectStore } from '../project';
const project = useProjectStore()
</script>

<template>
    <main>
        <div class="card">
            <div class="card-header">
                Settings
            </div>
            <!-- <div class="card-body">
                <div class="mb-3">
                    <label for="openfastExecutable" class="form-label">OpenFAST Executable</label>
                    <div class="input-group">
                        <input type="text" :value="project.analyze.Exec.Path" class="form-control" id="openfastExecutable"
                            aria-describedby="openfastExecutableHelp" readonly>
                        <button class="btn btn-outline-primary" type="button" id="openfastExecutable"
                            @click="project.selectExec">Browse</button>
                    </div>
                    <div id="openfastExecutableHelp" class="form-text">Path to OpenFAST executable</div>
                </div>
            </div>
            <div class="card-body">
                <h5 class="card-title">OpenFAST Version</h5>
                <pre><code>{{ project.analyze.Exec.Info }}</code></pre>
            </div> -->

        </div>
    </main>
</template>

<style scoped></style>
