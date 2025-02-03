#!/bin/bash -xe


REPO=Intel-Tiber-Broadcast-Suite
git clone https://github.com/OpenVisualCloud/${REPO}.git
cd ${REPO}
git submodule update --init --recursive