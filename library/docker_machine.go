/*
Copyright (c) 2016, Fabrizio Soppelsa
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

* Neither the name of ansible-swarm nor the names of its
  contributors may be used to endorse or promote products derived from
  this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

=============================================================================

Ansible docker_machine module
Requires Ansible 2.2+

Sample snippet for playbooks:

    - name: Create a Machine
	  docker_machine:
	    provider: "generic"
	    name: "node0"
	    ip_address: "192.168.1.1"
	    ssh_key: "~/.conf/keys/nodes.pem"
	    ssh_user: "ubuntu"
	  register: machine_result
*/

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

// Json utilities
//

// Playbook variables
type ModuleArgs struct {
	Provider   string
	Name       string
	Ip_address string
	Ssh_user   string
	Ssh_key    string
}

// Format response
type Response struct {
	Msg        string `json:"msg"`
	ConnString string `json:"connstring"`
	Cmd        string `json:"command"`
	Changed    bool   `json:"changed"`
	Failed     bool   `json:"failed"`
}

func ExitJson(responseBody Response) {
	returnResponse(responseBody)
}

func FailJson(responseBody Response) {
	responseBody.Failed = true
	returnResponse(responseBody)
}

func returnResponse(responseBody Response) {
	var response []byte
	var err error
	response, err = json.Marshal(responseBody)
	if err != nil {
		response, _ = json.Marshal(Response{Msg: "Invalid response object"})
	}
	fmt.Println(string(response))
	if responseBody.Failed {
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

// Machine functions
//

// Machine def
type Machine struct {
	provider   string
	name       string
	ip_address string
	ssh_user   string
	ssh_key    string
}

// Create machine
func createMachine(params *Machine) {
	var response Response
	var command string

	if params.provider == "generic" {
		command = "docker-machine create -d generic --generic-ip-address " + params.ip_address + " --generic-ssh-key " + params.ssh_key + " --generic-ssh-user " + params.ssh_user + " " + params.name
	}

	cmd := exec.Command(command)
	err := cmd.Run()
	if err != nil {
		response.Msg = "Machine creation error"
		response.Cmd = command
		FailJson(response)
	}
}

func main() {
	var response Response

	if len(os.Args) != 2 {
		response.Msg = "Not enough arguments provided!"
		FailJson(response)
	}

	argsFile := os.Args[1]

	text, err := ioutil.ReadFile(argsFile)
	if err != nil {
		response.Msg = "Could not read configuration file: " + argsFile
		FailJson(response)
	}

	var moduleArgs ModuleArgs
	err = json.Unmarshal(text, &moduleArgs)

	if err != nil {
		response.Msg = "Configuration file not valid JSON: " + argsFile
		FailJson(response)
	}

	var provider string = "generic"
	var name string = ""
	var ip_address = ""
	var ssh_user = ""
	var ssh_key = ""

	if moduleArgs.Provider != "" {
		provider = moduleArgs.Provider
	}
	if moduleArgs.Name != "" {
		name = moduleArgs.Name
	}
	if moduleArgs.Ip_address != "" {
		ip_address = moduleArgs.Ip_address
	}
	if moduleArgs.Ssh_user != "" {
		ssh_user = moduleArgs.Ssh_user
	}
	if moduleArgs.Ssh_key != "" {
		ssh_key = moduleArgs.Ssh_key
	}

	machine := Machine{provider, name, ip_address, ssh_user, ssh_key}
	createMachine(&machine)

	response.ConnString = "docker-machine create -d " + provider + " --generic-ip-address " + ip_address + " --generic--ssh-key " + ssh_key + " --generic-user " + ssh_user + " " + name
	ExitJson(response)
}
