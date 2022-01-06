#!/usr/bin/env bash

# This script is used to perform a nightly cleanup of any terraform resources that might still be present on the Cloud Native account. 

# Filter Requirements
# -------------------
# tf-
# Terraform

sudo apt-get -y install curl jq

#####################
### RegEx Filters ###
#####################
filters=("tf-" "Terraform")
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

####################
### VKE CLUSTERS ###
####################

CLUSTERS=$(curl --silent --location --request GET 'https://api.vultr.com/v2/kubernetes/clusters' --header "Authorization: Bearer $VULTR_API_KEY" | jq -r '.vke_clusters[] | select(.).id')

clean_vke(){
  for CLUSTER in $CLUSTERS; do
    label=$(curl --silent --location --request GET "https://api.vultr.com/v2/kubernetes/clusters/$CLUSTER" --header "Authorization: Bearer $VULTR_API_KEY" | jq -r '.vke_cluster | select(.).label')

    for element in "${filters[@]}"; do
        if [[ $label =~ $element* ]]; then
            echo -e "${RED}- Deleting $label - $CLUSTER ${NC} "
            curl --silent --location --request DELETE "https://api.vultr.com/v2/kubernetes/clusters/$CLUSTER" --header "Authorization: Bearer $VULTR_API_KEY"
            break
        fi
    done

  done
}

#################
### INSTANCES ###
#################

INSTANCES=$(curl --silent --location --request GET 'https://api.vultr.com/v2/instances' --header "Authorization: Bearer $VULTR_API_KEY" | jq -r '.instances[] | select(.).id')

clean_instance(){
  for INSTANCE in $INSTANCES; do
    label=$(curl --silent --location --request GET "https://api.vultr.com/v2/instances/$INSTANCE" --header "Authorization: Bearer $VULTR_API_KEY" | jq -r '.instance | select(.).label')

    for element in "${filters[@]}"; do
        if [[ $label =~ $element* ]]; then
            echo -e "${RED}- Deleting $label - $INSTANCE ${NC} "
            curl --silent --location --request DELETE "https://api.vultr.com/v2/instances/$INSTANCE" --header "Authorization: Bearer $VULTR_API_KEY"
            break
        fi
    done

  done
}

##################
### BARE METAL ###
##################

BAREMETALS=$(curl --silent --location --request GET 'https://api.vultr.com/v2/bare-metals' --header "Authorization: Bearer $VULTR_API_KEY" | jq -r '.bare_metals[] | select(.).id')

clean_baremetal(){
  for BAREMETAL in $BAREMETALS; do
    label=$(curl --silent --location --request GET "https://api.vultr.com/v2/bare-metals/$BAREMETAL" --header "Authorization: Bearer $VULTR_API_KEY" | jq -r '.bare_metal | select(.).label')

    for element in "${filters[@]}"; do
        if [[ $label =~ $element* ]]; then
            echo -e "${RED}- Deleting $label - $BAREMETAL ${NC} "
            curl --silent --location --request DELETE "https://api.vultr.com/v2/bare-metals/$BAREMETAL" --header "Authorization: Bearer $VULTR_API_KEY"
            break
        fi
    done

  done
}

######################
### LOAD BALANCERS ###
######################

LOADBALANCERS=$(curl --silent --location --request GET 'https://api.vultr.com/v2/load-balancers' --header "Authorization: Bearer $VULTR_API_KEY" | jq -r '.load_balancers[] | select(.).id')

clean_loadbalancers(){
  for LOADBALANCER in $LOADBALANCERS; do
    label=$(curl --silent --location --request GET "https://api.vultr.com/v2/load-balancers/$LOADBALANCER" --header "Authorization: Bearer $VULTR_API_KEY" | jq -r '.load_balancer | select(.).label')

    for element in "${filters[@]}"; do
        if [[ $label =~ $element* ]]; then
            echo -e "${RED}- Deleting $label - $LOADBALANCER ${NC} "
            curl --silent --location --request DELETE "https://api.vultr.com/v2/load-balancers/$LOADBALANCER" --header "Authorization: Bearer $VULTR_API_KEY"
            break
        fi
    done

  done
}

#################
### SNAPSHOTS ###
#################

SNAPSHOTS=$(curl --silent --location --request GET 'https://api.vultr.com/v2/snapshots' --header "Authorization: Bearer $VULTR_API_KEY" | jq -r '.snapshots[] | select(.).id')

