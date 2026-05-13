package endpoint

import "regexp"

// nextjsAppRouterMethods lists the HTTP methods exported from App Router route files.
var nextjsAppRouterMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}

// nextjsExportPattern matches named exports like: export async function GET(
var nextjsExportPattern = regexp.MustCompile(
	`export\s+(?:async\s+)?function\s+(GET|POST|PUT|PATCH|DELETE|HEAD|OPTIONS)\s*\(`,
)

// nextjsDefaultExportPattern matches: export default function or export default async function
var nextjsDefaultExportPattern = regexp.MustCompile(
	`export\s+default\s+(?:async\s+)?function\s*(\w*)`,
)
