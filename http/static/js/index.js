function validateEmail(email) {
    const domain = '{{domain}}';
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

function WsRetry(url, onMessageCallback, onOpenCallback, onCloseCallback, logEnabled = false) {
    this.url = url;
    this.onMessageCallback = onMessageCallback;
    this.onOpenCallback = onOpenCallback;
    this.onCloseCallback = onCloseCallback;
    this.logEnabled = logEnabled

    this.reconnectInterval = 2000;
    this.maxReconnectAttempts = 100;
    this.reconnectAttemptCount = 0;

    this.websocket = undefined;

    this.connect = () => {
        this.websocket = new WebSocket(this.url);
        this.websocket.onopen = () => {
            if (logEnabled) console.log('WebSocket connected');
            this.onOpenCallback();
            this.reconnectAttemptCount = 0;
        };
        this.websocket.onmessage = (event) => {
            this.onMessageCallback(event);
        };
        this.websocket.onclose = () => {
            if (this.logEnabled) console.warn('WebSocket disconnected');
            this.onCloseCallback();
            if (this.reconnectAttemptCount < this.maxReconnectAttempts) {
                setTimeout(() => {
                    if (this.logEnabled) console.log('WebSocket reconnecting...');
                    this.reconnectAttemptCount++;
                    this.connect();
                }, this.reconnectInterval);
            } else {
                if (this.logEnabled) console.error('WebSocket reconnect attempts exceeded');
            }
        };
    };

    this.send = (message) => {
        if (this.websocket.readyState === WebSocket.OPEN) {
            this.websocket.send(message);
        } else {
            if (this.logEnabled) console.log('WebSocket not connected');
        }
    };

    this.disconnect = () => {
        if (this.websocket.readyState !== WebSocket.CLOSED) {
            this.websocket.close();
        }
    };
}

function Inbox(section, address) {
    this.section = section // section element
    this.address = address; // string
    this.ws = undefined;
    this.destroyed = false
    this.connected = false

    // Crate inbox title
    this.title = document.createElement('h3')
    // Create inbox ul to store emails
    this.list = document.createElement('ul')
    // Add elements to inbox section
    this.section.appendChild(this.title)
    this.section.appendChild(this.list)
}

Inbox.prototype.wsUrl = function () {
    // Get current URL
    const url = window.location.origin;
    // Check if the current protocol is https or http
    const wsProtocol = (window.location.protocol === 'https:') ? 'wss' : 'ws';
    // Add WebSocket protocol
    return `${wsProtocol}://${url.split('//')[1]}/sync/${this.address}`;
};

Inbox.prototype.connect = function () {

    const eventHandler = (event) => {
        // Parse email
        const email = JSON.parse(event.data);
        // Add email to list
        const listItem = document.createElement('li');
        const emailLink = document.createElement('a')
        emailLink.innerText = 'From: ' + email.From + ', Subject: ' + email.Subject;
        emailLink.href = `/${email.ID}`
        emailLink.target = '_blank'
        listItem.appendChild(emailLink)
        this.list.appendChild(listItem);
    }

    const onConn = () => {
        this.title.innerText = `Waiting for mails from: ${this.address}`
        // Clean the email list
        while (this.list.firstChild) {
            this.list.removeChild(this.list.firstChild);
        }
    }

    const onConnClose = () => {
        if (!this.title.innerText.includes('Disconnected'))
            this.title.innerText = `${this.title.innerText} (Disconnected)`;
    }

    this.ws = new WsRetry(this.wsUrl(), eventHandler.bind(this), onConn.bind(this), onConnClose.bind(this))

    this.ws.connect()
};

Inbox.prototype.destroy = function () {
    this.ws.disconnect()
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
