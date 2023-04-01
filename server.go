package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

type httpServeDir struct {
	URL  string
	Path string
}

// serverConf defines the server configuration file.
type serverConf struct {
	Host             string
	Port             int
	DocumentRootDir  string
	IndexFile        string
	BackupDir        string
	BackupFileFormat string
	ServeDirs        []httpServeDir
	LogFileName      string
	LogMaxSize       int // MB
	LogMaxBackups    int
	LogMaxAge        int // Days
	LogCompress      bool
}

// config holds the server configuration file values.
var config = serverConf{
	Host:             "",
	Port:             8080,
	DocumentRootDir:  "./www",
	IndexFile:        "index.html",
	BackupDir:        "./backup",
	BackupFileFormat: ":name:.:date:.html",
	ServeDirs:        []httpServeDir{},
	LogFileName:      "./logs/twserver.log",
	LogMaxSize:       4,
	LogMaxBackups:    16,
	LogMaxAge:        28,
	LogCompress:      true,
}

// initLog initializes the logging system.
func initLog() error {
	var logDir = filepath.Dir(config.LogFileName)

	// Creates the log directory if it doesn't exist
	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			fmt.Printf("[ERROR]initLog: Unable to create the log directory: %v\n", err)
			return err
		}
	}

	log.SetOutput(&lumberjack.Logger{
		Filename:   config.LogFileName,
		MaxSize:    config.LogMaxSize,
		MaxBackups: config.LogMaxBackups,
		MaxAge:     config.LogMaxAge,
		Compress:   config.LogCompress,
	})
	return nil
}

// readConfig reads the given configuration file.
func readConfig(filename string) error {
	var cfgFile, err = os.Open(filename)

	if err != nil {
		return err
	}

	var decoder = json.NewDecoder(cfgFile)

	if err := decoder.Decode(&config); err != nil {
		cfgFile.Close()
		return err
	}
	return cfgFile.Close()
}

// makeBackup makes a unique compressed backup of the index file.
func makeBackup(buf *[]byte) {
	// Replaces :name: with the index file without extension and :date: with the current timestamp
	var fileNameStrReplacer = strings.NewReplacer(
		":name:", strings.TrimSuffix(filepath.Base(config.IndexFile), filepath.Ext(config.IndexFile)),
		":date:", strconv.FormatInt(time.Now().Unix(), 10),
	)

	// Defines the index, backup and compressed backup file names
	var indexFilename = filepath.Join(config.DocumentRootDir, config.IndexFile)
	var backupFilename = filepath.Join(config.BackupDir, fileNameStrReplacer.Replace(config.BackupFileFormat))
	var zipFilename = (backupFilename + ".zip")

	// Writes the buf content to the backup file
	if err := writeFile(backupFilename, buf); err != nil {
		log.Println(fmt.Sprintf("[ERROR]makeBackup: Unable to write the file: %v", err))
		return
	}

	// Replaces the index file with the newly created backup file
	if err := copyFile(backupFilename, indexFilename); err != nil {
		log.Println(fmt.Sprintf("[ERROR]makeBackup: Unable to copy the file: %v", err))
		return
	}

	// Compress the backup file
	if err := zipFile(backupFilename, zipFilename); err != nil {
		log.Println(fmt.Sprintf("[ERROR]makeBackup: Unable to compress the file: %v", err))
		return
	}

	// Removes the original backup file
	if err := deleteFile(backupFilename); err != nil {
		log.Println(fmt.Sprintf("[ERROR]makeBackup: Unable to delete the file: %v", err))
	}
}

// httpHandleReq handles any HTTP request on the document root directory.
func httpHandleReq(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Methods", "HEAD, OPTIONS, GET, PUT")
	(w).Header().Set("dav", "'tw5/put")

	if r.URL.Path != "/" {
		http.Error(w, fmt.Sprintf("Invalid URL: %s", r.URL.Path), http.StatusNotFound)
		log.Println(fmt.Sprintf("[WARN]httpHandleReq: Requested resource '%s' not found", r.URL.Path))
		return
	}

	switch r.Method {
	case "HEAD", "OPTIONS":
		break
	case "GET":
		http.ServeFile(w, r, filepath.Join(config.DocumentRootDir, config.IndexFile))
	case "PUT":
		var bodyBytes, err = ioutil.ReadAll(r.Body)

		if err != nil {
			log.Println(fmt.Sprintf("[ERROR]httpHandleReq: Unable to read the request body content: %v", err))
			break
		}

		if err := r.Body.Close(); err != nil {
			log.Println(fmt.Sprintf("[ERROR]httpHandleReq: Unable to close the request body stream: %v", err))
			break
		}

		go makeBackup(&bodyBytes)
	default:
		log.Println("[WARN]httpHandleReq: Unsupported HTTP request method. Only GET/PUT methods are supported.")
	}
}

// init initializes the app by initializing the log system, reading the configuration file and creating the backup folder.
func init() {
	var binPath, err = os.Executable()

	if err != nil {
		fmt.Printf("[ERROR]init: Unable to get the app binary path: %v\n", err)
		return
	}

	var binName = strings.TrimSuffix(filepath.Base(binPath), filepath.Ext(binPath))
	var binDir = filepath.Dir(binPath)
	var confFilename = filepath.Join(binDir, (binName + ".json"))

	// Reads the config file if present
	if _, err := os.Stat(confFilename); err == nil {
		if err := readConfig(confFilename); err != nil {
			fmt.Printf("[ERROR]init: Unable to read the configuration file: %v\n", err)
		}
	}

	// Inits the logging system
	if err := initLog(); err != nil {
		fmt.Printf("[ERROR]init: Unable to initialize the logs: %v\n", err)
	}

	// Creates the backup directory if it doesn't exist
	if _, err := os.Stat(config.BackupDir); os.IsNotExist(err) {
		if err := os.MkdirAll(config.BackupDir, os.ModePerm); err != nil {
			log.Fatal(fmt.Sprintf("[ERROR]init: Unable to create the backup directory: %v", err))
		}
	}
}

// main starts the HTTP server.
func main() {
	var addr = fmt.Sprintf("%s:%s", config.Host, strconv.Itoa(config.Port))
	var router = http.NewServeMux()
	var srv = &http.Server{
		Addr:    addr,
		Handler: router,
	}

	router.HandleFunc("/", httpHandleReq)

	// Serves user-defined directories
	if len(config.ServeDirs) > 0 {
		for _, sDir := range config.ServeDirs {
			router.Handle(sDir.URL, http.StripPrefix(sDir.URL, http.FileServer(http.Dir(sDir.Path))))
		}
	}

	var sigint = make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)

	go func() {
		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("[ERROR]main: Unable to start the HTTP Server: %v", err)
		}
	}()

	log.Println("[INFO]main: HTTP Server started")
	fmt.Printf("TW HTTP Server Listening on %s\n", addr)

	<-sigint

	log.Println("[INFO]main: HTTP Server stopped")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("[ERROR]main: HTTP Server shutdown failed: %v", err)
	}
}
