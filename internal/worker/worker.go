package worker

import (
	"fmt"
	"importerapi/internal/services"
	"importerapi/internal/util"
	"os"
	"time"
)

type ImportJob struct {
	FilePath string
	ImportID string
}

type Worker struct {
	Service *services.ImportService
}

func (w *Worker) StartImportWorker(jobQueue <-chan ImportJob, workerID int) {
	for job := range jobQueue {
		fmt.Printf("[Worker %d] Processing import job for file: %s\n", workerID, job.FilePath)
		startTime := time.Now()

		f, err := os.Open(job.FilePath)
		if err != nil {
			fmt.Printf("[Worker %d] Error opening file: %s\n", workerID, err.Error())
			continue
		}
		defer f.Close()

		w.Service.ImportStatusRepo.UpdateStatus(job.ImportID, "processing")
		records, err := util.ReadExcelFromReader(f)
		if err != nil {
			fmt.Printf("[Worker %d] Error reading excel file: %s\n", workerID, err.Error())
			continue
		}

		// err = w.Service.ImportFromXML(records)
		// if err != nil {
		// 	fmt.Printf("[Worker %d] Failed to import: %s\n", workerID, err)
		// 	continue
		// }

		fmt.Printf("[Worker %d] Processed %d records in %s\n", workerID, len(records), time.Since(startTime))
		w.Service.ImportStatusRepo.UpdateStatus(job.ImportID, "completed")

		if err := os.Remove(job.FilePath); err != nil {
			fmt.Printf("[Worker %d] Error removing file: %s\n", workerID, err)
		}
	}
}
