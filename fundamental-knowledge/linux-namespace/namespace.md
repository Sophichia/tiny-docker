# Linux Namespace

## Overall

Docker uses linux namespace to achieve the functionality of resource isolation.

At a high level, linux namespace allows for isolation of global system resources between independent processes.

Note that namespaces do not restrict access to physical resources such as CPU, memory and disk. That access is metered and restricted by a kernel feature called `cgroups`

## Linux Namespaces Types

| Namespace | Constant        | Isolates                             |
|-----------|-----------------|--------------------------------------|
| Cgroup    | CLONE_NEWCGROUP | Cgroup root directory                |
| IPC       | CLONE_NEWIPC    | System V IPC, POSIX message queues   |
| Network   | CLONE_NEWNET    | Network devices, stacks, ports, etc. |
| Mount     | CLONE_NEWNS     | Filesystem Mount points              |
| PID       | CLONE_NEWPID    | Process IDs                          |
| User      | CLONE_NEWUSER   | User and group IDs                   |
| UTS       | CLONE_NEWUTS    | Hostname and NIS domain name         |

## Namespace APIs

Namespace APIs are including `clone()`, `setns()`, `unshare()` and some files under `/proc`.

When using these APIs, we need to assign namespaces by `|`.

### using `clone()` to create a namespace together with creating a new process.
Using `clone()` to create an isolated namespace process is the common way. Docker also uses this way.

`int clone(int (*child_func)(void *), void *child_stack, int flags, void *arg)`

- `child_func` parses in the main function in sub-process
- `child_stack` parses in the stack space used in sub-prrocess
- `flags` means `CLONE_*` flags, like CLONE_NEWIPC, CLONE_NEWNS, etc.
- `args` parses in the user arguments

### `/proc/[pid]/ns` file
After kernel v3.8, users can find files pointing to namespaces ID under this directory.

By the command
`ls -l /proc/$$/ns`
we can see all the namespaces' ID.

```bash
lrwxrwxrwx 1 mengjial mengjial 0 Dec  4 15:12 ipc -> ipc:[4026531839]
lrwxrwxrwx 1 mengjial mengjial 0 Dec  4 15:12 mnt -> mnt:[4026531840]
lrwxrwxrwx 1 mengjial mengjial 0 Dec  4 15:12 net -> net:[4026531957]
lrwxrwxrwx 1 mengjial mengjial 0 Dec  4 15:12 pid -> pid:[4026531836]
lrwxrwxrwx 1 mengjial mengjial 0 Dec  4 15:12 user -> user:[4026531837]
lrwxrwxrwx 1 mengjial mengjial 0 Dec  4 15:12 uts -> uts:[4026531838]
```

if two processes are pointing to same namespace ID, that means they are under the same namespace.

Note: if the link file under this directory opened and the fd still is existing, even all the processes already terminated inside this namespace, the namespace won't be deleted.

### using `setns()` to join in a existed namespace

Using `setns()`, the process will adding an existed namespace into original namespace. In order to enable the new pid namespace take effect, it will usually call `clone()` after `setns()` to create a new sub-preocess to continue and terminate the previous process.

`int setns(int fd, int nstype)`

- `fd`: represents the file descriptor for the adding namespace.
- `nstype`: lets caller to check whether the type of namespace of fd meets the actual situation. 0 for no checking.

### using `unshare()` to isolate namespace on an existing process

`unshare()` is very similar as the `clone()`, the difference is that it doesn't need to start a new process.

Also, compared with `setns()`, `unshare()` doesn't need to link to a existing namespace, only need to specify the namespace need to be isolated. It will automatically create a new namespace.

`int unshare(int flags)`

Since docker doesn't use this API, will ignore the details about this API.

Note: The caller process for `unshare()` and `setns()` will not trap into the new PID namespace, only the child process created after will enter into the new PID namespace. This is because `getpid()` returns the PID depends on the caller's PID namespace. If caller traps into the new PID namespace, its PID will change. For the program running in user mode, they are thinking PID should be a constant value, if PID changed will make them crash.

