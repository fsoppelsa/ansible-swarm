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

Ansible docker_swarm module
Requires Ansible 2.2+
*/

package main

import (
	"encoding/json"
	"github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types/swarm"
	"github.com/docker/go-connections/tlsconfig"
	"github.com/fsoppelsa/ansible"
	"golang.org/x/net/context"
	"net/http"
	"os"
	"path/filepath"
)

/*
  Sample:

    - name: Operate swarm clusters and nodes
      docker_swarm:
        role: "manager"|"slave"
        operation: "init"|"join"|"leave"|"update"|"promote"|"demote"
        docker_url: "tcp://192.168.99.101:2376"
        join_url: ["tcp://192.168.99.100:3376"] # array of strings
        tls_path: "/path/to/"
      register: swarm_result

*/

type ModuleArgs struct {
	Role       string
	Operation  string
	Docker_url string
	Join_url   []string
	Tls_path   string
}

func connectEngine(args *ModuleArgs) *client.Client {
	var c *http.Client
	var response ansible.Response

	securityOptions := tlsconfig.Options{
		CAFile:             filepath.Join(args.Tls_path, "ca.pem"),
		CertFile:           filepath.Join(args.Tls_path, "cert.pem"),
		KeyFile:            filepath.Join(args.Tls_path, "key.pem"),
		InsecureSkipVerify: true,
	}

	tlsc, err := tlsconfig.Client(securityOptions)
	if err != nil {
		response.Msg = "TLS configuration error: " + args.Tls_path
		ansible.FailJson(response)
	}

	c = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsc,
		},
	}

	defaultHeaders := map[string]string{"User-Agent": "engine-api-cli-1.0"}
	// v1.24 = engine 1.12
	cli, err := client.NewClient(args.Docker_url, "v1.24", c, defaultHeaders)
	if err != nil {
		response.Msg = "Impossible to connect to Docker Engine: " + args.Docker_url
		ansible.FailJson(response)
	}

	return cli
}

func initSwarm(cli *client.Client) string {
	var response ansible.Response

	swarm, err := cli.SwarmInit(context.Background(), swarm.InitRequest{
		ListenAddr:      "0.0.0.0:2377",
		ForceNewCluster: true,
		Spec: {
			AcceptancePolicy: {
				Policies: []Policy{
					Policy{
						Role:       "manager",
						Autoaccept: true,
					},
					Policy{
						Role:       "slave",
						Autoaccept: true,
					},
				},
			},
		},
	})
	if err != nil {
		response.Msg = "ERROR: Swarm init: " + swarm
		ansible.FailJson(response)
	}
	return "ok init " + swarm
}

func joinSwarm(cli *client.Client, addr []string, role string) string {
	var response ansible.Response
	var isManager bool

	if role == "manager" {
		isManager = true
	} else {
		isManager = false
	}

	err := cli.SwarmJoin(context.Background(), swarm.JoinRequest{
		ListenAddr:  "0.0.0.0:2377",
		RemoteAddrs: addr,
		Manager:     isManager,
	})
	if err != nil {
		response.Msg = "ERROR: Swarm join: " + role
		ansible.FailJson(response)
	}

	return "ok join"
}

func leaveSwarm(cli *client.Client, force bool) string {
	var response ansible.Response

	err := cli.SwarmLeave(context.Background(), force)
	if err != nil {
		response.Msg = "ERROR: Swarm leave"
		ansible.FailJson(response)
	}

	return "ok leave"
}

func promoteNode(cli *client.Client) string {
	return "ok promote"
}

func demoteNode(cli *client.Client) string {
	return "ok demote"
}

func updateSwarm(cli *client.Client) string {
	return "ok update"
}

func main() {
	var response ansible.Response
	var moduleArgs ModuleArgs

	text := ansible.ParseVariables(os.Args)

	if err := json.Unmarshal(text, &moduleArgs); err != nil {
		response.Msg = "Configuration file not valid JSON: " + os.Args[1]
		ansible.FailJson(response)
	}

	cli := connectEngine(&moduleArgs)

	// Operations implementation
	// Init cluster
	if moduleArgs.Role == "manager" && moduleArgs.Operation == "init" {
		swarm := initSwarm(cli)
		response.Msg = swarm
		ansible.ExitJson(response)
	}

	// Join a node
	if moduleArgs.Operation == "join" {
		swarm := joinSwarm(cli, moduleArgs.Join_url, moduleArgs.Role)
		response.Msg = swarm
		ansible.ExitJson(response)
	}

	// A node leaves
	if moduleArgs.Operation == "leave" {
		swarm := leaveSwarm(cli, true)
		response.Msg = swarm
		ansible.ExitJson(response)
	}
}
