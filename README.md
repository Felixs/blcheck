# blcheck
Broken links check, looking for unavailable links on websites. Written in golang.

## What are we doing here?
After getting back into golang with [learn go with tests (link)](https://quii.gitbook.io/learn-go-with-tests/) i decided to start a project that i came up with durion one of the tasks.
How about a tool that checks every href/link reference on a website/html-page if the content its linking to is still accessable. Or at least returns a 200, content check might be a bit out of scope. 

>So here goes nothing!

[TODOs here](TODO.md)

## Requirements
- go >= 1.22

## Make and use
```shell
git clone https://github.com/Felixs/blcheck
cd blcheck
make build
./bin/blcheck https://www.only-on-pages-own-by-you.con
```

## Usage output
```shell
blcheck (0.0.1)- A simple tool to check which links on your websites are broken.

Usage: blcheck <URL>
  -max-parallel-requests uint
        Setting a maximum how many requests get executed in parallel (default 20)
  -v    Displays version of blcheck
  -version
        Displays version of blcheck
```


## Reliability
Make sure to only use this on your own websites or websites you have permission to check. This tool is not meant to be used for malicious purposes. It tries to use sane and safe defaults, but you should always be careful when running tools like this. **Use at own risk!**