package session

import (
	"encoding/json"
	"os"

	"golang.org/x/exp/slog"
)

type Session struct {
	Name string `json:"name"`
	Host string `json:"host"`
	Port string `json:"port"`
}

func ListSession() []Session {
	file, _ := os.OpenFile("session.json", os.O_RDWR|os.O_CREATE, 0666)
	defer file.Close()

	sessions := make([]Session, 0)

	decoder := json.NewDecoder(file)
	err := decoder.Decode(&sessions)
	if err != nil {
		slog.Warn("read session", err)
		return nil
	}

	return sessions
}

func AddSession(name, host, port string) {
	sessions := ListSession()

	session := Session{
		Name: name,
		Host: host,
		Port: port,
	}

	sessions = append(sessions, session)

	saveFile(sessions)
}

func UpdateSession(name, host, port string, id int) {
	sessions := ListSession()

	if id >= len(sessions) || id < 0 {
		slog.Warn("id is out of range")
		return
	}

	sessions[id].Name = name
	sessions[id].Host = host
	sessions[id].Port = port

	saveFile(sessions)
}

func DeleteSession(id int) {
	sessions := ListSession()

	if id >= len(sessions) || id < 0 {
		slog.Warn("id is out of range")
		return
	}

	sessions = append(sessions[:id], sessions[id+1:]...)

	saveFile(sessions)
}

func saveFile(sessions []Session) {
	file, _ := os.OpenFile("session.json", os.O_RDWR|os.O_CREATE, 0666)
	defer file.Close()

	encoder := json.NewEncoder(file)
	err := encoder.Encode(sessions)

	if err != nil {
		slog.Warn("write session", err)
	}
}
