import { defineStore } from 'pinia'
import { ref, reactive } from 'vue'
import { GenerateDiagram, LoadConfig, SaveConfig } from "../wailsjs/go/main/App"
import { OpenProjectDialog, SaveProjectDialog, OpenProject } from '../wailsjs/go/main/App'
import { FetchModel, UpdateModel, ImportModelDialog } from "../wailsjs/go/main/App"
import { FetchAnalysis, UpdateAnalysis, AddAnalysisCase, RemoveAnalysisCase } from "../wailsjs/go/main/App"
import { FetchEvaluate, UpdateEvaluate, SelectExec, EvaluateCase, GetEvaluateLog, CancelEvaluate } from "../wailsjs/go/main/App"
import { FetchResults, OpenCaseDirDialog } from "../wailsjs/go/main/App"
import { GetModeViz } from "../wailsjs/go/main/App"
// import { EvaluateLinearization, CancelEvaluate } from "../wailsjs/go/main/App"
import { main, diagram as diag, viz } from "../wailsjs/go/models"
import { EventsOn } from "../wailsjs/runtime/runtime"
import { LogError } from '../wailsjs/runtime/runtime'
import Project from './components/Project.vue'

export const NOT_LOADED = 0;
export const LOADING = 1;
export const LOADED = 2;

// Contains loading status of the various project components
class Loading {
    project: number = 0;
    config: number = 0;
    info: number = 0;
    model: number = 0;
    analysis: number = 0;
    evaluate: number = 0;
    results: number = 0;
    diagram: number = 0;
    viz: number = 0;
}

