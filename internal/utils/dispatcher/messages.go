package dispatcher

import (
	"context"
	"fmt"

	"github.com/go-telegram/bot"
)

type Sendable interface {
	Send(ctx context.Context, b *bot.Bot) error
}

// Message https://core.telegram.org/bots/api#sendmessage
type TextMessage struct {
	Params *bot.SendMessageParams
}

func (tm TextMessage) Send(ctx context.Context, b *bot.Bot) error {
	_, err := b.SendMessage(ctx, tm.Params)
	if err != nil {
		return fmt.Errorf("dispatcher.text_mesage.send %w", err)
	}

	return nil
}

// PhotoMessage https://core.telegram.org/bots/api#sendphoto
type PhotoMessage struct {
	Params *bot.SendPhotoParams
}

func (pm PhotoMessage) Send(ctx context.Context, b *bot.Bot) error {
	_, err := b.SendPhoto(ctx, pm.Params)
	if err != nil {
		return fmt.Errorf("dispatcher.photo_mesage.send %w", err)
	}

	return nil
}

// DocumentMessage https://core.telegram.org/bots/api#senddocument
type DocumentMessage struct {
	Params *bot.SendDocumentParams
}

func (dm DocumentMessage) Send(ctx context.Context, b *bot.Bot) error {
	_, err := b.SendDocument(ctx, dm.Params)
	if err != nil {
		return fmt.Errorf("dispatcher.document_mesage.send %w", err)
	}

	return nil
}

// EditMessageText https://core.telegram.org/bots/api#editmessagetext
type EditMessage struct {
	Params *bot.EditMessageTextParams
}

func (em EditMessage) Send(ctx context.Context, b *bot.Bot) error {
	_, err := b.EditMessageText(ctx, em.Params)
	if err != nil {
		return fmt.Errorf("dispatcher.edit_mesage.send %w", err)
	}

	return nil
}
