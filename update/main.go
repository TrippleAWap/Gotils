package update

import (
	"crypto/sha256"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

// UpdateToURL updates the application to the latest version from the given URL.
// the hashes of the downloaded file and the current version are compared to determine if an update is needed.
// If an update is needed, the application will spawn a new process to perform the update.
// The update process will replace the current process with the new version.
func UpdateToURL(URL string) error {
	req, err := http.NewRequest("GET", URL, nil)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fileHandle, err := os.Open(os.Args[0])
	defer fileHandle.Close()
	currentBytes, err := io.ReadAll(fileHandle)
	if err != nil {
		return err
	}
	bodyHash := sha256.Sum256(bodyBytes)
	currentHash := sha256.Sum256(currentBytes)
	if currentHash != bodyHash {
		pathSplit := strings.Split(os.Args[0], "\\")
		log.Printf("%x != %x\n", currentHash, bodyHash)
		tempPath := "new_" + pathSplit[len(pathSplit)-1] + ".tmp"
		err := os.WriteFile(tempPath, bodyBytes, 0644)
		if err != nil {
			return err
		}
		log.Printf("Update available. Spawning new process to perform update.")
		time.Sleep(1 * time.Second)
		nullFile, err := os.Open(os.DevNull)
		if err != nil {
			panic(err)
		}
		defer nullFile.Close()
		attr := &os.ProcAttr{
			Dir:   os.Getenv("CWD"),
			Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
		}
		_, err = os.StartProcess("C:\\Windows\\system32\\cmd.exe", []string{"/c", "timeout", "/t", "1", "/nobreak", "&&", "move", tempPath, os.Args[0], "&&", "cls", "&&", "start", os.Args[0]}, attr)
		if err != nil {
			return err
		}
		os.Exit(0)
		return nil
	}
	return nil
}
