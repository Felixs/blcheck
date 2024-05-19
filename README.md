# blcheck
Broken links check, looking for unavailable links on websites. Written in golang.

## What are we doing here?
After getting back into golang with [learn go with tests (link)](https://quii.gitbook.io/learn-go-with-tests/) i decided to start a project that i came up with durion one of the tasks.
How about a tool that checks every href/link reference on a website/html-page if the content its linking to is still accessable. Or at least returns a 200, content check might be a bit out of scope. 

>So here goes nothing!

# Tasks
- [x] create a cmd tool, that takes an url as input
- [x] validate url input
- [ ] add missing http(s) protocol prefix if missing
- [x] fetch the html content of the url
- [ ] parse the html content, collect all hrefs/links
- [ ] check if the links are still accessable
- [ ] create a report with all checked links and their status

# Features for the future
- recursive mode, that checks all links on the same domain as the first given url
- exclude/include regex parameter that can filter which links should be checked
- create an nice csv/html output of link-report

