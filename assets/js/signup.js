$(document).ready(function() {
    $('#join').bind('click', join);
    $('#cancel').bind('click', cancel);
});

// The function join takes all of the variables and submits them via
// an ajax request to /post/newuser.
function join() {
    $.ajax({
	type: 'POST',
	url: '/post/newuser',
	data: {
	    'fname': $('#fname').val(),
	    'lname': $('#lname').val(),
	    'desc': $('#desc').val(),
	    'website': $('#website').val(),
	    'address': $('#address').val()
	},
	success: function(response) {
	    alert(response);
	},
	error: function(response) {
	    // The only time this should error is if the server stops.
	    // TODO write error
	}
    });
}

// The function cancel erases everything in the form and makes it
// blank again.
function cancel() {
    $('#fname').val("");
    $('#lname').val("");
    $('#desc').val("");
    $('#website').val("");
}