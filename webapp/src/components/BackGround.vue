<script lang="ts">
import { checkCipher } from '@/api/cipher'
import { defineComponent } from 'vue'
import { ElMessage as elMessage } from 'element-plus'

export default defineComponent({
  emits: ['VALIDCIPHER'],
  methods: {
    check() {
      checkCipher({cipher: this.cipher})
        .then(res => {
          console.log(res)
          this.$emit('VALIDCIPHER')
        })
        .catch(err => {
          console.log(err)
          elMessage('permission deny')
        })
    }
  },
  data() {
    return {
      cipher: '',
      imgPath: '/api/background'
    }
  }
})
</script>

<template>
  <div class="bg">
    <img :src="imgPath" alt=""/>
  </div>

  <div class="cipher">
    <el-input size="small" type="password" v-model="cipher" placeholder="your cipher" @keypress.enter="check">
      <template #append>
        <el-button size="small" type="primary" @click="check" >check</el-button>
      </template>
    </el-input>
  </div>
</template>

<style scoped>
.cipher {
      top: 50%;
    left: 41%;
    position: absolute;
}
.bg > img {
  object-fit:fill;
  width:100%;
  height: 100%;
  overflow:hidden;
}
</style>
