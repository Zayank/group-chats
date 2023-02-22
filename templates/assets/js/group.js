$(document).ready(function() {
    $("#login").click(function(){
        if($("#name").val() == "") {
            alert("all fields are required");return;
        }
        
        $.ajax({
            type: "POST",
            url:  location.origin + "/group",
            data: $('#create-group-form').serialize(),
            cache: false,
            success: function (result)
            {
                window.location.href= location.origin + '/group/' + result.message.Id;
            }
        });
        
      });
});