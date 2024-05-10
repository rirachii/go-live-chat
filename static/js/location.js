
function detectUserLocation() {

    console.log("asking for location");

    let location = navigator.geolocation;

    if ( !location ){
        return null;
    }

    // set session storage
    function locationSuccess(pos) {

        const { coords, timestamp } = pos;
        const { latitude, longitude, accuracy } = coords;
    
        // console.log(latitude, longitude, accuracy);


        let locationData = {
            userLatitude: latitude,
            userLongitude: longitude,
            locAccuracy: accuracy
        };

        sessionStorage.setItem("location-data", JSON.stringify(locationData));
    
        return

    };

    
    function locationError(err) {
        console.log(err);

        sessionStorage.removeItem("location-data")

        return err;
    };


    const options = { };
    location.getCurrentPosition(locationSuccess, locationError, options);

    return 

}

function getLocationData() {

    let locationData = sessionStorage.getItem("location-data");

    if ( !locationData ){
        return {locationData: null};
    }
    
    console.log("location", locationData);

    return JSON.parse(locationData);


}