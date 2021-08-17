package scanner

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
)

var EICAR = []byte(`X5O!P%@AP[4\PZX54(P^)7CC)7}$EICAR-STANDARD-ANTIVIRUS-TEST-FILE!$H+H*`)
var version = "master"

type ScanInterceptor interface {
	Scanner() Scanner
	ScanFiles(files []*multipart.FileHeader) *ScanResult
	ScanUrls(urls []string) *ScanResult
	ScanForVirus(filename string, reader io.Reader) FileScanResult
	Info() *Info
}
type scanInterceptor struct {
	scanner Scanner
}

func NewScanInterceptor(scannerObj Scanner, url string) ScanInterceptor {
	scannerObj.SetAddress(url)
	return &scanInterceptor{scanner: scannerObj}
}

type Info struct {
	Version             string `json:"file_uploader_version"`
	Address             string `json:"scan_server_url"`
	PingResult          string `json:"ping_result"`
	ScannerVersion      string `json:"scan_server_version"`
	TestScanVirusResult string `json:"test_scan_virus"`
	TestScanCleanResult string `json:"test_scan_clean"`
}

type FileScanResult struct {
	*Result
	Error    bool   `json:"error"`
	FileName string `json:"file_name"`
	Message  string `json:"message"`
}
type ScanResult struct {
	Success bool             `json:"success"`
	Files   []FileScanResult `json:"files"`
	Message string           `json:"message"`
}

func (i *scanInterceptor) Scanner() Scanner {
	return i.scanner
}
func (i *scanInterceptor) ScanUrls(urls []string) *ScanResult {
	scanResults := &ScanResult{Files: []FileScanResult{}, Success: true}
	for _, u := range urls {
		response, err := http.Get(u)
		if err != nil {
			scanResults.Files = append(scanResults.Files, i.prepareErrorResponse(u, err))
			continue
		}
		scanResults.Files = append(scanResults.Files, i.ScanForVirus(u, response.Body))
		_ = response.Body.Close()
	}
	i.prepareResponse(scanResults)
	return scanResults
}

func (i *scanInterceptor) ScanFiles(files []*multipart.FileHeader) *ScanResult {
	scanResults := &ScanResult{Files: []FileScanResult{}, Success: true}
	for _, httpFile := range files {
		f, err := httpFile.Open()
		if err != nil {
			scanResults.Files = append(scanResults.Files, i.prepareErrorResponse(httpFile.Filename, err))
			continue
		}
		scanResults.Files = append(scanResults.Files, i.ScanForVirus(httpFile.Filename, f))
		_ = f.Close()
	}
	i.prepareResponse(scanResults)
	return scanResults
}

func (i *scanInterceptor) ScanForVirus(filename string, reader io.Reader) FileScanResult {
	fileScanResult := FileScanResult{FileName: filename, Message: ""}
	result, err := i.scanner.Scan(reader)
	if err != nil {
		fileScanResult.Error = true
		fileScanResult.Message = err.Error()
		return fileScanResult
	}
	fileScanResult.Result = result
	return fileScanResult
}

func (i *scanInterceptor) Info() *Info {
	info := &Info{
		Address: i.scanner.Address(),
		Version: version,
	}
	if err := i.scanner.Ping(); err != nil {
		info.PingResult = err.Error()
	} else {
		info.PingResult = "Connected to server OK"
		if response, err := i.scanner.Version(); err != nil {
			info.ScannerVersion = err.Error()
		} else {
			info.ScannerVersion = response
		}
		reader := bytes.NewReader(EICAR)
		if result, err := i.scanner.Scan(reader); err != nil {
			info.TestScanVirusResult = err.Error()
		} else {
			info.TestScanVirusResult = result.String()
		}
		// Validate the Clamd response for a non-viral string
		reader = bytes.NewReader([]byte("foo bar mcgrew"))
		if result, err := i.scanner.Scan(reader); err != nil {
			info.TestScanCleanResult = err.Error()
		} else {
			info.TestScanCleanResult = result.String()
		}
	}
	return info
}
func (i *scanInterceptor) prepareErrorResponse(fileName string, err error) FileScanResult {
	return FileScanResult{
		FileName: fileName,
		Error:    true,
		Message:  err.Error(),
		Result: &Result{
			Status:      "",
			Virus:       false,
			Description: "",
		},
	}
}
func (i *scanInterceptor) prepareResponse(scanResults *ScanResult) {
	for _, file := range scanResults.Files {
		if file.Error || file.Virus == true {
			scanResults.Success = false
			break
		}
	}
}
