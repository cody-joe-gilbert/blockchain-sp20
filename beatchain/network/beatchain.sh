#!/bin/bash
# Blockchain & Applications Project 2: Beatchain
# Owner(s): Cody Gilbert

####
# This shell file was orginially used within the IBM TradeChannel example as `trade.sh`,
# and has been modified from the original version to support this new network.
####


# This script will orchestrate a sample end-to-end execution of the Hyperledger
# Fabric network.
#
# The end-to-end verification provisions a sample Fabric network consisting of
# two organizations, each maintaining two peers, and a “solo” ordering service.
#
# This verification makes use of two fundamental tools, which are necessary to
# create a functioning transactional network with digital signature validation
# and access control:
#
# * cryptogen - generates the x509 certificates used to identify and
#   authenticate the various components in the network.
# * configtxgen - generates the requisite configuration artifacts for orderer
#   bootstrap and channel creation.
#
# Each tool consumes a configuration yaml file, within which we specify the topology
# of our network (cryptogen) and the location of our certificates for various
# configuration operations (configtxgen).  Once the tools have been successfully run,
# we are able to launch our network.  More detail on the tools and the structure of
# the network will be provided later in this document.  For now, let's get going...

# prepending $PWD/../bin to PATH to ensure we are picking up the correct binaries
# this may be commented out to resolve installed version of tools if desired

export PATH=${PWD}/../bin:${PWD}:$PATH
export FABRIC_CFG_PATH=${PWD}

# Print the usage message
function printHelp () {
  echo "Usage: "
  echo "  beatchain.sh up|down|restart|generate|reset|clean [-c <channel name>] [-f <docker-compose-file>] [-i <imagetag>] [-o <logfile>] [-dev]"
  echo "  beatchain.sh -h|--help (print this message)"
  echo "    <mode> - one of 'up', 'down', 'restart' or 'generate'"
  echo "      - 'up' - bring up the network with docker-compose up"
  echo "      - 'down' - clear the network with docker-compose down"
  echo "      - 'restart' - restart the network"
  echo "      - 'generate' - generate required certificates and genesis block"
  echo "      - 'reset' - delete chaincode containers while keeping network artifacts" 
  echo "      - 'clean' - delete network artifacts"
  echo "    -c <channel name> - channel name to use (defaults to \"beatchainchannel\")"
  echo "    -f <docker-compose-file> - specify which docker-compose file use (defaults to docker-compose.yaml)"
  echo "    -i <imagetag> - the tag to be used to launch the network (defaults to \"latest\")"
  echo "    -d - Apply command to the network in dev mode."
  echo
  echo "Typically, one would first generate the required certificates and "
  echo "genesis block, then bring up the network. e.g.:"
  echo
  echo "	beatchain.sh generate -c fullchannel"
  echo "	beatchain.sh up -c fullchannel -o logs/network.log"
  echo "        beatchain.sh up -c fullchannel -i 1.1.0-alpha"
  echo "	beatchain.sh down -c fullchannel"
  echo
  echo "Taking all defaults:"
  echo "	beatchain.sh generate"
  echo "	beatchain.sh up"
  echo "	beatchain.sh down"
}

# Keeps pushd silent
pushd () {
    command pushd "$@" > /dev/null
}

# Keeps popd silent
popd () {
    command popd "$@" > /dev/null
}

# Ask user for confirmation to proceed
function askProceed () {
  read -p "Continue? [Y/n] " ans
  case "$ans" in
    y|Y|"" )
      echo "proceeding ..."
    ;;
    n|N )
      echo "exiting..."
      exit 1
    ;;
    * )
      echo "invalid response"
      askProceed
    ;;
  esac
}

# Obtain CONTAINER_IDS and remove them
# TODO Might want to make this optional - could clear other containers
function clearContainers () {
  CONTAINER_IDS=$(docker ps -aq)
  if [ -z "$CONTAINER_IDS" -o "$CONTAINER_IDS" == " " ]; then
    echo "---- No containers available for deletion ----"
  else
    docker rm -f $CONTAINER_IDS
  fi
}

