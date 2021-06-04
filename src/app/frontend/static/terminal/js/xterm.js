import { Terminal } from 'xterm';
import { FitAddon } from 'xterm-addon-fit';
import { WebLinksAddon } from 'xterm-addon-web-links';
const webLinksAddon = new WebLinksAddon();
const fitAddon = new FitAddon();
export class Xterm {
    constructor(elem) {
        this.resizeListener = () => {
            //console.log("resizeListener called");
            fitAddon.fit();
            this.term.scrollToBottom();
            this.showMessage(String(this.term.cols) + "x" + String(this.term.rows), this.messageTimeout);
        };
        this.elem = elem;
        this.term = new Terminal();
        this.term.loadAddon(webLinksAddon);
        this.term.loadAddon(fitAddon);
        this.message = elem.ownerDocument.createElement("div");
        this.message.className = "xterm-overlay";
        this.messageTimeout = 2000;
        this.term.open(elem);
        this.resizeListener();
        //터미널 내부에 선택된 텍스트가 있을경우 우클릭 복사/붙여넣기 기능 
        elem.addEventListener('contextmenu', (e) => {
            if (this.term.hasSelection()) {
                e.preventDefault();
                this.term.paste(this.term.getSelection());
                this.term.clearSelection();
            }
        });
    }
    ;
    info() {
        return { columns: this.term.cols, rows: this.term.rows };
    }
    ;
    output(data) {
        // this.term.write(this.decoder.decode(data));
        this.term.write(data);
    }
    ;
    showMessage(message, timeout) {
        this.message.textContent = message;
        this.elem.appendChild(this.message);
        if (this.messageTimer) {
            clearTimeout(this.messageTimer);
        }
        if (timeout > 0) {
            this.messageTimer = window.setTimeout(() => {
                this.elem.removeChild(this.message);
            }, timeout);
        }
    }
    ;
    removeMessage() {
        if (this.message.parentNode == this.elem) {
            this.elem.removeChild(this.message);
        }
    }
    ;
    setWindowTitle(title) {
        document.title = title;
    }
    ;
    setPreferences(value) {
    }
    ;
    onInput(callback) {
        this.term.onData((data) => {
            callback(data);
        });
    }
    ;
    onResize(callback) {
        this.term.onResize((data) => {
            callback(data.cols, data.rows);
        });
    }
    ;
    deactivate() {
        // this.term.off("data", () => { });
        // this.term.off("resize", () => { });
        this.term.blur();
    }
    ;
    reset() {
        this.removeMessage();
        this.term.clear();
    }
    ;
    close() {
        window.removeEventListener("resize", this.resizeListener);
        this.term.dispose();
    }
    ;
    focus() {
        this.term.focus();
    }
    ;
}
//# sourceMappingURL=xterm.js.map