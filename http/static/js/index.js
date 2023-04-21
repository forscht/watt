function validateEmail(email) {
    const domain = '{{domain}}'
    const regex = new RegExp(`^[\\w-.]+@${domain}$`);
    return regex.test(email);
}

function copyTextToClipboard(text) {
    if (navigator.clipboard) {
        navigator.clipboard.writeText(text).catch(() => {
            copyTextToClipboardFallback(text);
        });
    } else {
        copyTextToClipboardFallback(text);
    }
}

function copyTextToClipboardFallback(text) {
    const textArea = document.createElement('textarea');
    textArea.value = text;
    textArea.style.position = 'fixed';
    textArea.style.top = '0';
    textArea.style.left = '0';
    document.body.appendChild(textArea);
    textArea.focus();
    textArea.select();
    document.execCommand('copy');
    document.body.removeChild(textArea);
}


function Inbox(section, address) {
    this.section = section // section element
    this.address = address; // string
    this.ws = undefined;
}

Inbox.prototype.wsUrl = function () {
    // Get current URL
    const url = window.location.href;
    // Check if the current protocol is https or http
    const wsProtocol = (window.location.protocol === 'https:') ? 'wss' : 'ws';
    // Add WebSocket protocol
    return `${wsProtocol}://${url.split('//')[1]}sync/${this.address}`;
};

Inbox.prototype.connect = function () {

    // Create new websocket connection
    this.ws = new WebSocket(this.wsUrl());

    // Handle websocket onopen event
    this.ws.onopen = () => {
        // Crate inbox title
        this.title = document.createElement('h3')
        this.title.innerText = `Waiting for mails from: ${this.address}`
        // Create inbox ul to store emails
        this.list = document.createElement('ul')
        // Add elements to inbox section
        this.section.appendChild(this.title)
        this.section.appendChild(this.list)
    };

    // Handle new message event
    this.ws.onmessage = (event) => {
        // Parse email
        const email = JSON.parse(event.data);
        // Add email to list
        const listItem = document.createElement('li');
        const emailLink = document.createElement('a')
        emailLink.innerText = 'From: ' + email.From + ', Subject: ' + email.Subject;
        emailLink.href = `/${email.ID}`
        emailLink.target = '_blank'
        // listItem.classList.add("email-li")
        listItem.appendChild(emailLink)
        this.list.appendChild(listItem);
    };

    // Handle onclose event
    this.ws.onclose = () => {
        this.title.innerText = `${this.title.innerText} (Disconnected)`;
    };
};

Inbox.prototype.destroy = function () {
    if (this.ws.readyState !== WebSocket.CLOSED) {
        this.ws.close();
    }
    this.title.remove()
    this.list.remove()
};

function getRecentItemsFromLS(key, limit = 8) {
    const recent = localStorage.getItem(key) || '';
    const items = recent.split(',').filter(i => i !== '');
    const uniqueItems = removeDuplicates(items);
    if (limit && limit < uniqueItems.length) {
        return uniqueItems.slice(-limit);
    }
    return uniqueItems;
}

function getRecentInboxesFromLS() {
    return removeDuplicates(getRecentItemsFromLS('recent'));
}

function addInboxToLS(addr) {
    const recentItems = getRecentItemsFromLS('recent');
    recentItems.push(addr);
    localStorage.setItem('recent', recentItems.join(','));
}

function existInLS(addr) {
    return getRecentItemsFromLS('recent').includes(addr);
}

function removeDuplicates(arr) {
    return [...new Set(arr)];
}
