twserver-go
=====

A [TiddlyWiki](https://github.com/Jermolene/TiddlyWiki5) HTTP Server written in Go.



## Features

- Automatic backup on each save.
- No external dependencies.
- Customizable settings.
- Easy to deploy.
- Works on macOS, Linux and Windows.



## Getting started

### Prerequisites

- [TiddlyWiki empty HTML file](https://tiddlywiki.com/#GettingStarted).


### Setup

- __[Download the latest build](https://github.com/k4rian/twserver-go/releases)__ for your platform.
- Extract the archive.
- `cd` into the directory that was just created.
- Put the TeddlyWiki's `empty.html` file inside the `www` subfolder.
- Rename the `empty.html` file to `index.html`.
- Run the server: `./twserver`.
- Open your web browser and browse to: [http://localhost:8080](http://localhost:8080).


### Customizing

All server settings can be tweaked using a configuration file located beside the server binary.
The configuration file must use the same name as the binary and must be saved with the `.json` extension.

- Create a configuration file in `JSON` format (on Linux/macOS):

```bash
touch twserver.json
```

- Open the file with your favorite text editor and write the following values *(default)*:

```json
{
  "Port": 8080,
  "DocumentRootDir": "./www",
  "IndexFile": "index.html",
  "BackupDir": "./backup",
  "BackupFileFormat": ":name:.:date:.html",
  "ServeDirs": []
}
```

- Tweak the values as needed, save the file and restart the server.

By default, the HTTP server only serves the index file and rejects any other request. To serve one or more custom directories containing extra resources (such as images), you have to add them by tweaking the `ServeDirs` value in the configuration file.

- Example: Add the `images` directory located in `./images` and accessed via the URL `<wiki_url>/img/`:

```json
{
  "ServeDirs": [
    {
      "URL": "/img/",
      "Path": "./images"
    }
  ]
}
```

- Any image from the `./images` directory can now be reached via `<wiki_url>/img/` and can be displayed inside any wiki post.



## Building

Building is done with the `go` tool. If you have setup your `GOPATH` correctly, the following should work:

```bash
go get github.com/k4rian/twserver-go
go build -ldflags "-w -s" github.com/k4rian/twserver-go
```



## Docker Image

A Docker image is available on Docker Hub under [k4rian/twserver](https://hub.docker.com/r/k4rian/twserver) and its corresponding [source repository](https://github.com/K4rian/docker-twserver) on GitHub.



## License

[MIT](LICENSE)