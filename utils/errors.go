package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/ztrue/tracerr"
)

func LogSuccess(success interface{}) {
	log.Println("🟩: " + fmt.Sprint(success))
}

func LogWarning(warning interface{}) {
	log.Println("🟨: " + fmt.Sprint(warning))
}

func LogError(err interface{}) {
	log.Println("🟥: " + fmt.Sprint(err))
}

func StacktraceError(err error) {
	err = tracerr.Wrap(err)
	tracerr.Print(err)
}

func StacktraceErrorAndExit(err error) {
	err = tracerr.Wrap(err)
	tracerr.Print(err)
	os.Exit(1)
}
