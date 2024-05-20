# blcheck
Broken links check, looking for unavailable links on websites. Written in golang.

## What are we doing here?
After getting back into golang with [learn go with tests (link)](https://quii.gitbook.io/learn-go-with-tests/) i decided to start a project that i came up with durion one of the tasks.
How about a tool that checks every href/link reference on a website/html-page if the content its linking to is still accessable. Or at least returns a 200, content check might be a bit out of scope. 

>So here goes nothing!

# Tasks
- [x] create a cmd tool, that takes an url as input
- [x] validate url input
- [x] add missing http(s) protocol prefix if missing
- [x] fetch the html content of the url
- [x] parse the html content, collect all unique lowercase hrefs/links
  - [x] remove ancor from links like www.example.com/#about -> www.example.com
- [x] check if the links are still accessable (from current machine)
- [x] create a report with all checked links and their status
- [x] move from sequential url check to parallel url check
- [ ] limit the number of parallel requests (add a programm flag aswell)
- [x] cleanup how the url report is displayed
- [ ] create a presentable output format csv/html
- [ ] add a --help --version flag to print out help and version text

# maybe features for the future
- add a counter how often an unique url appeared
- timeout and non 200 result should be distinguishable
- maybe check of if an anchor is given this anchor is still present on the page
- recursive mode, that checks all links on the same domain as the first given url
- exclude/include regex parameter that can filter which links should be checked
- ~~create an nice csv/html output of link-report~~

