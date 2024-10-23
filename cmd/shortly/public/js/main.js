const BASE_URL = '/encode';
const DEFAULT_HEADERS = {'Content-Type': 'application/json'}
const NOTIFICATION_TIMEOUT = 3000; // 5 seconds

const pageState = {
    isNotificationDisplayed: false,
    hasValidationErrors: false,
}

// TODO: fix styles
function createNotification(message) {
    if (pageState.isNotificationDisplayed) {
        return null;
    }

    const notification = document.createElement('div');
    notification.className = `p-4 mb-3 flex items-center justify-between notification-enter`;
    
    notification.innerHTML = `
      <div class="flex items-center">
        <span>${message}</span>
      </div>
      <button 
        onclick="removeNotification(this.parentElement)" 
        class="ml-4 hover:opacity-75"
      >
        ×
      </button>
    `;

    return notification;
}

function removeNotification(notification) {
    notification.classList.remove('notification-enter');
    notification.classList.add('notification-exit');
    notification.addEventListener('animationend', () => {
      if (notification.classList.contains('notification-exit')) {
        notification.remove();
      }
    });
    pageState.isNotificationDisplayed = false;
}

function showNotification(message) {
    const container = document.getElementById('notification-container');
    const notification = createNotification(message);
    if (notification === null) {
        return;
    }

    container.appendChild(notification);
    pageState.isNotificationDisplayed = true;

    setTimeout(() => {
      if (notification.parentElement) {
        removeNotification(notification);
      }
    }, NOTIFICATION_TIMEOUT);
}

function copyToClipboard() {
    const textToCopy = document.getElementById("copy-text");
    navigator.clipboard.writeText(textToCopy.value);
    showNotification('Copied to clipboard');
}

async function handleSuccess(response) {
    const json = await response.json();
    const resultBlock = document.getElementById("result");
    const copyTextArea = document.getElementById("copy-text");
    resultBlock.classList.remove("hidden");
    copyTextArea.value = json.short_url;
}

async function handleValidationError(response) {
    pageState.hasValidationErrors = true;
    const urlErrHintArea = document.getElementById("url-err-hint");
    const aliasErrHintArea = document.getElementById("alias-err-hint");
    const json = await response.json();
    const {url, alias} = json.errors;
    if (url) {
        urlErrHintArea.classList.remove("hidden");
        urlErrHintArea.innerText = url;
    }

    if (alias) {
        aliasErrHintArea.classList.remove("hidden");
        aliasErrHintArea.innerText = alias;
    }
}

function cleanUpErrorHints() {
    const urlErrHintArea = document.getElementById("url-err-hint");
    const aliasErrHintArea = document.getElementById("alias-err-hint");
    urlErrHintArea.innerText = "";
    aliasErrHintArea.innerText = "";
    urlErrHintArea.classList.add("hidden");
    aliasErrHintArea.classList.add("hidden");
    pageState.hasValidationErrors = false;
}

async function makeItShort() {
    if (pageState.hasValidationErrors) {
        cleanUpErrorHints();
    }

    const longUrl = document.getElementById("long-url-input").value;
    const alias = document.getElementById("alias-input").value;
    const payload = JSON.stringify({url: longUrl, alias});

    try {
        document.body.style.cursor = "wait";
        const response = await fetch(BASE_URL, {method: "POST", body: payload, headers: DEFAULT_HEADERS});
        console.log('DEBUG', response);
        if (response.ok) {
            await handleSuccess(response);
        } else {
            console.log('DEBUG validation error');
            await handleValidationError(response);
        }
        document.body.style.cursor = "default";
    } catch (_e) {
        document.body.style.cursor = "default";
        console.log(_e);
        showNotification("Service is currently unavailable. Please try again later.");
    }
}

// Setup event listeners
document.getElementById("submit-button")?.addEventListener("click", makeItShort);
document.getElementById("copy-button")?.addEventListener("click", copyToClipboard);