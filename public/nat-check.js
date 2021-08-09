var infoprint = ""
var userIP = ""
var listenOnPort = ""

function syntaxHighlight(json) {
    if (typeof json != 'string') {
        json = JSON.stringify(json, undefined, 2);
    }
    json = json.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
    return json.replace(/("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g, function (match) {
        var cls = 'number';
        if (/^"/.test(match)) {
            if (/:$/.test(match)) {
                cls = 'key';
            } else {
                cls = 'string';
            }
        } else if (/true|false/.test(match)) {
            cls = 'boolean';
        } else if (/null/.test(match)) {
            cls = 'null';
        }
        return '<span class="' + cls + '">' + match + '</span>';
    });
}

function parseCandidate(line) {
    var parts;
    // Parse both variants.
    if (line.indexOf('a=candidate:') === 0) {
        parts = line.substring(12).split(' ');
    } else {
        parts = line.substring(10).split(' ');
    }

    var candidate = {
        foundation: parts[0],
        component: parts[1],
        protocol: parts[2].toLowerCase(),
        priority: parseInt(parts[3], 10),
        ip: parts[4],
        port: parseInt(parts[5], 10),
        // skip parts[6] == 'typ'
        type: parts[7]
    };

    for (var i = 8; i < parts.length; i += 2) {
        switch (parts[i]) {
            case 'raddr':
                candidate.relatedAddress = parts[i + 1];
                break;
            case 'rport':
                candidate.relatedPort = parseInt(parts[i + 1], 10);
                break;
            case 'tcptype':
                candidate.tcpType = parts[i + 1];
                break;
            default: // Unknown extensions are silently ignored.
                break;
        }
    }
    return candidate;
};

var candidates = {};
var totalCandidates = 0;
var pc = new RTCPeerConnection({iceServers: [
        {urls: 'stun:stun1.l.google.com:19302'},
        {urls: 'stun:stun2.l.google.com:19302'},
        {urls: 'stun:stun3.l.google.com:19302'},
        {urls: 'stun:stun4.l.google.com:19302'},
        {urls: 'stun:stun.schlund.de'},
        {urls: 'stun:stun.voipstunt.com'}
    ]});
pc.createDataChannel("nat-checker");
pc.onicecandidate = function(e) {
    if (e.candidate && e.candidate.candidate.indexOf('srflx') !== -1) {
        totalCandidates++
        var cand = parseCandidate(e.candidate.candidate);
        userIP = cand.ip;
        console.log(e)
        infoprint += "Candidate: " + cand.ip + ":" + syntaxHighlight(cand.port) + " | " + syntaxHighlight(e.candidate.usernameFragment) + "<br/>"
        if (!candidates[cand.relatedPort]) candidates[cand.relatedPort] = [];
        candidates[cand.relatedPort].push(cand.port);
        // All ICE candidates have been sent
    } else {
        // all candidates should have the same relative port, ie. the local ip address and local port
        if (Object.keys(candidates).length === 1) {

            var ports = candidates[Object.keys(candidates)[0]];
            infoprint += "PORTS: " + syntaxHighlight(ports)
            // but if outside ports are different for every candidate, then NAT is symmetric
            infoprint += (ports.length === 1 ? '<h4 style="color:green;">NON-SYMMETRIC NAT</h4>' : '<br/><h4 style="color:red;">SYMMETRIC NAT</h4>');
        }
    }

    document.getElementById("status-bar-ip-addr-port").innerHTML = userIP + ":" + listenOnPort
    document.getElementById("log").innerHTML = "Checking NAT Type" + "<br/><br/>"
    document.getElementById("log").innerHTML += "Total Candidates:" + totalCandidates + "<br/>" +infoprint + "<br/>"
};