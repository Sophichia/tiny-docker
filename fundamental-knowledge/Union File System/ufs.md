# Union File System

Union file systems implement a union mount and operate by creating layers.

Docker uses union file systems in conjunction with copy-on-write techniques to provide the building blocks for containers, making them very lightweight and fast.

## Reading Materials
1. https://medium.com/@paccattam/drooling-over-docker-2-understanding-union-file-systems-2e9bf204177c
2. https://blog.csdn.net/xftony/article/details/80569777