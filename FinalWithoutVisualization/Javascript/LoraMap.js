var mapVar;
var markerArray = [];
var length;
var flag;
// var ws = new WebSocket("wss://majestic-lake-clark-39709.herokuapp.com/ws")
var ws = new WebSocket("ws://localhost:8050/ws")
ws.onmessage = function (event) {
    info = JSON.parse(event.data);
    console.log(info);
    settings(info);
}

function settings(properties) {
    var i;
    console.log("Inside addMarker function")
    for (i = 0; i < properties.Slice.length; i++) {
        if ((properties.Slice[i].UpdateValue != 0) && (length == properties.Slice.length)) {
            console.log("Previous Length", length)
            console.log("Length is same");
            if (properties.Slice[i].UpdateValue != 0) {
                removeMarker(i);
                addMarker(properties, i)
            }
            continue;
        }
        addMarker(properties, i)
        console.log("MarkerArray Length" ,markerArray.length)
    }
    length = properties.Slice.length;
}

function addMarker(properties, i) {
    var Color;
    if (properties.Slice[i].Amount <= 25) {
        Color = "Images/dustYellow.png"
    } else if (properties.Slice[i].Amount <= 50) {
        Color = "Images/dustBlue.png"
    } else if (properties.Slice[i].Amount <= 75) {
        Color = "Images/dustBlack.png"
    } else if (properties.Slice[i].Amount <= 100) {
        Color = "Images/dustRed.png"
    } else {
        console.log("No Color Match :: Wrong Amount")
    }
    var markerVar = new google.maps.Marker({
        position: {
            lat: properties.Slice[i].Lat,
            lng: properties.Slice[i].Lnd
        },
        map: mapVar,
        icon: Color
    });

    markerArray[i] = markerVar;

    var infoWindow = new google.maps.InfoWindow({
        content: properties.Slice[i].Message
    });

    markerVar.addListener('mouseover', function () {
        infoWindow.open(mapVar, markerVar);
    });

    markerVar.addListener('mouseout', function () {
        infoWindow.close(mapVar, markerVar);
    });
    markerVar.addListener('click', function () {
        //  document.getElementById('myModal').modal('show');
        console.log(properties)
        $("#myModal").modal();
        document.getElementById('bar').style.width = properties.Slice[i].Amount + "%";
        document.getElementById("label").innerHTML = (properties.Slice[i].Amount * 1) + '%';

    });

}

function removeMarker(pos) {
    console.log("Removing Marker of Position :", pos);
    console.log("Marker to Remove", markerArray[pos])
    markerArray[pos].setMap(null);

}
// ws.send("Say hi to sever");
function initMap() {
    var options = {
        zoom: 13,
        center: {
            lat: 19.022,
            lng: 72.856
        }
    }
    mapVar = new google.maps.Map(document.getElementById('map'), options);
}