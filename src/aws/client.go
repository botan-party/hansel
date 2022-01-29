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
func GetIPAddress() (string, StatusError) {
	statusOutputJSON, err := exec.Command("aws", "ec2", "describe-instances", "--instance-ids", os.Getenv("INSTANCE_ID"), "--query", "Reservations[].Instances[].{publicip:PublicIpAddress}").Output()
	if err != nil {
		return "", NewStatusError(ERR_FAILED_GET_IP_ADDRESS, err)
	}

	ssResponse := []ServerStatusResponse{}
	if err := json.Unmarshal(statusOutputJSON, &ssResponse); err != nil {
		return "", NewStatusError(ERR_INVALID_RESPONSE_GET_IP_ADDRESS, err)
	}

	ipaddress := ssResponse[0].PublicIP
	if ipaddress != "" {
		log.Println("IPアドレス : ", ipaddress)
	}

	return ipaddress, StatusError{}
}

func StartInstance() StatusError {
	outputJSON, err := exec.Command("aws", "ec2", "start-instances", "--instance-ids", os.Getenv("INSTANCE_ID")).Output()
	if err != nil {
		return NewStatusError(ERR_FAILED_START_INSTANCE, err)
	}

	startResponse := StartResponse{}
	if err := json.Unmarshal(outputJSON, &startResponse); err != nil {
		return NewStatusError(ERR_INVALID_RESPONSE_START_INSTANCE, err)
	}

	currentState := startResponse.StartingInstances[0].CurrentState.Name
	if currentState == "running" {
		return NewStatusError(ERR_INSTANCE_ALREADY_STARTED, nil)
	}

	previousState := startResponse.StartingInstances[0].PreviousState.Name
	if currentState == "pending" && previousState == "pending" {
		return NewStatusError(ERR_STARTING_INSTANCE, nil)
	}

	// 開始待ち
	if _, err := exec.Command("aws", "ec2", "wait", "instance-running", "--instance-ids", os.Getenv("INSTANCE_ID")).Output(); err != nil {
		return NewStatusError(ERR_FAILED_WAIT_START_INSTANCE, err)
	}

	return StatusError{}
}

func StopInstance() StatusError {
	outputJSON, err := exec.Command("aws", "ec2", "stop-instances", "--instance-ids", os.Getenv("INSTANCE_ID")).Output()
	if err != nil {
		return NewStatusError(ERR_FAILED_STOP_INSTANCE, err)
	}

	stopResponse := StopResponse{}
	if err := json.Unmarshal(outputJSON, &stopResponse); err != nil {
		return NewStatusError(ERR_INVALID_RESPONSE_STOP_INSTANCE, err)
	}

	currentState := stopResponse.StoppingInstances[0].CurrentState.Name
	if currentState == "stopped" {
		return NewStatusError(ERR_INSTANCE_ALREADY_STOPPED, nil)
	}

	previousState := stopResponse.StoppingInstances[0].PreviousState.Name
	if currentState == "stopping" && previousState == "stopping" {
		return NewStatusError(ERR_STOPPING_INSTANCE, nil)
	}

	// 停止待ち
	if _, err := exec.Command("aws", "ec2", "wait", "instance-stopped", "--instance-ids", os.Getenv("INSTANCE_ID")).Output(); err != nil {
		return NewStatusError(ERR_FAILED_WAIT_STOP_INSTANCE, err)
	}

	return StatusError{}
}
