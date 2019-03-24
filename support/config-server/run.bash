#!/bin/bash

timestamp=$(date +%s)
mv server.jks server.jks.${timestamp}

#create server.jks
keytool -genkeypair -alias aliaskey -keyalg RSA -dname "CN=My Microservices,OU=Unit,O=Organization,L=City,S=State,C=SE" -keypass secret -keystore server.jks -storepass password -validity 730

./gradlew build
