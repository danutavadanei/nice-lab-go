<template>
  <StudentIndexView v-if="userIsStudent"/>
  <ProfessorIndexView v-else-if="userIsProfessor" />
</template>

<script setup>
import { computed } from 'vue'
import { useStore } from 'vuex'
import StudentIndexView from './StudentIndexView.vue'
import ProfessorIndexView from './ProfessorIndexView.vue'

const store = useStore()
const user = store.getters.user
const userIsStudent = computed(() => user.type === 'student')
const userIsProfessor = computed(() => user.type === 'professor')
</script>
