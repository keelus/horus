if (window.location.search.indexOf('added') !== -1) {
	showPopup("Color added.", 3000, "success")
	window.history.replaceState(null, null, window.location.href.split("?")[0]);
}

$(".checkLine").on("click", (e) => {
	e.preventDefault()
	$(".color.selected").removeClass("selected")

	mode = $(e.target).closest(".option").attr("mode")

	
	$.ajax({
		type: "POST",
		url: `/api/ledControl/activate/${mode}`,
		// data: postData,
		success: function (r) {
			$(".option.active").removeClass("active")
			$(e.target).closest(".option").addClass("active")

			if (mode == "StaticGradient") {
				$($(e.target).closest(".option").find(".gradient")[0]).addClass("selected")
			} else if (mode == "StaticColor" || mode == "BreathingColor"){
				$($(e.target).closest(".option").find(".color")[0]).addClass("selected")
			}

			showPopup(`Led mode applied.`, 3000, "success")
		},
		error: function(r) {
			showPopup(r.responseJSON.details, 3000, "error")
		}
	});
})


$(".color:not(.color.new)").on("click", (e) => {
	e.preventDefault()

	if($(e.target).hasClass("delete") || $(e.target).is("svg")){
		return false
	}

	$(".color.selected").removeClass("selected")

	mode = $(e.target).closest(".option").attr("mode")
	hex = $(e.target).closest(".color").attr("hex")

		
	$.ajax({
		type: "POST",
		url: `/api/ledControl/activate/${mode}`,
		data: {
			"hexValue": hex
		},
		success: function (r) {
			$(e.target).closest(".color").addClass("selected")
			showPopup(`Led color applied.`, 3000, "success")
		},
		error: function(r) {
			showPopup(r.responseJSON.details, 3000, "error")
		}
	});
})

$(".gradient:not(.gradient.new)").on("click", (e) => {
	e.preventDefault()

	if($(e.target).hasClass("delete") || $(e.target).is("svg")){
		return false
	}

	$(".gradient.selected").removeClass("selected")

	rawGradient = $(e.target).closest(".gradient").attr("raw-gradient")

	
	$.ajax({
		type: "POST",
		url: `/api/ledControl/activate/StaticGradient`,
		data: {
			"rawGradient":rawGradient
		},
		success: function (r) {
			$(e.target).closest(".gradient").addClass("selected")
			showPopup(`Led gradient applied.`, 3000, "success")
		},
		error: function(r) {
			showPopup(r.responseJSON.details, 3000, "error")
		}
	});
})

$(".color .delete").on("click", (e) => {
	mode = $(e.target).closest(".option").attr("mode")
	hex = $(e.target).closest(".color").attr("hex")

	$.ajax({
		type: "POST",
		url: `/api/ledControl/delete/${mode}`,
		data:  {
			"hexValue":hex
		},
		success: function (r) {
			$(".color.selected").removeClass("selected")
			option = $(e.target).closest(".option")
			$(e.target).closest(".color").remove()
			if(mode != "StaticGradient")
				$($(option).find(".color")[0]).addClass("selected")
			showPopup(`Color removed.`, 3000, "success")
		},
		error: function(r) {
			showPopup(r.responseJSON.details, 3000, "error")
		}
	});
})

$(".gradient .delete").on("click", (e) => {
	rawGradient = $(e.target).closest(".gradient").attr("raw-gradient")
	
	$.ajax({
		type: "POST",
		url: `/api/ledControl/delete/StaticGradient`,
		data: {
			"rawGradient":rawGradient
		},
		success: function (r) {
			$(e.target).closest(".gradient").remove()
			showPopup(`Gradient deleted.`, 3000, "success")
		},
		error: function(r) {
			showPopup(r.responseJSON.details, 3000, "error")
		}
	});
})

$(".color.new").on("click", (e) => {
	mode = $(e.target).closest(".option").attr("mode")
	$("#newColorModal > .modal").attr("activeMode", mode)
	$("#newColorModal").addClass("show")
})

$(".gradient.new").on("click", (e) => {
	mode = $(e.target).closest(".option").attr("mode")
	$("#newGradientModal > .modal").attr("activeMode", mode)
	$("#newGradientModal").addClass("show")
})

