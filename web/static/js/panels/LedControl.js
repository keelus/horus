if (window.location.search.indexOf('added') !== -1) {
	showPopup("Color added.", 3000, "success")
	window.history.replaceState(null, null, window.location.href.split("?")[0]);
}

$(".checkLine").on("click", (e) => {
	e.preventDefault()
	$(".color.selected").removeClass("selected")
	$(".gradient.selected").removeClass("selected")

	let mode = $(e.target).closest(".option").attr("mode")

	
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

	let mode = $(e.target).closest(".option").attr("mode")
	let hex = $(e.target).closest(".color").attr("hex")

		
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
	
	if($(e.target).hasClass("edit") || $(e.target).hasClass("delete") || $(e.target).is("svg") || $(e.target).is("path")){
		return false
	}

	$(".gradient.selected").removeClass("selected")

	let rawGradient = $(e.target).closest(".gradient").attr("raw-gradient")

	
	$.ajax({
		type: "POST",
		url: `/api/ledControl/activate/StaticGradient`,
		data: {
			"rawGradient": rawGradient
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
	let mode = $(e.target).closest(".option").attr("mode")
	let hex = $(e.target).closest(".color").attr("hex")

	$.ajax({
		type: "POST",
		url: `/api/ledControl/delete/${mode}`,
		data:  {
			"hexValue":hex
		},
		success: function (r) {
			let option = $(e.target).closest(".option")
			
			$(".color.selected").removeClass("selected")
			$($(".color")[0]).addClass("selected")
			$(e.target).closest(".color").remove()
			
			
			showPopup(`Color removed.`, 3000, "success")
		},
		error: function(r) {
			showPopup(r.responseJSON.details, 3000, "error")
		}
	});
})

$(".gradient .delete").on("click", (e) => {
	let rawGradient = $(e.target).closest(".gradient").attr("raw-gradient")
	
	$.ajax({
		type: "POST",
		url: `/api/ledControl/delete/StaticGradient`,
		data: {
			"rawGradient": rawGradient
		},
		success: function (r) {
			$(e.target).closest(".gradient").remove()
			$($(".gradient")[0]).addClass("selected")
			showPopup(`Gradient deleted.`, 3000, "success")
		},
		error: function(r) {
			showPopup(r.responseJSON.details, 3000, "error")
		}
	});
})

$(".color.new").on("click", (e) => {
	let mode = $(e.target).closest(".option").attr("mode")
	$("#newColorModal > .modal").attr("activeMode", mode)
	$("#newColorModal").addClass("show")
})

$(".gradient.new").on("click", (e) => {
	$("#newGradientModal > .modal").attr("editing", "false")
	$(".hexCode").remove()
	$("#addGradient").text("Add gradient")

	addColorToGradient("#000000")
	addColorToGradient("#FFFFFF")

	let mode = $(e.target).closest(".option").attr("mode")
	$("#newGradientModal > .modal").attr("activeMode", mode)
	$("#newGradientModal").addClass("show")
})

$(document).on("click", "#addColor", (e) => {
	let mode = $("#newColorModal > .modal").attr("activeMode")
	let hsv =  hexPicker.getColor()
	let hex = hsvToHex(hsv)
	
	$.ajax({
		type: "POST",
		url: `/api/ledControl/add/${mode}`,
		data:  {
			"hexValue": hex
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
	let hexValues = []
	let previousHexValues = []

	$(".hexCodes").children().slice(0, -1).each((i, e) => {
	  hexValues.push(hexFromGradients(e))
	})
	
	if ($(".modal[activemode='StaticGradient']").attr("editing") == "true") {
		previousHexValues = $(".modal[activemode='StaticGradient']").attr("previousRawGradient").replaceAll("#", "").split(",")
	}
	
	$.ajax({
		type: "POST",
		url: `/api/ledControl/add/StaticGradient`,
		data: {
			"hexValues": JSON.stringify(hexValues),
			"previousHexValues": JSON.stringify(previousHexValues),
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
	$("#newColorModal > .modal").attr("activeMode", "")
	$("#newColorModal").removeClass("show")
})

$("#cancelAddGradient").on("click", (e) => {
	$("#newGradientModal > .modal").attr("activeMode", "")
	$("#newGradientModal").removeClass("show")
	$(".hexCode").remove()
})

$(".addColorToGradient").on("click", () => {
	drawGradientPreview()
	addColorToGradient("#000000")
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
	if ($(".hexCode").length == 2) {
		showPopup("Gradient must have 2 color at least.", 3000, "error")
		return
	}
	$(e.target).closest(".hexCode").remove()
	drawGradientPreview()
})

$(document).on("key keyup keydown value input focus click change", ".hexCode", () => {
	drawGradientPreview()
})

function drawGradientPreview() {
	let gradientHexesStr = ""
	console.log($(".hexCodes").children())
	$(".hexCodes").children().slice(0, -1).each((i, e) => {
	  gradientHexesStr += "#" + hexFromGradients(e)
	  if (i != $(".hexCodes").children().length - 1 - 1){
		  gradientHexesStr += ","
	  }
	})

	$(".gradientPreview").css("background-image", `linear-gradient(to right, ${gradientHexesStr})`)
}

$("#brightness").on("change", () => {
	let brightness = $("#brightness").val() || "0"
 	
	$.ajax({
		type: "POST",
		url: `/api/ledControl/brightness/${brightness}`,
		success: function (r) {
			showPopup(`Brightness applied.`, 3000, "success")
		},
		error: function(r) {
			console.log(r.responseJSON)
			if ("brightness" in r.responseJSON) {
				$("#brightnessVisual").text(r.responseJSON.brightness + "%")
				$("#brightness").val(r.responseJSON.brightness)
			}
			showPopup(r.responseJSON.details, 3000, "error")
		}
	});
})

$("#brightness").on("input", () => {
	let brightness = $("#brightness").val() || "0"
	$("#brightnessVisual").text(brightness + "%")
})

$("#setCooldownFadingRainbow").on("click", () => {
	let amount = $("#cooldownFadingRainbow").val() || "0"
	let mode = "FadingRainbow"
	
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
	let amount = $("#cooldownBreathingColor").val() || "0"
	let mode = "BreathingColor"
	
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

// StaticColor & breathing color modal's picker & hsv converter
const hexPicker = new Pickr({
	el: '#hexPicker',
	default: '#FFFFFF',
	components: {
		preview: true,
		opacity: false,
		hue: true,
		output: {
			hex: true,
			rgba: false,
			hsva: false,
			input: true
		},
	},
});

function hsvToHex(hsv) {
	let color = tinycolor(hsv)
	return color.toHexString().replace("#", "")
}

// Gradient color modal's picker & others
function hexFromGradients(pickerDiv) {
	let rgb = $(pickerDiv).children(".color-picker").children(".button").css("background").split(") ")[0].replace("rgb(", "").split(", ")
	return tinycolor({r:rgb[0],g:rgb[1],b:rgb[2]}).toHexString().replace("#", "")
}

let lastGradientAdded = 0

function addColorToGradient(initialHex) {
	let length = $(".hexCodes").children().length - 1
	$($(".hexCodes").children()[length]).before(`
		<li class="hexCode">
			<div class="drag"><svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path d="M13 6h4l-5-6-5 6h4v12h-4l5 6 5-6h-4z"/></svg></div>
			<div class='gradientPicker${lastGradientAdded}'></div>
			<div class="remove"><svg clip-rule="evenodd" fill-rule="evenodd" stroke-linejoin="round" stroke-miterlimit="2" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path d="m20.015 6.506h-16v14.423c0 .591.448 1.071 1 1.071h14c.552 0 1-.48 1-1.071 0-3.905 0-14.423 0-14.423zm-5.75 2.494c.414 0 .75.336.75.75v8.5c0 .414-.336.75-.75.75s-.75-.336-.75-.75v-8.5c0-.414.336-.75.75-.75zm-4.5 0c.414 0 .75.336.75.75v8.5c0 .414-.336.75-.75.75s-.75-.336-.75-.75v-8.5c0-.414.336-.75.75-.75zm-.75-5v-1c0-.535.474-1 1-1h4c.526 0 1 .465 1 1v1h5.254c.412 0 .746.335.746.747s-.334.747-.746.747h-16.507c-.413 0-.747-.335-.747-.747s.334-.747.747-.747zm4.5 0v-.5h-3v.5z" fill-rule="nonzero"/></svg></div>
		</li>`);
	let newPicker = new Pickr({
	el: '.gradientPicker' + lastGradientAdded,
	default: initialHex,
	components: {
		preview: true,
		opacity: false,
		hue: true,
		output: {
			hex: true,
			rgba: false,
			hsva: false,
			input: true
			},
		},
	})
	lastGradientAdded++;
}

$(".gradientNote").on("click", drawGradientPreview)


$(".gradient .edit").on("click", (e) => {
	$(".hexCode").remove()
	$("#addGradient").text("Edit gradient")
	
	let rawGradient = $(e.target).closest(".gradient").attr("raw-gradient")
	let gradientColors = rawGradient.split(",")
	for(i=0;i<gradientColors.length;i++){
		addColorToGradient(gradientColors[i])
	}

	let mode = $(e.target).closest(".option").attr("mode")
	$("#newGradientModal > .modal").attr("activeMode", mode)
	$("#newGradientModal > .modal").attr("previousRawGradient", rawGradient)
	$("#newGradientModal").addClass("show")
	$("#newGradientModal > .modal").attr("editing", "true")

	drawGradientPreview()
})

$("#applyLedAmount").on("click", () => {
	let ledAmount = $("#ledAmount").val() || "0";

	$.ajax({
		type: "POST",
		url: `/api/ledControl/ledAmount/${ledAmount}`,
		success: function (r) {
			showPopup(`Led amount applied.`, 3000, "success")
		},
		error: function (r) {
			showPopup(r.responseJSON.details, 3000, "error")
		}
	});
})
