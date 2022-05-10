package minecraft

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

type Properties map[string]string

// WriteProperties will merge the configuration properties into the server.properties for the instance
func (s *Server) WriteProperties() error {
	filePath := filepath.Join(s.Path, "server.properties")

	existing := Properties{}
	if err := existing.Open(filePath); err != nil {
		return err
	}

	// merge from hcl file
	for k, v := range s.Properties {
		existing[k] = v
	}

	return existing.Save(filePath)
}

// Open a pre-existing properties file and load it
func (p Properties) Open(file string) error {
	propertiesFile, err := os.Open(file)
	if err != nil {
		if !os.IsNotExist(err) {
			return errors.Wrapf(err, "failed to open properties file %q", file)
		}

		return nil
	}

	defer propertiesFile.Close()

	scanner := bufio.NewScanner(propertiesFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// comment
		if strings.HasPrefix(line, "#") {
			continue
		}

		// empty or not valid value
		kv := strings.SplitN(line, "=", 2)
		if len(kv) != 2 {
			continue
		}

		// set the value
		p[kv[0]] = kv[1]
	}

	return nil
}

// Save the properties to disk
func (p Properties) Save(path string) error {
	sb := strings.Builder{}
	keys := make([]string, len(p))
	current := 0
	for k := range p {
		keys[current] = k
		current++
	}

	sort.Strings(keys)

	for _, k := range keys {
		if _, err := sb.WriteString(fmt.Sprintf("%s=%s\n", k, p[k])); err != nil {
			return errors.Wrapf(err, "failed to write key value %q", k)
		}
	}

	return errors.Wrapf(os.WriteFile(path, []byte(sb.String()), 0644), "failed to write properties to file %q", path)
}

/*
https://minecraft.fandom.com/wiki/Server.properties

motd=A Minecraft Server

#World
level-name=world
level-type=default
level-seed=
max-world-size=29999984

allow-nether=true
generate-structures=true
##generator-settings=

#GamePlay
difficulty=easy
gamemode=survival
force-gamemode=false
hardcore=false
pvp=true
allow-flight=false

#GamePlay->Spawning
spawn-animals=true
spawn-monsters=true
spawn-npcs=true
spawn-protection=16



#Minecraft 1.17 Snapshot server properties
#Sun May 16 14:13:20 PDT 2021


broadcast-console-to-ops=true
broadcast-rcon-to-ops=true

enable-command-block=false
enable-jmx-monitoring=false
enable-query=false
enable-rcon=false
enable-status=true
enforce-whitelist=false
entity-broadcast-range-percentage=100

function-permission-level=2







max-players=20
max-tick-time=60000


network-compression-threshold=256
online-mode=true
op-permission-level=4
player-idle-timeout=0
prevent-proxy-connections=false

query.port=25565
rate-limit=0
rcon.password=
rcon.port=25575
require-resource-pack=false
resource-pack-prompt=
resource-pack-sha1=
resource-pack=
server-ip=
server-port=25565
snooper-enabled=true


sync-chunk-writes=true
text-filtering-config=
use-native-transport=true
view-distance=10
white-list=false
*/