## 1. UTS Namespace
UTS(UNIX Time-sharing System) namespaces provides the isolation for host name and domain name. So that each single docker container will own its own host name and domain name, it will be treated as a single node instead of a process in host machine from network perspective.

We will have an example shows how it works.

 ```c
#define _GNU_SOURCE
#include <sys/types.h>
#include <sys/wait.h>
#include <stdio.h>
#include <sched.h>
#include <signal.h>
#include <unistd.h>

#define STACK_SIZE (1024 * 1024)

static char child_stack[STACK_SIZE];
char* const child_args[] = {
    "/bin/bash",
    NULL
};

int child_main(void* args) {
    printf("I'm in the child process!\n");
    execv(child_args[0], child_args);
    return 1;
}

int main() {
    printf("Program starts now: \n");
    int child_pid = clone(child_main, child_stack + STACK_SIZE, SIGCHLD, NULL);
    waitpid(child_pid, NULL, 0);
    printf("exited\n");
    return 0;
}
```

Compile and run the above program, it will shows like
```bash
gcc -Wall uts.c -o uts.o && ./uts.o
Program starts now: 
I'm in the child process!
mengjial@ubuntu:/tmp$ exit
exit
exited
mengjial@ubuntu:/tmp
```
From the output we can see that when we run bash inside the sub process, we are still under the same user and same namespace

If we try to add an UTS isolation like the following:
```c
int child_main(void* args) {
    printf("I'm in the child process!\n");
    sethostname("NewNamespace", 12);
    execv(child_args[0], child_args);
    return 1;
}

int main() {
    //[...]
    int child_pid = clone(child_main, child_stack + STACK_SIZE, CLONE_NEWUTS | SIGCHLD, NULL);
    //[...]
}
```

Re-compile and run:
```bash
mengjial@ubuntu  ~/tiny_docker  gcc -Wall uts.c -o uts.o && sudo ./uts.o
Program starts now: 
I'm in the child process!
root@NewNamespace:~/tiny_docker# exit
exit
exited
mengjial@ubuntu  ~/tiny_docker 
```
We can see the hostname inside the sub process changed.
If we don't add `CLONE_NEWUTS` here, we can still see the hostname changed inside the sub process, but the actual current hostname is changed by the sub process.
We can use `hostname` command to check.

## 2. IPC Namespace
IPC(Inter-Process Communication) namespaces provides isolation for message queue, shared memory and other IPC resources. Processes in same IPC namespace can see each other, while processes in different IPC namespaces cannot.

If we add `CLONE_NEWIPC` flag when clone, like
```c
int child_pid = clone(child_main, child_stack + STACK_SIZE, CLONE_NEWIPC | CLONE_NEWUTS | SIGCHLD, NULL);
```

We can use `ipcmk -Q` to create a message queue, and use `ipcs -q` to check opened message queue.

If we create the message queue in main process, and check the message queue inside the child process, you will not find the message queue inside the parent.

```bash
mengjial@ubuntu  ~/tiny_docker  ipcmk -Q
Message queue id: 0
 mengjial@ubuntu  ~/tiny_docker  ipcs -q

------ Message Queues --------
key        msqid      owner      perms      used-bytes   messages    
0x29481ec1 0          mengjial   644        0            0           

mengjial@ubuntu  ~/tiny_docker  gcc -Wall ipc.c -o ipc.o && sudo ./ipc.o
Program starts now: 
I'm in the child process!
root@NewNamespace:~/tiny_docker# ipcs -q

------ Message Queues --------
key        msqid      owner      perms      used-bytes   messages    

root@NewNamespace:~/tiny_docker# exit
exit
exited

```

Docker uses IPC namespace to achieve IPC isolation between container with host, and container with container.

