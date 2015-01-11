package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type StaticServer struct {
	// real os path to directory
	BaseDir string
	// the base html file --  gets served with no path
	Index string
	// server port
	Port int
	// Cache filename
	Files []string
	// force no cache
	NoCache bool
	// Error Message
	Error string
}

// new
func NewStaticServer(basePath, index string, port int, cache bool, errMsg string) *StaticServer {
	return &StaticServer{
		BaseDir: basePath,
		Index:   index,
		Port:    port,
		Files:   []string{},
		NoCache: cache,
		Error:   errMsg,
	}
}

/**
 * returns filename for raw path
 * @param path req path
 * @return path with in context
 */
func (s *StaticServer) resolvePath(path string) string {
	if path == "/" {
		return filepath.Join(s.BaseDir, s.Index)
	} else {
		return filepath.Join(s.BaseDir, path)
	}
}

// checkFile
// saves the filename in the `Files` slice property
func (s *StaticServer) checkFile(filename string) bool {
	for _, f := range s.Files {
		if f == filename {
			return true
		}
	}
	return false
}

// push a file into the served file array
func (s *StaticServer) addFile(filename string) {
	s.Files = append(s.Files, filename)
}

// get Mod time
// try to get the modification time
// so that ServeContent
// can send a 304 if the file has been cached
// if s.Cache is true
// Note: it will get logged as a 200 regardless
// Check for err,
// if global no-caching option is true,
// and if the file has been served on this instance yet
func (s *StaticServer) getModTime(filename string) time.Time {
	if stats, err := os.Stat(filename); err != nil || (!s.checkFile(filename) || s.NoCache) {
		return time.Now()
	} else {
		return stats.ModTime()
	}
}

func (s *StaticServer) sendErrorMsg(w http.ResponseWriter) error {
	_, err := w.Write([]byte(s.Error))
	return err
}

// whats the status?
// @return status code, file
// @note: the status code 600 signifies that there is no `index.html`
// so we send the error message
func (s *StaticServer) getStatus(req *http.Request, filename string) (int, *os.File) {
	var status int

	if req.Method != "GET" {
		return 405, nil
	}

	file, err := os.Open(filename)

	if err == nil {
		return 200, file
	}

	if os.IsNotExist(err) {
		if "/" == req.URL.Path {
			status = 600
		} else {
			status = 404
		}
	} else {
		status = 500
	}

	return status, nil
}

// handle the response
func (s *StaticServer) serveFile(w http.ResponseWriter, req *http.Request, filename, path string) {
	status, file := s.getStatus(req, filename)

	switch status {
	case 200:
		modTime := s.getModTime(filename)
		http.ServeContent(w, req, filename, modTime, file)
		// add the file
		s.addFile(filename)
		file.Close()
	case 600:
		if err := s.sendErrorMsg(w); err != nil {
			http.Error(w, http.StatusText(500), 500)
		}
		// don't log a status of 600
		status = 404
	case 405, 404, 500:
		fmt.Println("GOT HERE!")
		http.Error(w, http.StatusText(status), status)
	}
	// log the request
	log(req.Method, req.URL.Path, filename, status)
}

func (s *StaticServer) Serve() {
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		filename := s.resolvePath(req.URL.Path)
		s.serveFile(w, req, filename, req.URL.Path)
	})
	fmt.Printf("Listening on port \x1b[35;2m%d\x1b[0m\n", s.Port)
	if err := http.ListenAndServe(":"+strconv.Itoa(s.Port), nil); err != nil {
		panic(err)
	}
}

// helper
func log(method, path, filename string, status int) {
	fmt.Printf("\x1b[32;1m%s\x1b[0m \x1b[30;1m%s\x1b[0m \x1b[31;4m%s\x1b[0m \x1b[36;1m%d\x1b[0m \x1b[30;1m%s\x1b[0m\n", method, path, http.StatusText(status), status, filename)
}
