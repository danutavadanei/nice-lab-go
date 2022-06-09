<template>
  <div class="px-4 sm:px-6 lg:px-8">
    <div class="sm:flex sm:items-center">
      <div class="sm:flex-auto">
        <h1 class="text-xl font-semibold text-gray-900">Files from {{ bucket }}</h1>
        <p class="mt-2 text-sm text-gray-700">A list of all the files from this bucket</p>
      </div>
    </div>
    <div class="mt-8 grid gap-6 lg:grid-flow-col-dense lg:grid-cols-3">
      <div class="space-y-6 lg:col-start-1 lg:col-span-2 flex flex-col">
        <div class="-my-2 -mx-4 overflow-x-auto sm:-mx-6 lg:-mx-8">
          <div class="inline-block min-w-full py-2 align-middle md:px-6 lg:px-8">
            <div class="overflow-hidden shadow ring-1 ring-black ring-opacity-5 md:rounded-lg">
              <table class="min-w-full divide-y divide-gray-300">
                <thead class="bg-gray-50">
                  <tr>
                    <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-900 sm:pl-6">Name</th>
                    <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Last Modified</th>
                    <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Size</th>
                    <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Storage Class</th>
                    <th scope="col" class="relative py-3.5 pl-3 pr-4 sm:pr-6">
                      <span class="sr-only">Download</span>
                    </th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-gray-200 bg-white">
                  <tr v-for="file in files" :key="file.Key">
                    <td class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-6">{{ file.Key }}</td>
                    <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{{ file.LastModified }}</td>
                    <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{{ file.Size }} Bytes</td>
                    <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{{ file.StorageClass }}</td>
                    <td class="relative whitespace-nowrap py-4 pl-3 pr-4 text-right text-sm font-medium sm:pr-6">
                      <a :href="`http://localhost:8080/s3/buckets/${bucket}/${file.Key}`" target="_blank" class="text-indigo-600 hover:text-indigo-900"
                        >Download<span class="sr-only">, {{ file.Key }}</span></a
                      >
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
      <div class="bg-white shadow sm:rounded-lg lg:col-span-1">
        <div class="px-4 py-5 sm:p-6">
          <h3 class="text-lg leading-6 font-medium text-gray-900">Upload a new file</h3>
          <form class="mt-5 sm:flex sm:items-center" @submit.prevent="onSubmit">
            <div class="flex items-end">
              <div class="w-96">
                <label for="file" class="form-label inline-block mb-2 text-gray-700">Default file input example</label>
                <input class="form-control
                  block
                  w-full
                  px-3
                  py-1.5
                  text-base
                  font-normal
                  text-gray-700
                  bg-white bg-clip-padding
                  border border-solid border-gray-300
                  rounded
                  transition
                  ease-in-out
                  m-0
                  focus:text-gray-700 focus:bg-white focus:border-blue-600 focus:outline-none" type="file" id="file" @change="onFileChange" ref="file">
              </div>
              <button type="submit" class="mt-3 w-full inline-flex items-center justify-center px-4 py-3 border border-transparent shadow-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm">Upload</button>
            </div>
          </form>
        </div>
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
const bucket = route.params.bucket
const file = ref(null)
const apiBaseUrl = inject('apiBaseUrl')
const apiEndpoint = `${apiBaseUrl}/s3/buckets/${bucket}`;

const fetch = async () => {
  return await axios.get(apiEndpoint)
    .then(response => (files.value = response.data.Contents))
}

const onFileChange = (e) => {
  file.value = e.target.files[0]
}

const onSubmit = async (e) => {
  const formData = new FormData()
  formData.append('file', file.value, file.value.name)

  await axios.post(apiEndpoint, formData, { headers: { "Content-Type": "multipart/form-data" } })
    .then(fetch)
}

onMounted(fetch)
</script>
