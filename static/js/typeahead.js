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
    source: substringMatcher(investments)
  }).on('typeahead:select', function(ev, suggestion) {
      listInv = [];
      $.each(investments, function (i, investment) {
        listInv.push(investment.invName);
      });
    $('#hiddenInputElement').val(investments[listInv.indexOf(suggestion)].invCode);
  });
});

/**
  Matcher function - This function search to the match between what the user
  typed and the invName at the investments Object.
**/
var substringMatcher = function(strs) {
  return function findMatches(q, cb) {
    var matches, substringRegex;

    // an array that will be populated with substring matches
    matches = [];

    // regex used to determine if a string contains the substring `q`
    substringRegex = new RegExp(q, 'i');

    // iterate through the pool of strings and for any string that
    // contains the substring `q`, add it to the `matches` array
    $.each(strs, function(i, str) {
      if (substringRegex.test(str.invName)) {
        matches.push(str.invName);
      }
    });

    cb(matches);
  };
};
