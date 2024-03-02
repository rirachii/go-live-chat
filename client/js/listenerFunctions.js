
function formInjectDisplayName(event) {
    const eventTarget = event.target


    if (eventTarget.tagName === "FORM"){

        // inject display-name into form
        const displayName = localStorage.getItem("displayName") ?
            localStorage.getItem("displayName") : "no name";

        const displayNameInput = eventTarget.querySelector('input[name="display-name"]');
        if (displayNameInput) {
            displayNameInput.value = displayName;
        }

    }

}