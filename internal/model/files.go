package model

import (
	"fmt"
	"os"
	"path/filepath"
)

type Files struct {
	ParamsFile   *os.File
	StatesFile   *os.File
	RequestsFile *os.File
}

const (
	ParamsInitString   = "Lambda Seed,Mu Seed,Number of iterations,Number of workers,Lambda,Mu"
	StatesInitString   = "Processing,Waiting,Average response time"
	RequestsInitString = "ArrivedAt,Processed time,Processing started,Processing finished"
)

func NewFiles(directoryName, paramsFileName, statesFileName, requestFileName string) (*Files, error) {

	if directoryName != "" {
		err := os.MkdirAll(directoryName, 0755)
		if err != nil {
			return nil, err
		}
	}

	params, err := os.OpenFile(filepath.Join(directoryName, paramsFileName), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(params, ParamsInitString)

	states, err := os.OpenFile(filepath.Join(directoryName, statesFileName), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(states, StatesInitString)

	requests, err := os.OpenFile(filepath.Join(directoryName, requestFileName), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return nil, err
	}
	fmt.Fprintln(requests, RequestsInitString)

	return &Files{
		ParamsFile:   params,
		StatesFile:   states,
		RequestsFile: requests,
	}, nil
}

func (f *Files) DumpParams(lseed, mseed, iterationNum, workersNum int64, lambda, mu float64) {
	fmt.Fprintf(f.ParamsFile, "%d,%d,%d,%d,%f,%f\n", lseed, mseed, iterationNum, workersNum, lambda, mu)
}

func (f *Files) DumpState(state *State, averateResponseTime float64) {
	fmt.Fprintf(f.StatesFile, "%d,%d,%f\n", state.Processing, state.Waiting, averateResponseTime)
}

func (f *Files) DumpRequests(requests []*Request) {
	for _, request := range requests {
		fmt.Fprintf(f.RequestsFile, "%d,%d,%d,%d\n",
			request.ArrivedAt, request.ProcessedTime, request.ProcessingStarted, request.ProcessingFinished)
	}
}

func (f *Files) Close() {
	f.ParamsFile.Close()
	f.RequestsFile.Close()
	f.StatesFile.Close()
}
