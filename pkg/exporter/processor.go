package exporter

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"k8s.io/klog"
)

// ProcessingPhase is the current phase being processed.
type ProcessingPhase string

// Processor holds the fields needed to process data from a data provider.
type Processor struct {
	// currentPhase is the phase the processing is in currently.
	currentPhase ProcessingPhase
	// destination to export
	destination DataDestinationInterface
	// sourcePath
	sourcePath string
	// expotDir
	exportDir string
	// expotPath
	exportPath string
}

// DataDestinationInterface is the interface all data destination should implement.
type DataDestinationInterface interface {
	// Transfer is called to transfer the image to the Destination
	Transfer() (ProcessingPhase, error)
	// Close closes any senders or other open resources.
	Close()
}

const (
	// ProcessingPhaseConvert is the phase in which Source Data is converted to the qcow2 image format and transferred to the exportDir.
	ProcessingPhaseConvert ProcessingPhase = "Convert"
	// ProcessingPhaseTransfer is the phase in which qcow2 image in exportDir is transferred to the destination
	ProcessingPhaseTransfer ProcessingPhase = "Transfer"
	// ProcessingPhaseComplete is the phase where the entire process completed successfully and we can exit gracefully.
	ProcessingPhaseComplete ProcessingPhase = "Complete"
	// ProcessingPhaseError is the phase in which we encountered an error and need to exit ungracefully.
	ProcessingPhaseError ProcessingPhase = "Error"
)

// NewProcessor create a new instance of a data processor using the passed in data provider.
func NewProcessor(dataDestination DataDestinationInterface, sourcePath, exportDir, exportPath string) *Processor {
	return &Processor{
		currentPhase: ProcessingPhaseConvert,
		destination:  dataDestination,
		sourcePath:   sourcePath,
		exportDir:    exportDir,
		exportPath:   exportPath,
	}
}

// cleanDir cleans the contents of a directory including its sub directories, but does NOT remove thedirectory itself.
func cleanDir(dest string) error {
	dir, err := ioutil.ReadDir(dest)
	if err != nil {
		klog.Errorf("Unable read directory to clean: %s, %v", dest, err)
		return err
	}
	for _, d := range dir {
		klog.V(1).Infoln("deleting file: " + filepath.Join(dest, d.Name()))
		err = os.RemoveAll(filepath.Join(dest, d.Name()))
		if err != nil {
			klog.Errorf("Unable to delete file: %s, %v", filepath.Join(dest, d.Name()), err)
			return err
		}
	}
	return nil
}

// Process is the main synchronous processing loop
func (p *Processor) Process() error {
	if err := cleanDir(p.exportDir); err != nil {
		return errors.Wrap(err, "Failure cleaning up export space")
	}
	return p.processDataWithPause()
}

func (p *Processor) processDataWithPause() error {
	var err error
	for p.currentPhase != ProcessingPhaseComplete {
		switch p.currentPhase {
		case ProcessingPhaseConvert:
			p.currentPhase, err = p.convert()
			if err != nil {
				err = errors.Wrap(err, "Unable to convert source data to Qcow2 format")
			}
		case ProcessingPhaseTransfer:
			p.currentPhase, err = p.destination.Transfer()
			if err != nil {
				err = errors.Wrap(err, "Unable to transfer Qcow2 image to destination")
			}
		default:
			return errors.Errorf("Unknown processing phase %s", p.currentPhase)
		}
		if err != nil {
			klog.Errorf("%+v", err)
			return err
		}
		klog.V(1).Infof("New phase: %s\n", p.currentPhase)
	}
	return err
}

func (p *Processor) convert() (ProcessingPhase, error) {
	klog.V(1).Infoln("Converting to Qcow2")
	err := convertToQcow2(p.sourcePath, p.exportPath)
	if err != nil {
		return ProcessingPhaseError, errors.Wrap(err, "Conversion to Qcow2 failed")
	}

	return ProcessingPhaseTransfer, nil
}
