function showPopup(text, duration, type) {
	$(".popupMessage").remove()
	let popup = $(`<div class='popupMessage type-${type} '>${text}</div>`)
	$(document.body).append(popup)
	
	setTimeout(() => {
		popup.addClass("show");
	  }, 50);
	
	  setTimeout(() => {
		popup.addClass("hide");
		setTimeout(() => {
			popup.remove();
		}, 500);
	  }, duration);
}