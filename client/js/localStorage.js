"use strict";

function setDisplayName(name){

    const displayNameElement = document.getElementById("user-display-name") 

    localStorage.setItem("displayName", name);
    displayNameElement.textContent = getDisplayName();

}

function getDisplayName(){


    const displayName = 
        localStorage.getItem("displayName") ?  
        localStorage.getItem("displayName") : "[No Name]"

    return displayName
}

function injectDisplayName(){

    const displayNameElement = document.getElementById("user-display-name") 
    displayNameElement.innerHTML = getDisplayName();

}
