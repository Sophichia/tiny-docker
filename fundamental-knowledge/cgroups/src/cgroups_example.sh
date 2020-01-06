#!/usr/bin/env bash

# Create a hierarchy mount point
mkdir cgroup-test

# Mount a hierarchy
sudo mount -t cgroup -o none,name=cgroup-test cgroup-test ./cgroup-test

# After mounting, there're some files generated under this foler.
# Those files are the configuration in this root cgroup.
ls ./cgroup-test

cd cgroup-test

# Create two sub-cgroups in this root cgroup
sudo mkdir cgroup-1
sudo mkdir cgroup-2

# Add current process into cgroup1
cd cgroup1
sudo sh -c "echo $$ >> tasks"
cat /proc/<pid>/cgroup

# Applying limitations on cgroup via subsystem
