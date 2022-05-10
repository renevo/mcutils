package java

// Download links taken from Adoptium
// https://adoptium.net/

type VersionMap map[int]map[string]map[string]string

var (
	// VersionPaths of the installed JRE versions
	VersionPaths = map[int]string{
		8:  "java/8/jdk8u332-b09/",
		16: "java/16/jdk-16.0.2+7/",
		17: "java/17/jdk-17.0.3+7/",
	}

	// Versions of Java supported for download by version:runtime.GOOS:runtime.GOARCH
	Versions = VersionMap{
		8: {
			"linux": {
				"amd64":   "https://github.com/adoptium/temurin8-binaries/releases/download/jdk8u332-b09/OpenJDK8U-jdk_x64_linux_hotspot_8u332b09.tar.gz",
				"aarch64": "https://github.com/adoptium/temurin8-binaries/releases/download/jdk8u332-b09/OpenJDK8U-jdk_aarch64_linux_hotspot_8u332b09.tar.gz",
			},
			"darwin": {
				"amd64": "https://github.com/adoptium/temurin8-binaries/releases/download/jdk8u332-b09/OpenJDK8U-jdk_x64_mac_hotspot_8u332b09.tar.gz",
			},
			"windows": {
				"amd64": "https://github.com/adoptium/temurin8-binaries/releases/download/jdk8u332-b09/OpenJDK8U-jdk_x64_windows_hotspot_8u332b09.zip",
			},
		},
		16: {
			"linux": {
				"amd64":   "https://github.com/adoptium/temurin16-binaries/releases/download/jdk-16.0.2%2B7/OpenJDK16U-jdk_x64_linux_hotspot_16.0.2_7.tar.gz",
				"aarch64": "https://github.com/adoptium/temurin16-binaries/releases/download/jdk-16.0.2%2B7/OpenJDK16U-jdk_aarch64_linux_hotspot_16.0.2_7.tar.gz",
			},
			"darwin": {
				"amd64": "https://github.com/adoptium/temurin16-binaries/releases/download/jdk-16.0.2%2B7/OpenJDK16U-jdk_x64_mac_hotspot_16.0.2_7.tar.gz",
			},
			"windows": {
				"amd64": "https://github.com/adoptium/temurin16-binaries/releases/download/jdk-16.0.2%2B7/OpenJDK16U-jdk_x64_windows_hotspot_16.0.2_7.zip",
			},
		},
		17: {
			"linux": {
				"amd64":   "https://github.com/adoptium/temurin17-binaries/releases/download/jdk-17.0.3%2B7/OpenJDK17U-jdk_x64_linux_hotspot_17.0.3_7.tar.gz",
				"aarch64": "https://github.com/adoptium/temurin17-binaries/releases/download/jdk-17.0.3%2B7/OpenJDK17U-jdk_aarch64_linux_hotspot_17.0.3_7.tar.gz",
			},
			"darwin": {
				"amd64": "https://github.com/adoptium/temurin17-binaries/releases/download/jdk-17.0.3%2B7/OpenJDK17U-jdk_x64_mac_hotspot_17.0.3_7.tar.gz",
			},
			"windows": {
				"amd64": "https://github.com/adoptium/temurin17-binaries/releases/download/jdk-17.0.3%2B7/OpenJDK17U-jdk_x64_windows_hotspot_17.0.3_7.zip",
			},
		},
	}
)
