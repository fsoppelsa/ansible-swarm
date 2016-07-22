ansible-swarm
=============
Tested on <a href="https://github.com/swarm2k/swarm2k">Swarm2k</a>, see in **swarm2k** dir.

In **library**, Ansible modules :
* `docker_machine.go` - <a href="https://github.com/docker/machine">machine</a>
* `docker_swarm.go` - <a href="https://github.com/docker/swarm">swarm (mode)</a>

In **playbooks/** some samples to instantiate on AWS

Software requirements
* Ansible 2.2+
* Docker Machine

Provider requirements, for example (EC2):
* pip install boto
* export AWS_ACCESS_KEY
* export AWS_SECRET_KEY

For example (OpenStack)
* export OS_USERNAME
* export OS_PASSWORD
* export OS_TENANT_NAME
* export OS_AUTH_URL

Steps:

1. Compile `library/docker_machine.go` and `library/docker_swarm.go`

2. Setup env manually or through a play (i.e. `ansible-playbook *setup.yml`)

3. Instanciate managers (i.e. `forloop: ansible-playbook -M library *swarm_manager.yml`)

4. Instanciate workers (i.e. `forloop: ansible-playbook -M library *swarm_worker.yml`)

5. Have fun
