package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/easymirror/easymirror-backend/internal/common"
	"github.com/easymirror/easymirror-backend/internal/db"
	"github.com/easymirror/easymirror-backend/internal/mirrorlink"
)

const (
	query_limit = 25
)

// Returns a list of items a user has uploaded
func (u user) MirrorLinks(ctx context.Context, db *db.Database, pageNum int) ([]mirrorlink.MirrorLink, error) {
	links, err := mirrorlink.GetUserLinks(ctx, db, u.ID().String(), common.GetPageOffset(query_limit, pageNum))
	if err != nil {
		return nil, fmt.Errorf("GetUserLinks error: %w", err)
	}
	return links, nil
}

func (u user) UpdateMirrorLinkName(
	ctx context.Context,
	db *db.Database,
	linkID, name string,
) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("name cannot be empty")
	}

	if err := mirrorlink.UpdateName(ctx, db, linkID, u.ID().String(), name); err != nil {
		log.Println("Error updating name:", err)
		return fmt.Errorf("UpdateName error: %w", err)
	}
	return nil
}

func (u user) DeleteMirrorLink(ctx context.Context, db *db.Database, linkID string) error {
	if err := mirrorlink.Delete(ctx, db, u.ID().String(), linkID); err != nil {
		log.Println("Error deleting:", err)
		return fmt.Errorf("delete error: %w", err)
	}
	return nil
}

func (u user) GetFiles(ctx context.Context, db *db.Database, mirrorLinkID string) ([]mirrorlink.File, error) {
	files, err := mirrorlink.GetFilesFromMirror(ctx, db, mirrorLinkID, u.ID().String())
	if err != nil {
		return nil, fmt.Errorf("GetFilesFromMirror error: %w", err)
	}
	return files, nil
}
