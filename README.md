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
./bin/blcheck --help
blcheck (0.0.2)- A simple tool to check which links on your websites are broken.

Usage: blcheck <URL>
  -c    Export output as csv format (default if no other format given)
  -csv
        Export output as csv format (default if no other format given)
  -j    Export output as json format
  -json
        Export output as json format
  -max-parallel-requests int
        Maximum number of parallel requests executed (default 5)
  -max-response-timeout int
        Maximum timeout wait on requests in seconds (default 5)
  -mpr int
        Maximum number of parallel requests executed (default 5)
  -mrt int
        Maximum timeout wait on requests in seconds (default 5)
  -v    Displays version of blcheck
  -version
        Displays version of blcheck
```

## Example output*
```shell
./bin/blcheck www.google.com
2024/05/26 20:37:49 Checking URL:  www.google.com
2024/05/26 20:37:49 infered https:// prefix, because given url did not have an protocol
Started: 2024-05-26 20:37:49.796592218 +0200 CEST m=+0.487706897 , took: 657.857335ms, urlcount: 12
Meta information:
initial_parsing_duration: 487.463374ms
#1      true    OK      https://www.google.com/imghp      1
#2      true    OK      https://www.google.com/setprefdomain        1
#3      true    OK      https://www.youtube.com 1
#4      true    OK      https://news.google.com 1
#5      true    OK      https://maps.google.de/maps        1
#6      true    OK      https://play.google.com   1
#7      true    OK      https://accounts.google.com/servicelogin  1
#8      false   Not Found       http://schema.org/webpage       1
#9      true    OK      http://www.google.de/history/optout       1
#10     true    OK      https://www.google.de/intl/de/about/products     1
#11     true    OK      https://mail.google.com/mail    1
#12     true    OK      https://drive.google.com        1
```
(*) removed get parameter outputs from displayed urls

## Reliability
Make sure to only use this on your own websites or websites you have permission to check. This tool is not meant to be used for malicious purposes. It tries to use sane and safe defaults, but you should always be careful when running tools like this. **Use at own risk!**