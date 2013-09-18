$(document).ready(function() {
    $('#join').bind('click', join);
    $('#cancel').bind('click', cancel);
});

// The function join takes all of the variables and submits them via
// an ajax request to /api/newuser.
function join() {
    
}

// The function cancel erases everything in the form and makes it
// blank again.
function cancel() {
    $('#fname').val("");
    $('#lname').val("");
    $('#desc').val("");
    $('#website').val("");
}