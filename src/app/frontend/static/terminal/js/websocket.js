export class ConnectionFactory {
    constructor(url, protocols) {
        this.url = url;
        this.protocols = protocols;
    }
    ;
    create() {
        return new Connection(this.url, this.protocols);
    }
    ;
}
export class Connection {
    constructor(url, protocols) {
        this.bare = new WebSocket(url, protocols);
    }
    open() {
        // nothing todo for websocket
    }
    ;
    close() {
        this.bare.close();
    }
    ;
    send(data) {
        this.bare.send(data);
    }
    ;
    isOpen() {
        if (this.bare.readyState == WebSocket.CONNECTING ||
            this.bare.readyState == WebSocket.OPEN) {
            return true;
        }
        return false;
    }
    onOpen(callback) {
        this.bare.onopen = (event) => {
            callback();
        };
    }
    ;
    onReceive(callback) {
        this.bare.onmessage = (event) => {
            callback(event.data);
        };
    }
    ;
    onClose(callback) {
        this.bare.onclose = (event) => {
            callback();
        };
    }
    ;
}
//# sourceMappingURL=websocket.js.map