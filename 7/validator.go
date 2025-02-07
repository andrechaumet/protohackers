package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var types = map[string]bool{
	"connect": true,
	"data":    true,
	"ack":     true,
	"close":   true,
}

var validations = []func([]string) error{
	validPrefixSuffix,
	validStructure,
	validMessageType,
	validSessionID,
	validThirdField,
	validDataField,
}

func Validate(data string) error {
	data = data[1 : len(data)-1]
	parts := strings.Split(data, "/")
	for _, validate := range validations {
		if err := validate(parts); err != nil {
			return err
		}
	}
	return nil
}

func validPrefixSuffix(parts []string) error {
	msg := parts[0]
	if !strings.HasPrefix(msg, "/") || !strings.HasSuffix(msg, "/") {
		return errors.New("the message must start and end with '/'\n")
	}
	return nil
}

func validStructure(parts []string) error {
	if len(parts) < 3 {
		return errors.New("the message must have at least 3 parts")
	}
	return nil
}

func validMessageType(parts []string) error {
	if _, valid := types[parts[0]]; !valid {
		return fmt.Errorf("invalid message type")
	}
	return nil
}

func validSessionID(parts []string) error {
	if _, err := strconv.Atoi(parts[1]); err != nil {
		return fmt.Errorf("invalid session id")
	}
	return nil
}

func validThirdField(parts []string) error {
	if parts[0] == "data" || parts[0] == "ack" {
		if _, err := strconv.Atoi(parts[2]); err != nil {
			return fmt.Errorf("invalid pos or length")
		}
	}
	return nil
}

func validDataField(parts []string) error {
	if parts[0] == "data" && len(parts) > 3 && parts[3] == "" {
		return errors.New("the DATA field cannot be empty for 'data' message type")
	}
	return nil
}
