package model

import (
	"fmt"

	"gorm.io/gorm"
)

type ExecEntry struct {
	gorm.Model
	SubmissionId uint
	Submission   Submission
	Status       ExecutionStatus
	Result       string
	ExitCode     int
	CompilerOpt  string
	TimeLimit    uint
	MemoryLimit  uint
}

type ExecutionStatus int

const (
	Running ExecutionStatus = iota + 1
	Finished
	Interrupted
	Failed
)

func (e ExecutionStatus) String() string {
	s := [...]string{"running", "finished", "interrupted", "failed"}
	if e < Running || e > Failed {
		return fmt.Sprintf("ExecutionStatus(%d)", int(e))
	}
	return s[(e)-1]
}

func (e ExecutionStatus) IsValid() bool {
	switch e {
	case Running, Finished, Interrupted, Failed:
		return true
	}
	return false
}
