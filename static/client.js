// The socket connection to the server.
var socket = new WebSocket(
  "wss://" + window.location.hostname + ":" + window.location.port + "/im");

// Who is online ?
var online_users = {};

// Who am I talking to ?
var recipient = "";

function selectRecepient(event) {
  var new_recipient = $(event.target).text();
  if (new_recipient == recipient) return;

  recipient = new_recipient;
  $("h2#recipient").html(recipient);
  online_users[recipient].new_message = false;
  DisplayUsers();

  $("div.messages").empty();
  for (i = 0; i < online_users[recipient].messages.length; i++) {
    var msg = online_users[recipient].messages[i];
    PrintTextMessage(msg.Text, msg.From == me);
  }
}

function DisplayUsers() {
  var list = "";
  for (user in online_users) {
    var notification = "";
    if (online_users[user].new_message) {
      notification = "<font color='red'>New Messsage!</font>";
    }
    if (user == recipient) {
      user = "<font color='#3b5998'>" + user + "</font>";
    }
    var line = "<tr>";
    line += "<td><p onclick='selectRecepient(event)'>" + user + "</p></td>";
    line += "<td>" + notification + "</td>";
    line += "</tr>";
    list += line;
  }
  $("table#users").html(list);
}

function PrintTextMessage(text, fromSelf) {
  if (!fromSelf) {
    text = "<span id='other'>" + text + "</span>";
  }
  $("div.messages").append(text + "<br><hr>");
  $("div.messages").animate({scrollTop: $("div.messages")[0].scrollHeight}, "slow");
}

function ShowIncomingMessage(msg) {
  // Save the message.
  online_users[msg.From].messages.push(msg);

  if (!recipient) {
    recipient = msg.From;
    $("h2#recipient").html(recipient);
  }

  if (recipient == msg.From) {
    PrintTextMessage(msg.Text, false);
  } else {
    online_users[msg.From].new_message = true;
    DisplayUsers();
  }
}

function SendPresence(socket) {
  var msg = {From: me, To: "", Text: ""};
  socket.send(JSON.stringify(msg));
}

// Wrap and send the message.
function SendMessage(socket) {
  if (socket.readyState === 1) {
    var text = $('input#send_box').val();
    if (text) {
      if (recipient && recipient in online_users) {
        var msg = {From: me, To: recipient, Text: text};
        socket.send(JSON.stringify(msg));
        $('input#send_box').val("");
        online_users[recipient].messages.push(msg);
        PrintTextMessage(text, true);
      } else {
        $("div.messages").append("Who are you talking to bro?<br><hr>");
      }
    }
  } else {
    $("div.messages").append("Socket is not ready :( <br><hr>");
  }
}

function ReceiveMessage(e) {
  var msg = JSON.parse(e.data);
  if (msg["Text"]) {  // This is a real message
    ShowIncomingMessage(msg);
  } else if (msg["To"]) { // This indicates a new online presence.
    online_users[msg.From] = {};
    online_users[msg.From].messages = [];
    online_users[msg.From].online = true;
    online_users[msg.From].new_message = false;
    DisplayUsers();
  } else {  // The sender is going offline.
    delete online_users[msg.From];
    DisplayUsers();
  }
}

$(document).ready(function() {
  socket.onopen = function() {
    SendPresence(socket);
  }
	socket.onmessage = function(e) {
	  ReceiveMessage(e);
	}
  socket.onclose = function() {
    $("div.messages").append(
      "<hr>The connection with the server has been closed!<br><hr>");
  }

  // Send Message on clicking the button or on pressing enter.
	$("input#send_button").click(function() {SendMessage(socket);});
  $("input#send_box").keydown(function(e) {
      if(e.which == 13) {
        e.preventDefault();
        SendMessage(socket);
      }
  });
});