# Delete any images that were generated as a part of this setup
# specifically the following images are often left behind:
# TODO list generated image naming patterns
function removeUnwantedImages() {
  DOCKER_IMAGE_IDS=$(docker images | grep "dev\|none\|test-vp\|peer[0-9]-" | awk '{print $3}')
  if [ -z "$DOCKER_IMAGE_IDS" -o "$DOCKER_IMAGE_IDS" == " " ]; then
    echo "---- No images available for deletion ----"
  else
    docker rmi -f $DOCKER_IMAGE_IDS
  fi
}

# Do some basic sanity checking to make sure that the appropriate versions of fabric
# binaries/images are available.  In the future, additional checking for the presence
# of go or other items could be added.
function checkPrereqs() {
  # Note, we check configtxlator externally because it does not require a config file, and peer in the
  # docker image because of FAB-8551 that makes configtxlator return 'development version' in docker
  LOCAL_VERSION=$(configtxlator version | sed -ne 's/ Version: //p')
  DOCKER_IMAGE_VERSION=$(docker run --rm hyperledger/fabric-tools:$IMAGETAG peer version | sed -ne 's/ Version: //p'|head -1)

  echo "LOCAL_VERSION=$LOCAL_VERSION"
  echo "DOCKER_IMAGE_VERSION=$DOCKER_IMAGE_VERSION"

  if [ "$LOCAL_VERSION" != "$DOCKER_IMAGE_VERSION" ] ; then
     echo "=================== WARNING ==================="
     echo "  Local fabric binaries and docker images are  "
     echo "  out of  sync. This may cause problems.       "
     echo "==============================================="
  fi
}

# Generate the needed certificates, the genesis block and start the network.
function networkUp () {
  checkPrereqs
  # generate artifacts if they don't exist
  if [ ! -d "crypto-config" ]; then
    generateCerts
    replacePrivateKey
    generateChannelArtifacts
  fi
  # Create folder for docker network logs
  LOG_DIR=$(dirname $LOG_FILE)
  if [ ! -d $LOG_DIR ]
  then
    mkdir -p $LOG_DIR
  fi
  IMAGE_TAG=$IMAGETAG docker-compose -f $COMPOSE_FILE up >$LOG_FILE 2>&1 &

  if [ $? -ne 0 ]; then
    echo "ERROR !!!! Unable to start network"
    exit 1
  fi
}

# Bring down running network
function networkDown () {

  docker-compose -f $COMPOSE_FILE down --volumes

  for PEER in peer0.creatororg.beatchain.com peer0.customerorg.beatchain.com peer0.appdevorg.beatchain.com; do
    # Remove any old containers and images for this peer
    CC_CONTAINERS=$(docker ps -a | grep dev-$PEER | awk '{print $1}')
    if [ -n "$CC_CONTAINERS" ] ; then
      docker rm -f $CC_CONTAINERS
    fi
  done

}


# Delete network artifacts
function networkClean () {
  #Cleanup the chaincode containers
  clearContainers
  #Cleanup images
  removeUnwantedImages
  # remove orderer block and other channel configuration transactions and certs
  rm -rf channel-artifacts crypto-config add_org/crypto-config
  # remove the docker-compose yaml file that was customized to the example
  rm -f docker-compose.yaml
  # remove client certs 
  rm -rf client-certs
}

