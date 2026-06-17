# Description

**This is a project to display the previously rendered ascii-art project to a web browser instead of a terminal.**

**ASCII-ART-WEB** allows users input text and choose a banner type through a web form and see the result rendered in their browser.

## Authors
Hamza Musa

## Usage: how to run

Go to the project's home directory in your terminal and run `go run .`

You'll see a message that looks like `2026/06/12 10:51:03 Server started on http://localhost:8080`

Go to that link `http://localhost:8080` in your browser and type in your desired text (***Note: Use \n for newline***) and see the ascii-art representation printed out.

***Example: Hello\nWorld***. Go ahead and test it.

--

## Implementation details: algorithm

### Server

The server was built using the `net/http` package in Go to register the port, create routes and assign handlers to those routes.

The handlers verify request method types, parses the html template then write responses back to the browser.


### ASCII Art Generation

When a POST request is received from the browser;

1. The text and banner's name are gotten from that request

2. The banner's name and text are passed through `Loadbanner()` and `ParseInput()`

3. Their results are passed to the render which processes it and generates the desired ascii-art string.

### Template Rendering

The string is passed into the parsed HTML `.Result` placeholder through a struct then sent back to the browser using the response writer.