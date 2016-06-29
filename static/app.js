// Loading language.
// Extracted from github.com/adamwdraper/Numeral-js/blob/master/languages/pt-br.js
numeral.language('pt-br', {
  delimiters: {
    thousands: '.',
    decimal: ','
  },
  abbreviations: {
    thousand: 'mil',
    million: 'milhões',
    billion: 'b',
    trillion: 't'
  },
  ordinal: function (number) {
    return 'º';
  },
  currency: {
    symbol: 'R$'
  },
})

// Setting language and default currency formatter.
numeral.language('pt-br')

function formatCurrency(n) {
  return numeral(n.toFixed(2)).format('$ 0.0,00').slice(0, -2)
}

function formatPercent(n) {
  return numeral(n.toFixed(2)).format('0.0,00').slice(0, -2) + ' %'
}

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

// Default value for date-picker.
$(document).ready(function() {
	var now = new Date();
	// The two-digit day and month are needed due to time conversion formatting templates.
	var today = now.getFullYear() + '-' + ('0' + (now.getMonth()+1)).slice(-2) + '-' + ('0' + now.getDate()).slice(-2);
	$('#datePicker').val(today);
});

$(document).ready(function(){
    $('[data-toggle="tooltip"]').tooltip();
});

$(document).ready(function(){
  $('#op').submit(function() {
    var txt = $('#opvalue');
    console.log(numeral().unformat(txt.val()))
    txt.val(numeral().unformat(txt.val()));
  });
});

$(function() {
     $("#datePicker").datepicker({ dateFormat: 'yy-mm-dd'});
});
