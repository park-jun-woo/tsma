//ff:func feature=coverage type=implementation control=sequence
//ff:what Runs the Java build tool with JaCoCo and computes per-function coverage from jacoco.xml
package coverage

import (
	"fmt"

	"github.com/park-jun-woo/tsma/internal/model"
)

// Check runs the project's tests with JaCoCo (via Maven or Gradle) and maps the
// resulting jacoco.xml onto the given function's line range. It branches on the
// detected build tool and requires a working JDK + build tool + JaCoCo plugin.
func (c *JavaChecker) Check(projectRoot, testFile string, fn *model.Function) (*Report, error) {
	buildTool := detectJavaBuildTool(projectRoot)
	if buildTool == "" {
		return nil, fmt.Errorf("no java build tool detected: expected pom.xml (Maven) or build.gradle(.kts) (Gradle) in %s", projectRoot)
	}

	bin, err := findJavaTool(projectRoot, buildTool)
	if err != nil {
		return nil, err
	}

	args := buildJavaCoverageArgs(buildTool)
	if err := runJavaCoverage(bin, projectRoot, args); err != nil {
		return nil, fmt.Errorf("java coverage failed: %w", err)
	}

	reportPath := jacocoReportPath(projectRoot, buildTool)
	cov, err := parseJacoco(reportPath)
	if err != nil {
		return nil, fmt.Errorf("parse jacoco xml: %w", err)
	}

	ranges := collectJavaRanges(fn)

	return buildJavaReport(ranges, cov, projectRoot), nil
}
