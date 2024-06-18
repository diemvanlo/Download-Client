package configs

type DownloadMode string

const (
	DownloadModeLocal DownloadMode = "local"
	DownloadModeS3    DownloadMode = "s3"
)

type DownloadConfig struct {
	Mode        DownloadMode `yaml:"mode"`
	DownloadDir string       `yaml:"download_dir"`
}
