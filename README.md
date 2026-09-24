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

## Docker

The project ships with a `Dockerfile` so the server can be built and run as a container instead of with `go run .`.

### Requirements

Docker installed and the daemon running. On Linux your user must be in the `docker` group, otherwise every command below needs `sudo`.

### Build and run

Build the image and check that it is there:

```
docker build -t aaw-app:1.0.0 .
docker images
```

Run it in the background and check that the container is up:

```
docker run -d -p 8081:8081 --name aaw-container --label project=ascii-art-web --rm aaw-app:1.0.0
docker ps
```

The app is then at `http://localhost:8081`, exactly as it is when run with `go run .`.

### Verify

Read the labels off the image, read them off the running container, and confirm which user the server runs as:

```
docker inspect --format '{{json .Config.Labels}}' aaw-app:1.0.0
docker inspect --format '{{json .Config.Labels}}' aaw-container
docker exec aaw-container whoami
```

### Clean up

Stop the container, confirm nothing of this project is left lying around, then drop the dangling images and the build cache:

```
docker stop aaw-container
docker ps -a --filter label=project=ascii-art-web
docker image prune -f --filter label=org.opencontainers.image.title=ascii-art-web
docker builder prune -f
```

### How the image is built

The `Dockerfile` is a two-stage build.

The **builder** stage compiles the project from the official `golang:1.22.2` image with `CGO_ENABLED=0`, so the binary is statically linked. That matters because the final stage runs on Alpine, which uses a different C library from the Debian-based Go image; a dynamically linked binary built there would not run on Alpine.

The **final** stage starts from `alpine:3.20` and receives only the compiled `aaw_app` binary plus `templates/`, `static/` and `banners/`. Those three directories have to be in the image because the server reads them at runtime through relative paths (`templates/index.html`, `http.Dir("static")`, `banners/<name>.txt`), which are resolved against `WORKDIR /app`. The Go source and the whole Go toolchain stay behind in the discarded builder stage, which keeps the final image small.

### Metadata

The image carries OCI labels set with `LABEL` in the final stage: `org.opencontainers.image.title`, `.description`, `.authors` and `.source`. A separate `project=ascii-art-web` label is attached to the container at run time with `--label`, which is what `docker ps -a --filter label=project=ascii-art-web` matches on. A container inherits the labels of its image and adds its own on top, so the container inspect above shows both the OCI labels and `project`.

### Security

The final stage creates an unprivileged user with `adduser -D not_archer` and switches to it with `USER not_archer` before the server starts, so the process does not run as root. `docker exec aaw-container whoami` prints `not_archer`.

### Garbage collection

- `--rm` on `docker run` removes the container as soon as it stops, so no stopped containers pile up.
- `docker image prune` with `--filter label=org.opencontainers.image.title=ascii-art-web` removes only this project's dangling images, leaving images from other projects alone.
- `docker builder prune` clears the build cache, which is where the layers of the discarded builder stage are held.

### `.dockerignore`

`.dockerignore` keeps things out of the build context that the image does not need, which makes the context smaller and avoids needless cache invalidation:

- `.git` and `.gitignore` — version control history, of no use inside the image.
- `README.md` — documentation, not needed at runtime.
- `Dockerfile` — the build instructions themselves.
- `ascii-art-web` and `aaw_app` — binaries compiled locally. The builder stage compiles its own inside the image, and a host binary could be built for a different platform.

--

## Behaviour notes

- `homeHandler` returns `404 Not Found` for any path other than `/`, because the `"/"` pattern in Go's default `ServeMux` is a prefix match rather than an exact one and would otherwise serve the home page for every unmatched URL.
- `/export` is an endpoint, not a page. Visiting it without the query parameters returns `400 Bad Request` rather than a form, because there is nothing to render.
- The banner's name is checked against a fixed whitelist (`standard`, `shadow`, `thinkertoy`) on both `/ascii-art` and `/export`. The name is only turned into a `banners/<name>.txt` path after it matches one of those three, so a crafted banner name cannot be used to walk the file system through the banner file path.
