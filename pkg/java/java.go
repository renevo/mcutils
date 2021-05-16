package java

// Download links taken from AdoptOpenJDK
// https://adoptopenjdk.net/archive.html?variant=openjdk16&jvmVariant=hotspot

type VersionMap map[int]map[string]map[string]string

var (
	// VersionPaths of the installed JRE versions
	VersionPaths = map[int]string{
		8:  "java/8/jdk8u292-b10-jre/",
		16: "java/16/jdk-16.0.1+9-jre/",
	}

	// Versions of Java supported for download by version:runtime.GOOS:runtime.GOARCH
	Versions = VersionMap{
		8: {
			"linux": {
				"amd64":   "https://github.com/AdoptOpenJDK/openjdk8-binaries/releases/download/jdk8u292-b10/OpenJDK8U-jre_x64_linux_hotspot_8u292b10.tar.gz",
				"aarch64": "https://github.com/AdoptOpenJDK/openjdk8-binaries/releases/download/jdk8u292-b10/OpenJDK8U-jre_aarch64_linux_hotspot_8u292b10.tar.gz",
			},
			"darwin": {
				"amd64": "https://github.com/AdoptOpenJDK/openjdk8-binaries/releases/download/jdk8u292-b10/OpenJDK8U-jre_x64_mac_hotspot_8u292b10.tar.gz",
			},
			"windows": {
				"amd64": "https://github.com/AdoptOpenJDK/openjdk8-binaries/releases/download/jdk8u292-b10/OpenJDK8U-jre_x64_windows_hotspot_8u292b10.zip",
			},
		},
		16: {
			"linux": {
				"amd64":   "https://github.com/AdoptOpenJDK/openjdk16-binaries/releases/download/jdk-16.0.1%2B9/OpenJDK16U-jre_x64_linux_hotspot_16.0.1_9.tar.gz",
				"aarch64": "https://github.com/AdoptOpenJDK/openjdk16-binaries/releases/download/jdk-16.0.1%2B9/OpenJDK16U-jre_aarch64_linux_hotspot_16.0.1_9.tar.gz",
			},
			"darwin": {
				"amd64": "https://github.com/AdoptOpenJDK/openjdk16-binaries/releases/download/jdk-16.0.1%2B9/OpenJDK16U-jre_x64_mac_hotspot_16.0.1_9.tar.gz",
			},
			"windows": {
				"amd64": "https://github.com/AdoptOpenJDK/openjdk16-binaries/releases/download/jdk-16.0.1%2B9/OpenJDK16U-jre_x64_windows_hotspot_16.0.1_9.zip",
			},
		},
	}
)
