document.addEventListener("DOMContentLoaded", function () {
    // Custom JavaScript code here
    const devicesSelect = document.getElementById("devices");
    const recordBtn = document.getElementById("recordBtn");
    const stopBtn = document.getElementById("stopBtn");
    const resultDiv = document.getElementById("result");
    const postToSlackBtn = document.getElementById("postToSlack");
    const statusDiv = document.getElementById("status");
    const recordingIndicator = document.getElementById("recordingIndicator");

    // Update the status message
    function updateStatus(message) {
        statusDiv.textContent = message;
    }

    // Show or hide the recording indicator
    function setRecordingIndicator(visible) {
        recordingIndicator.style.display = visible ? "block" : "none";
    }

    let conclusion = "";

    async function fetchDevices() {
        const response = await fetch("/devices");
        const devices = await response.json();
        devices.forEach((device) => {
            const option = document.createElement("option");
            option.value = device.index;
            option.textContent = device.name;
            devicesSelect.appendChild(option);
        });
    }

    async function selectDevice(index) {
        const response = await fetch(`/select-device/${index}`, {
            method: "POST",
        });
        return response.ok;
    }

    async function startRecording() {
        const response = await fetch("/record", { method: "POST" });
        return response.ok;
    }

    async function stopRecording() {
        document.getElementById("status").innerText = "Stopping recording...";
        try {
            const response = await fetch("/stop", {
                method: "POST",
            });
            const data = await response.json();
            if (response.ok) {
                document.getElementById("status").innerText = data.message;
            } else {
                document.getElementById("status").innerText =
                    "Error: " + data.error;
            }
            return data;
        } catch (error) {
            console.error("Error stopping recording:", error);
            document.getElementById("status").innerText =
                "Error: Failed to stop recording";
        }
    }

    async function getConclusion() {
        try {
            const response = await fetch("/conclusion");
            if (response.ok) {
                const data = await response.json();
                document.getElementById("conclusion").innerText =
                    data.conclusion;
            } else {
                document.getElementById("conclusion").innerText =
                    "Error: Conclusion not available";
            }
        } catch (error) {
            console.error("Error fetching conclusion:", error);
            document.getElementById("conclusion").innerText =
                "Error: Conclusion not available";
        }
    }

    async function postConclusionToSlack() {
        try {
            const response = await fetch("/post-to-slack", {
                method: "POST",
            });

            if (response.ok) {
                updateStatus("Conclusion posted to Slack");
            } else {
                const data = await response.json();
                updateStatus("Error posting to Slack: " + data.error);
            }
        } catch (error) {
            console.error("Error posting conclusion to Slack:", error);
            updateStatus("Error posting conclusion to Slack");
        }
    }

    // Update the status message and recording indicator accordingly
    devicesSelect.addEventListener("change", async () => {
        if (await selectDevice(devicesSelect.value)) {
            updateStatus("Device selected successfully");
        } else {
            updateStatus("Error selecting device");
        }
    });

    recordBtn.addEventListener("click", async () => {
        if (await startRecording()) {
            updateStatus("Recording started");
            setRecordingIndicator(true);
            recordBtn.disabled = true;
            stopBtn.disabled = false;
        } else {
            updateStatus("Error starting recording");
        }
    });

    postToSlackBtn.addEventListener("click", async () => {
        updateStatus("Posting conclusion to Slack...");
        await postConclusionToSlack();
    });

    stopBtn.addEventListener("click", async () => {
        const data = await stopRecording();
        if (data.message) {
            updateStatus("Recording stopped by user");
            setRecordingIndicator(false);
            recordBtn.disabled = false;
            stopBtn.disabled = true;
            conclusion = data.conclusion;
            resultDiv.textContent = conclusion;
            postToSlackBtn.disabled = false;
        } else {
            updateStatus("Error stopping recording");
        }
    });

    // Initialize the recording indicator as hidden
    setRecordingIndicator(false);

    // Call fetchDevices at the end
    fetchDevices();
});
