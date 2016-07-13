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
	"github.com/fsoppelsa/ansible"
)

/*
  Sample:

    - name: Operate swarm nodes
      docker_swarm:
        role: "master"|"slave"
        operation: "init"|"join"|"update"|"rm"|"promote"|"demote"
        detach: yes
        docker_url: "192.168.1.1"
        use_tls: encrypt
        tls_ca_cert: "/path/ca.pem"
        tls_client_cert: "/path/cert.pem"
        tls_client_key: "/path/key.pem"
    register: swarm_result

*/

type ModuleArgs struct {
	Role          string
	Operation     string
	Detach        bool
	DockerUrl     string
	UseTls        string
	TlsCaCert     string
	TlsClientCert string
	TlsClientKey  string
}

type Swarm struct {
}

type SwarmNode struct {
}

func main() {
}
