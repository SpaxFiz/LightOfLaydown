<template>
  <div id="calculator">
    <div class="calculator" @click="showDarwer = true">
      <div class="display">
        <svg
          class="icon"
          width="25"
          height="25"
          viewBox="0 0 1024 1024"
          xmlns="http://www.w3.org/2000/svg"
          data-v-365b8594
        >
          <path
            fill="currentColor"
            d="M764.416 254.72a351.68 351.68 0 0186.336 149.184H960v192.064H850.752a351.68 351.68 0 01-86.336 149.312l54.72 94.72-166.272 96-54.592-94.72a352.64 352.64 0 01-172.48 0L371.136 936l-166.272-96 54.72-94.72a351.68 351.68 0 01-86.336-149.312H64v-192h109.248a351.68 351.68 0 0186.336-149.312L204.8 160l166.208-96h.192l54.656 94.592a352.64 352.64 0 01172.48 0L652.8 64h.128L819.2 160l-54.72 94.72zM704 499.968a192 192 0 10-384 0 192 192 0 00384 0z"
          />
        </svg>
        <span>头铁计算器</span>
      </div>
    </div>

    <el-dialog
      ref="drawer"
      v-model="showDarwer"
      title="头铁计算器"
      :before-close="handleClose"
      width="70%"
    >
      <div class="dialog">
        <el-form :model="form">
          <el-row :gutter="24">
            <el-col class="col" :span="8">
              <el-form-item label="代码" label-width="auto">
                <el-input class="code" size="small" v-model="code" @blur="fetchFundInfo"></el-input>
                <span class="name">{{ form.name }}</span>
              </el-form-item>
            </el-col>
            <el-col class="col" :span="8">
              <el-form-item label="当前净值" label-width="auto">
                <el-input-number
                  :controls="false"
                  size="small"
                  v-model="form.price"
                ></el-input-number>
              </el-form-item>
            </el-col>
            <el-col class="col" :span="8">
              <el-form-item label="当前盈利" label-width="auto">
                <span :style="{ color: (form.currentProfit || 0) > 0 ? 'red' : 'green' }">{{
                  form.currentProfit?.toFixed(2)
                }}</span>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="24">
            <el-col class="col" :span="8">
              <el-form-item label="预期下跌幅度" label-width="auto">
                <el-input-number
                  @blur="calculate"
                  :controls="true"
                  :step="0.05"
                  size="small"
                  v-model="form.expectDropOff"
                  id="expectedDropOff"
                  ><template #append>%</template></el-input-number
                >
              </el-form-item>
            </el-col>
            <el-col class="col" :span="8">
              <el-form-item label="底仓" label-width="auto">
                <el-input-number
                  @blur="
                    calculate();
                    currentProfit()
                  "
                  :controls="false"
                  size="small"
                  v-model="form.baseAmount"
                ></el-input-number>
              </el-form-item>
            </el-col>
            <el-col class="col" :span="8">
              <el-form-item label="成本价" label-width="auto">
                <el-input-number
                  @blur="
                    calculate();
                    currentProfit()
                  "
                  size="small"
                  :precision="4"
                  v-model="form.costPrice"
                  :controls="false"
                ></el-input-number>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="24">
            <el-col class="col" :span="8">
              <el-form-item label="每下跌" label-width="auto">
                <el-input-number
                  @blur="calculate"
                  :controls="false"
                  size="small"
                  v-model="form.addOnPercentage"
                  id="addOnPercentage"
                  ><template #append>%</template></el-input-number
                >
                <span>&nbsp; 加仓</span>
              </el-form-item>
            </el-col>
            <el-col class="col" :span="8">
              <el-form-item label="每次加仓金额" label-width="auto">
                <el-input-number
                  @blur="calculate"
                  :controls="false"
                  size="small"
                  v-model="form.addOnBase"
                ></el-input-number>
              </el-form-item>
            </el-col>
            <el-col class="col" :span="8">
              <el-form-item label="加仓系数" label-width="auto">
                <el-input-number
                  @blur="calculate"
                  size="small"
                  :precision="2"
                  :step="0.05"
                  v-model="form.addOnCoefficient"
                ></el-input-number>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="24">
            <el-col class="col" :span="8">
              <el-form-item label="最终成本价" label-width="auto">
                <span>{{ finalCostPrice?.toFixed(4) }}</span>
              </el-form-item>
            </el-col>


            <el-col class="col" :span="8">
              <el-form-item label="价格回到当前净值盈利" label-width="auto">
                <el-input
                  size="small"
                  v-model="form.priceBackToNowProfit"
                  disabled="disabled"
                ></el-input>
              </el-form-item>
            </el-col>
            <el-col class="col" :span="8">
              <el-form-item label="谷底上涨" label-width="auto">
                <span>{{ (form.valleyPercentage || 0) * 100 }} %</span>
                <span>回本</span>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row :gutter="24">
             <el-col class="col" :span="18">
              <el-form-item label="每次加仓前总亏损" label-width="auto">
                <el-table size="mini" :data="form.lossAmount || []" :border="true">
                  <el-table-column label="加仓价格" prop="price"></el-table-column>
                  <el-table-column label="该价格时亏损" prop="loss"></el-table-column>
                  <el-table-column label="加仓前亏损比例" prop="totalRatio" :formatter="(row:any ,col: any, val:number)=>{return val.toFixed(4)}"></el-table-column>
                  <el-table-column label="当前加仓金额" prop="curAddOn"></el-table-column>
                  <el-table-column label="累计加仓金额" prop="cumAddOn"></el-table-column>
                  <el-table-column label="加仓后股数" prop="cumAmount"></el-table-column>
                </el-table>
              </el-form-item>
            </el-col>
          </el-row>
          <el-row>
            <el-col>
              <el-button size="small" type="primary" @click="calculate">计算</el-button>
            </el-col>
          </el-row>
        </el-form>
      </div>
    </el-dialog>
  </div>
