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
      if (genAxis.length > 5) {
          var midIndex = ~~(genAxis.length/2)
          var half = ~~(midIndex/2)
          return [genAxis[0], genAxis[midIndex-half], genAxis[midIndex], genAxis[midIndex+half], genAxis[genAxis.length-1]]
      }
      return genAxis
    },
  }

  aSummaryOptions = JSON.parse(JSON.stringify(options))
  aSummaryOptions.yaxis =  {
    tickFormatter: formatCurrency,
    axisLabel: "Valor Total (R$)",
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

/**
  Typeahead function - This function enables the typeahead input
  and writes the investment id in #hiddenInputElement to form submit
**/
$(document).ready(function(){
  $('#investment-field .typeahead').typeahead({
    hint: true,
    highlight: true,
    minLength: 1
  },
  {
    name: 'investments',
    source: function (query, process) {
      listInv = [];
      map = {};
      $.each(investments, function (i, inv) {
          map[inv.name] = inv.id;
          listInv.push(inv.name);
      });
      process(listInv);
    }
  }).on('typeahead:select', function(ev, suggestion) {
    $('#hiddenInputElement').val(investments[listInv.indexOf(suggestion)].id);
  });
});