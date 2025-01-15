Here’s the correctly formatted README.md for your project:

Page Analysis API

A simple web API that accepts a URL, fetches its HTML content, and submits it to an Ollama server for analysis—returning the response to the client.

Table of Contents

1.	Overview
2.	Prerequisites
3.	Installation
4.	Running
5.	Usage
6.	Notes

Overview

This API demonstrates how to:
1.	Accept a POST request with a JSON payload containing a URL.
2.	Fetch the HTML content from the provided URL.
3.	Send that content to an Ollama server for text analysis.
4.	Return the analysis response to the client.

Prerequisites
1. Go (1.18+ recommended)
2. Git
3. An Ollama server running locally or remotely. (Installation instructions at ollama.ai)

Installation

1.	Clone this repository:

```
git clone https://github.com/aodr3w/page-analysis-api
```

2.	Change to the project directory:

```
cd page-analysis-api
```

Running

Start the application:

```
go run main.go
```

This will launch the server on port 3000 by default (e.g., http://localhost:3000).

Usage

After the server starts, you can make a POST request to /find:

```
curl -X POST http://localhost:3000/find \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://www.readpoetry.com/10-vivid-haikus-to-leave-you-breathless/"
  }'

```

Example Response:

```
{
  "data": "This is a HTML code snippet for a webpage. It includes various JavaScript files and references to different plugins and features used on the website, such as the EasyPieChart plugin, jQuery, WordPress Simple Firewall, Page Links To, and more. The page also includes a Not Bot Shield protection feature to prevent automated attacks."
}

```


Note: The output will vary based on the LLM in use and the content of the provided page.

Notes

•	You can configure which LLM (model) Ollama uses by modifying the Go client code (llm.NewClient()).

•	If you want to extract only certain parts of the HTML (e.g., <body>), you can modify the fetching logic in fetchHTML or use an HTML parser for more robust parsing.

•	The above example shows a simplified prompt and response. Actual usage can involve more complex instructions or formatting for the LLM.

