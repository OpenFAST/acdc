import { defineStore } from 'pinia'
import { computed, ref, reactive } from 'vue'
import { OpenProject, OpenProjectDialog, SaveProject, SaveProjectDialog, UpdateAnalysis } from '../wailsjs/go/main/App'
import { LoadConfig, SaveConfig } from "../wailsjs/go/main/App"
import { SelectExec, ImportModelDialog, UpdateModel } from "../wailsjs/go/main/App"
import { AddAnalysisCase, RemoveAnalysisCase } from "../wailsjs/go/main/App"
import { EvaluateLinearization, CancelEvaluate } from "../wailsjs/go/main/App"
import { OpenCaseDirectoryDialog } from "../wailsjs/go/main/App"
import { main } from "../wailsjs/go/models"
import { File, Field } from "./types"
import { EventsOn } from "../wailsjs/runtime/runtime"

export const useProjectStore = defineStore('project', () => {

    const saving = ref(false)
    const loaded = ref(false)
    const info = reactive<main.Info>(new main.Info)
    const exec = reactive<main.Exec>(new main.Exec)
    const model = reactive<main.Model>(new main.Model)
    const analysis = reactive<main.Analysis>(new main.Analysis)
    const results = reactive<main.Results>(new main.Results)
    const statusMap = reactive<Map<number, main.EvalStatus>>(new Map)
    const config = reactive<main.Config>(new main.Config)

    // Load config when store is initialized
    LoadConfig().then(result => {
        Object.assign(config, result)
    }).catch(err => {
        console.log(err)
    })

    function $reset() {
        loaded.value = false
    }

    function open(path: string) {
        OpenProject(path).then(result => {
            Object.assign(info, result.Info)
            Object.assign(exec, result.Exec)
            Object.assign(model, result.Model)
            Object.assign(analysis, result.Analysis)
            updateRecentProjects(path)
            loaded.value = true
        }).catch(err => {
            console.log(err)
        })
    }

    function openDialog() {
        OpenProjectDialog().then(result => {
            Object.assign(info, result.Info)
            Object.assign(exec, result.Exec)
            Object.assign(model, result.Model)
            Object.assign(analysis, result.Analysis)
            updateRecentProjects(info.Path)
            loaded.value = true
        }).catch(err => {
            console.log(err)
        })
    }

    function save() {
        saving.value = true
        SaveProject(info.Path).then(result => {
            Object.assign(info, result.Info)
            saving.value = false
        }).catch(err => {
            console.log(err)
        })
    }

    function saveDialog() {
        SaveProjectDialog().then(result => {
            Object.assign(info, result.Info)
            Object.assign(exec, result.Exec)
            Object.assign(model, result.Model)
            Object.assign(analysis, result.Analysis)
            updateRecentProjects(info.Path)
            loaded.value = true
        }).catch(err => {
            console.log(err)
        })
    }

    function selectExec() {
        SelectExec().then(result => {
            Object.assign(info, result.Info)
            Object.assign(exec, result.Exec)
        }).catch(err => {
            console.log(err)
        })
    }

    function updateRecentProjects(path: string) {
        // If path in recent, remove it
        const index = config.RecentProjects.indexOf(path)
        if (index > -1) config.RecentProjects.splice(index, 1)
        config.RecentProjects.unshift(path) // Prepend new path
        config.RecentProjects = config.RecentProjects.slice(0, 5)      // Limit to 5 items
        // Save config
        SaveConfig(config).catch(err => {
            console.log(err)
        })
    }

    function importModel() {
        ImportModelDialog().then(result => {
            Object.assign(info, result.Info)
            Object.assign(model, result.Model)
        }).catch(err => {
            console.log(err)
        })
    }

    function openCaseDirectory() {
        OpenCaseDirectoryDialog().then(result => {
            Object.assign(info, result.Info)
            Object.assign(results, result.Results)
        }).catch(err => {
            console.log(err)
        })
    }


    function updateModel() {
        UpdateModel(model).then(result => {
            Object.assign(info, result.Info)
        }).catch(err => {
            console.log(err)
        })
    }

    function updateAnalysis() {
        UpdateAnalysis(analysis).then(result => {
            Object.assign(info, result.Info)
            Object.assign(analysis, result.Analysis)
        }).catch(err => {
            console.log(err)
        })
    }

    function addAnalysisCase() {
        return new Promise<main.Case>((resolve, reject) => {
            AddAnalysisCase().then(result => {
                Object.assign(info, result.Info)
                Object.assign(analysis, result.Analysis)
                resolve(analysis.Cases[analysis.Cases.length - 1])
            }).catch(err => {
                console.log(err)
                reject(err)
            })
        })
    }

    function removeAnalysisCase(id: number) {
        RemoveAnalysisCase(id).then(result => {
            Object.assign(info, result.Info)
            Object.assign(analysis, result.Analysis)
        }).catch(err => {
            console.log(err)
        })
    }

    function instanceOfField(obj: any): obj is Field {
        return typeof obj == 'object' && 'Name' in obj && 'Type' in obj;
    }

    const modelFileOptions = computed<File[]>(() => {
        const options: File[] = []
        if (!model.Files) return options
        for (const files of Object.values(model.Files) as File[][]) {
            for (const file of files) {
                options.push({
                    Name: file.Name,
                    Type: file.Type,
                    Fields: Object.values(file).filter(instanceOfField),
                } as File)
            }
        }
        return options
    })

    // Setup listener for evaluation status updates
    EventsOn("evalStatus", (status: main.EvalStatus) => {
        statusMap.set(status.ID, status)
    })

    function startEvaluate(caseID: number, numCPUs: number) {
        EvaluateLinearization(analysis.Cases[caseID - 1], numCPUs).then(result => {
            statusMap.clear()
            for (const status of result) {
                statusMap.set(status.ID, status)
            }
        }).catch(err => {
            console.log(err)
        })
    }

    function cancelEvaluate() {
        CancelEvaluate().catch(err => {
            console.log(err)
        })
    }

    return {
        loaded, saving, config, info, exec, model, analysis, results, statusMap, modelFileOptions,
        $reset, open, saveDialog, save, openDialog, selectExec,
        importModel, updateModel,
        openCaseDirectory, updateAnalysis, addAnalysisCase, removeAnalysisCase,
        startEvaluate, cancelEvaluate
    }
})
