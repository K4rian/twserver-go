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
- __[Download the latest build](/releases)__ for your platform.
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
    "BackupFileFormat": ":name:.:date:.html"
}
```

- Tweak the values as needed, save the file and restart the server.



## Building

Building is done with the `go` tool. If you have setup your `GOPATH` correctly, the following should work:

```bash
go get github.com/k4rian/twserver-go
go build -ldflags "-w -s" github.com/k4rian/twserver-go
```



## License

[MIT](LICENSE)