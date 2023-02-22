$(document).ready(function () {
    
    $('#myForm').ajaxForm(function () {
        $('#message_form').val('');
        $('#message_form').focus();

    });

    currentGroup = $("#groupuuid").val();
    currentGroupName = $("#roomName").val();
    currentUserId = $("#userid").val();
    conversation_history = $("#history").val();

    //show chat history
    myArray = conversation_history ? JSON.parse(conversation_history) : [];

    ApplyClass = '';

    for (var j = myArray.length - 1; j >= 0; j--) {

        if (!usericons[myArray[j].user_id]) {

            if (counter != 0 && counter % 3 == 0) {
                counter = 0;
            }
            usericons[myArray[j].user_id] = "user-logo-" + counter;
            counter++;
        }

        user_icon = usericons[myArray[j].user_id];

        d = new Date(myArray[j].created_at);

        if (currentUserId == myArray[j].user_id) {
            ApplyClass = 'pull-right sb5';
            ApplyStyle = '';
            ApplyClassOnName = 'hidden';
        } else {
            ApplyClass = 'sb2';
            ApplyStyle = '';
            ApplyClassOnName = '';
        }

        $('#messages').append("<div class='row block-view '><div class='row message-box " + ApplyClass + "' style='" + ApplyStyle + "'><i class='text-left'><span class='" + ApplyClassOnName + "'><i class='" + user_icon + " pull-left'></i><span class='pull-left' style='font-size:13px;font-weight:50'>" + myArray[j].sender_name + '</span></span></i>' + (ApplyClassOnName ? '' : '<br>') + '<p style="font-size:20px;font-weight:100;margin-bottom:0px;margin-bottom:0px">' + myArray[j].message + "</p><h6 class='text-right' style='font-size:10px;font-weight:100;margin-bottom:0px;position:absolute;bottom:0;margin-top:0px;margin-top:0px;color:grey;" + ((ApplyClassOnName != '') ? "left:2;" : "right:2;") + "'> " + d.toTimeString().split(" ")[0] + "</h6></div></div>");
    }
    document.getElementById('messages').scrollIntoView(false);



    if (!!window.EventSource) {

        var source = new EventSource('/events/conversations/' + currentGroup);
        source.addEventListener('message', function (e) {

            $('#' + currentGroup).addClass('hidden');

            $('#chats').prepend("<div style='cursor: pointer;' class='current-chat-box'  id='" + currentGroup + "'><span><i class='group-logo-mini pull-left'></i>" + currentGroupName + "</span></div>");

            arr = e.data.split("|");

            d = new Date(arr[2]);
            if (!usericons[arr[3]]) {

                usericons[arr[3]] = "user-logo-0";

            }

            user_icon = usericons[arr[3]];

            if (currentUserId == arr[3]) {
                ApplyClass = 'pull-right sb5';
                ApplyStyle = '';
                ApplyClassOnName = 'hidden';
            } else {
                ApplyClass = 'sb2';
                ApplyStyle = '';
                ApplyClassOnName = '';
            }

            $('#messages').append("<div class='row block-view'><div class='row message-box " + ApplyClass + "' style='" + ApplyStyle + "'><i class='text-left'><span class='" + ApplyClassOnName + "'><i class='" + user_icon + " pull-left'></i><span class='pull-left' style='font-size:13px;font-weight:50'>" + arr[0] + '</span></span></i>' + (ApplyClassOnName ? '' : '<br>') + '<p style="font-size:20px;font-weight:100;margin-bottom:0px;margin-bottom:0px">' + arr[1] + "</p><h6 class='text-right' style='font-size:10px;font-weight:100;margin-bottom:0px;position:absolute;bottom:0;margin-top:0px;margin-top:0px;color:grey;" + ((ApplyClassOnName != '') ? "left:2;" : "right:2;") + "'> " + d.toTimeString().split(" ")[0] + "</h6></div></div>");

            document.getElementById('messages').scrollIntoView({ behavior: 'smooth', block: 'end' });

        }, false);

        window.addEventListener('beforeunload', (event) => {
            // Cancel the event as stated by the standard.
            event.preventDefault();
            // Chrome requires returnValue to be set.
            source.close();
            return;
          });

    } else {
        alert("NOT SUPPORTED");
    }
});

function push_new_chat(group_id, group_name) {
    $('#' + group_id).addClass('hidden');
    count = $('#' + group_id).attr('data-new_messages');
    if (isNaN(count) || count == '') {
        count = 1;
    }
    else {
        count++;
    }
    applyFunction = "onclick='window.location.href=`" + group_id + "`'";
    $('#chats').prepend("<div style='cursor: pointer;' id='" + group_id + "' class='new-chat-box' " + applyFunction + " data-new_messages= '" + count + "'><span><i class='group-logo-mini pull-left'></i>" + group_name + " (" + count + " new messages)</span></div>");

}