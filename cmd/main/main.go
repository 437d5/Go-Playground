package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type Request struct {
	Code string `json:"code"`
}

type Response struct {
	Output string `json:"output,omitempty"`
	Error  string `json:"error,omitempty"`
}

func main() {
	s := &http.Server{
		Addr: ":5050",
	}

	http.HandleFunc("/run", runHandler)
	http.HandleFunc("/frmt", frmtHandler)

	log.Fatal(s.ListenAndServe())
}

func runHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("JSON DECODE ERROR:", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	res, err := processCode(req.Code)
	if err != nil {
		log.Println("PROCESS ERROR:", err)
	}

	response := Response{
		Output: res,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("JSON ENCODE ERROR:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func processCode(code string) (string, error) {
	ts := time.Now().Unix()
	tmpDir := filepath.Join("tmp", fmt.Sprintf("%d", ts))
	err := os.MkdirAll(tmpDir, 0755)
	if err != nil {
		return "", fmt.Errorf("error creating folder: %v", err)
	}

	filePath := filepath.Join(tmpDir, "main.go")
	err = os.WriteFile(filePath, []byte(code), 0644)
	if err != nil {
		return "", fmt.Errorf("error writing to file: %v", err)
	}

	cmd := exec.Command("go", "run", filePath)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func frmtHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Println("JSON DECODE ERROR:", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	log.Println("Received code:", req.Code)

	res, err := frmtCode(req.Code)
	if err != nil {
		log.Println("FORMAT ERROR:", err)
		http.Error(w, "Failed to process code", http.StatusInternalServerError)
	}

	response := Response{
		Output: res,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("JSON ENCODE ERROR:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func frmtCode(code string) (string, error) {
	ts := time.Now().Unix()
	tmpDir := filepath.Join("tmp", fmt.Sprintf("%d", ts))
	err := os.MkdirAll(tmpDir, 0755)
	if err != nil {
		return "", fmt.Errorf("error creating folder: %v", err)
	}

	filePath := filepath.Join(tmpDir, "main.go")
	err = os.WriteFile(filePath, []byte(code), 0644)
	if err != nil {
		return "", fmt.Errorf("error writing to file: %v", err)
	}

	cmd := exec.Command("go", "fmt", filePath)
	_, err = cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error formatting file: %v", err) 
	}

	res, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed read file after fmt: %v", err)
	}

	return string(res), err
}
