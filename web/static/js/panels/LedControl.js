if (window.location.search.indexOf('added') !== -1) {
	showPopup("Color added.", 3000, "success")
	window.history.replaceState(null, null, window.location.href.split("?")[0]);
}

// $(".color > .delete").on("click", (e) => {
// 	e.preventDefault()
// 	$(e.target).closest(".color").remove()
// })

// $(".color").on("click", (e) => {
// 	$(".color.selected").removeClass("selected")
// 	$(e.target).closest(".color").addClass("selected")
// })



$(".checkLine").on("click", (e) => {
	e.preventDefault()
	$(".color.selected").removeClass("selected")

	mode = $(e.target).closest(".option").attr("mode")

	
	$.ajax({
		type: "POST",
		url: `/back/ledControl/activate/${mode}`,
		// data: postData,
		success: function (r) {
			$(".option.active").removeClass("active")
			$(e.target).closest(".option").addClass("active")
			if (mode != "FadingColors" && mode != "FadingRainbow") {
				$($(e.target).closest(".option").find(".color")[0]).addClass("selected")
			}
			showPopup(`Led mode applied.`, 3000, "success")
			console.log(r)
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

	if (mode == "FadingColors" || mode == "FadingRainbow") // On fading colors all created are activated. On fading rainbow entire spectrum.
		return false

	
	$.ajax({
		type: "POST",
		url: `/back/ledControl/activate/${mode}/${hex}`,
		// data: postData,
		success: function (r) {
			$(e.target).closest(".color").addClass("selected")
			showPopup(`Led color applied.`, 3000, "success")
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
		url: `/back/ledControl/delete/${mode}/${hex}`,
		// data: postData,
		success: function (r) {
			$(".color.selected").removeClass("selected")
			option = $(e.target).closest(".option")
			$(e.target).closest(".color").remove()
			if(mode != "FadingColors")
				$($(option).find(".color")[0]).addClass("selected")
			showPopup(`Color removed.`, 3000, "success")
		},
		error: function(r) {
			showPopup(r.responseJSON.details, 3000, "error")
		}
	});
})

$(".color.new").on("click", (e) => {
	mode = $(e.target).closest(".option").attr("mode")
	$(".modal").attr("activeMode", mode)
	$(".darker").addClass("show")
})

$(document).on("click", "#addColor", (e) => {
	mode = $(".modal").attr("activeMode")
	hex =  $("#newHex").val() ? $("#newHex").val() : "000000"
	
	$.ajax({
		type: "POST",
		url: `/back/ledControl/add/${mode}/${hex}`,
		// data: postData,
		success: function (r) {
			showPopup(`Color added.`, 3000, "success")
			
			window.location.href = window.location.href + "?added";
		},
		error: function(r) {
			showPopup(r.responseJSON.details, 3000, "error")
		}
	});
})

$("#cancelAddColor").on("click", (e) => {
	mode = $(".modal").attr("activeMode", "")
	$(".darker").removeClass("show")
})

$("#brightness").on("change", () => {
	brightness = parseInt($("#brightness").val())
 	
	$.ajax({
		type: "POST",
		url: `/back/ledControl/brightness/${brightness}`,
		success: function (r) {
			showPopup(`Led color applied.`, 3000, "success")
		},
		error: function(r) {
			showPopup(r.responseJSON.details, 3000, "error")
		}
	});
})
$("#brightness").on("input", () => {
	brightness = parseInt($("#brightness").val())
	$("#brightnessVisual").text(brightness + "%")
})
$("#applyMSFadingRainbow").on("click", () => {
	amount = $("#MSamountRainbow").val()
	mode = "FadingRainbow"
	
	$.ajax({
		type: "POST",
		url: `/back/ledControl/cooldown/${mode}/${amount}`,
		success: function (r) {
			showPopup(`MS amount applied.`, 3000, "success")
		},
		error: function(r) {
			showPopup(r.responseJSON.details, 3000, "error")
		}
	});
})
$("#applyMSFadingColors").on("click", () => {
	amount = $("#MSamountColors").val()
	mode = "FadingColors"
	
	$.ajax({
		type: "POST",
		url: `/back/ledControl/cooldown/${mode}/${amount}`,
		success: function (r) {
			showPopup(`MS amount applied.`, 3000, "success")
		},
		error: function(r) {
			showPopup(r.responseJSON.details, 3000, "error")
		}
	});
})