$(document).on("click", "#addColor", (e) => {
	mode = $("#newColorModal > .modal").attr("activeMode")
	hex =  $("#newHex").val() ? $("#newHex").val() : "000000"
	
	$.ajax({
		type: "POST",
		url: `/api/ledControl/add/${mode}`,
		data:  {
			"hexValue":hex
		},
		success: function (r) {
			showPopup(`Color added.`, 3000, "success")
			
			window.location.href = window.location.href + "?added";
		},
		error: function(r) {
			showPopup(r.responseJSON.details, 3000, "error")
		}
	});
})

$(document).on("click", "#addGradient", (e) => {
	hexValues = []
	
	$(".hexCodes").children().slice(0, -1).each((i, e) => {
	  hexValues.push($(e).find("input").val())
	})
	console.log(hexValues)
	
	$.ajax({
		type: "POST",
		url: `/api/ledControl/add/StaticGradient`,
		data: {
			"hexValues":JSON.stringify(hexValues)
		},
		success: function (r) {
			showPopup(`Gradient added.`, 3000, "success")
			
			window.location.href = window.location.href + "?added";
		},
		error: function(r) {
			showPopup(r.responseJSON.details, 3000, "error")
		}
	});
})

$("#cancelAddColor").on("click", (e) => {
	mode = $("#newColorModal > .modal").attr("activeMode", "")
	$("#newColorModal").removeClass("show")
})

$("#cancelAddGradient").on("click", (e) => {
	mode = $("#newGradientModal > .modal").attr("activeMode", "")
	$("#newGradientModal").removeClass("show")
})

$(".addColorToGradient").on("click", () => {
	lastIndex = $(".hexCodes").children().length - 1 - 1
	$($(".hexCodes").children()[lastIndex]).after(`
		<li class="hexCode">
			<div class="drag"></div>
			<div>#<input type="text" placeholder="XXXXXX" maxlength="6" value="000000"></div>
			<div class="remove">Remove</div>
		</li>`);
	drawGradientPreview()
})

$(document).ready(function() {
    $(".hexCodes").sortable({
		handle:".drag",
		update: function(event, ui) {
			drawGradientPreview()
		}
	});
	$(".hexCodes").disableSelection();
});


$(document).on("click", ".remove", (e) => {
	$(e.target).closest(".hexCode").remove()
	drawGradientPreview()
})

$(document).on("key keyup keydown value input focus click change", ".hexCode input", () => {
	console.log("update")
	drawGradientPreview()
})

function drawGradientPreview() {
	gradientHexesStr = ""
	console.log($(".hexCodes").children())
	$(".hexCodes").children().slice(0, -1).each((i, e) => {
	  gradientHexesStr += "#" + $(e).find("input").val()
	  if (i != $(".hexCodes").children().length - 1 - 1){
		  gradientHexesStr += ","
	  }
	})

	$(".gradientPreview").css("background-image", `linear-gradient(to right, ${gradientHexesStr})`)
}


$("#brightness").on("change", () => {
	brightness = parseInt($("#brightness").val())
 	
	$.ajax({
		type: "POST",
		url: `/api/ledControl/brightness/${brightness}`,
		success: function (r) {
			showPopup(`Brightness applied.`, 3000, "success")
		},
		error: function(r) {
			if (r.responseJSON.brightness != null) {
				$("#brightness").val(r.responseJSON.brightness)
			}
			showPopup(r.responseJSON.details, 3000, "error")
		}
	});
})
$("#brightness").on("input", () => {
	brightness = parseInt($("#brightness").val())
	$("#brightnessVisual").text(brightness + "%")
})
$("#setCooldownFadingRainbow").on("click", () => {
	amount = $("#cooldownFadingRainbow").val()
	mode = "FadingRainbow"
	
	$.ajax({
		type: "POST",
		url: `/api/ledControl/cooldown/${mode}/${amount}`,
		success: function (r) {
			showPopup(`Cooldown applied.`, 3000, "success")
		},
		error: function(r) {
			showPopup(r.responseJSON.details, 3000, "error")
		}
	});
})
$("#setCooldownBreathingColor").on("click", () => {
	amount = $("#cooldownBreathingColor").val()
	mode = "BreathingColor"
	
	$.ajax({
		type: "POST",
		url: `/api/ledControl/cooldown/${mode}/${amount}`,
		success: function (r) {
			showPopup(`Cooldown applied.`, 3000, "success")
		},
		error: function(r) {
			showPopup(r.responseJSON.details, 3000, "error")
		}
	});
})
