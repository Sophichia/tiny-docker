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

## UTS Namespace
UTS(UNIX Time-sharing System) namespaces provides the isolation for host name and domain name. So that each single docker container will own its own host name and domain name, it will be treated as a single node instead of a process in host machine from network perspective.

We will have an example shows how it works.

 ```c
#define _GNU_SOURCE
#incude <sys/types.h>
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
    int child_pod = clone(child_main, child_stack + STACK_SIZE, CLONE_NEWUTS | SIGCHLD, NULL);
    //[...]
}
```

Re-compile and run:
```bash
gcc -Wall uts.c -o uts.o && sudo ./uts.o
sudo: unable to resolve host NewNamespace
Program starts now: 
I'm in the child process!
root@NewNamespace:/tmp# exit
exit
exited!
```
We can see the hostname inside the sub process changed.
If we don't add `CLONE_NEWUTS` here, we can still see the hostname changed inside the sub process, but the actual current hostname is changed by the sub process.
We can use `hostname` command to check.

## Reference
1. 《自己动手写Docker》
2. [Linux Namespaces](https://medium.com/@teddyking/linux-namespaces-850489d3ccf)
3. [浅谈Linux Namespace机制（一）](https://zhuanlan.zhihu.com/p/73248894)