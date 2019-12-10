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

If there's no logic code for how to handle a signal inside init process, all this signal sent to init process from processes inside same namespace will be blocked even from root permission process. This feature is helping on preventing init process killed by mistake.

If process in parent PID namespace sends a signal to child namespace's init process, if the signal is not `SIGKILL` or `SIGSTOP`, it will be blocked also. Once the init process gets killed, all processes inside same PID namespace will receives `SIGKILL` signal to be destroyed. And ideally the namespace gets destroyed together. But if `/proc/[pid]/ns/pid` is in open state, namespace will be kept. But the kept namespace cannot create new process.

### Mount proc File System


## Reference
1. 《自己动手写Docker》
2. [Linux Namespaces](https://medium.com/@teddyking/linux-namespaces-850489d3ccf)
3. [浅谈Linux Namespace机制（一）](https://zhuanlan.zhihu.com/p/73248894)