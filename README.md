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
  -c    Export output as csv format (default if no other format given) (default true)
  -csv
        Export output as csv format (default if no other format given) (default true)
  -d    Only gets urls from initial webpage and does not check the status of other urls
  -dry
        Only gets urls from initial webpage and does not check the status of other urls
  -ex string
        Parsed urls need to not contain this string to get checked
  -exclude string
        Parsed urls need to not contain this string to get checked
  -in string
        Parsed urls need to contain this string to get checked
  -include string
        Parsed urls need to contain this string to get checked
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
  -o string
        Writes output to given location. If directory is given, writes to blcheck.log in directory.
  -out string
        Writes output to given location. If directory is given, writes to blcheck.log in directory.
  -v    Displays version of blcheck
  -version
        Displays version of blcheck
```

## Example output*
```shell
./bin/blcheck --show-reachable www.google.com
Infered https:// prefix, because given url did not have a protocol
Checking URL:  https://www.google.com
Started: 2024-06-01T10:57:54+02:00 , took: 872ms, urlcount: 12
Meta information:
        initial_parsing_duration: 285.398457ms
        total_extracted_urls: 12

url                                             is_reachable    status_message  content_length  response_time   num_occured
https://www.google.com/imghp                    true            OK              -1              52.767129ms     1
https://www.google.com/setprefdomain            true            OK              -1              55.834137ms     1
https://maps.google.de/maps                     true            OK              -1              145.421355ms    1
https://accounts.google.com/servicelogin        true            OK              170442          222.108459ms    1
https://www.youtube.com                         true            OK              499519          257.611805ms    1
https://news.google.com                         true            OK              1636600         356.455411ms    1
https://www.google.de/intl/de/about/products    true            OK              235684          249.307575ms    1
http://www.google.de/history/optout             true            OK              246539          375.696233ms    1
http://schema.org/webpage                       false           Not Found       284             438.883645ms    1
https://play.google.com                         true            OK              2686403         697.007397ms    1
https://mail.google.com                         true            OK              170850          439.482341ms    1
https://drive.google.com                        true            OK              170785          423.845873ms    1
```
(*) removed get parameter outputs from displayed urls

## Reliability
Make sure to only use this on your own websites or websites you have permission to check. This tool is not meant to be used for malicious purposes. It tries to use sane and safe defaults, but you should always be careful when running tools like this. **Use at own risk!**