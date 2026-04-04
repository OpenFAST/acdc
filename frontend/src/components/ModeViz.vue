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
    showNodeOrientation: boolean
}>()

watch(
    () => props.ModeData,
    (md) => { createFrames(md) },
    { deep: true }
)

watch(() => props.showNodePaths, (snp) => (nodePaths.visible = snp))
watch(() => props.showNodeOrientation, (sno) => {
    for (const ofr of orientationFrames) {
        ofr.visible = sno && lineFrames.indexOf(ofr) === frameNum
    }
})

let lineFrames = new Array<THREE.Group>;
let orientationFrames = new Array<THREE.Group>;
let nodePaths = new THREE.Group;
let frameNum = 0;
let frameCenter = new THREE.Vector3;
let frameSize = new THREE.Vector3;
let clock = new THREE.Clock();
let delta = 0
const FOV = 10

function createFrames(modeData: viz.ModeData) {
    if (modeData.Frames == null) return
    scene.clear()

    const geometry = new THREE.BufferGeometry();
    geometry.setAttribute('position', new THREE.Float32BufferAttribute([0, 0, 0], 3));
    const material = new THREE.PointsMaterial({ color: 0x888888 });
    const origin = new THREE.Points(geometry, material);
    origin.visible = false

    // Clear existing frames
    lineFrames = [] as THREE.Group[];
    orientationFrames = [] as THREE.Group[];

    const allFramesGroup = new THREE.Group()
    allFramesGroup.add(origin)

    const xMaterial = new THREE.LineBasicMaterial({ color: 0xff0000, linewidth: 2 });
    const yMaterial = new THREE.LineBasicMaterial({ color: 0x00ff00, linewidth: 2 });
    const zMaterial = new THREE.LineBasicMaterial({ color: 0x0000ff, linewidth: 2 });

    // Loop through frames
    for (const f of modeData.Frames) {

        // Lines
        const lineFrameGroup = new THREE.Group()
        for (const c of Object.values(f.Components)) {
            const curve = new THREE.CatmullRomCurve3(
                c.Line.map((p) => new THREE.Vector3(p.XYZ[0], p.XYZ[1], p.XYZ[2])))
            const points = curve.getPoints(50);
            const geometry = new THREE.BufferGeometry().setFromPoints(points);
            const material = new THREE.LineBasicMaterial({ color: 0xffffff, linewidth: 1 });
            const curveObject = new THREE.Line(geometry, material);
            lineFrameGroup.add(curveObject)
            allFramesGroup.add(curveObject.clone()) // Add clone of object to be used for view sizing
        }
        lineFrameGroup.visible = false // Initialize each group to not visible for animation
        lineFrames.push(lineFrameGroup)

        // Orientations
        const orientationFrameGroup = new THREE.Group()
        for (const c of Object.values(f.Components)) {
            const indices = new Uint16Array(c.Line.map((_, i) => i * 2).flatMap(i => [i, i + 1]));

            const pointsX = new Float32Array(c.Line.flatMap(p => [p.XYZ[0], p.XYZ[1], p.XYZ[2], p.XYZ[0] + p.OrientationX[0] * 4, p.XYZ[1] + p.OrientationX[1] * 4, p.XYZ[2] + p.OrientationX[2] * 4]));
            const geometryX = new THREE.BufferGeometry();
            geometryX.setAttribute('position', new THREE.BufferAttribute(pointsX, 3));
            geometryX.setIndex(new THREE.BufferAttribute(indices, 1));
            const lineX = new THREE.LineSegments(geometryX, xMaterial);
            orientationFrameGroup.add(lineX);

            const pointsY = new Float32Array(c.Line.flatMap(p => [p.XYZ[0], p.XYZ[1], p.XYZ[2], p.XYZ[0] + p.OrientationY[0] * 4, p.XYZ[1] + p.OrientationY[1] * 4, p.XYZ[2] + p.OrientationY[2] * 4]));
            const geometryY = new THREE.BufferGeometry();
            geometryY.setAttribute('position', new THREE.BufferAttribute(pointsY, 3));
            geometryY.setIndex(new THREE.BufferAttribute(indices, 1));
            const lineY = new THREE.LineSegments(geometryY, yMaterial);
            orientationFrameGroup.add(lineY);

            // const pointsZ = new Float32Array(c.Line.flatMap(p => [p.XYZ[0], p.XYZ[1], p.XYZ[2], p.XYZ[0] + p.OrientationZ[0] * 4, p.XYZ[1] + p.OrientationZ[1] * 4, p.XYZ[2] + p.OrientationZ[2] * 4]));
            // const geometryZ = new THREE.BufferGeometry();
            // geometryZ.setAttribute('position', new THREE.BufferAttribute(pointsZ, 3));
            // geometryZ.setIndex(new THREE.BufferAttribute(indices, 1));
            // const lineZ = new THREE.LineSegments(geometryZ, zMaterial);
            // orientationFrameGroup.add(lineZ);
        }
        orientationFrameGroup.visible = false // Initialize each group to not visible for animation
        orientationFrames.push(orientationFrameGroup)

        scene.add(lineFrameGroup)
        scene.add(orientationFrameGroup)
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
    if (delta > 1.5 / lineFrames.length) {
        delta = 0

        // Hide current frame
        lineFrames[frameNum].visible = false;
        orientationFrames[frameNum].visible = false;

        // Increment frame number
        frameNum++
        if (frameNum >= lineFrames.length) frameNum = 0

        // Show next frame
        lineFrames[frameNum].visible = true;
        if (props.showNodeOrientation) {
            orientationFrames[frameNum].visible = true;
        }

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