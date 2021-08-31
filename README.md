# somewhat-emulate-skopeo
Proof of Concept for B/R of images in order to eliminate dependency on registry deployment

## quickstart
- Install dependencies like `yum install libassuan-devel gpgme-devel device-mapper-devel`
- Install `dnf install btrfs-progs-devel`
- Build the binary as root and then execute as the program needs root acccess
