package refactoring

import (
	"chast.io/core/internal/model/run_models/refactoring"
	"log"
	"os/exec"
)

func BuildRunPipeline(runModel *refactoring.RunModel) {
	log.Printf("Refactoring BuildRunPipeline")
	log.Printf("Run Command: %s", runModel.Run[0].Command.Cmd)

	runModel.Run[0].Command.Cmd[3] = "/shared/home/rjenni/Projects/mse-repos/master-thesis/chast/chast-refactoring-antlr/base/src/main/java/CSharpLexerBase.java"
	runModel.Run[0].Command.Cmd[4] = "/shared/home/rjenni/Projects/mse-repos/master-thesis/chast/chast-refactoring-antlr/refactorings/rearrange_class_members/src/test/resources/rearrange_class_members/default_config.yaml"

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
