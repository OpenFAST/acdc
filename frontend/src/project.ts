import { defineStore } from 'pinia'
import { computed, ref, reactive } from 'vue'
import { OpenProject, OpenProjectDialog, SaveProject, SaveProjectDialog } from '../wailsjs/go/main/App'
import { LoadConfig, SaveConfig } from "../wailsjs/go/main/App"
import { SelectExec, ImportModelDialog, UpdateModel } from "../wailsjs/go/main/App"
import { main } from "../wailsjs/go/models"
import { File, Field } from "./types"

export const useProjectStore = defineStore('project', () => {

    const saving = ref(false)
    const loaded = ref(false)
    const info = reactive<main.Info>(new main.Info)
    const exec = reactive<main.Exec>(new main.Exec)
    const model = reactive<main.Model>(new main.Model)
    const analyze = reactive<main.Analyze>(new main.Analyze)
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
            Object.assign(analyze, result.Analyze)
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
            Object.assign(analyze, result.Analyze)
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
        config.RecentProjects.slice(0, 5)      // Limit to 5 items
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

    function updateModel() {
        UpdateModel(model).then(result => {
            Object.assign(info, result.Info)
        }).catch(err => {
            console.log(err)
        })
    }

    function instanceOfField(obj: any): obj is Field {
        return typeof obj == 'object' && 'Name' in obj && 'Type' in obj;
    }

    const modelFileOptions = computed<File[]>(() => {
        const options: File[] = []
        for (const files of Object.values(model) as File[][]) {
            for (const file of files) {
                options.push({
                    Name: file.Name,
                    Type: file.Type,
                    // Text: file.Text,
                    Fields: Object.values(file).filter(instanceOfField),
                } as File)
            }
        }
        return options
    })

    return {
        loaded, saving, config, info, exec, model, analyze,
        modelFileOptions,
        $reset, open, saveDialog, save, openDialog, selectExec,
        importModel, updateModel
    }
})