export const useProjectStore = defineStore('project', () => {

    const config = ref<main.Config | null>(null)
    const info = ref<main.Info | null>(null)
    const evalStatus = reactive<Array<main.EvalStatus>>(new Array)
    const diagram = reactive<diag.Diagram>(new diag.Diagram)
    const status = reactive<Loading>(new Loading)
    const modeViz = reactive<Array<viz.ModeData>>(new Array)

    // Load config when store is initialized
    LoadConfig().then(result => {
        config.value = result
    }).catch(err => {
        LogError(err)
        console.log(err)
    })

    //--------------------------------------------------------------------------
    // Project
    //--------------------------------------------------------------------------

    function openDialog() {
        OpenProjectDialog().then(result => {
            info.value = result
            updateRecentProjects(info.value.Path)
            status.project = LOADED
        }).catch(err => {
            LogError(err)
            console.log(err)
        })
    }

    function saveDialog() {
        SaveProjectDialog().then(result => {
            info.value = result
            updateRecentProjects(info.value.Path)
            status.project = LOADED
        }).catch(err => {
            LogError(err)
            console.log(err)
        })
    }

    function open(path: string) {
        OpenProject(path).then(result => {
            info.value = result
            updateRecentProjects(info.value.Path)
            status.project = LOADED
        }).catch(err => {
            LogError(err)
            console.log(err)
        })
    }

    function updateRecentProjects(path: string) {

        if (config.value == null) return

        // If path in recent, remove it
        const index = config.value.RecentProjects.indexOf(path)
        if (index > -1) config.value.RecentProjects.splice(index, 1)

        // Prepend new path
        config.value.RecentProjects.unshift(path)

        // Limit to 5 items
        config.value.RecentProjects = config.value.RecentProjects.slice(0, 5)

        // Save config
        SaveConfig(config.value).catch(err => {
            LogError(err)
            console.log(err)
        })
    }

    //--------------------------------------------------------------------------
    // Model
    //--------------------------------------------------------------------------

    const model = reactive<main.Model>(new main.Model)

    function fetchModel() {
        return new Promise<main.Model>((resolve, reject) => {
            FetchModel().then(result => {
                Object.assign(model, result)
                resolve(model)
            }).catch(err => {
                LogError(err)
                console.log(err)
                reject(err)
            })
        })
    }

    function importModelDialog() {
        ImportModelDialog().then(result => {
            Object.assign(model, result)
        }).catch(err => {
            LogError(err)
            console.log(err)
        })
    }

    function updateModel() {
        UpdateModel(model).then(result => {
            Object.assign(model, result)
        }).catch(err => {
            LogError(err)
            console.log(err)
        })
    }

    //--------------------------------------------------------------------------
    // Analysis
    //--------------------------------------------------------------------------

    const analysis = ref<main.Analysis | null>(null)
    const currentCaseID = ref<number>(1)

    function fetchAnalysis() {
        return new Promise<main.Analysis>((resolve, reject) => {
            FetchAnalysis().then(result => {
                analysis.value = result
            }).catch(err => {
                LogError(err)
                console.log(err)
                reject(err)
            })
        })
    }

    function updateAnalysis() {
        if (analysis.value == null) return
        UpdateAnalysis(analysis.value).then(result => {
            analysis.value = result
        }).catch(err => {
            LogError(err)
            console.log(err)
        })
    }

    function addAnalysisCase() {
        return new Promise<main.Case>((resolve, reject) => {
            AddAnalysisCase().then(result => {
                analysis.value = result
                currentCaseID.value = result.Cases.length - 1
            }).catch(err => {
                LogError(err)
                console.log(err)
                reject(err)
            })
        })
    }

    function removeAnalysisCase(id: number) {
        RemoveAnalysisCase(id).then(result => {
            Object.assign(analysis, result)
        }).catch(err => {
            LogError(err)
            console.log(err)
        })
    }

    //--------------------------------------------------------------------------
    // Evaluate
    //--------------------------------------------------------------------------

    const evaluate = ref<main.Evaluate | null>(null)

    function fetchEvaluate() {
        FetchEvaluate().then(result => {
            evaluate.value = result
        }).catch(err => {
            LogError(err)
            console.log(err)
        })
    }

    function updateEvaluate() {
        if (evaluate.value == null) return
        UpdateEvaluate(evaluate.value).then(result => {
            evaluate.value = result
        }).catch(err => {
            LogError(err)
            console.log(err)
        })
    }

    function selectExec() {
        SelectExec().then(result => {
            evaluate.value = result
        }).catch(err => {
            LogError(err)
            console.log(err)
        })
    }

    // Setup listener for evaluation status updates
    EventsOn("evalStatus", (status: main.EvalStatus) => {
        Object.assign(evalStatus[status.ID - 1], status)
    })

    function startEvaluate(caseID: number) {
        results
        EvaluateCase(caseID).then(result => {
            Object.assign(evalStatus, result)
        }).catch(err => {
            LogError(err)
            console.log(err)
        })
    }

    function cancelEvaluate() {
        CancelEvaluate().catch(err => {
            LogError(err)
            console.log(err)
        })
    }

    //--------------------------------------------------------------------------
    // Results
    //--------------------------------------------------------------------------

    const results = ref<main.Results | null>(null)

    function fetchResults() {
        FetchResults().then(result => {
            Object.assign(results, result)
            console.log(result)
        }).catch(err => {
            LogError(err)
            console.log(err)
        })
    }

    function openCaseDirDialog() {
        status.results = LOADING
        OpenCaseDirDialog().then(result => {
            results.value = result
            status.diagram = NOT_LOADED
            status.results = LOADED
        }).catch(err => {
            LogError(err)
            console.log(err)
            status.results = NOT_LOADED
        })
    }

    //--------------------------------------------------------------------------
    // Diagram
    //--------------------------------------------------------------------------

    function generateDiagram(doCluster: boolean) {
        if (results.value == null) return
        status.diagram = LOADING
        GenerateDiagram(results.value.MinFreq, results.value.MaxFreq, doCluster).then(result => {
            Object.assign(diagram, result)
            status.diagram = LOADED
        }).catch(err => {
            LogError(err)
            console.log(err)
            status.diagram = NOT_LOADED
        })
    }

    //--------------------------------------------------------------------------
    // Visualization
    //--------------------------------------------------------------------------

    function getModeViz(opID: number, modeID: number, scale: number) {
        status.viz = LOADING
        return new Promise<viz.ModeData>((resolve, reject) => {
            GetModeViz(opID, modeID, scale).then(result => {
                console.log(result)
                const found = modeViz.find((md) => result.OPID == md.OPID && result.ModeID == md.ModeID)
                if (found !== undefined) {
                    Object.assign(found, result)
                } else {
                    modeViz.push(result);
                }
                status.viz = LOADED
                resolve(result)
            }).catch(err => {
                LogError(err)
                console.log(err)
                status.viz = NOT_LOADED
                reject(err)
            })
        })
    }

    function clearModeViz() {
        modeViz.splice(0);
    }

    //--------------------------------------------------------------------------
    // Other
    //--------------------------------------------------------------------------

    return {
        status: status,
        // Project
        config,
        info,
        saveDialog,
        openDialog,
        open,
        // Model
        model,
        fetchModel,
        importModelDialog,
        updateModel,
        // Analysis
        analysis,
        currentCaseID,
        fetchAnalysis,
        updateAnalysis,
        addAnalysisCase,
        removeAnalysisCase,
        // Evaluate
        evaluate,
        evalStatus,
        fetchEvaluate,
        updateEvaluate,
        startEvaluate,
        cancelEvaluate,
        selectExec,
        // Results
        results,
        fetchResults,
        openCaseDirDialog,
        // Diagram
        diagram,
        generateDiagram,
        // Visualization
        getModeViz,
        clearModeViz,
        modeViz,
    }
})
