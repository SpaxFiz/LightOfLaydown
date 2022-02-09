<script lang="ts">
// This starter template is using Vue 3 <script setup> SFCs
// Check out https://v3.vuejs.org/api/sfc-script-setup.html#sfc-script-setup
import EmData from '@/components/EmData.vue'
import PEData from '@/components/PEData.vue'
import Warn from '@/components/Warn.vue'
import IndustryPEDataNew from './components/IndustryPEDataNew.vue'
import LossCaculator from './components/LossCaculator.vue'
import BackGround from './components/BackGround.vue'

import { defineComponent } from 'vue'

export default defineComponent({
  components: {EmData, PEData, Warn, IndustryPEDataNew, LossCaculator, BackGround},
  methods: {
    VALIDCIPHER() {
      this.showBG = false
    }
  },
  data() {
    return {
      showDarwer: false,
      showBG: true
    }
  }
})
</script>

<template>
  <div v-if="showBG">
    <BackGround @VALIDCIPHER="VALIDCIPHER"></BackGround>
  </div>
  <div v-else>
    <Warn></Warn>
    <LossCaculator></LossCaculator>
    <PEData></PEData>
    <EmData></EmData>
    <IndustryPEDataNew></IndustryPEDataNew>
  </div>
</template>
