import { defineStore } from 'pinia'
import { ref, reactive } from 'vue'
import { GenerateDiagram, LoadConfig, SaveConfig, UpdateDiagram } from "../wailsjs/go/main/App"
import { OpenProjectDialog, SaveProjectDialog, OpenProject } from '../wailsjs/go/main/App'
import { FetchModel, UpdateModel, ImportModelDialog } from "../wailsjs/go/main/App"
import { FetchAnalysis, UpdateAnalysis, AddAnalysisCase, RemoveAnalysisCase } from "../wailsjs/go/main/App"
import { FetchEvaluate, UpdateEvaluate, SelectExec, EvaluateCase, CancelEvaluate } from "../wailsjs/go/main/App"
import { FetchResults, OpenCaseDirDialog } from "../wailsjs/go/main/App"
import { GetModeViz } from "../wailsjs/go/main/App"
import { main, diagram as diag, viz } from "../wailsjs/go/models"
import { EventsOn } from "../wailsjs/runtime/runtime"
import { LogError } from '../wailsjs/runtime/runtime'

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

    const status = reactive<Loading>(new Loading)
    const config = ref<main.Config | null>(null)
    const info = ref<main.Info | null>(null)
    const model = ref<main.Model | null>(null)
    const analysis = ref<main.Analysis | null>(null)
    const evaluate = ref<main.Evaluate | null>(null)
    const results = ref<main.Results | null>(null)
    const diagram = ref<diag.Diagram | null>(null)

    const currentCaseID = ref<number>(1)
    const evalStatus = reactive<Array<main.EvalStatus>>(new Array)
    const modeViz = reactive<Array<viz.ModeData>>(new Array)

    function $reset() {
        info.value = null
        model.value = null
        analysis.value = null
        evaluate.value = null
        results.value = null
        diagram.value = null
        currentCaseID.value = 1
        clearEvalStatus()
        clearModeViz()
    }

    //--------------------------------------------------------------------------
    // Project
    //--------------------------------------------------------------------------

    // Load config when store is initialized
    LoadConfig().then(result => {
        config.value = result
    }).catch(err => {
        LogError(err)
        console.log(err)
    })

    function openDialog() {
        OpenProjectDialog().then(result => {
            $reset()
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
            $reset()
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

    function fetchModel() {
        if (model.value != null) return
        FetchModel().then(result => {
            model.value = result
        }).catch(err => {
            LogError(err)
            console.log(err)
        })
    }

    function importModelDialog() {
        ImportModelDialog().then(result => {
            model.value = result
        }).catch(err => {
            LogError(err)
            console.log(err)
        })
    }

    function updateModel() {
        if (model.value == null) return
        UpdateModel(model.value).then(result => {
            model.value = result
        }).catch(err => {
            LogError(err)
            console.log(err)
        })
    }

    //--------------------------------------------------------------------------
    // Analysis
    //--------------------------------------------------------------------------

    function fetchAnalysis() {
        if (analysis.value != null) return
        FetchAnalysis().then(result => {
            analysis.value = result
        }).catch(err => {
            LogError(err)
            console.log(err)
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
        AddAnalysisCase().then(result => {
            analysis.value = result
            currentCaseID.value = result.Cases.length - 1
        }).catch(err => {
            LogError(err)
            console.log(err)
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

    function fetchEvaluate() {
        if (evaluate.value != null) return
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

    // Setup listener for evaluation status updates
    EventsOn("evalStatus", (status: main.EvalStatus) => {
        Object.assign(evalStatus[status.ID], status)
    })

    function clearEvalStatus() {
        evalStatus.splice(0)
    }

    //--------------------------------------------------------------------------
    // Results
    //--------------------------------------------------------------------------

    function fetchResults() {
        if (results.value != null) return
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

    function generateDiagram(minFreq: number, maxFreq: number, doCluster: boolean, filterStructural: boolean) {
        status.diagram = LOADING
        GenerateDiagram(minFreq, maxFreq, doCluster, filterStructural).then(result => {
            diagram.value = result
            status.diagram = LOADED
        }).catch(err => {
            LogError(err)
            console.log(err)
            status.diagram = NOT_LOADED
        })
    }

    function updateDiagram() {
        if (diagram.value == null) return
        UpdateDiagram(diagram.value).catch(err => {
            LogError(err)
            console.log(err)
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
        $reset,
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
        clearEvalStatus,
        // Results
        results,
        fetchResults,
        openCaseDirDialog,
        // Diagram
        diagram,
        generateDiagram,
        updateDiagram,
        // Visualization
        getModeViz,
        clearModeViz,
        modeViz,
    }
})
