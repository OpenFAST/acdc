<script lang="ts" setup>
import { reactive, ref, onMounted, computed, watch } from 'vue'
import { main, viz } from '../../wailsjs/go/models';
import { useProjectStore } from '../project';
import * as THREE from 'three'
import chroma from 'chroma-js';

const project = useProjectStore()

const props = defineProps<{
    ModeData: viz.ModeData
    showNodePaths: boolean
}>()

watch(
    () => props.ModeData,
    (md) => { createFrames(md) },
    { deep: true }
)

watch(() => props.showNodePaths, (snp) => (nodePaths.visible = snp))

const showNodePaths = ref(true)

let frames = new Array<THREE.Group>;
let nodePaths = new THREE.Group;
let frameNum = 0;
let capturing = false;
let frameCenter = new THREE.Vector3;
let frameSize = new THREE.Vector3;
let clock = new THREE.Clock();
let delta = 0
const FOV = 10

function createFrames(modeData: viz.ModeData) {
    if (modeData.Frames == null) return
    scene.clear()
    frames = [] as THREE.Group[];
    const geometry = new THREE.BufferGeometry();
    geometry.setAttribute('position', new THREE.Float32BufferAttribute([0, 0, 0], 3));
    const material = new THREE.PointsMaterial({ color: 0x888888 });
    const origin = new THREE.Points(geometry, material);
    origin.visible = false
    const allFramesGroup = new THREE.Group()
    allFramesGroup.add(origin)
    for (const f of modeData.Frames) {
        const frameGroup = new THREE.Group()
        for (const c of Object.values(f.Components)) {
            const curve = new THREE.CatmullRomCurve3(
                c.Line.map((p) => new THREE.Vector3(p.XYZ[0], p.XYZ[1], p.XYZ[2])))
            const points = curve.getPoints(50);
            const geometry = new THREE.BufferGeometry().setFromPoints(points);
            const material = new THREE.LineBasicMaterial({ color: 0xffffff, linewidth: 1 });
            const curveObject = new THREE.Line(geometry, material);
            frameGroup.add(curveObject)
            allFramesGroup.add(curveObject.clone()) // Add clone of object to be used for view sizing
        }
        frameGroup.visible = false // Initialize each group to not visible for animation
        frames.push(frameGroup)
        scene.add(frameGroup)
    }
    // Node paths
    const componentNames = Object.keys(modeData.Frames[0].Components)
    const curves = new Array<THREE.CatmullRomCurve3>
    const curveLengths = new Array<number>
    for (const compName of componentNames) {
        const numNodes = modeData.Frames[0].Components[compName].Line.length
        for (let i = 0; i < numNodes; i++) {
            let vectors = [] as THREE.Vector3[]
            for (const f of modeData.Frames) {
                const line = f.Components[compName].Line
                const p = line[i]
                vectors.push(new THREE.Vector3(p.XYZ[0], p.XYZ[1], p.XYZ[2]))
            }
            const curve = new THREE.CatmullRomCurve3(vectors)
            curve.closed = true
            curveLengths.push(curve.getLength())
            curves.push(curve)
        }
    }
    nodePaths.clear()
    const cs = chroma.scale(['008ae5', 'yellow']).domain([Math.min(...curveLengths), Math.max(...curveLengths)])
    for (let i = 0; i < curves.length; i++) {
        const points = curves[i].getPoints(50);
        const geometry = new THREE.BufferGeometry().setFromPoints(points);
        const material = new THREE.LineBasicMaterial({ color: cs(curveLengths[i]).hex(), linewidth: 1, transparent: true })
        material.opacity = 0.8
        const line = new THREE.Line(geometry, material)
        line.computeLineDistances()
        nodePaths.add(line)
    }
    scene.add(nodePaths)
    const bb = new THREE.Box3().setFromObject(allFramesGroup);
    frameCenter = bb.getCenter(new THREE.Vector3())
    frameSize = bb.getSize(new THREE.Vector3())
    frameNum = 0
    const axesHelper = new THREE.AxesHelper(frameSize.x / 2);
    scene.add(axesHelper)
}

let scene: THREE.Scene;
let renderer: THREE.WebGLRenderer;

