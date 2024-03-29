#!/bin/bash
#instruction=$1
#version=$2

set -ev

#chaincode install
docker exec cli peer chaincode install -n bc_game -v 2 -p github.com/bc_game
#chaincode instatiate
docker exec cli peer chaincode instantiate -n bc_game -v 2 -C mychannel -c '{"Args":[]}' -P 'OR ("Org1MSP.member", "Org2MSP.member","Org3MSP.member")'
sleep 5
#chaincode invoke user1
docker exec cli peer chaincode invoke -n bc_game -C mychannel -c '{"Args":["addUser","kim"]}'
sleep 5
#chaincode query user1
docker exec cli peer chaincode query -n bc_game -C mychannel -c '{"Args":["readUserInfo","kim"]}'

#chaincode invoke add rating
#docker exec cli peer chaincode invoke -n donation -C mychannel -c '{"Args":["addRating","user1","p1","5.0"]}'
#sleep 5

echo '-------------------------------------END-------------------------------------'
