package refactoring

import (
	"chast.io/core/internal/model/run_models/refactoring"
	"log"
	"os/exec"
)

func BuildRunPipeline(runModel *refactoring.RunModel) {
	log.Printf("Refactoring BuildRunPipeline")
	log.Printf("Run Command: %s", runModel.Run[0].Command.Cmd)

	runModel.Run[0].Command.Cmd[3] = "" // TODO load params from cli arguments
	runModel.Run[0].Command.Cmd[4] = "" // TODO load params from cli arguments

	cmd := exec.Command(runModel.Run[0].Command.Cmd[0], runModel.Run[0].Command.Cmd[1:]...)
	cmd.Dir = runModel.Run[0].Command.WorkingDirectory

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed with %s. \n%s", err, string(out))
	} else {
		log.Printf(string(out))
	}
}

func buildDockerRunPipeline(runModel *refactoring.Run) {

}
