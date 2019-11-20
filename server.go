package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// ServerConf defines the server configuration file
type ServerConf struct {
	Port             int
	DocumentRootDir  string
	IndexFile        string
	BackupDir        string
	BackupFileFormat string
}

// AppConf holds the server configuration file values
var AppConf = ServerConf{
	Port:             8080,
	DocumentRootDir:  "./www",
	IndexFile:        "index.html",
	BackupDir:        "./backup",
	BackupFileFormat: ":name:.:date:.html",
}

func initLog(filename *string) {
	logFile, err := os.OpenFile(*filename, os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf(fmt.Sprintf("[ERROR]initLog -> Unable to open/create the log file: '%s'.", *filename), err)
		return
	}
	log.SetOutput(logFile)
}

func readConfig(filename *string) error {
	cfgFile, err := os.Open(*filename)
	if err != nil {
		log.Fatal(fmt.Sprintf("[ERROR]readConfig -> Unable to read the config file: '%s'.", *filename), err)
		return err
	}

	decoder := json.NewDecoder(cfgFile)
	err = decoder.Decode(&AppConf)

	if err != nil {
		cfgFile.Close()
		log.Fatal("[ERROR]readConfig -> Unable to parse the config file to JSON format.", err)
		return err
	}
	return cfgFile.Close()
}

func makeBackup(buf *[]byte) {
	indexHTMLReplacer := strings.NewReplacer(".html", "", ".htm", "")
	fileNameStrReplacer := strings.NewReplacer(
		":name:", indexHTMLReplacer.Replace(AppConf.IndexFile),
		":date:", strconv.FormatInt(time.Now().Unix(), 10),
	)

	indexFilePath := filepath.Join(AppConf.DocumentRootDir, AppConf.IndexFile)
	backupFilePath := filepath.Join(
		AppConf.BackupDir,
		fileNameStrReplacer.Replace(AppConf.BackupFileFormat),
	)
	zipFilePath := backupFilePath + ".zip"

	if writeFile(&backupFilePath, buf) == nil {
		copyFile(&backupFilePath, &indexFilePath)

		if zipFile(&backupFilePath, &zipFilePath) == nil {
			delFile(&backupFilePath)
		}
	}
}

func httpHandleReq(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Methods", "HEAD, OPTIONS, GET, PUT")
	(w).Header().Set("dav", "'tw5/put")

	if r.URL.Path != "/" {
		http.Error(w, "Invalid URL.", http.StatusNotFound)
		log.Print(fmt.Sprintf("[WARN]httpHandleReq -> The page URL '%s' doesn't exists.", r.URL.Path))
		return
	}

	switch r.Method {
	case "HEAD":
	case "OPTIONS":
		break
	case "GET":
		http.ServeFile(w, r, filepath.Join(AppConf.DocumentRootDir, AppConf.IndexFile))
	case "PUT":
		bodyBytes, err := ioutil.ReadAll(r.Body)

		if err != nil {
			log.Print("[ERROR]httpHandleReq -> Unable to read the request body content.", err)
			break
		}

		if err := r.Body.Close(); err != nil {
			log.Print("[ERROR]main -> Unable to close the request body stream.", err)
			break
		}

		go makeBackup(&bodyBytes)
	default:
		log.Print("[WARN]httpHandleReq -> Unsupported HTTP request method. Only GET/PUT methods are supported.")
	}
}

func init() {
	binPath, err := os.Executable()
	if err != nil {
		fmt.Printf("[ERROR]init -> Unable to get the binary's path.")
		return
	}

	binName := filepath.Base(binPath)
	binDirPath := filepath.Dir(binPath)
	logFilePath := filepath.Join(binDirPath, (strings.Replace(binName, ".exe", "", 1) + ".log"))
	confFilePath := filepath.Join(binDirPath, (strings.Replace(binName, ".exe", "", 1) + ".json"))

	// Inits the logs
	initLog(&logFilePath)

	// Reads the config file if present
	if _, err := os.Stat(confFilePath); err == nil {
		readConfig(&confFilePath)
	}

	// Creates the backup directory if it doesn't exist
	if _, err := os.Stat(AppConf.BackupDir); os.IsNotExist(err) {
		err := os.MkdirAll(AppConf.BackupDir, os.ModePerm)
		if err != nil {
			log.Print(fmt.Sprintf("[ERROR]init -> Unable to create the backup directory: '%s'.", AppConf.BackupDir), err)
		}
	}
}

func main() {
	http.HandleFunc("/", httpHandleReq)

	bindStr := fmt.Sprintf(":%s", strconv.Itoa(AppConf.Port))
	fmt.Printf(fmt.Sprintf("TW HTTP Server Listening on %s\n", bindStr))

	if err := http.ListenAndServe(bindStr, nil); err != nil {
		log.Fatal("[ERROR]main -> Unable to start the TW HTTP Server: ", err)
	}
}
