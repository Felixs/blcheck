# Refactoring Session
- [ ] make sure marshalling erros get properly propergated and maybe add tests (lets see how it works)
- [ ] split the main into several part. Dunno about parsing and execution in differen files, but the main processing routine becomes quite large
- [x] move flag parsing out of main
- [ ] add acceptence test for the main, after splitting into part.

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
- [x] limit the number of parallel requests 
  - [x] add a programm flag aswell
- [x] cleanup how the url report is displayed
- [x] add a --help --version flag to print out help and version text
- [x] fix basic info output for all http status codes and errors
- [x] make http get timeout configurable with flag
- [x] add advances info for UrlResults like response time, content length
- [x] create a presentable string output format, as default output
- [ ] create a presentable output format csv, make accessible via flag
- [x] create a presentable output format json, make accessible via flag
- [x] while parsing the first html content, add a counter to unique urls how often they appear
- [x] add github action to run test on push
- [x] move time to parse and number of urls to UrlReport meta data
- [x] add a flag to exclude/include certain urls with regex (and the functionality ofc)
- [x] add a dry run flag to check how many unique urls are found on a webpage
- [x] write output of blckeck to file with an -o/-out flag
- [x] limit report output to only broken or timed out urls and give an exit code of != 0 if broken links found
- [ ] adda flag to add all checked urls to output report
- [ ] create a presentable output format html, make accessible via flag
- [ ] validate programm flags for sane inputs
- [ ] *experimental* add a flag to use a certain proxy server or maybe dns resolver
- [ ] add CHANGELOG.md by autochangelog
- [ ] add a method to retry timed out requests if wanted (flag)
- [ ] check how http.Head/Get handles redirects and how it can be tested in unit tests
- [ ] serve output html als webserver
- [ ] urls parser need to find relativ links to
- [ ] make urls parser give infor about found link, is it a href, src, relativ link or text search url
- [ ] check out cobra-cli for advanced cli argument parsing (https://www.kosli.com/blog/understanding-golang-command-line-arguments/)

# Fixes
- [ ] Check JSON urls encoding in url strings (e.g. "url": "https://maps.google.de/maps?hl=de\u0026tab=wl",)


# maybe features for the future
- check also urls with anchor and if this anchor is still present on the page
- recursive mode, that checks all links on the same domain as the first given url
- ~~add a counter how often an unique url appeared~~
- ~~exclude/include regex parameter that can filter which links should be checked~~
- ~~create an nice csv/html output of link-report~~
- ~~timeout and non 200 result should be distinguishable~~