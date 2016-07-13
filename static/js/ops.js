function operationType(type){
  switch(type) {
    case "Balance":
        return "Saldo";
        break;
    case "Deposit":
        return "Aporte";
        break;
    case "Withdrawal":
        return "Retirada";
        break;
  }
}

// Default value for date-picker.
$(document).ready(function() {
	var now = new Date();
	// The two-digit day and month are needed due to time conversion formatting templates.
  var today = ('0' + now.getDate()).slice(-2) + '/' + ('0' + (now.getMonth()+1)).slice(-2) + '/' + now.getFullYear();
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
     $("#datePicker").datetimepicker({locale: 'pt-br', format: 'DD/MM/YYYY'});
});