# Using docker-compose-template.yaml, replace constants with private key file names
# generated by the cryptogen tool and output a docker-compose.yaml specific to this
# configuration
function replacePrivateKey () {
  # Copy the template to the file that will be modified to add the private key
  cp docker-compose-template.yaml docker-compose.yaml
  cp config-template.json config.json

  # The next steps will replace the template's contents with the
  # actual values of the private key file names for the CAs.
  CURRENT_DIR=$PWD
  CUR_ORG="appdevorg"
  ORG_KEY="APPDEVORG"
  CA_DIR=crypto-config/peerOrganizations/${CUR_ORG}.beatchain.com/ca/
  cp fabric-ca-server-config.yaml ${CA_DIR}/fabric-ca-server-config.yaml
  cd ${CA_DIR}
  PRIV_KEY=$(ls *_sk)
  cp ca.${CUR_ORG}.beatchain.com-cert.pem ca-cert.pem
  cp ${PRIV_KEY} ca-key.pem
  cd "$CURRENT_DIR"
  if [ $(uname -s) == 'Darwin' ] ; then
    sed -i '' "s/${ORG_KEY}_CA_PRIVATE_KEY/${PRIV_KEY}/g" docker-compose.yaml
  else
    sed -i "s/${ORG_KEY}_CA_PRIVATE_KEY/${PRIV_KEY}/g" docker-compose.yaml
  fi
  cd crypto-config/peerOrganizations/${CUR_ORG}.beatchain.com/users/Admin@${CUR_ORG}.beatchain.com/msp/keystore/
  PRIV_KEY=$(ls *_sk)
  cp ca.${CUR_ORG}.beatchain.com-cert.pem ca-cert.pem
  cp ${PRIV_KEY} ca-key.pem
  cd "$CURRENT_DIR"
  if [ $(uname -s) == 'Darwin' ] ; then
    sed -i '' "s/${ORG_KEY}_ADMIN_PRIVATE_KEY/${PRIV_KEY}/g" config.json
  else
    sed -i "s/${ORG_KEY}_ADMIN_PRIVATE_KEY/${PRIV_KEY}/g" config.json
  fi

  CUR_ORG="creatororg"
  ORG_KEY="CREATORORG"
  CA_DIR=crypto-config/peerOrganizations/${CUR_ORG}.beatchain.com/ca/
  cp fabric-ca-server-config.yaml ${CA_DIR}/fabric-ca-server-config.yaml
  cd ${CA_DIR}
  PRIV_KEY=$(ls *_sk)
  cp ca.${CUR_ORG}.beatchain.com-cert.pem ca-cert.pem
  cp ${PRIV_KEY} ca-key.pem
  cd "$CURRENT_DIR"
  if [ $(uname -s) == 'Darwin' ] ; then
    sed -i '' "s/${ORG_KEY}_CA_PRIVATE_KEY/${PRIV_KEY}/g" docker-compose.yaml
  else
    sed -i "s/${ORG_KEY}_CA_PRIVATE_KEY/${PRIV_KEY}/g" docker-compose.yaml
  fi
  cd crypto-config/peerOrganizations/${CUR_ORG}.beatchain.com/users/Admin@${CUR_ORG}.beatchain.com/msp/keystore/
  PRIV_KEY=$(ls *_sk)
  cd "$CURRENT_DIR"
  if [ $(uname -s) == 'Darwin' ] ; then
    sed -i '' "s/${ORG_KEY}_ADMIN_PRIVATE_KEY/${PRIV_KEY}/g" config.json
  else
    sed -i "s/${ORG_KEY}_ADMIN_PRIVATE_KEY/${PRIV_KEY}/g" config.json
  fi

  CUR_ORG="customerorg"
  ORG_KEY="CUSTOMERORG"
  CA_DIR=crypto-config/peerOrganizations/${CUR_ORG}.beatchain.com/ca/
  cp fabric-ca-server-config.yaml ${CA_DIR}/fabric-ca-server-config.yaml
  cd ${CA_DIR}
  PRIV_KEY=$(ls *_sk)
  cp ca.${CUR_ORG}.beatchain.com-cert.pem ca-cert.pem
  cp ${PRIV_KEY} ca-key.pem
  cd "$CURRENT_DIR"
  if [ $(uname -s) == 'Darwin' ] ; then
    sed -i '' "s/${ORG_KEY}_CA_PRIVATE_KEY/${PRIV_KEY}/g" docker-compose.yaml
  else
    sed -i "s/${ORG_KEY}_CA_PRIVATE_KEY/${PRIV_KEY}/g" docker-compose.yaml
  fi
  cd crypto-config/peerOrganizations/${CUR_ORG}.beatchain.com/users/Admin@${CUR_ORG}.beatchain.com/msp/keystore/
  PRIV_KEY=$(ls *_sk)
  cd "$CURRENT_DIR"
  if [ $(uname -s) == 'Darwin' ] ; then
    sed -i '' "s/${ORG_KEY}_ADMIN_PRIVATE_KEY/${PRIV_KEY}/g" config.json
  else
    sed -i "s/${ORG_KEY}_ADMIN_PRIVATE_KEY/${PRIV_KEY}/g" config.json
  fi


}


