<script lang="ts" setup>
import { reactive, onMounted } from 'vue'
import { Greet } from '../../wailsjs/go/main/App'
import * as THREE from 'three'

const data = reactive({
  name: "",
  resultText: "Please enter your name below 👇",
})

function greet() {
  Greet(data.name).then(result => {
    data.resultText = result
  })
}

onMounted(() => {

  const renderer = new THREE.WebGLRenderer();
  const container = document.getElementById('canvas');
  const w = container?.offsetWidth!;
  const h = container?.offsetHeight!;
  renderer.setSize(w, h);
  container?.appendChild(renderer.domElement);

  const scene = new THREE.Scene();
  const camera = new THREE.PerspectiveCamera(75, w / h, 0.1, 1000);

  const geometry = new THREE.BoxGeometry(1, 1, 1);
  const material = new THREE.MeshBasicMaterial({ color: 0x00ff00 });
  const cube = new THREE.Mesh(geometry, material);
  scene.add(cube);

  camera.position.z = 5;

  function animate() {
    requestAnimationFrame(animate);

    cube.rotation.x += 0.01;
    cube.rotation.y += 0.01;

    renderer.render(scene, camera);
  }

  animate();
})


</script>

<template>
  <main>
    <div id="result" class="result">{{ data.resultText }}</div>
    <div id="input" class="input-box">
      <input id="name" v-model="data.name" autocomplete="off" class="input" type="text" />
      <button class="btn" @click="greet">Greet</button>
    </div>
    <div id="canvas"></div>
  </main>
</template>

<style scoped>
#canvas {
  height: 300px
}
</style>
