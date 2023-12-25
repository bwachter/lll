#!/bin/bash

if ! [ -d rootfs ]; then
   echo "You need to generate a rootfs in ./rootfs. See README"
   exit 1
fi

if ! [ -f rootfs/bzImage ]; then
    echo "./rootfs/bzImage not found. Build a kernel and put it there. See README."
    exit 1
fi

if [ "`which kvm`" ]; then
    kvm -cpu host -nographic -virtfs local,id=linux,path=rootfs,security_model=none,mount_tag=linux-root -net nic  -append "root=linux-root rootfstype=9p rootflags=trans=virtio,version=9p2000.L,nodevmap console=ttyS0,115200n8" -kernel rootfs/bzImage
elif [ "`which qemu-system-x86_64`" ]; then
    qemu-system-x86_64 -nographic -virtfs local,id=linux,path=rootfs,security_model=none,mount_tag=linux-root -net nic  -append "root=linux-root rootfstype=9p rootflags=trans=virtio,version=9p2000.L,nodevmap console=ttyS0,115200n8" -kernel rootfs/bzImage
else
    echo "Neither kvm nor qemu-system-x86_64 found, exiting"
    exit 1
fi