## 3. PID Namespace
PID namespace will re-number process, which means two processes can have same PID in different namespaces.

Kernel maintains a tree structure to hold all namespaces. Root of the tree is created when the system starting, which is called namespace. The namespace created by root namespace is call child namespace, and will also be the child node of the root.

In this way, parent namespace can have the view of child process, and can affect child process by signal. But child namespace cannot have any effect on parent.

- Each initial process in namespace (PID=1) will just like init process in Linus, has special capability.
- Processes in a namespace, cannot kill the process in its parent or sibling namespace.
- Root namespace has the view for all the processes.

So if you want to monitor the program running inside Docker, you can monitor all the processes under Docker daemon's PID namespace.

```
//[...]
int child_pid = clone(child_main, child_stack + STACK_SIZE, CLONE_NEWPID | CLONE_NEWIPC | CLONE_NEWUTS | SIGCHLD, NULL);
//[...]
```

```bash
mengjial@ubuntu  ~/tiny_docker  gcc -Wall pid.c -o pid.o && sudo ./pid.o
[sudo] password for mengjial: 
Program starts now: 
I'm in the child process!
root@NewNamespace:~/tiny_docker# echo $$
1
root@NewNamespace:~/tiny_docker# exit
exit
exited
 mengjial@ubuntu  ~/tiny_docker  echo $$
2414
 mengjial@ubuntu  ~/tiny_docker  
```
From the log we can find that the shell PID inside child process changed to 1.

But if we use `ps aux/top` commands inside the child process, we still can see all PIDs for parent process. That is because we didn't isolate the file system mount point. `ps aux/top` are actually calling actually file content under /proc. So compared with other namespace, in order to make container secure and stable, PID namespace needs more works as following.

### Init Process in PID Namespace
In the tradition UNIX operating system, process with PID=1 is the init process. It is the parent for all process, and it will maintain a process table to check them periodically. Once a child process becomes orphaned, init process will adopt it until it gets killed.

In hence, if it will run multiple process in docker container, the very first command process to run should be something can manage resources like monitoring and recycle, like bash.

### Signal and Init Process
Kernel offers another privilege to init process - signal shielding. 

If there's no logic code for how to handle a signal inside init process, all this signal sent to init process from processes inside same namespace will be blocked even from r oot permission process. This feature is helping on preventing init process killed by mistake.

If process in ancestor PID namespace sends a signal to child namespace's init process, if the signal is not `SIGKILL` or `SIGSTOP`, it will be blocked also if there's no handler for it. Once the init process gets killed, all processes inside same PID namespace will receives `SIGKILL` signal to be destroyed. And ideally the namespace gets destroyed together. But if `/proc/[pid]/ns/pid` is in open state, namespace will be kept. But the kept namespace cannot create new process.

### Mount proc File System
As what talked above, if you want to only have the view of processes in current namespace, you need to re-mount procfs
```bash
root@NewNamespace:~/tiny_docker# mount -t proc proc /proc
root@NewNamespace:~/tiny_docker# ps a
   PID TTY      STAT   TIME COMMAND
     1 pts/8    S      0:00 /bin/bash
    15 pts/8    R+     0:00 ps a
root@NewNamespace:~/tiny_docker# 
```

Note: Since we are not having mount namespace isolation yet, this operation actually affected root namespace's file system. If exit from child process, and `ps a` again, there will show error now. `mount -t proc proc /proc` can restore it.

### `unshare()` and `setns()`
Using `unshare()` and `setns()` for PID namespace will need more attentions. 

`unshare()` allows user to create new namespace in original namespace to do the isolation. But after create the new PID namespace, the caller process of `unshare()` won't get into new namespace. The following new child process will get into the new namespace, and it will become init process in new namespace. Similar process happens for `setns()`.

That's why `docker exec` uses `setns()` to join existing namespace, but will eventually call `clone()`.