# We will use the cryptogen tool to generate the cryptographic material (x509 certs)
# for our various network entities.  The certificates are based on a standard PKI
# implementation where validation is achieved by reaching a common trust anchor.
#
# Cryptogen consumes a file - ``crypto-config.yaml`` - that contains the network
# topology and allows us to generate a library of certificates for both the
# Organizations and the components that belong to those Organizations.  Each
# Organization is provisioned a unique root certificate (``ca-cert``), that binds
# specific components (peers and orderers) to that Org.  Transactions and communications
# within Fabric are signed by an entity's private key (``keystore``), and then verified
# by means of a public key (``signcerts``).  You will notice a "count" variable within
# this file.  We use this to specify the number of peers per Organization; in our
# case it's two peers per Org.  The rest of this template is extremely
# self-explanatory.
#
# After we run the tool, the certs will be parked in a folder titled ``crypto-config``.

# Generates Org certs using cryptogen tool
function generateCerts (){
  which cryptogen
  if [ "$?" -ne 0 ]; then
    echo "cryptogen tool not found. exiting"
    exit 1
  fi
  echo
  echo "##########################################################"
  echo "##### Generate certificates using cryptogen tool #########"
  echo "##########################################################"

  if [ -d "crypto-config" ]; then
    rm -Rf crypto-config
  fi
  set -x
  cryptogen generate --config=./crypto-config.yaml
  res=$?
  set +x
  if [ $res -ne 0 ]; then
    echo "Failed to generate certificates..."
    exit 1
  fi
  echo
}

# The `configtxgen tool is used to create four artifacts: orderer **bootstrap
# block**, fabric **channel configuration transaction**, and two **anchor
# peer transactions** - one for each Peer Org.
#
# The orderer block is the genesis block for the ordering service, and the
# channel transaction file is broadcast to the orderer at channel creation
# time.  The anchor peer transactions, as the name might suggest, specify each
# Org's anchor peer on this channel.
#
# Configtxgen consumes a file - ``configtx.yaml`` - that contains the definitions
# for the sample network. There are five members - one Orderer Org (``beatchainOrdererOrg``)
# and four Peer Orgs (``appdevorg``, ``creatororg``, ``customerorg`` & ``RegulatorOrg``)
# each managing and maintaining one peer node.
# This file also specifies a consortium - ``beatchainConsortium`` - consisting of our
# four Peer Orgs.  Pay specific attention to the "Profiles" section at the top of
# this file.  You will notice that we have two unique headers. One for the orderer genesis
# block - ``FourOrgsbeatchainOrdererGenesis`` - and one for our channel - ``FourOrgsbeatchainChannel``.
# These headers are important, as we will pass them in as arguments when we create
# our artifacts.  This file also contains two additional specifications that are worth
# noting.  Firstly, we specify the anchor peers for each Peer Org
# (``peer0.appdevorg.beatchain.com`` & ``peer0.creatororg.beatchain.com``).  Secondly, we point to
# the location of the MSP directory for each member, in turn allowing us to store the
# root certificates for each Org in the orderer genesis block.  This is a critical
# concept. Now any network entity communicating with the ordering service can have
# its digital signature verified.
#
# This function will generate the crypto material and our four configuration
# artifacts, and subsequently output these files into the ``channel-artifacts``
# folder.
#
# If you receive the following warning, it can be safely ignored:
#
# [bccsp] GetDefault -> WARN 001 Before using BCCSP, please call InitFactories(). Falling back to bootBCCSP.
#
# You can ignore the logs regarding intermediate certs, we are not using them in
# this crypto implementation.

