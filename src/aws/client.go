package aws

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
)

// ServerStatusResponse EC2起動後のステータス確認レスポンス
type ServerStatusResponse struct {
	PublicIP string `json:"publicip"`
}

// StartResponse EC2起動指示時のレスポンス
type StartResponse struct {
	StartingInstances []InstanceStatus `json:"StartingInstances"`
}

// StopResponse EC2停止指示時のレスポンス
type StopResponse struct {
	StoppingInstances []InstanceStatus `json:"StoppingInstances"`
}

// InstanceStatus EC2指示時の共通レスポンス
type InstanceStatus struct {
	InstanceID   string `json:"InstanceId"`
	CurrentState struct {
		Code int    `json:"Code"`
		Name string `json:"Name"`
	} `json:"CurrentState"`
	PreviousState struct {
		Code int    `json:"Code"`
		Name string `json:"Name"`
	} `json:"PreviousState"`
}

// GetIPAddress インスタンスのIPアドレス取得
func GetIPAddress() (string, error) {
	statusOutputJSON, err := exec.Command("aws", "ec2", "describe-instances", "--instance-ids", os.Getenv("INSTANCE_ID"), "--query", "Reservations[].Instances[].{publicip:PublicIpAddress}").Output()
	if err != nil {
		return "", WrapError(err, ErrFailedGetIpAddress)
	}

	ssResponse := []ServerStatusResponse{}
	if err := json.Unmarshal(statusOutputJSON, &ssResponse); err != nil {
		return "", WrapError(err, ErrInvalidResponseGetIpAddress)
	}

	ipaddress := ssResponse[0].PublicIP
	if ipaddress != "" {
		log.Println("IPアドレス : ", ipaddress)
	}

	return ipaddress, nil
}

func StartInstance() error {
	outputJSON, err := exec.Command("aws", "ec2", "start-instances", "--instance-ids", os.Getenv("INSTANCE_ID")).Output()
	if err != nil {
		return WrapError(err, ErrFailedStartInstance)
	}

	startResponse := StartResponse{}
	if err := json.Unmarshal(outputJSON, &startResponse); err != nil {
		return WrapError(err, ErrInvalidResponseStartInstance)
	}

	currentState := startResponse.StartingInstances[0].CurrentState.Name
	if currentState == "running" {
		return WrapError(nil, ErrInstanceAlreadyStarted)
	}

	previousState := startResponse.StartingInstances[0].PreviousState.Name
	if currentState == "pending" && previousState == "pending" {
		return WrapError(nil, ErrStartingInstance)
	}

	// 開始待ち
	if _, err := exec.Command("aws", "ec2", "wait", "instance-running", "--instance-ids", os.Getenv("INSTANCE_ID")).Output(); err != nil {
		return WrapError(err, ErrFailedWaitStartInstance)
	}

	return nil
}

func StopInstance() error {
	outputJSON, err := exec.Command("aws", "ec2", "stop-instances", "--instance-ids", os.Getenv("INSTANCE_ID")).Output()
	if err != nil {
		return WrapError(err, ErrFailedStopInstance)
	}

	stopResponse := StopResponse{}
	if err := json.Unmarshal(outputJSON, &stopResponse); err != nil {
		return WrapError(err, ErrInvalidResponseStopInstance)
	}

	currentState := stopResponse.StoppingInstances[0].CurrentState.Name
	if currentState == "stopped" {
		return WrapError(nil, ErrInstanceAlreadyStopped)
	}

	previousState := stopResponse.StoppingInstances[0].PreviousState.Name
	if currentState == "stopping" && previousState == "stopping" {
		return WrapError(nil, ErrStoppingInstance)
	}

	// 停止待ち
	if _, err := exec.Command("aws", "ec2", "wait", "instance-stopped", "--instance-ids", os.Getenv("INSTANCE_ID")).Output(); err != nil {
		return WrapError(err, ErrFailedWaitStopInstance)
	}

	return nil
}
