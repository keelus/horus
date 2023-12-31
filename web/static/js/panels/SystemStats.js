temperatureUnit = temperatureUnit || "C"

let gTemperature = Gauge(document.getElementById("graph0"),{max: 1,value: 0, extension: temperatureUnit});
let gCpu = Gauge(document.getElementById("graph1"),{max: 1,value: 0, extension: " %"});
let gRam = Gauge(document.getElementById("graph2"),{max: 1,value: 0, extension: " %"});
let gDisk = Gauge(document.getElementById("graph3"),{max: 1,value: 0, extension: " GB"});
let gSystem = Gauge(document.getElementById("graph4"),{max: 1,value: 0, extension: " s"});


function getStats(animate) {
	$.ajax({
		type: "GET",
		url: "/api/getStats",
		success: function (r) {
			$("#graph0").removeClass("good warning danger")
			$("#graph1").removeClass("good warning danger")
			$("#graph2").removeClass("good warning danger")
			$("#graph3").removeClass("good warning danger")
			$("#graph4").removeClass("good warning danger")

			let tempMax = 85
			let margin = [50, 55] // TODO: Separate margins for each Stat (if needed)

			gTemperature.setMaxValue(85)
			gCpu.setMaxValue(100)
			gRam.setMaxValue(100)
			gDisk.setMaxValue(r.DiskMax)
			gSystem.setMaxValue(100000)
			
			let temperatureConverted = r.Temperature // Celsius
			if (temperatureUnit == "F") {
				temperatureConverted = (r.Temperature * 1.8) + 32
				gTemperature.setMaxValue((85 * 1.8) + 32)
			}

			if(animate){
				gTemperature.setValueAnimated(temperatureConverted)
				gCpu.setValueAnimated(r.CPU)
				gRam.setValueAnimated(r.RAM)
				gDisk.setValueAnimated(r.Disk)
				gSystem.setValueAnimated(r.Uptime)
			} else {
				gTemperature.setValue(temperatureConverted)
				gCpu.setValue(r.CPU)
				gRam.setValue(r.RAM)
				gDisk.setValue(r.Disk)
				gSystem.setValue(r.Uptime)
			}
			
			
			if (r.Temperature * 100 / tempMax < margin[0]) {$("#graph0").addClass("good")} else if (r.Temperature * 100 / tempMax >= margin[0] && r.Temperature * 100 / tempMax < margin[1]) {$("#graph0").addClass("warning")} else {$("#graph0").addClass("danger")}
			if (r.CPU < margin[0]) {$("#graph1").addClass("good")} else if (r.CPU >= margin[0] && r.CPU < margin[1]) {$("#graph1").addClass("warning")} else {$("#graph1").addClass("danger")}
			if (r.RAM < margin[0]) {$("#graph2").addClass("good")} else if (r.RAM >= margin[0] && r.RAM < margin[1]) {$("#graph2").addClass("warning")} else {$("#graph2").addClass("danger")}
			if (r.Disk * 100 / r.DiskMax < margin[0]) {$("#graph3").addClass("good")} else if (r.Disk * 100 / r.DiskMax >= margin[0] && r.Disk * 100 / r.DiskMax < margin[1]) {$("#graph3").addClass("warning")} else {$("#graph3").addClass("danger")}
		},
		error: function(r) {
			showPopup("Error getting stats. Check browser console.", 3000, "error")
			console.log(r)
		}
	});
}
getStats(false) // Don't animate for the first time
setInterval(() => {getStats(true)}, 1000)