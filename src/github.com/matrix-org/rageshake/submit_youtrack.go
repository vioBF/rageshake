package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

var youTrackProjectID = ""
var youTrackUrl = ""
var youTrackToken = ""

type Payload struct {
	Project     Project `json:"project"`
	Summary     string  `json:"summary"`
	Description string  `json:"description"`
}
type Project struct {
	ID string `json:"id"`
}

type YouTrackProjectInformation []struct {
	ShortName string `json:"shortName"`
	Name      string `json:"name"`
	ID        string `json:"id"`
	Type      string `json:"$type"`
}

//func youTrackMainProcess(project string, summary string, description string) {
//	if err := curlYouTrackProjectID(project); err != nil {
//		os.Exit(1)
//	}
//	createYouTrackIssue(summary, description)
//}

func curlYouTrackProjectID(project string) error {
	req, err := http.NewRequest("GET", youTrackUrl+"/api/admin/projects?fields=id,name,shortName&query="+project, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+youTrackToken)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	youtrackJson := YouTrackProjectInformation{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&youtrackJson); err != nil {
		return err
	}
	defer resp.Body.Close()

	youTrackProjectID = youtrackJson[0].ID

	return nil
}

func createYouTrackIssue(summary string, description string) error {

	data := Payload{
		Project:     Project{ID: youTrackProjectID},
		Summary:     summary,
		Description: description,
	}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", youTrackUrl+"/api/issues", body)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+youTrackToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
