package types

var InvalidPrefixes = []string{
	"/System/",
	"/usr/libexec/",
	"/Library/",
	"/Applications/",
}

var DevRuntimes = []string{
	"python",
	"node",
	"bun",
	"ruby",
	"vite",
	"next",
	"uvicorn",
	"flask",
	"fastapi",
}

var DevPorts = map[uint16]bool{
	3000: true,
	8080: true,
	5173: true,
	4200: true,
	5000: true,
	8000: true,
	4000: true,
}
