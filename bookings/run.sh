#!/bin/bash

go build -o bookings cmd/web/*.go && 
./bookings -dbname=bookings -dbuser=kennethakor -cache=false -production=false