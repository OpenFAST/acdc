<script lang="ts" setup>
import HelloWorld from './components/HelloWorld.vue'
import MainNav from './components/MainNav.vue'
</script>

<template>
  <MainNav />
  <div class="container-fluid">
    <router-view />
  </div>
</template>

<style></style>
