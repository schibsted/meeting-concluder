document.getElementById('record-btn').addEventListener('click', async () => {
  const slackChannel = document.getElementById('slack-channel').value;
  const slackAPIKey = document.getElementById('slack-api-key').value;
  const openAIAPIKey = document.getElementById('openai-api-key').value;

  const recordRequest = {
    slack_channel: slackChannel,
    slack_api_key: slackAPIKey,
    openai_api_key: openAIAPIKey,
  };

  try {
    const response = await fetch('/record', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(recordRequest),
    });

    if (response.ok) {
      const jsonResponse = await response.json();
      const summary = jsonResponse.summary;
      alert(`Summary: ${summary}`);
    } else {
      alert(`Error: ${response.status} - ${response.statusText}`);
    }
  } catch (error) {
    alert(`Error: ${error.message}`);
  }
});
