package control

import (
	"context"
	"fmt"

	"github.com/renevo/mcutils/pkg/minecraft"
)

type Minecraft struct {
	Server *minecraft.Server
}

func (m *Minecraft) Execute(ctx context.Context, cmd string, reply *string) error {
	if err := m.Server.ExecuteCommand(cmd); err != nil {
		return err
	}

	*reply = fmt.Sprintf("Executed command: %q", cmd)

	return nil
}

func (m *Minecraft) WhitelistAdd(ctx context.Context, user string, reply *string) error {
	if err := m.Server.ExecuteCommand(fmt.Sprintf("whitelist add %s", user)); err != nil {
		return err
	}

	if err := m.Server.ExecuteCommand("whitelist reload"); err != nil {
		return err
	}

	*reply = fmt.Sprintf("Whitelist added %q and reloaded", user)

	return nil
}
