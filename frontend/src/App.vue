<script lang="ts" setup>

import MainNav from './components/MainNav.vue'
import { storeToRefs } from 'pinia'
import { ref, reactive, onMounted, watch } from 'vue'
import { useProjectStore } from './project';
import { Modal } from 'bootstrap';

const project = useProjectStore()

const { errMsg } = storeToRefs(project)
const modalInstance = ref<Modal | null>();

onMounted(() => {
  modalInstance.value = new Modal(document.getElementById("errorModal") as Element);
})

watch(errMsg, () => {
  if (errMsg != null) modalInstance.value?.show()
})

function closeModal() {
  project.errMsg = null;
  modalInstance.value?.hide()
}

</script>

<template>
  <MainNav />
  <div class="container-fluid">
    <router-view />
  </div>

  <div class="modal fade" id="errorModal" data-bs-backdrop="static" data-bs-keyboard="false" tabindex="-1"
    aria-labelledby="staticBackdropLabel" aria-hidden="true">
    <div class="modal-dialog modal-lg">
      <div class="modal-content">
        <div class="modal-header text-bg-danger">
          <h1 class="modal-title fs-5" id="staticBackdropLabel">Error</h1>
          <button type="button" class="btn-close" @click="closeModal()" aria-label="Close"></button>
        </div>
        <div class="modal-body">
          {{ project.errMsg }}
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-primary" @click="closeModal()">Close</button>
        </div>
      </div>
    </div>
  </div>

</template>

<style></style>
