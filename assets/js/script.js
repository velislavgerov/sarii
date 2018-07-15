function ChooseAChore() {
	var chooseAChoreBtn = document.getElementById("choose_a_chore_btn");
  chooseAChoreBtn.setAttribute('class', 'siimple-spinner siimple-spinner--primary');

	var xmlhttp = new XMLHttpRequest();
	xmlhttp.onreadystatechange=function() {
	  if (xmlhttp.readyState==4 && xmlhttp.status==200) {
		chooseAChoreBtn.remove();
	    var response = xmlhttp.responseText; //if you need to do something with the returned value
	  	console.log(response);
	  	var element = document.getElementById("chore");
	  	element.innerHTML = response;
	  }
	}

	setTimeout(function() {
		xmlhttp.open("GET","/chore",true);
		xmlhttp.send();
	}, 2000)
}