package analyzer

import "path/filepath"

// detectLanguage determines the programming language from file extension
func detectLanguage(filePath string) string {
	ext := filepath.Ext(filePath)

	languageMap := map[string]string{
		".go":    "go",
		".py":    "python",
		".js":    "javascript",
		".ts":    "typescript",
		".jsx":   "javascript",
		".tsx":   "typescript",
		".java":  "java",
		".rb":    "ruby",
		".rs":    "rust",
		".c":     "c",
		".cpp":   "cpp",
		".cc":    "cpp",
		".cxx":   "cpp",
		".cs":    "csharp",
		".php":   "php",
		".swift": "swift",
		".kt":    "kotlin",
		".scala": "scala",
		".sh":    "bash",
	}

	if lang, ok := languageMap[ext]; ok {
		return lang
	}

	return "unknown"
}
