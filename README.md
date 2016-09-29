# swarm-util
docker-machine wrapper for swarm mode(from docker 1.12) 

## create grouped machines
### GCE
swarm-util -name GROUP_NAME -manager 1 -worker 2 create --driver google --google-project PROJECT_ID --google-zone asia-east1-c --google-machine-type n1-standard-1 --google-preemptible

GROUP_NAME-m0

GROUP_NAME-w0

GROUP_NAME-w1

## remove grouped machines
swarm-util -name GROUP_NAME rm

## create swarm (init/join)
swarm-util -name GROUP_NAME init

## remove swarm (leave)
swarm-util -name GROUP_NAME leave
