<template>
  <div class="px-4 sm:px-6 lg:px-8">
    <div class="sm:flex sm:items-center">
      <div class="sm:flex-auto">
        <h1 class="text-xl font-semibold text-gray-900">Connecting to lab #{{ lab }}</h1>
      </div>
    </div>
  </div>
</template>

<script setup>
import { inject, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const axios = inject('axios')
const files = ref([])
const lab = route.params.lab
const file = ref(null)
const apiBaseUrl = inject('apiBaseUrl')
const apiEndpoint = `${apiBaseUrl}/s3/buckets/${lab}`;

const fetch = async () => {
  return await axios.get(apiEndpoint)
    .then(response => (files.value = response.data.Contents))
}

const onFileChange = (e) => {
  file.value = e.target.files[0]
}

// const onSubmit = async (e) => {
//   const formData = new FormData()
//   formData.append('file', file.value, file.value.name)
//
//   await axios.post(apiEndpoint, formData, { headers: { "Content-Type": "multipart/form-data" } })
//     .then(fetch)
// }
//
// onMounted(fetch)
</script>
