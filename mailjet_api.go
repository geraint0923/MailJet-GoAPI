package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type MailJet struct {
	username string
	password string
}

type Mail struct {
	Locale      string
	Sender      string
	SenderEmail string
	Subject     string
	ListId      int
}

func InitMailJet(username string, password string) *MailJet {
	mail := new(MailJet)
	mail.username = username
	mail.password = password
	return mail
}

func parseJSON(data string) map[string]interface{} {
	jsonData := []byte(data)
	//	fmt.Println(data)
	u := map[string]interface{}{}
	err := json.Unmarshal(jsonData, &u)
	if err != nil {
		fmt.Println("heihei")
		panic(err)
	}
	return u
}

func (m *MailJet) BuildGroup(name string) int {
	var jsonStr = []byte("{\"name\":" + name + "}")
	req, _ := http.NewRequest("POST", "https://api.mailjet.com/v3/REST/contactslist", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(m.username, m.password)
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//	fmt.Println(string(body))
	jsonObj := parseJSON(string(body))
	if jsonObj["Data"] == nil {
		return -1
	}
	retVal := jsonObj["Data"].([]interface{})[0].(map[string]interface{})["ID"].(float64)
	return int(retVal)
}

func (m *MailJet) addContact(addr string) int {
	var jsonStr = []byte("{\"Email\":" + addr + "}")
	req, _ := http.NewRequest("POST", "https://api.mailjet.com/v3/REST/contact", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(m.username, m.password)
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	jsonObj := parseJSON(string(body))
	if jsonObj["Data"] == nil {
		return -1
	}
	retVal := jsonObj["Data"].([]interface{})[0].(map[string]interface{})["ID"].(float64)
	return int(retVal)
}

func (m *MailJet) getContactID(addr string) int {
	req, _ := http.NewRequest("GET", "https://api.mailjet.com/v3/REST/contact/"+addr, nil)
	req.SetBasicAuth(m.username, m.password)
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	jsonObj := parseJSON(string(body))
	if jsonObj["Data"] == nil {
		return -1
	}
	retVal := jsonObj["Data"].([]interface{})[0].(map[string]interface{})["ID"].(float64)
	return int(retVal)
}

func (m *MailJet) AddToGroup(listId int, addr string) {
	res := m.addContact(addr)
	if res < 0 {
		return
	}
	res = getContactID(addr)
	if res < 0 {
		return
	}
	var jsonStr = []byte("{\"ContactID\":" + strconv.Itoa(res) + ",\"ListID\":" + strconv.Itoa(listId) + ",\"IsActive\": true}")
	req, _ := http.NewRequest("POST", "https://api.mailjet.com/v3/REST/listrecipient", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(m.username, m.password)
	client := &http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body))
	/*
		jsonObj := parseJSON(string(body))
		if jsonObj["Data"] == nil {
			return
		}
	*/
}

func (m *MailJet) SendToGroup(mail *Mail) {
}

func (m *MailJet) SendToUser(mail *Mail) {
}

func main() {
	mail := InitMailJet("3bfc81a965fbc505da2950df6e2a5a4b", "6a26367d687abc1c4d543483d1870365")
	res := mail.getContactID("0923@gmail.com")
	fmt.Println("Result: " + strconv.Itoa(res))
}