clean_snapshots(){
  for SNAPSHOT in $SNAPSHOTS; do
    label=$(curl --silent --location --request GET "https://api.vultr.com/v2/snapshots/$SNAPSHOT" --header "Authorization: Bearer $VULTR_API_KEY" | jq -r '.snapshot | select(.).description')

    for element in "${filters[@]}"; do
        if [[ $label =~ $element* ]]; then
            echo -e "${RED}- Deleting $label - $SNAPSHOT ${NC} "
            curl --silent --location --request DELETE "https://api.vultr.com/v2/snapshots/$SNAPSHOT" --header "Authorization: Bearer $VULTR_API_KEY"
            break
        fi
    done

  done
}

#####################
### BLOCK STORAGE ###
#####################

BLOCKS=$(curl --silent --location --request GET 'https://api.vultr.com/v2/blocks' --header "Authorization: Bearer $VULTR_API_KEY" | jq -r '.blocks[] | select(.).id')

clean_blocks(){
  for BLOCK in $BLOCKS; do
    label=$(curl --silent --location --request GET "https://api.vultr.com/v2/blocks/$BLOCK" --header "Authorization: Bearer $VULTR_API_KEY" | jq -r '.block | select(.).label')

    for element in "${filters[@]}"; do
        if [[ $label =~ $element* ]]; then
            echo -e "${RED}- Deleting $label - $BLOCK ${NC} "
            curl --silent --location --request DELETE "https://api.vultr.com/v2/blocks/$BLOCK" --header "Authorization: Bearer $VULTR_API_KEY"
            break
        fi
    done

  done
}

#####################
### RESERVED IP'S ###
#####################

RIPS=$(curl --silent --location --request GET 'https://api.vultr.com/v2/reserved-ips' --header "Authorization: Bearer $VULTR_API_KEY" | jq -r '.reserved_ips[] | select(.).id')

clean_rips(){
  for RIP in $RIPS; do
    label=$(curl --silent --location --request GET "https://api.vultr.com/v2/reserved-ips/$RIP" --header "Authorization: Bearer $VULTR_API_KEY" | jq -r '.reserved_ip | select(.).label')

    for element in "${filters[@]}"; do
        if [[ $label =~ $element* ]]; then
            echo -e "${RED}- Deleting $label - $RIP ${NC} "
            curl --silent --location --request DELETE "https://api.vultr.com/v2/reserved-ips/$RIP" --header "Authorization: Bearer $VULTR_API_KEY"
            break
        fi
    done

  done
}


main(){
  # VKE Clusters
  if [ "$CLUSTERS" != "" ]; then
    echo "+ Cleaning vke-clusters"
    clean_vke
    echo "+ VKE-clusters cleaned"
  else
    echo -e "${GREEN}+ No vke-clusters to remove${NC}"
  fi

  # Instances
  if [ "$INSTANCES" != "" ]; then
    echo "+ Cleaning instances"
    clean_instance
    echo "+ Instances cleaned"
  else
    echo -e "${GREEN}+ No instances to remove${NC}"
  fi

  # Bare-Metal
  if [ "$BAREMETALS" != "" ]; then
    echo "+ Cleaning bare metals"
    clean_baremetal
    echo "+ Bare metals cleaned"
  else
    echo -e "${GREEN}+ No bare metals to remove${NC}"
  fi

  # Load-Balancers
  if [ "$LOADBALANCERS" != "" ]; then
    echo "+ Cleaning load balancers"
    clean_loadbalancers
    echo "+ Load balancers cleaned"
  else
    echo -e "${GREEN}+ No load balancers to remove${NC}"
  fi

  # Snapshots
  if [ "$SNAPSHOTS" != "" ]; then
    echo "+ Cleaning snapshots"
    clean_snapshots
    echo "+ Snapshots cleaned"
  else
    echo -e "${GREEN}+ No snapshots to remove${NC}"
  fi

    # Blocks
  if [ "$BLOCKS" != "" ]; then
    echo "+ Cleaning block storage"
    clean_blocks
    echo "+ Block Storage cleaned"
  else
    echo -e "${GREEN}+ No block storages to remove${NC}"
  fi

    # Reserved IP's
  if [ "$RIPS" != "" ]; then
    echo "+ Cleaning reserved ip's"
    clean_rips
    echo "+ Reserved ip's cleaned"
  else
    echo -e "${GREEN}+ No reserved ip's to remove${NC}"
  fi

}

main
