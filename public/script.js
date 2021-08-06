// Urls
    var sock = null;
    var wsuri = "ws://" + window.location.host + "/api/socket";
    var api = "http://" + window.location.host + "/api/";

// Get the modal
    var modal = document.getElementById("modal");
    var modalKeyring = document.getElementById("modal-keyring");

// Get the button that opens the modal
    var editConfigBtn = document.getElementById("edit-config-link");
    var editKeyringBtn = document.getElementById("edit-keyring-link");

// Get the <span> element that closes the modal
    var closeBtnElements = document.getElementsByClassName("close");
    var cancelBtnElements = document.getElementsByClassName("cancelBtn");

// Window onload events
window.onload = function() {
    getConfig();
    getKeyring();
    pc.createOffer()
        .then(offer => pc.setLocalDescription(offer))
    // fetch node status every three seconds
    setInterval(checkNodeStatus, 3000);
    // Open websocket for monitoring logs
    sock = new WebSocket(wsuri);
    sock.onopen = function() {
        console.log("connected to " + wsuri);
    }
    sock.onclose = function(e) {
        console.log("connection closed (" + e.code + ")");
        console.log(e.data);
    }
    sock.onmessage = function(e) {
        //console.log("message received: " + e.data);
        var stdOut = document.getElementById('log');
        stdOut.innerHTML += "<br/>" + ansi2html_string(ansiconf, e.data)
        stdOut.scrollTop = stdOut.scrollHeight;
    }
    sock.onerror = function(e) {
        console.log(e);
    }
}

// Config control
    function getConfig() {
        const Http = new XMLHttpRequest();

        Http.onreadystatechange = function() {
            if (this.readyState == 4 && this.status == 200) {
                var resp = JSON.parse((this.responseText))
                console.log(resp)
                document.getElementById("conf-gas-adj").value = resp.Chain.GasAdjustment
                document.getElementById("conf-gas").value = resp.Chain.Gas
                document.getElementById("config-section-gas").innerHTML = resp.Chain.Gas
                document.getElementById("conf-gas-prices").value = resp.Chain.GasPrices
                document.getElementById("config-section-gas-price").innerHTML = resp.Chain.GasPrices
                document.getElementById("conf-chain-id").value = resp.Chain.ID
                document.getElementById("config-section-chain-id").innerHTML = resp.Chain.ID
                document.getElementById("conf-chain-rpc-addr").value = resp.Chain.RPCAddress
                document.getElementById("conf-chain-sim-exec").checked = resp.Chain.SimulateAndExecute
                document.getElementById("conf-node-interval-set-sessions").value = resp.Node.IntervalSetSessions
                document.getElementById("config-section-set-sessions").innerHTML = resp.Node.IntervalSetSessions
                document.getElementById("conf-node-interval-update-sessions").value = resp.Node.IntervalUpdateSessions
                document.getElementById("config-section-update-sessions").innerHTML = resp.Node.IntervalUpdateSessions
                document.getElementById("conf-node-interval-update-status").value = resp.Node.IntervalUpdateStatus
                document.getElementById("config-section-update-status").innerHTML = resp.Node.IntervalUpdateStatus
                document.getElementById("conf-node-listen-on").value = resp.Node.ListenOn
                document.getElementById("config-section-listen-on").innerHTML = resp.Node.ListenOn
                document.getElementById("conf-node-moniker").value = resp.Node.Moniker
                document.getElementById("config-section-moniker").innerHTML = resp.Node.Moniker
                document.getElementById("conf-node-price").value = resp.Node.Price
                document.getElementById("conf-node-provider").value = resp.Node.Provider
                document.getElementById("conf-node-remote-url").value = resp.Node.RemoteURL
                document.getElementById("status-bar-ip-addr-port").innerHTML = resp.Node.RemoteURL.match(RegExp(/[0-9].[0.9].*/))[0]
                document.getElementById("config-section-handshake-peers").innerHTML = resp.Handshake.Peers
                document.getElementById("config-section-handshake").innerHTML = (resp.Handshake.Enable == true ? "ENABLED" : "DISABLED")

                listenOnPort = resp.Node.ListenOn.split(":")[1]
            }
        };

        const url=api + 'config';
        Http.open("GET", url);
        Http.send();
    }

    function getKeyring() {
        const Http = new XMLHttpRequest();

        Http.onreadystatechange = function() {
            if (this.readyState == 4 && this.status == 200) {
                var resp = JSON.parse((this.responseText))
                resp.Keys.forEach(keyring => {
                    console.log(keyring)
                    document.getElementById("keyrings").innerHTML = document.getElementById("keyrings").innerHTML +
                        `<div class="keyring-block">
                            <div class="keyring">
                                <div class="keyring-name" style="font-weight: bold;">Name:</div>
                                <div class="keyring-operator" style="font-weight: bold;">Address:</div>
                                <div class="keyring-address" style="font-weight: bold;">Operator:</div>
                            </div>
                            <div class="keyring">
                                <div class="keyring-name">${keyring.Name}</div>
                                <div class="keyring-operator">${keyring.Operator}</div>
                                <div class="keyring-address">${keyring.Address}</div>
                            </div>
                        </div>`
                })
            }
        }

        const url=api + 'keys';
        Http.open("GET", url);
        Http.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
        Http.send();
    }

// Node
    function checkNodeStatus() {
        const Http = new XMLHttpRequest();

        Http.onreadystatechange = function() {
            if (this.readyState == 4 && this.status == 200) {
                var resp = JSON.parse((this.responseText))
                if (resp.Online == true) {
                    document.getElementById("node-status").innerHTML = "Running"
                } else {
                    document.getElementById("node-status").innerHTML = "Inactive"
                }
            }
        };

        const url = api + 'node';
        Http.open("GET", url);
        Http.send();
    }

    function startNode() {
        const Http = new XMLHttpRequest();

        Http.onreadystatechange = function() {
            checkNodeStatus();
        }

        const url = api + 'node/start/stream';
        Http.open("GET", url);
        Http.send();
    };

    function kill() {
        const Http = new XMLHttpRequest();

        Http.onreadystatechange = function() {
            checkNodeStatus();
        }

        const url= api + '/node/kill';
        Http.open("POST", url);
        Http.send();
    }

// When the user clicks on the button, open the modal
    editConfigBtn.onclick = function() {
        modal.style.display = "block";
    }

    editKeyringBtn.onclick = function() {
        modalKeyring.style.display = "block";
    }

// When the user clicks on <span> (x), close the modal
    for (let closeBtn of closeBtnElements) {
        closeBtn.onclick = function() {
            modal.style.display = "none";
            modalKeyring.style.display = "none";
        }
    }
    for (let cancelBtn of cancelBtnElements) {
        cancelBtn.onclick = function() {
            modal.style.display = "none";
            modalKeyring.style.display = "none";
        }
    }

// When the user clicks anywhere outside of the modal, close it
    window.onclick = function(event) {
        if (event.target == modal) {
            modal.style.display = "none";
        }
        if (event.target == modalKeyring) {
            modalKeyring.style.display = "none";
        }
    }

