var service_name = document.getElementById("service_name").value;

function sendToBackend() {
    fetch("https://bot-dev-domain.com/service", {
        method: 'POST',
        headers: {
            'Accept' : 'application/json',
            'Content-Type' : 'application/json'
        },
        body : JSON.stringify({"name" : service_name})
    }).then(response => response.json())
      .then(response => console.log(JSON.stringify(response)))
}