//ff:func feature=runner type=helper control=selection
//ff:what Locates the build-tool binary (maven/gradle), preferring a project wrapper, else PATH
package runner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// findJavaTool resolves the executable for the given build tool ("maven" or
// "gradle") in projectRoot. A project wrapper (mvnw/gradlew) is preferred when
// present; otherwise the tool is looked up on PATH. A clear error is returned
// when neither a wrapper nor a PATH binary is available.
func findJavaTool(projectRoot, buildTool string) (string, error) {
	var wrapper, cmd string
	switch buildTool {
	case "maven":
		wrapper, cmd = "mvnw", "mvn"
	case "gradle":
		wrapper, cmd = "gradlew", "gradle"
	default:
		return "", fmt.Errorf("unknown java build tool %q", buildTool)
	}

	wrapperPath := filepath.Join(projectRoot, wrapper)
	if _, err := os.Stat(wrapperPath); err == nil {
		return wrapperPath, nil
	}

	path, err := exec.LookPath(cmd)
	if err != nil {
		return "", fmt.Errorf("%s not found on PATH: install the JDK and %s, or add a %s wrapper", cmd, buildTool, wrapper)
	}
	return path, nil
}
