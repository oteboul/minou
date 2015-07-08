online_users = {};
online_users["tonton"] = false;
online_users["papi"] = false;
online_users["lartiste"] = true;

function DisplayUsers() {
  for (user in online_users) {
    var notification = online_users[user] ? "--> new message!" : "";
    var line = "<tr>";
    line += "<td><p>" + user + "</p></td>";
    line += "<td>" + notification + "</td>"
    line += "</tr>"
    $("table#users").append(line);
  }
}


$(document).ready(function() {
  DisplayUsers();

	$("table#users tr td p").click(function(event) {
    var user = $(event.target).text().toLowerCase().replace(" ", "_");
    $(event.target).css("color", "#ef3a3a");
    //window.location = "/client.html";
  });

});