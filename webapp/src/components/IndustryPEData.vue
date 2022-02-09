<template>
  <el-button-group>
    <el-button type="primary" size="small" @click="switchMode('close')">收盘价</el-button>
    <el-button type="primary" size="small" @click="switchMode('volume')">交易量</el-button>
  </el-button-group>
  <div v-loading="loading" id="industry"> </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'
import * as echarts from 'echarts'
import { getIndustryPEData } from '@/api/industryPEData'

  enum mode {
    close = 'close',
    volume = 'volume'
  }

  enum title {
    close = '收盘价',
    volume = '交易量'
  }

  interface industryDataStruct {
    dataset: echarts.DatasetComponentOption[]
    series: echarts.SeriesOption[]
  }

// interface backendData {
//   data: Record<string, string|number>[],
//   industryCodeMapping: Record<string, string>
// }

  interface iData {
    loading: boolean
    rawData: Record<string, Record<string, string | number>[]>
    industryCodeMappinng: Record<string, string>
    options: echarts.EChartsOption
    currentMode: mode
    currentTitle: string
    data?: industryDataStruct
    myChart?: echarts.ECharts
  }

export default defineComponent({
  methods: {
    getIndustryPEData() {
      getIndustryPEData().then((res) => {
        let d = res.data

        this.industryCodeMappinng = d.industryCode
        this.rawData = d.data
        this.data = this.transformRawData()
        this.options = this.constructData()
        this.render()

        this.loading = false
      })
    },
    transformRawData(): industryDataStruct {
      let dataset: echarts.DatasetComponentOption[] = [],
        series: echarts.SeriesOption[] = [],
        result: industryDataStruct

      for (const [key, val] of Object.entries(this.rawData)) {
        let value = val,
          currentSeries: echarts.SeriesOption

        value.sort(function (a, b) {
          return a.date < b.date ? -1 : 1
        })
        dataset.push({
          source: value
        })
        currentSeries = {
          name: this.industryCodeMappinng[key],
          datasetIndex: dataset.length - 1,
          type: 'line',
          smooth: true,
          symbol: 'none',
          dimensions: ['date'].concat([this.currentMode]),
          encode: {
            tooltip: [this.currentMode],
            x: 'date',
            y: [this.currentMode]
          },
          sampling: 'lttb'
        }

        series.push(currentSeries)
      }
      result = {
        dataset: dataset,
        series: series
      }
      return result
    },
    constructData(): echarts.EChartsOption {
      let option: echarts.EChartsOption

      option = {
        title: {
          text: this.currentTitle,
          left: '9%'
        },
        toolbox: {
          feature: {
            dataZoom: {
              yAxisIndex: 'none'
            },
            restore: {}
          }
        },
        dataZoom: [
          {
            type: 'inside',
            start: 0,
            end: 100
          },
          {
            start: 0,
            end: 100
          }
        ],
        tooltip: {
          trigger: 'axis',
          axisPointer: {
            snap: true
          }
        },
        legend: {},
        xAxis: {
          splitLine: { show: false },
          data: this.rawData.sh000922?.map((item) => {
            return item.date
          }),
          type: 'category',
          axisTick: {
            alignWithLabel: true
          }
        },
        yAxis: {},
        dataset: this.data?.dataset,
        series: this.data?.series
      }
      return option
    },
    render() {
      var chartDom = document.getElementById('industry')

      if (chartDom === null) {
        return
      }

      this.myChart = echarts.init(chartDom)
      if (this.options === null) {
        return
      }
      this.myChart.setOption(this.options)
      this.myChart.on('click', function (params) {
        console.log(params)
      })
      let c = this.myChart

      window.onresize = function () {
        c.resize()
      }


    },
    switchMode(m: string) {
      switch (m) {
      case 'close':
        this.currentMode = mode.close
        this.currentTitle = title.close
        break
      case 'volume':
        this.currentMode = mode.volume
        this.currentTitle = title.volume
        break
      default:
        break
      }
      this.data = this.transformRawData()
      this.options = this.constructData()

      this.myChart?.setOption(this.options, true, false)
    }
  },
  mounted() {
    this.getIndustryPEData()
  },
  data(): iData {
    return {
      loading: true,
      rawData: {},
      industryCodeMappinng: {},
      options: {},
      currentMode: mode.close,
      currentTitle: title.close
    }
  }
})
</script>

<style scoped>
  #industry {
    height: 400px;
  }
</style>
