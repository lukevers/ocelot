$(document).ready(function() {
    $('#join').bind('click', join);
    $('#cancel').bind('click', cancel);
    $('#fname').bind('keyup', function() {
	if ($('#fname').val() == "") { 
	    $('#fname').addClass('empty'); 
	} else $('#fname').removeClass('empty'); 
    });
    $('#lname').bind('keyup', function() {
	if ($('#lname').val() == "") {
	    $('#lname').addClass('empty');
	} else $('#lname').removeClass('empty'); 
    });
});

// The function join takes all of the variables and submits them via
// an ajax request to /api/newuser.
function join() {
    $.ajax({
	type: 'POST',
	url: '/api/token',
	data: {
	    'address': $('#address').val()
	},
	success: function(token) {
	    $.ajax({
		type: 'POST',
		url: '/api/newuser',
		data: {
		    'fname': $('#fname').val(),
		    'lname': $('#lname').val(),
		    'desc': $('#desc').val(),
		    'website': $('#website').val(),
		    'address': $('#address').val(),
		    'token': token
		},
		success: function(response) {
		    notifyMe(response);
		},
		error: function(response) {
		    // The only time this should error is if the server stops.
		    // TODO write error
		}
	    });
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
    $('#lname, #fname').removeClass('empty');
}

// The function notifyMe creates a notification in the top right
// corner of the screen telling you what happened with your request.
function notifyMe(what) {
    $('body').append('<div id="notify">'+what+'</div>');
    $('#notify').fadeIn(500);
    setTimeout(function() {
	$('#notify').fadeOut(500, function() {
	    $('#notify').remove();
	    if (what.indexOf("Success") != -1) {
		location.reload();
	    }
	});
    }, 3000);
    
}