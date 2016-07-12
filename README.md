ansible-swarm
=============

!!! WIP

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

1. Compile library/docker_machine.go

2. Setup env (ex. `ansible-playbook *setup.yml`)

3. Instanciate masters (ex `forloop: ansible-playbook -M library *swarm_master.yml`)

4. Instanciate slaves (ex `forloop: ansible-playbook -M library *swarm_slave.yml`)
