package main

import (
	"html/template"
	"log"
	"net/http"
)

const htmlTemplate = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Meeting Concluder</title>
    <style>
        {{.CSS}}
    </style>
</head>
<body>
    <h1>Meeting Concluder</h1>
    <p>
        <button id="start">Start Recording</button>
        <button id="stop" disabled>Stop Recording</button>
    </p>
    <script>
        {{.JS}}
    </script>
</body>
</html>
`

const css = `
body {
    font-family: Arial, sans-serif;
    margin: 40px;
}

button {
    padding: 10px 20px;
    font-size: 16px;
    cursor: pointer;
}
`

const js = `
document.getElementById("start").addEventListener("click", async () => {
    try {
        const response = await fetch("/start", { method: "POST" });
        if (!response.ok) {
            throw new Error("Failed to start recording");
        }
        document.getElementById("start").disabled = true;
        document.getElementById("stop").disabled = false;
    } catch (err) {
        alert(err.message);
    }
});

document.getElementById("stop").addEventListener("click", async () => {
    try {
        const response = await fetch("/stop", { method: "POST" });
        if (!response.ok) {
            throw new Error("Failed to stop recording");
        }
        document.getElementById("start").disabled = false;
        document.getElementById("stop").disabled = true;
    } catch (err) {
        alert(err.message);
    }
});
`

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.New("index").Parse(htmlTemplate)
		if err != nil {
			log.Fatal(err)
		}

		err = tmpl.Execute(w, struct {
			CSS string
			JS  template.JS
		}{
			CSS: css,
			JS:  template.JS(js),
		})
		if err != nil {
			log.Fatal(err)
		}
	})

	addr := ":8080"
	log.Printf("Starting server on %s...\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
