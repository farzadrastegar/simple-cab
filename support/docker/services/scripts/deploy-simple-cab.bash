#!/bin/bash
homeDir=`pwd`

cd gateway
nohup make all >${homeDir}/simple-cab-gateway.out 2>&1 &
echo $! > ${homeDir}/run.pid
cd ${homeDir}
        
cd zombie_driver
nohup make all >${homeDir}/simple-cab-zombie_driver.out 2>&1 &
echo $! >> ${homeDir}/run.pid
cd ${homeDir}
        
cd driver_location 
nohup make all >${homeDir}/simple-cab-driver_location.out 2>&1 &
echo $! >> ${homeDir}/run.pid
cd ${homeDir}
