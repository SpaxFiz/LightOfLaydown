<script  lang="ts">
import { getEmAccountData } from '@/api/emData'
import * as echarts from 'echarts'
import { defineComponent } from 'vue'

type EChartsOption = echarts.EChartsOption
interface EMAccountData {
  date: string,
  end_investor: number,
  market_avg_value: number,
  market_total_value: number,
  new_investor: number,
  new_investor_yoy: number,
  new_investor_mom: number,
  sh_index: number,
  [param: string]: any,
}

export default defineComponent({
  methods: {
    setup() {
      getEmAccountData({})
        .then(res => {
          let options = this.constructData(JSON.parse(res.data))
          this.render(options)
        })
        .catch(err => {
          console.log(err)
        })
    },
    constructData(data: EMAccountData[]): EChartsOption {
      let xData: any[] = []
      let yData: any[] = []
      let newInvester: any[] = []
      let cycleBasis: any[] = []
      data.forEach(item => {
        xData.push(item.date)
        yData.push(item.sh_index)
        newInvester.push(item.new_investor / 10)
        cycleBasis.push(item.new_investor_mom * 100)
      })
      let option: EChartsOption = {
        title: {
          text: '股东开户月度数据',
          left: 'center'
        },
        tooltip: {
          trigger: 'axis'
        },
        xAxis: {
          type: 'category',
          data: xData,
        },
        yAxis: [{
          id: 'sh_index',
          type: 'value'
        }, {
          id: 'new_investor',
          type: 'value',
          position: 'right',
          axisLabel: {
            formatter: '{value} 万'
          },
          axisLine: {
            show: true,
            lineStyle: {
              color: '#916400'
            }
          },
        }, {
          id: 'cycleBasis',
          min: -50,
          max: 100,
          offset: 80,
          position: 'right',
          axisLine: {
            show: true,
            lineStyle: {
              color: '#91CC75'
            }
          },
          axisLabel: {
            formatter: '{value} %'
          }
        }],
        series: [
          {
            name: '上证指数',
            data: yData,
            type: 'line',
            yAxisId: 'sh_index',
          }, {
            name: '新增投资者(万)',
            data: newInvester,
            type: 'bar',
            yAxisId: 'new_investor',

          }, {
            name: '环比',
            data: cycleBasis,
            type: 'line',
            yAxisId: 'cycleBasis',
            smooth: true
          }
        ]
      }
      return option
    },
    render(option: EChartsOption) {
      var chartDom = document.getElementById('chart')!
      var myChart = echarts.init(chartDom)
      myChart.setOption(option)
      window.onresize = function () {
        myChart.resize()
      }
    }
  },
  mounted() {
    this.setup()
  },
  data() {
    return {
      data: [],
      chartOption: {},
    }
  }
})

</script>

<template>
  <div id="chart"></div>
</template>

<style scoped>
#chart {
  height: 400px;
}
</style>
