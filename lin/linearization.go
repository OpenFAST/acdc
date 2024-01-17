package lin

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
)

type FileGroup struct {
	Name  string
	Files []string
}

type OPResults struct {
	Name   string        `json:"Name"`
	MBC    *MBC          `json:"MBC"`
	EigRes *EigenResults `json:"EigRes"`
	err    error         `json:"-"`
}

func ProcessFiles(LinFiles []string) ([]OPResults, error) {

	// Organize linearization files by operating point
	linFileGroups := map[string][]string{}
	for _, filePath := range LinFiles {
		tmp := strings.Split(filePath, ".")
		filePathNoLinExt := strings.Join(tmp[:len(tmp)-2], ".")
		fileSlice := linFileGroups[filePathNoLinExt]
		fileSlice = append(fileSlice, filePath)
		linFileGroups[filePathNoLinExt] = fileSlice
	}

	linFileGroupChan := make(chan FileGroup, len(linFileGroups))
	linFileResultChan := make(chan OPResults, len(linFileGroups))

	// Launch workers
	wg := &sync.WaitGroup{}
	go func() {
		for i := 0; i < min(len(linFileGroups), 1+2*runtime.NumCPU()/3); i++ {
			wg.Add(1)
			go linFileWorker(wg, linFileGroupChan, linFileResultChan)
		}
		wg.Wait()
		close(linFileResultChan)
	}()

	// Collect results
	groupResults := []OPResults{}

	// Loop through linearization files sent on channel
	for groupName, files := range linFileGroups {

		linFileGroup := FileGroup{Name: groupName, Files: files}

		if len(linFileGroup.Files) == 0 {
			continue
		}

		// Read linearization files
		linFileData := make([]*LinData, len(linFileGroup.Files))
		var err error
		for i, linFilePath := range linFileGroup.Files {
			linFileData[i], err = ReadLinFile(linFilePath)
			if err != nil {
				return nil, fmt.Errorf("error reading '%s': %w", linFilePath, err)
			}
		}

		// Extract matrix data from linearization file data
		matData := NewMatData(linFileData)

		// Perform multi-blade coordinate transform
		mbc, err := matData.MBC3()
		if err != nil {
			return nil, err
		}

		// Perform Eigenanalysis to get modes
		eigRes, err := mbc.EigenAnalysis()
		if err != nil {
			return nil, err
		}

		// Write MBC data to file
		bs, err := json.MarshalIndent(mbc, "", "\t")
		if err != nil {
			return nil, err
		}
		err = os.WriteFile(linFileGroup.Name+"_mbc.json", bs, 0777)
		if err != nil {
			return nil, err
		}

		// Write Eigen analysis results data to file
		bs, err = json.MarshalIndent(eigRes.Modes, "", "\t")
		if err != nil {
			return nil, err
		}
		err = os.WriteFile(linFileGroup.Name+"_modes.json", bs, 0777)
		if err != nil {
			return nil, err
		}

		// Send MBC and mode results
		groupResults = append(groupResults, OPResults{
			Name:   linFileGroup.Name,
			MBC:    mbc,
			EigRes: eigRes,
		})
	}

	// Sort results by group name
	sort.Slice(groupResults, func(i, j int) bool {
		return groupResults[i].Name < groupResults[j].Name
	})

	// Add linearization file data to results
	for i, gr := range groupResults {

		// Set operating point identifier for modes
		for j := range gr.EigRes.Modes {
			gr.EigRes.Modes[j].OP = i
		}

	}

	return groupResults, nil
}

func linFileWorker(wg *sync.WaitGroup, linFilesChan <-chan FileGroup, resultsChan chan<- OPResults) {
	defer wg.Done()

	// Loop through linearization files sent on channel
	for linFileGroup := range linFilesChan {

		if len(linFileGroup.Files) == 0 {
			continue
		}

		// Read linearization files
		linFileData := make([]*LinData, len(linFileGroup.Files))
		var err error
		for i, linFilePath := range linFileGroup.Files {
			linFileData[i], err = ReadLinFile(linFilePath)
			if err != nil {
				resultsChan <- OPResults{err: fmt.Errorf("error reading '%s': %w", linFilePath, err)}
				return
			}
		}

		// Extract matrix data from linearization file data
		matData := NewMatData(linFileData)

		// Perform multi-blade coordinate transform
		mbc, err := matData.MBC3()
		if err != nil {
			resultsChan <- OPResults{err: err}
			return
		}

		// Perform Eigenanalysis to get modes
		eigRes, err := mbc.EigenAnalysis()
		if err != nil {
			resultsChan <- OPResults{err: err}
			return
		}

		// Write MBC data to file
		bs, err := json.MarshalIndent(mbc, "", "\t")
		if err != nil {
			resultsChan <- OPResults{err: err}
			return
		}
		err = os.WriteFile(linFileGroup.Name+"_mbc.json", bs, 0777)
		if err != nil {
			resultsChan <- OPResults{err: err}
			return
		}

		// Write Eigen analysis results data to file
		bs, err = json.MarshalIndent(eigRes.Modes, "", "\t")
		if err != nil {
			resultsChan <- OPResults{err: err}
			return
		}
		err = os.WriteFile(linFileGroup.Name+"_modes.json", bs, 0777)
		if err != nil {
			resultsChan <- OPResults{err: err}
			return
		}

		// Send MBC and mode results
		resultsChan <- OPResults{
			Name:   linFileGroup.Name,
			MBC:    mbc,
			EigRes: eigRes,
		}
	}
}
