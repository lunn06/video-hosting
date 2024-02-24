package models

type Role struct {
	Id                    uint32 `yaml:"id" db:"id"`
	Name                  string `yaml:"name" db:"name"`
	CanRemoveUsers        bool   `yaml:"can_remove_users" db:"can_remove_users"`
	CanRemoveOthersVideos bool   `yaml:"can_remove_others_videos" db:"can_remove_others_videos"`
}