# Generate orderer genesis block, channel configuration transaction and
# anchor peer update transactions
function generateChannelArtifacts() {
  which configtxgen
  if [ "$?" -ne 0 ]; then
    echo "configtxgen tool not found. exiting"
    exit 1
  fi
  mkdir -p channel-artifacts

  echo "###########################################################"
  echo "#########  Generating Orderer Genesis block  ##############"
  echo "###########################################################"

  PROFILE=BeatchainGenesis
  CHANNEL_PROFILE=fullchannel
  # Note: For some unknown reason (at least for now) the block file can't be
  # named orderer.genesis.block or the orderer will fail to launch!
  set -x
  configtxgen -profile $PROFILE -outputBlock ./channel-artifacts/genesis.block
  res=$?
  set +x
  if [ $res -ne 0 ]; then
    echo "Failed to generate orderer genesis block..."
    exit 1
  fi
  echo
  echo "###################################################################"
  echo "###  Generating channel configuration transaction  'channel.tx' ###"
  echo "###################################################################"
  set -x
  configtxgen -profile $CHANNEL_PROFILE -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID $CHANNEL_NAME
  res=$?
  set +x
  if [ $res -ne 0 ]; then
    echo "Failed to generate channel configuration transaction..."
    exit 1
  fi


  echo
  echo "#####################################################################"
  echo "#######  Generating anchor peer update for AppDevMSP  ##########"
  echo "#####################################################################"
  set -x
  configtxgen -profile $CHANNEL_PROFILE -outputAnchorPeersUpdate ./channel-artifacts/AppDevMSPanchors.tx -channelID $CHANNEL_NAME -asOrg AppDevOrg
  res=$?
  set +x
  if [ $res -ne 0 ]; then
    echo "Failed to generate anchor peer update for AppDevMSP..."
    exit 1
  fi

  echo
  echo "#####################################################################"
  echo "#######  Generating anchor peer update for CreatorMSP  ##########"
  echo "#####################################################################"
  set -x
  configtxgen -profile $CHANNEL_PROFILE -outputAnchorPeersUpdate \
  ./channel-artifacts/CreatorMSPanchors.tx -channelID $CHANNEL_NAME -asOrg CreatorOrg -channelID $CHANNEL_NAME
  res=$?
  set +x
  if [ $res -ne 0 ]; then
    echo "Failed to generate anchor peer update for creatororgMSP..."
    exit 1
  fi

  echo
  echo "####################################################################"
  echo "#######  Generating anchor peer update for CustomerMSP  ##########"
  echo "####################################################################"
  set -x
  configtxgen -profile $CHANNEL_PROFILE -outputAnchorPeersUpdate \
  ./channel-artifacts/CustomerMSPanchors.tx -channelID $CHANNEL_NAME -asOrg CustomerOrg -channelID $CHANNEL_NAME
  res=$?
  set +x
  if [ $res -ne 0 ]; then
    echo "Failed to generate anchor peer update for CustomerMSP..."
    exit 1
  fi

}

# channel name (overrides default 'testchainid')
CHANNEL_NAME="fullchannel"
# use this as the default docker-compose yaml definition
COMPOSE_FILE=docker-compose.yaml
# default image tag
IMAGETAG="latest"
# default log file
LOG_FILE="logs/network.log"
# Parse commandline args
MODE=$1;shift
# Determine whether starting, stopping, restarting or generating for announce
if [ "$MODE" == "up" ]; then
  EXPMODE="Starting"
elif [ "$MODE" == "down" ]; then
  EXPMODE="Stopping"
elif [ "$MODE" == "restart" ]; then
  EXPMODE="Restarting"
elif [ "$MODE" == "clean" ]; then
  EXPMODE="Cleaning"
elif [ "$MODE" == "generate" ]; then
  EXPMODE="Generating certs and genesis block"
else
  printHelp
  exit 1
fi

while getopts "h?m:c:f:i:o:d:" opt; do
  case "$opt" in
    h|\?)
      printHelp
      exit 0
    ;;
    c)  CHANNEL_NAME=$OPTARG
    ;;
    f)  COMPOSE_FILE=$OPTARG
    ;;
    i)  IMAGETAG=`uname -m`"-"$OPTARG
    ;;
    o)  LOG_FILE=$OPTARG
    ;;
    d)  DEV_MODE=$OPTARG 
    ;;
  esac
done

# Announce what was requested
echo "${EXPMODE} with channel '${CHANNEL_NAME}'"
# ask for confirmation to proceed
askProceed

#Create the network using docker compose
if [ "${MODE}" == "up" ]; then
  networkUp
elif [ "${MODE}" == "down" ]; then ## Clear the network
  networkDown
elif [ "${MODE}" == "generate" ]; then ## Generate Artifacts
  generateCerts
  replacePrivateKey
  generateChannelArtifacts
elif [ "${MODE}" == "restart" ]; then ## Restart the network
  networkDown
  networkUp
elif [ "${MODE}" == "reset" ]; then ## Delete chaincode containers while keeping network artifacts
  removeUnwantedImages
elif [ "${MODE}" == "clean" ]; then ## Delete network artifacts
  networkClean
else
  printHelp
  exit 1
fi
