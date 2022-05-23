package control

import (
	"context"
	"fmt"
	"sync"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/pkg/errors"
	"github.com/renevo/mcutils/pkg/minecraft"
)

type Minecraft struct {
	Controller *MinecraftController
}

func (m *Minecraft) Execute(ctx context.Context, cmd string, reply *bool) error {
	err := m.Controller.Server.ExecuteCommand(cmd)
	*reply = err == nil
	return err
}

func (m *Minecraft) WhitelistAdd(ctx context.Context, player string, reply *string) error {
	player, err := m.Controller.WhitelistAdd(ctx, player)
	*reply = player
	return err
}

func (m *Minecraft) WhitelistRemove(ctx context.Context, player string, reply *string) error {
	player, err := m.Controller.WhitelistRemove(ctx, player)
	*reply = player
	return err
}

type MinecraftController struct {
	mu         sync.Mutex
	Server     *minecraft.Server
	Subscriber message.Subscriber
}

func (m *MinecraftController) SaveGame(ctx context.Context) error {
	_, err := m.Execute(ctx, "save-all", minecraft.EventSaved)
	return errors.Wrap(err, "failed to save the game")
}

func (m *MinecraftController) Stop(ctx context.Context) error {
	_, err := m.Execute(ctx, "stop", minecraft.StateStopping)
	return errors.Wrap(err, "failed to stop the game")
}

func (m *MinecraftController) Say(ctx context.Context, msg string) error {
	_, err := m.Execute(ctx, fmt.Sprintf("say %s", msg), minecraft.EventServerChat)
	return errors.Wrap(err, "failed to say something")
}

func (m *MinecraftController) Emote(ctx context.Context, msg string) error {
	_, err := m.Execute(ctx, fmt.Sprintf("me %s", msg), minecraft.EventServerEmote)
	return errors.Wrap(err, "failed to emote something")
}

func (m *MinecraftController) WhitelistAdd(ctx context.Context, player string) (string, error) {
	result, meta, err := m.TryExecute(ctx, fmt.Sprintf("whitelist add %s", player), minecraft.EventWhitelistAdd, minecraft.EventWhitelistUnknown)
	if err != nil {
		return "", err
	}

	if result == minecraft.EventWhitelistUnknown {
		return "", errors.Errorf("player %q does not exist", player)
	}

	_, err = m.Execute(ctx, "whitelist reload", minecraft.EventWhitelistReloaded)

	return meta.Get("player"), err
}

func (m *MinecraftController) WhitelistRemove(ctx context.Context, player string) (string, error) {
	result, meta, err := m.TryExecute(ctx, fmt.Sprintf("whitelist remove %s", player), minecraft.EventWhitelistRemove, minecraft.EventWhitelistUnknown)
	if err != nil {
		return "", err
	}

	if result == minecraft.EventWhitelistUnknown {
		return "", errors.Errorf("player %q does not exist", player)
	}

	_, err = m.Execute(ctx, "whitelist reload", minecraft.EventWhitelistReloaded)

	return meta.Get("player"), err
}

func (m *MinecraftController) Execute(ctx context.Context, cmd, result string) (message.Metadata, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sub, err := m.Subscriber.Subscribe(ctx, result)
	if err != nil {
		return nil, err
	}

	if err := m.Server.ExecuteCommand(cmd); err != nil {
		return nil, err
	}

	select {
	case <-ctx.Done():
	case msg := <-sub:
		msg.Ack()
		return msg.Metadata, nil
	}

	return nil, ctx.Err()
}

func (m *MinecraftController) TryExecute(ctx context.Context, cmd, success, failure string) (string, message.Metadata, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	subSuccess, err := m.Subscriber.Subscribe(ctx, success)
	if err != nil {
		return "", nil, err
	}
	subFailure, err := m.Subscriber.Subscribe(ctx, failure)
	if err != nil {
		return "", nil, err
	}

	if err := m.Server.ExecuteCommand(cmd); err != nil {
		return "", nil, err
	}

	select {
	case <-ctx.Done():
	case msg := <-subSuccess:
		msg.Ack()
		return success, msg.Metadata, nil
	case msg := <-subFailure:
		msg.Ack()
		return failure, msg.Metadata, nil
	}

	return "", nil, ctx.Err()
}
