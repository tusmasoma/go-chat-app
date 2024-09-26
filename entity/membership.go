package entity

import (
	"fmt"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

type Membership struct {
	UserID          string
	WorkspaceID     string
	Name            string
	ProfileImageURL string
	IsAdmin         bool
}

func NewMembership(userID, workspaceID, name, profileImageURL string, isAdmin bool) (*Membership, error) {
	if userID == "" {
		log.Error("UserID is required", log.Fstring("userID", userID))
		return nil, fmt.Errorf("userID is required")
	}
	if workspaceID == "" {
		log.Error("WorkspaceID is required", log.Fstring("workspaceID", workspaceID))
		return nil, fmt.Errorf("workspaceID is required")
	}
	if name == "" {
		log.Error("Name is required", log.Fstring("name", name))
		return nil, fmt.Errorf("name is required")
	}
	if profileImageURL == "" {
		profileImageURL = "https://www.hoge.com/avatar.jpg"
	}
	return &Membership{
		UserID:          userID,
		WorkspaceID:     workspaceID,
		Name:            name,
		ProfileImageURL: profileImageURL,
		IsAdmin:         isAdmin,
	}, nil
}