</template>

<script lang="ts">
import { getFundInfo } from '@/api/lossCalculator'
import { defineComponent } from 'vue'
import { ElMessage } from 'element-plus'

  /**
   * 代码
   * (名称)
   * (当前净值)
   * 预计下跌幅度
   * 底仓
   * 成本价
   * 从现在起每跌百分之几加仓
   * 从现在起每次加仓基数
   * 从现在起每次加仓系数
   * 每次加仓前总亏损
   * 价格回到当前净值盈利
   * 价格回到谷底与当前净值中位数盈利
   */
  interface calculateForm {
    code?: number
    price?: number
    name?: string
    expectDropOff: number
    baseAmount: number
    costPrice: number
    addOnPercentage: number
    addOnBase: number
    addOnCoefficient: number
    lossAmount?: Array<loss>
    priceBackToNowProfit?: number
    valleyPercentage?: number
    currentProfit?: number
  }

  interface loss {
    loss: number
    price: number
    curAddOn: number
    cumAddOn: number
    cumAmount: number
    totalRatio: number
  }

  interface iData {
    showDarwer: boolean
    formLabelWidth: string
    form: calculateForm
    constructedForm?: calculateForm
    code?: string | number
    finalCostPrice?: number
  }

export default defineComponent({
  data(): iData {
    return {
      showDarwer: false,
      formLabelWidth: '80px',
      form: {
        expectDropOff: 0.5,
        baseAmount: 0,
        costPrice: 0,
        addOnPercentage: 0.1,
        addOnBase: 10000,
        addOnCoefficient: 1
      },
      finalCostPrice: 0,
      code: ''
    }
  },
  methods: {
    handleClose(done: () => void) {
      done()
    },
    cancelForm() {
      this.showDarwer = false
    },
    fetchFundInfo() {
      if (typeof this.code === 'undefined') {
        return
      }
      getFundInfo(this.code)
        .then((res) => {
          this.form.price = res.data.price
          this.form.name = res.data.name
        })
        .catch((err) => {
          console.log(err)
          // eslint-disable-next-line new-cap
          ElMessage({
            message: err.message,
            type: 'warning'
          })
        })
    },

    currentProfit() {
      this.form.currentProfit =
          (this.form.price || 0 - this.form.costPrice) * (this.form.baseAmount || 0)
    },
    calculate() {
      // if (!(typeof this.form.expectDropOff !== 'undefined' &&typeof this.form.baseAmount !== 'undefined' &&typeof this.form.costPrice !== 'undefined' &&typeof this.form.addOnPercentage !== 'undefined' &&typeof this.form.addOnBase !== 'undefined' &&typeof this.form.addOnCoefficient !== 'undefined')) {
      //   return
      // }
      if (typeof this.form.price === 'undefined') {
        alert('!!!')
        return
      }
      let times = Math.floor(this.form.expectDropOff / this.form.addOnPercentage),
        curPrice = this.form.price,
        costPrice = this.form.costPrice,
        total = this.form.baseAmount * this.form.costPrice,
        amount = this.form.baseAmount,
        cumAddOn = 0,
        loss = 0,
        curAddOn = this.form.addOnBase,
        eachAddOnLoss: Array<loss> = []

      for (let i = 1; i <= times; i++) {
        curPrice = this.form.price * (1 - this.form.addOnPercentage * i)
        cumAddOn += curAddOn
        loss = Math.round((curPrice - costPrice) * amount)

        eachAddOnLoss.push({
          loss: loss,
          price: parseFloat((this.form.price * (1 - this.form.addOnPercentage * i)).toFixed(4)),
          curAddOn: curAddOn,
          cumAddOn: cumAddOn,
          cumAmount: 0,
          totalRatio: -loss/ total
        })

        curAddOn *= this.form.addOnCoefficient
        amount += Math.floor(this.form.addOnBase * this.form.addOnCoefficient / curPrice)
        total += this.form.addOnBase * this.form.addOnCoefficient
        costPrice = total / amount
        eachAddOnLoss[eachAddOnLoss.length-1].cumAmount = amount
      }
      this.form.lossAmount = eachAddOnLoss
      this.form.priceBackToNowProfit =
          Math.round((this.form.costPrice - costPrice) * amount * 100) / 100
      this.form.currentProfit = (this.form.price - this.form.costPrice) * this.form.baseAmount
      this.form.valleyPercentage = Math.round((costPrice - curPrice) / curPrice * 100) / 100
      this.finalCostPrice = costPrice
    }
  }
})
</script>

<style lang="less" scoped>
  .calculator {
    text-align: right;
    margin-right: 10px;
    display: flex;

    .display {
      margin-left: auto;
      display: flex;
      cursor: pointer;
    }
  }
  .calculator .dialog .col {
    text-align: right;
  }
  .code {
    width: 80px;
  }
  .name {
    margin-left: 10px;
  }
</style>
