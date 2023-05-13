async function sendToBackend() {
    var name = document.getElementById("input_data").value;
    const response = await fetch("https://bot-dev-domain.com/service", {
        method: 'POST',
        headers: {
            'Accept' : 'application/json',
            'Content-Type' : 'application/json'
        },
        body : JSON.stringify({"name" : name})
    })
    const jsonData = response.json()
    console.log(JSON.stringify(jsonData))
}