## 4. Mount Namespace
Mount namespaces provide isolation of the list of mount points seen by the processes in each namespace instance.  Thus, the processes in each of the mount namespace instances will see distinct single-directory hierarchies.

When a process tries to create a new mount namespace, the current mount point list will be copied to the child namespace. 

The isolation provided by mount namespace sometimes is too great. For example, in order to make a newly loaded optical disk available in all mount namespaces, a mount operation was required in each namespace. For this use case, a feature called `mount propagation` introduced.

- `MS_SHARED`. This mount point shares events with members of a peer group. Mount and unmount events immediately under this mount point will propagate to other mount points that are members of the peer group.
- `MS_PRIVATE`. This mount point is private; it does not have a peer group. Mount and unmount events do not propagate into or out of this mount point.
- `MS_SLAVE`. Mount and unmount events propagate into this mount point from a master shared peer group. Mount and unmount events under this mount point do not propagate to any peer. Not that a mount point can be the slave of another peer group while at the same time sharing mount and unmount events with a peer group of which it is a member.
- `MS_UNBINDABLE`. This is like a private mount, and in addition this mount can't be bind mounted. Attempts to bind mount this mount with the `MS_BIND` flag will fail.

## 5. Network Namespace
Network namespaces provide isolation of the system resources associated with networking: network devices, IPv4 and IPv6 protocol stacks, IP routing tables, firewall rules, the /proc/net directory (which is a symbolic link to /proc/PID/net), the /sys/class/net directory, various files under /proc/sys/net, port numbers (sockets), and so on.  In addition, network namespaces isolate the UNIX domain abstract socket namespace (see unix(7)).

A physical network device can live in exactly one network namespace. When a network namespace is freed (i.e., when the last process in the namespace terminates), its physical network devices are moved back to the initial network namespace (not to the parent of the process).

A virtual network (veth(4)) device pair provides a pipe-like abstraction that can be used to create tunnels between network namespaces, and can be used to create a bridge to a physical network device in another namespace.  When a namespace is freed, the veth devices that it contains are destroyed.

## 6. User Namespace
User namespace is the last namespace kernel supports. As it is not fully stable, some linux version doesn't enable USER_NS when compiling the kernel.

It mainly isolates security related identifier and attribute, including user ID, user group ID, root directory, key and special permissions. In a general way, a normal user's process can create a new process which can has different user or user groups in the new user namespace.

In Linux, root user's ID is 0, we will have a demo that a process whose user ID is not 0 becomes to 0 after create a new user namespace.

```c
//[...]
#include <sys/capability.h>
//[...]

int child_main(void* args) {
    printf("I'm in the child process!\n");
    cap_t caps;
    printf("eUID = %ld; eGID = %ld; ", (long) geteuid(), (long) getegid());
    caps = cap_get_proc();
    printf("capabilities: %s\n", cap_to_text(caps, NULL));
    execv(child_args[0], child_args);
    return 1;
}

//[...]
int child_pid = clone(child_main, child_stack + STACK_SIZE, CLONE_NEWUSER | SIGCHLD, NULL);
```

```bash
mengjial@ubuntu  ~/tiny_docker  gcc -Wall userns.c -lcap -o userns.o && sudo ./userns.o
Program starts now: 
I'm in the child process!
eUID = 65534; eGID = 65534; capabilities: = cap_chown,cap_dac_override,cap_dac_read_search,cap_fowner,cap_fsetid,cap_kill,cap_setgid,cap_setuid,cap_setpcap,cap_linux_immutable,cap_net_bind_service,cap_net_broadcast,cap_net_admin,cap_net_raw,cap_ipc_lock,cap_ipc_owner,cap_sys_module,cap_sys_rawio,cap_sys_chroot,cap_sys_ptrace,cap_sys_pacct,cap_sys_admin,cap_sys_boot,cap_sys_nice,cap_sys_resource,cap_sys_time,cap_sys_tty_config,cap_mknod,cap_lease,cap_audit_write,cap_audit_control,cap_setfcap,cap_mac_override,cap_mac_admin,cap_syslog,cap_wake_alarm,cap_block_suspend,37+ep
nobody@ubuntu:~/tiny_docker$ exit
exit
exited

mengjial@ubuntu  ~/tiny_docker  id -u  
1000
 mengjial@ubuntu  ~/tiny_docker  id -g
1000

```

