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
*/

package main

import (
	"encoding/json"
	"github.com/fsoppelsa/ansible"
	"os"
	"os/exec"
)

/*
  Json utilities

  Playbook variables
  Sample:

  - name: Create a Machine
    docker_machine:
      provider: "generic"
      name: "node0"
      ip_address: "192.168.1.1"
      ssh_key: "~/.conf/keys/nodes.pem"
      ssh_user: "ubuntu"
    register: machine_result
*/
type ModuleArgs struct {
	Provider   string
	Name       string
	Ip_address string
	Ssh_user   string
	Ssh_key    string
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
	var response ansible.Response
	var cmdName = "docker-machine"
	var cmdArgs string

	if params.provider == "generic" {
		cmdArgs = "create -d generic --generic-ip-address " + params.ip_address + " --generic-ssh-key " + params.ssh_key + " --generic-ssh-user " + params.ssh_user + " " + params.name
	}

	// TODO: Check that Machine does not exist already
	if err := exec.Command("sh", "-c", cmdName+" "+cmdArgs).Run(); err != nil {
		response.Msg = "Machine creation error"
		response.Cmd = cmdName + " " + cmdArgs
		ansible.FailJson(response)
	}
}

func main() {
	var response ansible.Response
	var moduleArgs ModuleArgs
	var provider string = "generic"
	var name string = ""
	var ip_address = ""
	var ssh_user = ""
	var ssh_key = ""

	text := ansible.ParseVariables(os.Args)

	if err := json.Unmarshal(text, &moduleArgs); err != nil {
		response.Msg = "Configuration file not valid JSON: " + os.Args[1]
		ansible.FailJson(response)
	}

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
	ansible.ExitJson(response)
}
