<template>
  <div class="px-4 sm:px-6 lg:px-8">
    <div class="sm:flex sm:items-center">
      <div class="sm:flex-auto">
        <h1 class="text-xl font-semibold text-gray-900">Buckets</h1>
        <p class="mt-2 text-sm text-gray-700">A list of all buckets</p>
      </div>
    </div>
    <div class="mt-8 flex flex-col">
      <div class="-my-2 -mx-4 overflow-x-auto sm:-mx-6 lg:-mx-8">
        <div class="inline-block min-w-full py-2 align-middle md:px-6 lg:px-8">
          <div class="overflow-hidden shadow ring-1 ring-black ring-opacity-5 md:rounded-lg">
            <table class="min-w-full divide-y divide-gray-300">
              <thead class="bg-gray-50">
                <tr>
                  <th scope="col" class="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-900 sm:pl-6">Name</th>
                  <th scope="col" class="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Creation Date</th>
                  <th scope="col" class="relative py-3.5 pl-3 pr-4 sm:pr-6">
                    <span class="sr-only">View</span>
                  </th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-200 bg-white">
                <tr v-for="bucket in buckets" :key="bucket.Name">
                  <td class="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-6">{{ bucket.Name }}</td>
                  <td class="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{{ bucket.CreationDate }}</td>
                  <td class="relative whitespace-nowrap py-4 pl-3 pr-4 text-right text-sm font-medium sm:pr-6">
                    <router-link :to="{ name: 'show', params: { bucket: bucket.Name }}" class="text-indigo-600 hover:text-indigo-900">
                      View<span class="sr-only">, {{ bucket.Name }}</span>
                    </router-link>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, inject, onMounted } from 'vue'
import { RouterLink } from 'vue-router'

const axios = inject('axios')
const buckets = reactive([])
const apiBaseUrl = inject('apiBaseUrl')
const apiEndpoint = `${apiBaseUrl}/s3/buckets`;

onMounted(async () => {
  await axios.get(apiEndpoint)
    .then(response => buckets.push(...response.data.Buckets))
})

</script>
