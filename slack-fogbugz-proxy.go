package main

// Written by Kim Blomqvist <kim.blomqvist@lekane.com>
// Modified by Adam Krone <adam.krone@thirdwavellc.com>

/*
Copyright (c) 2014 Lekane Oy. All rights reserved.
Copyright (c) 2016 Adam Krone. All rights reserved.
Copyright (c) 2016 Thirdwave, LLC. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Lekane Oy nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

var config Config

type Message struct {
	Text     string `json:"text"`
	Channel  string `json:"channel"`
	Username string `json:"username"`
	IconUrl  string `json:"icon_url"`
}

type Config struct {
	Webhook         string                 `json:"webhook"`
	Port            string                 `json:"port"`
	FogbugzUrl      string                 `json:"fogbugz_url"`
	SlackUser       string                 `json:"slack_user"`
	DefaultChannel  string                 `json:"default_channel"`
	ChannelMappings map[string]interface{} `json:"channel_mappings"`
}

func findChannel(project_name string) string {
	channel := config.ChannelMappings[project_name]
	if channel == nil {
		return config.DefaultChannel
	} else {
		return channel.(string)
	}
}

func prepareMessageText(queryParams *url.Values) string {
	message := "<" + config.FogbugzUrl + "/default.asp?" + queryParams.Get("case_number")
	message = message + "|Case " + queryParams.Get("case_number") + ">: "
	message = message + queryParams.Get("title")
	return message
}

func prepareMessage(queryParams *url.Values) Message {
	return Message{
		Text:     prepareMessageText(queryParams),
		Channel:  findChannel(queryParams.Get("project_name")),
		Username: config.SlackUser,
		IconUrl:  "http://www.fogcreek.com/images/fogbugz/pricing/kiwi.png",
	}
}

func post(queryParams *url.Values) {
	message := prepareMessage(queryParams)

	binaryMessage, err := json.Marshal(message)
	if err != nil {
		fmt.Printf("ERROR '%s' marshalling message: %s\n", err, message)
		return
	}

	fmt.Printf("Posting: %s\n", binaryMessage)
	req, err := http.NewRequest("POST", config.Webhook, bytes.NewReader(binaryMessage))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("ERROR '%s' sending message\n", err)
	}

	if resp != nil {
		defer resp.Body.Close()
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	post(&queryParams)
}

func loadConfig(config *Config) {
	file, _ := ioutil.ReadFile("./config.json")

	json.Unmarshal(file, &config)
}

func loadDefaults(config *Config) {
	if config.Webhook == "" {
		fmt.Printf("Missing webhook url. You must include this in your config.json.")
		os.Exit(1)
	}

	if config.Port == "" {
		config.Port = "10333"
	}

	if config.SlackUser == "" {
		config.SlackUser = "fogbugz"
	}

	if config.DefaultChannel == "" {
		fmt.Printf("Missing default channel. You must include this in your config.json.")
		os.Exit(1)
	}
}

func setupServer() {
	fmt.Printf("Listening to port: %s\n", config.Port)
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(":"+config.Port, nil)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
}

func main() {
	loadConfig(&config)
	loadDefaults(&config)
	setupServer()
}