From the output, we can get these knowledge
- After a new user namespace created, the init process is granted all permissions inside this namespace.
- The UID and GID are different inside new namespace compare with outside.
- User namespace also maintains a tree structure just like pid namespace.

Here is the example how to mapping users.
```c
void set_uid_map(pid_t pid, int inside_id, int outside_id, int length) {
  char path[256];
  sprintf(path, "/proc/%d/uid_map", getpid());
  FILE* uid_map = fopen(path, "w");
  fprintf(uid_map, "%d %d %d", inside_id, outside_id, length);
  fclose(uid_map);
}

void set_gid_map(pid_t pid, int inside_id, int outside_id, int length) {
  char path[256];
  sprintf(path, "/proc/%d/gid_map", getpid());
  FILE* gid_map = fopen(path, "w");
  fprintf(gid_map, "%d %d %d", inside_id, outside_id, length);
  fclose(gid_map);
}

int child_main(void* args) {
    printf("I'm in the child process!\n");
    cap_t caps;
    set_uid_map(getpid(), 0, 1000, 1);
    set_gid_map(getpid(), 0, 1000, 1);
    printf("eUID = %ld; eGID = %ld; ", (long) geteuid(), (long) getegid());
    caps = cap_get_proc();
    printf("capabilities: %s\n", cap_to_text(caps, NULL));
    execv(child_args[0], child_args);
    return 1;
}
```

```bash
mengjial@ubuntu  ~/tiny_docker  gcc -Wall userns.c -lcap -o userns.o && ./userns.o
Program starts now: 
I'm in the child process!
eUID = 0; eGID = 65534; capabilities: = cap_chown,cap_dac_override,cap_dac_read_search,cap_fowner,cap_fsetid,cap_kill,cap_setgid,cap_setuid,cap_setpcap,cap_linux_immutable,cap_net_bind_service,cap_net_broadcast,cap_net_admin,cap_net_raw,cap_ipc_lock,cap_ipc_owner,cap_sys_module,cap_sys_rawio,cap_sys_chroot,cap_sys_ptrace,cap_sys_pacct,cap_sys_admin,cap_sys_boot,cap_sys_nice,cap_sys_resource,cap_sys_time,cap_sys_tty_config,cap_mknod,cap_lease,cap_audit_write,cap_audit_control,cap_setfcap,cap_mac_override,cap_mac_admin,cap_syslog,cap_wake_alarm,cap_block_suspend,37+ep
root@ubuntu:~/tiny_docker# exit
exit
exited
```

We can see we changed the user in new user namespace as a root, and outside that namespace it is just user with uid=1000.

## Reference and Recommend Blog
1. 《容器与容器云》
2. [Linux Namespaces](https://medium.com/@teddyking/linux-namespaces-850489d3ccf)
3. [浅谈Linux Namespace机制（一）](https://zhuanlan.zhihu.com/p/73248894)
4. [pid_namespaces - overview of Linux PID namespaces](http://man7.org/linux/man-pages/man7/pid_namespaces.7.html)
5. [Namespaces in operation, part 3: PID namespaces](https://lwn.net/Articles/531419/)
6. [mount_namespaces - overview of Linux mount namespaces](http://man7.org/linux/man-pages/man7/mount_namespaces.7.html#)
7. [network_namespaces - overview of Linux network namespaces](http://man7.org/linux/man-pages/man7/network_namespaces.7.html)
8. [Linux namespace 简介 part 6 - USER](http://blog.lucode.net/linux/intro-Linux-namespace-6.html)