package constants

const (
	// exit code definition
	ExitSuccess                          int = 0
	ExitNotAllReportReachable            int = 1
	ExitMissingParameter                 int = 2
	ExitInvalidUrlParameter              int = 3
	ExitUrlNotReachable                  int = 4
	ExitFailedToWriteReport              int = 5
	ExitFailedToCreateReport             int = 6
	ExitErrorInParameterEvaluation       int = 7
	ExitToManyOutputFormats              int = 8
	ExitInvalidNumberMaxParallelRequests int = 9
	ExitInlvaidNumberMaxTimeoutInSeconds int = 10
)
