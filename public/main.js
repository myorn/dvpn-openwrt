var sock = null;
var wsuri = "ws://" + document.location.host + "/api/socket";

window.onload = function() {
    getConfig();
    getKeyring();

    // fetch node status every three seconds
    setInterval(checkNodeStatus, 3000);

    sock = new WebSocket(wsuri);

    sock.onopen = function() {
        console.log("connected to " + wsuri);
    }

    sock.onclose = function(e) {
        console.log("connection closed (" + e.code + ")");
        console.log(e.data);
    }

    sock.onmessage = function(e) {
        var stdOut = document.getElementById('console');
        stdOut.innerHTML += "<br/>" + ansi2html_string(ansiconf, e.data)
        stdOut.scrollTop = stdOut.scrollHeight;
    }

    sock.onerror = function(e) {
        console.log(e);
    }
};

function startNode() {
    document.getElementById("btn-start-node").style.display = "none"
    document.getElementById("btn-stop-node").style.display = "none"

    const Http = new XMLHttpRequest();

    Http.onreadystatechange = function() {
        checkNodeStatus();
    }

    const url='http://localhost:9000/api/node/start/stream';
    Http.open("GET", url);
    Http.send();
};

function getNode() {
    console.log("getting node status")
    document.getElementById("btn-start-node").style.display = "none"
    document.getElementById("btn-stop-node").style.display = "none"
    checkNodeStatus();
}

function checkNodeStatus() {
    const Http = new XMLHttpRequest();

    Http.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            var resp = JSON.parse((this.responseText))
            console.log(resp)
            if (resp.Online == true) {
                document.getElementById("node-status").innerHTML = "Running";
                document.getElementById("process-id").innerHTML = resp.Pid;
                document.getElementById("node-status").style.color = "darkgreen"
                if (document.getElementById("btn-start-node").style.display == "block") {
                    document.getElementById("btn-start-node").style.display = "none"
                }
                if (document.getElementById("btn-stop-node").style.display == "none") {
                    document.getElementById("btn-stop-node").style.display = "block"
                }
            } else {
                document.getElementById("node-status").innerHTML = "Stopped";
                document.getElementById("process-id").innerHTML = "Unavailable";
                document.getElementById("node-status").style.color = "red"
                if (document.getElementById("btn-start-node").style.display == "none") {
                    document.getElementById("btn-start-node").style.display = "block"
                }
                if (document.getElementById("btn-stop-node").style.display == "block") {
                    document.getElementById("btn-stop-node").style.display = "none"
                }
            }
        }
    };

    const url='http://localhost:9000/api/node';
    Http.open("GET", url);
    Http.send();
}

function kill() {
    document.getElementById("btn-start-node").style.display = "none"
    document.getElementById("btn-stop-node").style.display = "none"
    const Http = new XMLHttpRequest();

    Http.onreadystatechange = function() {
        getNode();
    }

    const url='http://localhost:9000/api/node/kill';
    Http.open("POST", url);
    Http.send();
}

function getConfig() {
    const Http = new XMLHttpRequest();

    Http.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            var resp = JSON.parse((this.responseText))
            console.log(resp)
            document.getElementById("conf-gas-adj").value = resp.Chain.GasAdjustment
            document.getElementById("conf-gas").value = resp.Chain.Gas
            document.getElementById("conf-gas-prices").value = resp.Chain.GasPrices
            document.getElementById("conf-chain-id").value = resp.Chain.ID
            document.getElementById("conf-chain-rpc-addr").value = resp.Chain.RPCAddress
            document.getElementById("conf-chain-sim-exec").checked = resp.Chain.SimulateAndExecute
            document.getElementById("conf-node-interval-set-sessions").value = resp.Node.IntervalSetSessions
            document.getElementById("conf-node-interval-update-sessions").value = resp.Node.IntervalUpdateSessions
            document.getElementById("conf-node-interval-update-status").value = resp.Node.IntervalUpdateStatus
            document.getElementById("conf-node-listen-on").value = resp.Node.ListenOn
            document.getElementById("conf-node-moniker").value = resp.Node.Moniker
            document.getElementById("conf-node-price").value = resp.Node.Price
            document.getElementById("conf-node-provider").value = resp.Node.Provider
            document.getElementById("conf-node-remote-url").value = resp.Node.RemoteURL
        }
    };

    const url='http://localhost:9000/api/config';
    Http.open("GET", url);
    Http.send();
}

function saveConfig() {
    const Http = new XMLHttpRequest();

    Http.onreadystatechange = function() {
        if (this.readyState == 4 && this.status == 200) {
            getConfig();
        }
    }

    var config = {
        "Chain": {
            "GasAdjustment": parseFloat(document.getElementById("conf-gas-adj").value),
            "Gas": parseFloat(document.getElementById("conf-gas").value),
            "GasPrices": document.getElementById("conf-gas-prices").value,
            "ID": document.getElementById("conf-chain-id").value,
            "RPCAddress": document.getElementById("conf-chain-rpc-addr").value,
            "SimulateAndExecute": document.getElementById("conf-chain-sim-exec").checked
        },
        "Handshake": {
            "Enable": false,
            "Peers": 8
        },
        "Keyring": {
            "Backend": "test",
            "From": "nuage"
        },
        "Node": {
            "IntervalSetSessions": document.getElementById("conf-node-interval-set-sessions").value,
            "IntervalUpdateSessions": document.getElementById("conf-node-interval-update-sessions").value,
            "IntervalUpdateStatus": document.getElementById("conf-node-interval-update-status").value,
            "ListenOn": document.getElementById("conf-node-listen-on").value,
            "Moniker": document.getElementById("conf-node-moniker").value,
            "Price": document.getElementById("conf-node-price").value,
            "Provider": document.getElementById("conf-node-provider").value,
            "RemoteURL": document.getElementById("conf-node-remote-url").value
        },
        "Qos": {}
    }

    const url='http://localhost:9000/api/config';
    Http.open("POST", url);
    Http.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
    Http.send(JSON.stringify(config));
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


    const url='http://localhost:9000/api/keys';
    Http.open("GET", url);
    Http.setRequestHeader("Content-Type", "application/json;charset=UTF-8");
    Http.send();
}

function restoreKey() {

}