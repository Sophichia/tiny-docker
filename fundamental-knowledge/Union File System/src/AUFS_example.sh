#!/usr/bin/env bash

cd /tmp/aufs

# Create a mnt foler as the mount point
mkdir mnt

# Create a container-layer folder with a container-layer.txt file in it.
mkdir container-layer
echo "I am container layer" > container-layer/container-layer.txt

# Create 4 imager-layer<number> folders with image-layer<num>.txt
for i in {1..4}
do
    mkdir "image-layer${i}"
    echo "I am image layer ${i}" >> "image-layer${i}"/"image-layer${i}.txt"
done

# use aufs to mount all those folders to mnt.
# By default, the most left dir will hav read-write permission, all others will be read-only
sudo mount -t aufs -o dirs=./container-layer:./image-layer4:./image-layer3:./image-layer2:./image-layer1 none ./mnt

# Check mnt structure
tree mnt

# Append a line into image-layer4.txt
echo -e "\nwrite to mnt's image-layer4.txt" >> ./mnt/image-layer4.txt

# Line added in mnt/image-layer4.txt
cat ./mnt/image-layer4.txt

# Line not added in image-layer4/image-layer4.txt
cat ./image-layer4/image-layer4.txt

# It actually create a new file in container-layer
cat container-layer/image-layer4.txt