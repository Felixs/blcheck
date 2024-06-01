package arguments

import (
	"flag"
	"fmt"
	"os"

	"github.com/Felixs/blcheck/pkg/url"
)

var (
	// Tool comandline flags
	DisplayVersion bool

	// Report parameter
	URL            string
	RegexInclude   string
	RegexExclude   string
	ShowReachables bool
	ExecuteDryRun  bool

	// Constrains for url checks
	MaxParallelRequests int
	MaxTimeoutInSeconds int

	// Output parameter
	OutputAsJSON bool
	OutputAsCSV  bool
	OutputInFile string

	// Error message on flag errors/missmatch
	ErrorMessage string
)

// Parses the command line arguments and checks if they all needed arguments are present.
func init() {
	// TODO: if using 2 different flag names for a single values flag.Usage needs to be overwritten
	// Version
	flag.BoolVar(&DisplayVersion, "version", false, "Displays version of blcheck")
	flag.BoolVar(&DisplayVersion, "v", false, "Displays version of blcheck")
	// Ratelimit / parallel requests
	flag.IntVar(&MaxParallelRequests, "max-parallel-requests", url.MaxNumParallelQueries, "Maximum number of parallel requests executed")
	flag.IntVar(&MaxParallelRequests, "mpr", url.MaxNumParallelQueries, "Maximum number of parallel requests executed")
	// Timeout
	flag.IntVar(&MaxTimeoutInSeconds, "max-response-timeout", int(url.DefaultHttpGetTimeout.Seconds()), "Maximum timeout wait on requests in seconds")
	flag.IntVar(&MaxTimeoutInSeconds, "mrt", int(url.DefaultHttpGetTimeout.Seconds()), "Maximum timeout wait on requests in seconds")
	// Output as json flag
	flag.BoolVar(&OutputAsJSON, "json", false, "Export output as json format")
	flag.BoolVar(&OutputAsJSON, "j", false, "Export output as json format")
	// Output as csv flag
	flag.BoolVar(&OutputAsCSV, "csv", false, "Export output as csv format")
	flag.BoolVar(&OutputAsCSV, "c", false, "Export output as csv format")
	// Include flag for which string needs to be present in url to check
	flag.StringVar(&RegexInclude, "include", "", "Parsed urls need to contain this string to get checked")
	flag.StringVar(&RegexInclude, "in", "", "Parsed urls need to contain this string to get checked")
	// Exclude flag for which string can not be present in url to check
	flag.StringVar(&RegexExclude, "exclude", "", "Parsed urls need to not contain this string to get checked")
	flag.StringVar(&RegexExclude, "ex", "", "Parsed urls need to not contain this string to get checked")
	// Flag if tool should run in dry mode, only getting links from initial webpage
	flag.BoolVar(&ExecuteDryRun, "dry", false, "Only gets urls from initial webpage and does not check the status of other urls")
	flag.BoolVar(&ExecuteDryRun, "d", false, "Only gets urls from initial webpage and does not check the status of other urls")
	// Flag if output should be writen into file, gives path an name of file
	flag.StringVar(&OutputInFile, "o", "", "Writes output to given location. If directory is given, writes to blcheck.log in directory.")
	flag.StringVar(&OutputInFile, "out", "", "Writes output to given location. If directory is given, writes to blcheck.log in directory.")
	// Flag if reachable urls should be included into the output
	flag.BoolVar(&ShowReachables, "sr", false, "Includes reachable urls in report")
	flag.BoolVar(&ShowReachables, "show-reachable", false, "Includes reachable urls in report")

	// setting own print function, to handle positonal arguments
	flag.Usage = printUsage
	flag.Parse()

	if DisplayVersion {
		fmt.Println("blcheck " + Version + "\n2024 - Felix Sponholz")
		os.Exit(0)
	}
	if flag.NArg() != 1 {
		ErrorMessage = "URL is required"
		printUsage()
		os.Exit(3)
	}
	URL = flag.Arg(0)
}

// Prints how to use the tool to stdout, with an error message if present.
func printUsage() {
	if ErrorMessage != "" {
		fmt.Println("ERROR:" + ErrorMessage)
	}

	fmt.Printf(`blcheck (%s)- A simple tool to check which links on your websites are broken.
	
Usage: blcheck <URL>
`, Version)
	flag.PrintDefaults()
}
