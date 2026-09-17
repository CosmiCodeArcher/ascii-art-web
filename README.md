# Description

**This is a project to display the previously rendered ascii-art project to a web browser instead of a terminal.**

**ASCII-ART-WEB** allows users input text and choose a banner type through a web form and see the result rendered in their browser, then download that result as a `.txt` file.

## Authors
Hamza Musa

## Usage: how to run

Go to the project's home directory in your terminal and run `go run .`

You'll see a message that looks like `2026/06/12 10:51:03 Server started on http://localhost:8081`

Go to that link `http://localhost:8081` in your browser and type in your desired text (***Note: Use \n for newline***) and see the ascii-art representation printed out.

***Example: Hello\nWorld***. Go ahead and test it.

--

## Endpoints

| Method | Route | Parameters | Description |
| --- | --- | --- | --- |
| `GET` | `/` | none | Serves the form page where the text and banner are chosen. |
| `POST` | `/ascii-art` | form fields `text`, `banner` | Renders the ascii-art from the submitted form data and writes it back into the results page. |
| `GET` | `/export` | query parameters `text`, `banner` | Regenerates the same ascii-art and sends it back as a downloadable `.txt` file. |
| `GET` | `/static/` | none | Serves files from the `static` directory (currently the stylesheet). Registered with `http.Handle` using `http.FileServer` wrapped in `http.StripPrefix`, so `/static/style.css` maps to `static/style.css` on disk. |

`text` is the text to render and `banner` is the banner's name: `standard`, `shadow` or `thinkertoy`.

### Downloading the result

Once a result has been rendered, the results page shows a **Download** link. That link points at
`/export?text=...&banner=...`, carrying the text and banner name that were just used, so clicking it
downloads exactly the art that is on screen.

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

The `/export` route reuses that exact same generation step; it only reads the text and banner's name from the query string instead of the form, so the downloaded art always matches what the page showed.

### Template Rendering

The string is passed into the parsed HTML `.Result` placeholder through a struct then sent back to the browser using the response writer.

### Styling

The CSS lives in `static/style.css` and is linked from the `<head>` of `templates/index.html`. The browser fetches it through the `/static/` route above.

- The colour palette (`--bg`, `--text`, `--accent-color`, `--border-color`, `--error`) is defined once as CSS custom properties in `:root` and referenced with `var()` throughout, so each palette colour is declared in a single place.
- Every interactive control (the text input, the banner select, the submit button and the download link) has `:hover`, `:active` and `:focus-visible` states. `:focus-visible` is used rather than `:focus` so the keyboard focus ring appears during tab navigation but not on mouse clicks.
- The `<pre>` output element uses `overflow-x: auto` so wide ascii-art scrolls inside its own box instead of forcing the whole document sideways. Wrapping is deliberately not used, because a wrapped row would break the alignment of the art.
- The layout is a single column of inline form controls inside a `max-width: 1200px` body, which reflows naturally at narrow widths without media queries.
- The page requests `JetBrains Mono` and falls back to the generic `monospace` family, so the art stays correctly aligned on any machine but the typeface differs where JetBrains Mono is not installed.

### Export

No file is ever written to the server's disk. The generated string is streamed straight to the client
through the response writer, and the response carries three headers:

- `Content-Type: text/plain` — tells the browser it is plain text
- `Content-Length` — the byte length of the generated art
- `Content-Disposition: attachment; filename=ascii-art.txt` — asks the browser to save the response instead of displaying it

The browser then creates the file itself from the filename in `Content-Disposition`.

--

## Behaviour notes

- `homeHandler` returns `404 Not Found` for any path other than `/`, because the `"/"` pattern in Go's default `ServeMux` is a prefix match rather than an exact one and would otherwise serve the home page for every unmatched URL.
- `/export` is an endpoint, not a page. Visiting it without the query parameters returns `400 Bad Request` rather than a form, because there is nothing to render.
- The banner's name is checked against a fixed whitelist (`standard`, `shadow`, `thinkertoy`) on both `/ascii-art` and `/export`. The name is only turned into a `banners/<name>.txt` path after it matches one of those three, so a crafted banner name cannot be used to walk the file system through the banner file path.
