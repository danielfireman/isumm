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