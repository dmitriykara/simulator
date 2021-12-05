package main

import (
	"flag"
	"log"
	"math/rand"

	"github.com/dmitriykara/queuing/internal/model"
)

var (
	lseed            = flag.Int64("lseed", 1, "Lambda random seed")
	mseed            = flag.Int64("mseed", 2, "Mu random seed")
	iterationNum     = flag.Int64("n", 1000, "Number of iteration")
	workersNum       = flag.Int64("w", 5, "Number of workers")
	lambda           = flag.Float64("lambda", 0.2, "Lambda")
	mu               = flag.Float64("mu", 0.1, "Mu")
	noLog            = flag.Bool("noLog", false, "Do not log data")
	directoryName    = flag.String("dir", "results", "Output directory")
	paramsFileName   = flag.String("params", "params.csv", "File name for parameters")
	statesFileName   = flag.String("states", "states.csv", "File name for simulation states")
	requestsFileName = flag.String("requests", "requests.csv", "File name for request")
)

var (
	lambdaGenerator *rand.Rand
	muGenerator     *rand.Rand
	i               int64
	files           *model.Files
	err             error
)

func ArrivalFunc() float64 {
	return lambdaGenerator.ExpFloat64() / *lambda
}

func OperationFunc() float64 {
	return muGenerator.ExpFloat64() / *mu
}

func main() {
	flag.Parse()

	if !*noLog {
		files, err = model.NewFiles(*directoryName, *paramsFileName, *statesFileName, *requestsFileName)
		if err != nil {
			log.Fatalf("could not create files: %v", err)
		}
	}

	if !*noLog {
		files.DumpParams(*lseed, *mseed, *iterationNum, *workersNum, *lambda, *mu)
	}

	lambdaGenerator = rand.New(rand.NewSource(*lseed))
	muGenerator = rand.New(rand.NewSource(*mseed))

	model := model.NewModel(*workersNum, ArrivalFunc, OperationFunc)
	processedRequests := 0
	averageResponseTime := 0.0

	for i = 0; i < *iterationNum; i++ {
		state := model.GetCurrentState()
		for _, responseTime := range state.ResponseTimes {
			averageResponseTime = (float64(processedRequests)*averageResponseTime + float64(responseTime)) /
				float64(processedRequests+1)
			processedRequests++
		}
		if !*noLog {
			files.DumpState(state, averageResponseTime)
		}
	}

	if !*noLog {
		files.DumpRequests(model.Queue.Requests())
	}
}
