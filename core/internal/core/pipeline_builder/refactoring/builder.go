package refactoring

import (
	"chast.io/core/internal/core/version_control/overlay"
	"chast.io/core/internal/model/run_models/refactoring"
	"log"
)

func BuildRunPipeline(runModel *refactoring.RunModel) {
	log.Printf("Refactoring BuildRunPipeline")
	log.Printf("Run Command: %s", runModel.Run[0].Command.Cmds)

	var nsContext = overlay.NewNamespaceContext(
		"/",                            // This will be defined by the versioning system
		"/tmp/overlay-auto-test/upper", // This will be defined by the versioning system
		"/tmp/overlay-auto-test/operationDirectory", // This will be defined by the versioning system
		runModel.Run[0].Command.WorkingDirectory,
		runModel.Run[0].Command.Cmds[0]..., // TODO support multiple commands
	)

	if err := overlay.RunCommandInIsolatedEnvironment(nsContext); err != nil {
		log.Fatalf("Error running command in isolated environment - %s", err)
	}
}

func buildDockerRunPipeline(runModel *refactoring.Run) {

}
