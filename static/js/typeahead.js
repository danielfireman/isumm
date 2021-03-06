/**
 *  Matcher function - This function search to the match between what the user
 *  typed and the invName at the investments Object.
 *  @param {Object[]} strs -  array of objects to be matched
 */
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

/**
 *  Typeahead function - This function enables the typeahead input
 *  and writes the investment id in #hiddenInputElement to form submit
 *  @param {Object[]} invs - The investments list
 *  @param {string} invField - id of investment text field div
 *  @param {string} invHiddenField - id of hidden investment text field
 */
var typeaheadWrapper = (function(invs, invField, invHiddenField) {
  return function() {
     $('#' + invField + ' .typeahead').typeahead({
      hint: true,
      highlight: true,
      minLength: 1
    },
    {
      name: 'investments',
      source: substringMatcher(invs)
    }).on('typeahead:select', function(ev, suggestion) {
        listInv = [];
        $.each(invs, function (i, investment) {
          listInv.push(investment.invName);
        });
      $('#' + invHiddenField).val(invs[listInv.indexOf(suggestion)].invCode);
    });
  }
});
