#!/bin/bash

for i in `seq 1 6`
do
    rm -rf /home/tsy/tmp/$i/objects/*
    rm -rf /tmp/$i/temp/*
done
