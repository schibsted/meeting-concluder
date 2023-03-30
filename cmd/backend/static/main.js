document.getElementById("recordButton").addEventListener("click", async () => {
  const recordButton = document.getElementById("recordButton");
  const stopButton = document.getElementById("stopButton");
  const errorElement = document.getElementById("error");
  const summaryElement = document.getElementById("summary");

  recordButton.classList.add("hidden");
  stopButton.classList.remove("hidden");
  errorElement.classList.add("hidden");
  summaryElement.textContent = "";

  const clapDetection = window.ClapDetection;
  const maxDuration = window.MaxDuration;

  try {
    const response = await fetch("/record/start", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ ClapDetection: clapDetection, MaxDuration: maxDuration }),
    });

    if (!response.ok) {
      throw new Error(`${response.status} - ${response.statusText}`);
    }

    const result = await response.json();
    summaryElement.textContent = `Recorded for ${result.duration.toFixed(2)} seconds.`;
  } catch (error) {
    console.error(error);
    errorElement.textContent = `Error: ${error.message}`;
    errorElement.classList.remove("hidden");
  }
});

document.getElementById("stopButton").addEventListener("click", async () => {
  const recordButton = document.getElementById("recordButton");
  const stopButton = document.getElementById("stopButton");
  const errorElement = document.getElementById("error");
  const summaryElement = document.getElementById("summary");

  recordButton.classList.remove("hidden");
  stopButton.classList.add("hidden");
  errorElement.classList.add("hidden");

  try {
    const response = await fetch("/record/stop", {
      method: "POST",
    });

    if (!response.ok) {
      throw new Error(`${response.status} - ${response.statusText}`);
    }

    const result = await response.json();
    summaryElement.textContent += ` Transcription: ${result.transcription}`;
  } catch (error) {
    console.error(error);
    errorElement.textContent = `Error: ${error.message}`;
    errorElement.classList.remove("hidden");
  }
});
