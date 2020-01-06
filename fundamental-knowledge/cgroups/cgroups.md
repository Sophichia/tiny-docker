# Cgroups - Linux control groups

Control groups, usually referred to as cgroups, are a Linux kernel feature which allow processes to be organized into hierarchical groups whose usage of various types of resources can then be limited and monitored.

There's 4 features of cgroups to a developer:
1. The kernel's cgroup interface is provided through a pseudo-filesystem called cgroupsfs.
2. Granularity of cgroups can down to thread.
3. All resource management is achieved by sub-system.
4. All sub-tasks originally created within same group with its parent.

## Reading Materials
1. https://www.cnblogs.com/menkeyi/p/10941843.html
2. https://git.kernel.org/pub/scm/linux/kernel/git/torvalds/linux.git/tree/Documentation/cgroup-v1/cgroups.txt?h=v4.18-rc1