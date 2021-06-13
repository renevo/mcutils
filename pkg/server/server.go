package server

import (
	"os"
	"path/filepath"

	"github.com/renevo/mcutils/pkg/minecraft/version"
)

// Server for Minecraft
type Server struct {
	Name     string   `hcl:"name,label" property:"motd"`
	Path     string   `hcl:"path"`
	Version  string   `hcl:"version"`
	Snapshot bool     `hcl:"snapshot,optional"`
	JavaHome string   `hcl:"java_home,optional"`
	World    World    `hcl:"world,block"`
	GamePlay GamePlay `hcl:"game,block"`

	version version.Version
}

// World for Minecraft
type World struct {
	Name             string `hcl:"name,label" property:"level-name"`
	Type             string `hcl:"type,optional" property:"level-type"`
	Seed             string `hcl:"seed,optional" property:"level-seed"`
	Size             int    `hcl:"size,optional" property:"max-world-size"`
	Height           int    `hcl:"height,optional" property:"max-build-height"`
	EnableNether     bool   `hcl:"nether,optional" property:"allow-nether"`
	EnableStructures bool   `hcl:"structures,optional" property:"generate-structures"`

	// Current issue with Generator Settings on a server, will revisit this in the future...
	//Generator       string `hcl:"generator,optional" property:"generator-settings,json"`
}

type GamePlay struct {
	Difficulty string           `hcl:"difficulty,optional" property:"difficulty"`
	Mode       string           `hcl:"mode,optional" property:"gamemode"`
	ModeForced bool             `hcl:"mode_forced,optional" property:"force-gamemode"`
	Hardcore   bool             `hcl:"hardcore,optional" property:"hardcore"`
	PVP        bool             `hcl:"pvp,optional" property:"pvp"`
	Flying     bool             `hcl:"flying,optional" property:"allow-flight"`
	Spawning   GamePlaySpawning `hcl:"spawning,block"`
}

// TODO: This doesn't work for some reason :/
type GamePlaySpawning struct {
	Animals  bool `hcl:"animals,optional" property:"spawn-animals"`
	Monsters bool `hcl:"monsters,optional" property:"spawn-monsters"`
	NPCS     bool `hcl:"npcs,optional" property:"spawn-npcs"`
}

// Default will return a default configured Minecraft server
func Default() *Server {
	return &Server{
		Name:     "minecraft",
		Path:     "./.minecraft/",
		Version:  "latest",
		JavaHome: os.Getenv("JAVA_HOME"),
		GamePlay: GamePlay{
			Difficulty: "normal",
			Mode:       "survival",
			PVP:        true,
			Spawning: GamePlaySpawning{
				Animals:  true,
				Monsters: true,
				NPCS:     true,
			},
		},
		World: World{
			Type:             "default",
			Size:             29999984, // no idea why this is the default...
			Height:           256,
			EnableNether:     true,
			EnableStructures: true,
		},
	}
}

// Entrypoint returns the location of the Minecraft server jar
func (s *Server) Entrypoint() string {
	if s.version.ID == "" {
		return filepath.Join(s.Path, "server.jar")
	}

	return filepath.Join(s.Path, s.version.ID+".jar")
}
