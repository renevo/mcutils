package pubsub

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/renevo/mcutils/internal/command/minecraft/modules/ext"
	"github.com/sirupsen/logrus"
)

type logger struct {
	ctx context.Context
}

func (l *logger) Error(msg string, err error, fields watermill.LogFields) {
	log := ext.Logger(l.ctx)
	if len(fields) > 0 {
		log = log.WithFields(map[string]interface{}(fields))
	}

	log.Errorf(msg+": %s", err.Error())
}
func (l *logger) Info(msg string, fields watermill.LogFields) {
	log := ext.Logger(l.ctx)
	if len(fields) > 0 {
		log = log.WithFields(map[string]interface{}(fields))
	}

	log.Info(msg)
}

func (l *logger) Debug(msg string, fields watermill.LogFields) {
	log := ext.Logger(l.ctx)
	if len(fields) > 0 {
		log = log.WithFields(map[string]interface{}(fields))
	}

	log.Debug(msg)
}

func (l *logger) Trace(msg string, fields watermill.LogFields) {
	log := ext.Logger(l.ctx)
	if len(fields) > 0 {
		log = log.WithFields(map[string]interface{}(fields))
	}

	log.Trace(msg)
}

func (l *logger) With(fields watermill.LogFields) watermill.LoggerAdapter {
	log := ext.Logger(l.ctx)
	return &logger{ext.WithLogger(l.ctx, log.WithFields(logrus.Fields(map[string]interface{}((fields)))))}
}
