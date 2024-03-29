twserver-go
=====

A [TiddlyWiki][1] HTTP Server written in Go.

## Features
- Automatic backup on each save.
- Customizable settings.
- Easy to deploy.
- Works on macOS, Linux and Windows.

## Getting started
### Prerequisites
- [An empty copy of TiddlyWiki][2].

### Setup
- __[Download the latest build][3]__ for your platform.
- Extract the archive.
- `cd` into the directory that was just created.
- Put the TiddlyWiki's `empty.html` file inside the `www` subfolder.
- Rename the `empty.html` file to `index.html`.
- Run the server: `./twserver`.
- Open your web browser and browse to: [http://localhost:8080](http://localhost:8080)

### Customizing
All server settings can be tweaked using a configuration file located beside the server binary.
The configuration file must use the same name as the binary and saved with the `.json` extension.

- Create a configuration file in `JSON` format (on Linux/macOS):
```bash
touch twserver.json
```

- Open the file with your favorite text editor and write the following values *(default)*:
```json
{
  "Host": "",
  "Port": 8080,
  "DocumentRootDir": "./www",
  "IndexFile": "index.html",
  "BackupDir": "./backup",
  "BackupFileFormat": ":name:.:date:.html",
  "ServeDirs": [],
  "LogFileName": "./logs/twserver.log",
  "LogMaxSize": 4,
  "LogMaxBackups": 16,
  "LogMaxAge": 28,
  "LogCompress": true
}
```
- Tweak the values as needed, save the file and restart the server.

By default, the HTTP server only serves the index file and rejects any other request. To serve one or more custom directories containing extra resources (such as images), you have to add them by tweaking the `ServeDirs` value in the configuration file.

- __Example__: Add the `images` directory located in `./images` and accessed via the URL `<wiki_url>/img/`:
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
- Any image from the `./images` directory can now be reached via `<wiki_url>/img/` and displayed inside any wiki post.

## Building
Building is done with the `go` tool. If you have setup your `GOPATH` correctly, the following should work:
```bash
go get github.com/k4rian/twserver-go
go build -ldflags "-w -s" github.com/k4rian/twserver-go
```

## Dependencies
The rotating logging system is powered by [lumberjack][4]:
```
gopkg.in/natefinch/lumberjack.v2="v2.2.1"
```

## Docker Image
A Docker image is available on Docker Hub under [k4rian/twserver][5] and its corresponding [source repository][6] on GitHub.

## License
[MIT][7]

[1]: https://github.com/Jermolene/TiddlyWiki5 "TiddlyWiki 5 GitHub repository"
[2]: https://tiddlywiki.com/#GettingStarted "TiddlyWiki: Getting Started"
[3]: https://github.com/k4rian/twserver-go/releases "Latest twserver-go releases"
[4]: https://github.com/natefinch/lumberjack "Lumberjack GitHub repository"
[5]: https://hub.docker.com/r/k4rian/twserver "twserver-go Docker Image"
[6]: https://github.com/K4rian/docker-twserver "twserver-go Docker Image GitHub repository"
[7]: https://github.com/K4rian/twserver-go/blob/master/LICENSE