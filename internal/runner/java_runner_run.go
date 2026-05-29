//ff:func feature=runner type=implementation control=sequence
//ff:what Executes Java tests for the given test file via Maven or Gradle, branching on build tool
package runner

import (
	"fmt"
	"os/exec"
)

// E2E note: running this path requires a working JDK plus Maven or Gradle,
// which is not available in CI/sandbox environments. The build-tool detection
// and argument builders are split out (detectJavaBuildTool, buildJavaTestArgs)
// for environment-independent unit testing.

// Run executes the given Java test file against the project. It detects the
// build tool from project markers, locates the tool binary (preferring a
// project wrapper), and runs only the single test class.
func (r *JavaRunner) Run(projectRoot, testFile string) (*Result, error) {
	buildTool := detectJavaBuildTool(projectRoot)
	if buildTool == "" {
		return nil, fmt.Errorf("no java build tool detected: expected pom.xml (Maven) or build.gradle(.kts) (Gradle) in %s", projectRoot)
	}

	bin, err := findJavaTool(projectRoot, buildTool)
	if err != nil {
		return nil, err
	}

	args := buildJavaTestArgs(buildTool, javaTestClass(testFile))
	cmd := exec.Command(bin, args...)
	cmd.Dir = projectRoot

	output, err := cmd.CombinedOutput()
	result := &Result{
		Output: string(output),
		Pass:   err == nil,
	}
	return result, nil
}
