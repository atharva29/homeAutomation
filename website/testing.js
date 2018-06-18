
    function initMap() {


      var options = {
        zoom: 15,
        center: {
          lat: 19.022,
          lng: 72.856689
        },
        mapTypeControlOptions: {
          style: google.maps.MapTypeControlStyle.DROPDOWN_MENU,
          mapTypeIds: ['roadmap', 'terrain', 'hybrid', 'satellite']
        },
        mapTypeControlOptions: {
          style: google.maps.MapTypeControlStyle.HORIZONTAL_BAR,
          position: google.maps.ControlPosition.TOP_CENTER
        }
    }



      // New Map
      var map = new google.maps.Map(document.getElementById('map'), options);

      //Listen for click on map
      google.maps.event.addListener(map, 'click', function(event) {


        //Add marker
        addMarker({
          coords: event.latLng
        })
      });







        // Add Marker
/*
        var marker = new google.maps.Marker({
          position :{lat:info.Lat,lng:info.Lng},
          map:map
        });

        var infoWindow = new google.maps.InfoWindow({
          content :'<h1>Kachra Kundi  </h1>'
        });lat: myStr.loc[m].lat,
            lng: myS

        marker.addListener('click',function(){
          infoWindow.open(map,marker);
        });

*/
      //



      //Recieved String from server
      // var myString = '{ "content" :"<h1> VJTI chi Kachra Kundi </h1>" ,lat" :19.022 ,"lng" :72.856689},{ "content" :"<h1> DADAR chi Kachra Kundi </h1>" ,lat" :19.0213 ,"lng" :72.84243}' ;
      // var obj = JSON.parse(myString);





      //Recieved object by Parsing String
      //Following is JSON Object

      var myStr = {

        "loc": [{
            "content": "<h1> VJTI chi Kachra Kundi </h1>",
        //    "lat": 19.022,
        //    "lng": 72.856689,
        //    "iconImage": './dustIcon.png'
              "lat": info.Lat,
              "lng": info.Lng,
          },
          {
            "content": "<h1> DADAR chi Kachra Kundi </h1>",
            "lat": 19.0213,
            "lng": 72.84243,
          //  "iconImage": './dustIcon.png'
          },
          {
            "content": "<h1> DOMBIVLI chi Kachra Kundi </h1>",
            "lat": 19.0213,
            "lng": 72.86243,
          //  "iconImage": './dustIcon.png'
          }
        ]
      }




/*
      //Array of markers
      var marker = [{
          coords: {
            lat: myStr.loc[0].lat,
            lng: myStr.loc[0].lng
          },
          content: myStr.loc[0].content,
          iconImage: './dustIcon.png'

        },
        {
          coords: {
            lat: myStr.loc[1].lat,
            lng: myStr.loc[1].lng
          },
          content: myStr.loc[1].content,
        }
      ];
      // Adds Markers that are resulted by parsing a string
      for (var i = 0; i < marker.length; i++) {
    //    addMarker(marker[i]);
      }
*/

    //  Add markers by gettings data from JSON object
      var m = 0 ;
      for (var m in myStr.loc) {

        var marker1 = {
          coords: {
            lat: myStr.loc[m].lat,
            lng: myStr.loc[m].lng
          },
          content: myStr.loc[m].content,
          iconImage :myStr.loc[m].iconImage
        } ;

        addMarker(marker1);
      }


      //Add marker Function

      function addMarker(props) {
        var marker = new google.maps.Marker({
          position: props.coords,
          map: map,
          animation: google.maps.Animation.DROP

        });

        //Check content
        if (props.content) {
          var infoWindow = new google.maps.InfoWindow({
            content: props.content,
          });

          marker.addListener('click', function() {
            infoWindow.open(map, marker, content);
          });
        }


        //Check for custom icon image
        if (props.iconImage) {
          //Set Icon image
          marker.setIcon(props.iconImage);
        }
      }
  }
