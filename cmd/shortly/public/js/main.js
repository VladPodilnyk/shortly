const BASE_URL = '/encode';
const DEFAULT_HEADERS = {'Content-Type': 'application/json'}

async function handleSuccess(response) {
    const json = await response.json();
    const resultBlock = document.getElementById("result");
    const copyTextArea = document.getElementById("copy-text");
    resultBlock.classList.remove("hidden");
    copyTextArea.value = json.short_url;
}

function handleFailure(response) {
    // TODO: Handle error
    console.error(response);
}

function copyToClipboard() {
    // TODO: Implement copy to clipboard
}

async function makeItShort() {
    const longUrl = document.getElementById("long-url-input").value;
    const alias = document.getElementById("alias-input").value;
    const payload = JSON.stringify({url: longUrl, alias});

    try {
        console.time("Time taken"); 
        const response = await fetch(BASE_URL, {method: "POST", body: payload, headers: DEFAULT_HEADERS});
        await handleSuccess(response);
        console.timeEnd("Time taken");
    } catch (e) {
        handleFailure(e);
    }
}

// Setup event listeners
document.getElementById("submit-button")?.addEventListener("click", makeItShort);