const views = [
    {
        // Top View
        left: 0,
        bottom: 0.705,
        width: 0.4,
        height: 0.30,
        up: [1, 0, 0],
        updateCamera: function (camera: THREE.PerspectiveCamera) {
            // Calculate distance along Z axis to fit model in frame horizontally
            const fov = camera.fov * (Math.PI / 180);
            const fovh = 2 * Math.atan(Math.tan(fov / 2) * camera.aspect);
            let distance = 1.05 * (frameSize.y / 2 / Math.tan(fovh / 2) + frameSize.z)
            camera.position.fromArray([0, 0, distance]); // Looking along -Z (downward)
            camera.lookAt(frameCenter);
        },
        camera: new THREE.PerspectiveCamera,
    },
    {
        // Front View
        left: 0,
        bottom: 0,
        width: 0.4,
        height: 0.70,
        up: [0, 0, 1],
        updateCamera: function (camera: THREE.PerspectiveCamera) {
            // Calculate distance along -X axis to fit model in frame vertically
            // See https://wejn.org/2020/12/cracking-the-threejs-object-fitting-nut/ for equation
            let distance = 1.05 * (frameSize.z / 2 / Math.tan(camera.fov * Math.PI / 180 / 2) + frameSize.x / 2)
            camera.position.fromArray([-distance, 0, frameCenter.z]); // Looking along X (downwind)
            camera.lookAt(frameCenter);
        },
        camera: new THREE.PerspectiveCamera,
    },
    {
        // Side View
        left: 0.402,
        bottom: 0,
        width: 0.3,
        height: 0.70,
        up: [0, 0, 1],
        updateCamera: function (camera: THREE.PerspectiveCamera) {
            // Calculate distance along -Y axis to fit model in frame vertically
            let distance = 1.05 * (frameSize.z / 2 / Math.tan(camera.fov * Math.PI / 180 / 2) + frameSize.y / 2)
            camera.position.fromArray([0, -distance, frameCenter.z]); // Looking along -Y (side)
            camera.lookAt(frameCenter);
        },
        camera: new THREE.PerspectiveCamera,
    },

    {
        // Isometric View
        left: 0.704,
        bottom: 0,
        width: 0.3,
        height: 1.0,
        up: [0, 0, 1],
        updateCamera: function (camera: THREE.PerspectiveCamera) {
            // Calculate distance along Z axis to fit model in frame horizontally
            let distanceFront = 0.8 * (frameSize.z / 2 / Math.tan(camera.fov * Math.PI / 180 / 2) + frameSize.x / 2)
            let distanceSide = 0.8 * (frameSize.z / 2 / Math.tan(camera.fov * Math.PI / 180 / 2) + frameSize.y / 2)
            camera.position.fromArray([-distanceFront, -distanceSide, frameCenter.z + 3 * frameSize.z]); // Looking along -Z (downward)
            camera.lookAt(frameCenter);
        },
        camera: new THREE.PerspectiveCamera,
    }
];

function animate() {
    requestAnimationFrame(animate);
    delta += clock.getDelta()
    if (delta > 1.5 / frames.length) {
        delta = 0
        frames[frameNum].visible = false;
        frameNum++
        if (frameNum >= frames.length) frameNum = 0
        frames[frameNum].visible = true;
        render();
    }
}

function render() {

    const canvas = renderer.domElement;
    const canvasWidth = canvas.clientWidth;
    const canvasHeight = canvas.clientHeight;
    if (canvas.width !== canvasWidth || canvas.height !== canvasHeight) {
        renderer.setSize(canvasWidth, canvasHeight, false);
    }

    for (let ii = 0; ii < views.length; ++ii) {

        const view = views[ii];
        const camera = view.camera;

        view.updateCamera(camera);

        const left = Math.floor(canvasWidth * view.left);
        const bottom = Math.floor(canvasHeight * view.bottom);
        const width = Math.floor(canvasWidth * view.width);
        const height = Math.floor(canvasHeight * view.height);

        renderer.setViewport(left, bottom, width, height);
        renderer.setScissor(left, bottom, width, height);
        renderer.setScissorTest(true);

        camera.aspect = width / height;
        camera.updateProjectionMatrix();

        renderer.render(scene, camera);
    }
}

onMounted(() => {

    const canvas = <HTMLCanvasElement>document.getElementById('modeVizCanvas')!;

    for (let ii = 0; ii < views.length; ++ii) {
        const view = views[ii];
        const camera = new THREE.PerspectiveCamera(FOV, 2, 1, 10000);
        camera.up.fromArray(view.up);
        view.updateCamera(camera)
        view.camera = camera;
    }

    scene = new THREE.Scene();
    createFrames(props.ModeData)

    renderer = new THREE.WebGLRenderer({ antialias: true, canvas });
    renderer.setClearColor(0x3a3b3c);

    animate();
})

</script>

<template>
    <canvas id="modeVizCanvas" class="h-100 w-100"></canvas>
</template>

<style scoped></style>