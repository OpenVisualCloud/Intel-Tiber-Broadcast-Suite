#!/bin/bash -ex

SCRIPT_DIR="jenkins/scripts"
mkdir -p Trivy
mkdir -p Trivy/dockerfile
chmod a+w Trivy
chmod a+w Trivy/dockerfile

trivy conf  -o Trivy/dockerfile/source_config.txt ./Dockerfile 
