package appid

const (
	Name = "HoldToType"
	Slug = "holdtotype"

	Exe      = Slug + ".exe"
	SetupExe = Slug + "-setup.exe"
	LogFile  = Slug + ".log"
	Portable = Slug + "-portable.zip"

	InstallDirName = Name
	ShortcutName   = Name + ".lnk"
	RunValue       = Name
	UninstallKey   = `Software\Microsoft\Windows\CurrentVersion\Uninstall\` + Name
	MutexName      = `Local\` + Name + `TrayMutex`

	Repo      = "Vitalii-Yemets/" + Slug
	RepoURL   = "https://github.com/" + Repo
	AuthorURL = "https://github.com/Vitalii-Yemets"
	LatestAPI = "https://api.github.com/repos/" + Repo + "/releases/latest"

	LLMAlias = Slug + "-local"

	PrevSlug = "voxterminal"
)

var Version = "0.19.0"

func Class(suffix string) string { return Name + suffix }

func TempDirName(kind string, pid int) string {
	return Slug + "-" + kind + "-" + itoa(pid)
}

func TempDirPrefix(kind string) string { return Slug + "-" + kind }

func PrevTempPrefix(kind string) string { return PrevSlug + "-" + kind }

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}
