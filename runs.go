package openaiassistantsgo

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// Runs take a Thread id which could contain Messages and runs it against an Assistant

func CreateRun(assistant_id string, thread_id string) (string, error) {

	var run_id string
	client := &http.Client{}

	var data = strings.NewReader(fmt.Sprintf(`{
    "assistant_id": "%s"
  }`, assistant_id))

	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.openai.com/v1/threads/%s/runs", thread_id), data)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var b interface{}
	if err := json.Unmarshal(bodyText, &b); err != nil {
		log.Fatal(err)
	}

	m := b.(map[string]interface{})

	for k, v := range m {
		if k == "id" {
			run_id = v.(string)
			fmt.Println(v)
			return run_id, nil
		}
	}

	fmt.Printf("%s\n", bodyText)
	return "Run ID not found.", nil

}

func GetRunStatus(run_id string, thread_id string) (string, error) {

	client := &http.Client{}
	var run_status string

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.openai.com/v1/threads/%s/runs/%s", thread_id, run_id), nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))
	req.Header.Set("OpenAI-Beta", "assistants=v2")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var b interface{}
	if err := json.Unmarshal(bodyText, &b); err != nil {
		log.Fatal(err)
	}

	m := b.(map[string]interface{})

	fmt.Printf("%s\n", bodyText)

	for k, v := range m {
		if k == "status" {
			run_status = v.(string)
			fmt.Println(v)
			return run_status, nil
		}
	}

	return "Run Status not found.", nil
}
