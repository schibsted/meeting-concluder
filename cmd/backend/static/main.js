// Custom JavaScript code here
const devicesSelect = document.getElementById("devices");
const recordBtn = document.getElementById("recordBtn");
const stopBtn = document.getElementById("stopBtn");
const resultDiv = document.getElementById("result");
const postToSlackBtn = document.getElementById("postToSlack");

let conclusion = "";

async function fetchDevices() {
    const response = await fetch("/devices");
    const devices = await response.json();
    devices.forEach(device => {
        const option = document.createElement("option");
        option.value = device.index;
        option.textContent = device.name;
        devicesSelect.appendChild(option);
    });
}

async function selectDevice(index) {
    const response = await fetch(`/select-device/${index}`, { method: "POST" });
    return response.ok;
}

async function startRecording() {
    const response = await fetch("/record", { method: "POST" });
    return response.ok;
}

async function stopRecording() {
    const response = await fetch("/stop", { method: "POST" });
    const data = await response.json();
    return data;
}

async function postToSlack(conclusion) {
    // Implement this function to post the conclusion to Slack
    console.log("Posting to Slack:", conclusion);
}

devicesSelect.addEventListener("change", () => {
    selectDevice(devicesSelect.value);
});

recordBtn.addEventListener("click", async () => {
    const success = await startRecording();
    if (success) {
        recordBtn.disabled = true;
        stopBtn.disabled = false;
    }
});

stopBtn.addEventListener("click", async () => {
    const data = await stopRecording();
    if (data.message) {
        recordBtn.disabled = false;
        stopBtn.disabled = true;
        conclusion = data.message;
        resultDiv.textContent = `Conclusion: ${conclusion}`;
        postToSlackBtn.disabled = false;
    }
});

postToSlackBtn.addEventListener("click", () => {
    postToSlack(conclusion);
});

fetchDevices();
