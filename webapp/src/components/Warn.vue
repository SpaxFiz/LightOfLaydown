<template>
  <div>
    <el-alert :title="content" show-icon center type="warning" />
    <br>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

const copy = '本站所有内容仅供私人使用, 不构成任何投资建议, 使用本站内容造成的任何后果自行承担'

export default defineComponent({
  data() {
    return {
      'content': copy
    }
  }
})
</script>
