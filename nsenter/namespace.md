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

`int unshare(int flags)`

Docker doesn't use this API.
