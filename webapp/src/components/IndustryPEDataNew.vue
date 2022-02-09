<script lang="ts">
import { getIndustryPEData } from '@/api/industryPEData'
import * as echarts from 'echarts'
import { defineComponent } from 'vue'

const yAxis: echarts.YAXisComponentOption[] = [{
  'id': 'PE',
  'name': 'PE',
  'min': 0,
  'max': 130,
  'position': 'right',
  'axisLine': {
    'show': 'auto',
    'lineStyle': {
      'color': '#188773'
    }
  }
}]

type EChartsOption = echarts.EChartsOption

  interface seriesData {
    data: Array<Array<string|number>>,
    name: string
  }

  interface iData {
    chartOptions: EChartsOption,
    rawData: Array<seriesData>,
    myChart?: echarts.ECharts,
  }


export default defineComponent({
  methods: {
    setup() {
      getIndustryPEData()
        .then(res => {
          this.rawData = JSON.parse(res.data)
          this.chartOptions = this.constructData()
          this.render()
        })
        .catch(err => {
          console.log(err)
        })
    },
    constructData():EChartsOption {
      let options: EChartsOption

      options = {
        'title': {
          'text': '行业估值趋势',
          'left': '9%'
        },
        'toolbox': {
          'feature': {
            'dataZoom': {
              'yAxisIndex': 'none'
            },
            'restore': {}
          }
        },
        'dataZoom': [
          {
            'type': 'inside',
            'start': 30,
            'end': 100
          },
          {
            'start': 30,
            'end': 100
          }
        ],
        'tooltip': {
          'trigger': 'axis',
          'axisPointer': {
            'snap': true
          }
        },
        'legend': {
          'align': 'left',
          'right': '10%'
        },
        'xAxis': {
          'splitLine': { 'show': false },
          'type': 'time'
          // 'data': ['时间', 'PE'],
        },
        'yAxis': yAxis,
        'series': this.rawData
      }

      return options
    },
    render() {
      var chartDom = document.getElementById('chartPEIndustry')

      if (chartDom === null) {
        return
      }

      this.myChart = echarts.init(chartDom)
      this.myChart.setOption(this.chartOptions)
      let c = this.myChart

      window.onresize = function () {
        c.resize()
      }
    }
  },
  mounted() {
    this.setup()
  },
  data(): iData {
    return {
      'chartOptions': {},
      'rawData': []
    }
  }
})
</script>

<template>
  <div id="chartPEIndustry"></div>
</template>

<style scoped>
#chartPEIndustry {
  height: 400px;
}
</style>
