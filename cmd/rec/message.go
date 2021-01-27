package main

import "fmt"

type message struct {
	IP        string `json:"ip"`
	Name      string `json:"name"`
	Class     string `json:"class"`
	Timestamp int64  `json:"ts"`
}

func (m message) String() string {
	return fmt.Sprintf("%12d, %5s, %-16s, %s", m.Timestamp, m.Class, m.IP, m.Name)
}

func messageHeader() string { return "   timestamp, class, ip              , name" }
