#!/bin/bash
SWARM2KID="swarm2k"

#for machine in `docker-machine ls | grep $SWARM2KID | awk '{print $1;}'`;
#do
$machine="swarm2k0"
MACHINE_NAME=$machine
ansible-playbook -M library --extra-vars '{"machine_name": swarm2k0}' swarm2k/join_swarm2k.yml
#done
