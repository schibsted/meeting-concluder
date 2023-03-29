const body = document.body;
const slackChannelEnv = body.getAttribute("data-slack-channel");
const slackAPIKeyEnv = body.getAttribute("data-slack-api-key");
const openAIAPIKeyEnv = body.getAttribute("data-openai-api-key");

document.getElementById("record-btn").addEventListener("click", async () => {
    const resultMessage = document.getElementById("result-message").value;
    const slackChannel =
        document.getElementById("slack-channel").value || slackChannelEnv;
    const slackAPIKey =
        document.getElementById("slack-api-key").value || slackAPIKeyEnv;
    const openAIAPIKey =
        document.getElementById("openai-api-key").value || openAIAPIKeyEnv;

    const recordRequest = {
        slack_channel: slackChannel,
        slack_api_key: slackAPIKey,
        openai_api_key: openAIAPIKey,
    };

    try {
        const recordBtn = document.getElementById("record-btn");
        recordBtn.classList.add("bg-red-500", "animate-pulse");
        recordBtn.classList.remove("bg-green-500");
        recordBtn.textContent = "Recording...";

        const response = await fetch("/record", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(recordRequest),
        });

        const resultMessage = document.getElementById("result-message");
        if (response.ok) {
            resultMessage.textContent = "Summary sent to Slack.";
            resultMessage.classList.remove("text-red-500");
            resultMessage.classList.add("text-green-500");
        } else {
            resultMessage.textContent = `Error: ${response.status} - ${response.statusText}`;
            resultMessage.classList.remove("text-green-500");
            resultMessage.classList.add("text-red-500");
        }
    } catch (error) {
        resultMessage.textContent = `Error: ${error.message}`;
        resultMessage.classList.remove("text-green-500");
        resultMessage.classList.add("text-red-500");
    } finally {
        const recordBtn = document.getElementById("record-btn");
        recordBtn.classList.add("bg-green-500");
        recordBtn.classList.remove("bg-red-500", "animate-pulse");
        recordBtn.textContent = "Record";
    }
});
