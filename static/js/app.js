// Plotting charts
$(document).ready(function() {
  var options = {
    series: {
      shadowSize: 5,
      points: { show: true },
      lines: { show: true },
    },
    grid: {
      hoverable: true,
      clickable: true,
    },
    legend: {
      position: "nw",
    },
    tooltip : true,
  };

  xaxis = {
    mode: "time",
    timeformat: "%m/%Y",
    ticks: function(axis){
      genAxis = axis.tickGenerator(axis);
      if (genAxis.length > 2) {
          var midIndex = ~~(genAxis.length/2)
          return [genAxis[0], genAxis[midIndex], genAxis[genAxis.length-1]]
      }
      return genAxis
    },
  }

  aSummaryOptions = JSON.parse(JSON.stringify(options))
  aSummaryOptions.yaxis =  {
    tickFormatter: formatCurrency,
    axisLabel: "Valor Total (R$)",
  }
  aSummaryOptions.tooltipOpts = {
    content: function(label, xval, yval, flotItem){
        return $.plot.formatDate(new Date(xval), "%m/%Y") + ": " + formatCurrencyFull(yval)
    }
  }
  aSummaryOptions.xaxis = xaxis
  $.plot("#summChart", amountSummaryData, aSummaryOptions);

  irateChartOptions = JSON.parse(JSON.stringify(options))
  irateChartOptions.yaxis =  {
    tickFormatter: formatPercent,
    axisLabel: "Rentabilidade Mensal (%)",
  }
  irateChartOptions.xaxis = xaxis
  $.plot("#irateChart", interestRateData, irateChartOptions);
});
