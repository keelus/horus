$(".botonGuardar").on("click", (e) => {
	category = $(e.target).closest(".category").attr("category")
	$(`[category='${category}']`).find(".hasError").removeClass("hasError")

	postData = {}

	switch(category) {
	case "UserInfo":
		postData.Username = $("#userInfo0").val()
		postData.Password = $("#userInfo1").val()
		break
	case "SessionSettings":
		activeRadio = $('input[name="sessionDurationType"]:checked').attr("id")
		if(activeRadio == "sessionDuration0"){
			postData.Lifespan = -1
			postData.Unit = ""
		} else {
			postData.Lifespan = $("#sessionDuration2").val()
			postData.Unit = $("#sessionDuration3").val()
		}
		break
	case "LedControl":
		postData.LedControl = JSON.stringify([$("#ledControl0").prop("checked"),$("#ledControl0").prop("checked")])
		break
	case "SystemStats":
		postData.SystemStats = JSON.stringify([$("#systemStats0").prop("checked"), $("#systemStats1").prop("checked"),
								$("#systemStats2").prop("checked"), $("#systemStats3").prop("checked"), $("#systemStats4").prop("checked")])
		break
	case "Logging":
		postData.Logging = $("#logging0").prop("checked")
		break
	case "Security":
		postData.UserInput = $("#security0").prop("checked")
		break
	case "Units":
		postData.TimeMode = $("#units0").val()
		postData.Temperature = $("#units1").val()
		break
	case "Design":
		// TODO
		break
	case "Deletion":
		// TODO
		break
	}

	
	$.ajax({
		type: "POST",
		url: `/back/saveConfiguration/${category}`,
		data: postData,
		success: function (r) {
			showPopup(`${category} saved.`, 3000, "success")
		},
		error: function(r) {
			popupMessage = ""
			for(i=0;i<r.responseJSON.length;i++){
				console.error(r.responseJSON[i].details)
				popupMessage += `• ${r.responseJSON[i].details}<br>`
				

				field = r.responseJSON[i].field
				if(field != "") {
					$(`#${field}`).addClass("hasError")
				}
			}

			showPopup(popupMessage, 3000, "error")
		}
	});
})