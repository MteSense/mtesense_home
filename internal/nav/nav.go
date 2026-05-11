package nav

import (
	"database/sql"
	"fmt"
)

type Group struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	SortOrder int    `json:"sortOrder"`
	Visible   bool   `json:"visible"`
}

type Link struct {
	ID           int64  `json:"id"`
	GroupID      int64  `json:"groupId"`
	Title        string `json:"title"`
	URL          string `json:"url"`
	Icon         string `json:"icon"`
	IconType     string `json:"iconType"`
	Description  string `json:"description"`
	SortOrder    int    `json:"sortOrder"`
	Visible      bool   `json:"visible"`
	OpenInNewTab bool   `json:"openInNewTab"`
}

type GroupWithLinks struct {
	Group
	Links []Link `json:"links"`
}

type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) PublicNavigation() ([]GroupWithLinks, error) {
	return s.listNavigation(true)
}

func (s *Service) AdminNavigation() ([]GroupWithLinks, error) {
	return s.listNavigation(false)
}

func (s *Service) listNavigation(visibleOnly bool) ([]GroupWithLinks, error) {
	groupQuery := "SELECT id, title, sort_order, visible FROM nav_groups"
	linkQuery := "SELECT id, group_id, title, url, icon, icon_type, description, sort_order, visible, open_in_new_tab FROM nav_links"
	if visibleOnly {
		groupQuery += " WHERE visible = 1"
		linkQuery += " WHERE visible = 1"
	}
	groupQuery += " ORDER BY sort_order ASC, id ASC"
	linkQuery += " ORDER BY sort_order ASC, id ASC"

	rows, err := s.db.Query(groupQuery)
	if err != nil {
		return nil, fmt.Errorf("list groups: %w", err)
	}
	defer rows.Close()

	groups := make([]GroupWithLinks, 0)
	groupIndex := map[int64]int{}
	for rows.Next() {
		var group GroupWithLinks
		if err := rows.Scan(&group.ID, &group.Title, &group.SortOrder, &group.Visible); err != nil {
			return nil, fmt.Errorf("scan group: %w", err)
		}
		group.Links = []Link{}
		groupIndex[group.ID] = len(groups)
		groups = append(groups, group)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate groups: %w", err)
	}

	linkRows, err := s.db.Query(linkQuery)
	if err != nil {
		return nil, fmt.Errorf("list links: %w", err)
	}
	defer linkRows.Close()
	for linkRows.Next() {
		var link Link
		if err := linkRows.Scan(&link.ID, &link.GroupID, &link.Title, &link.URL, &link.Icon, &link.IconType, &link.Description, &link.SortOrder, &link.Visible, &link.OpenInNewTab); err != nil {
			return nil, fmt.Errorf("scan link: %w", err)
		}
		if idx, ok := groupIndex[link.GroupID]; ok {
			groups[idx].Links = append(groups[idx].Links, link)
		}
	}
	if err := linkRows.Err(); err != nil {
		return nil, fmt.Errorf("iterate links: %w", err)
	}
	return groups, nil
}

func (s *Service) CreateGroup(group Group) (Group, error) {
	result, err := s.db.Exec(
		"INSERT INTO nav_groups (title, sort_order, visible) VALUES (?, ?, ?)",
		group.Title,
		group.SortOrder,
		group.Visible,
	)
	if err != nil {
		return Group{}, fmt.Errorf("create group: %w", err)
	}
	group.ID, _ = result.LastInsertId()
	return group, nil
}

func (s *Service) UpdateGroup(id int64, group Group) (Group, error) {
	_, err := s.db.Exec(
		"UPDATE nav_groups SET title = ?, sort_order = ?, visible = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		group.Title,
		group.SortOrder,
		group.Visible,
		id,
	)
	if err != nil {
		return Group{}, fmt.Errorf("update group: %w", err)
	}
	group.ID = id
	return group, nil
}

func (s *Service) DeleteGroup(id int64) error {
	_, err := s.db.Exec("DELETE FROM nav_groups WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete group: %w", err)
	}
	return nil
}

func (s *Service) CreateLink(link Link) (Link, error) {
	result, err := s.db.Exec(
		`INSERT INTO nav_links
		(group_id, title, url, icon, icon_type, description, sort_order, visible, open_in_new_tab)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		link.GroupID,
		link.Title,
		link.URL,
		link.Icon,
		link.IconType,
		link.Description,
		link.SortOrder,
		link.Visible,
		link.OpenInNewTab,
	)
	if err != nil {
		return Link{}, fmt.Errorf("create link: %w", err)
	}
	link.ID, _ = result.LastInsertId()
	return link, nil
}

func (s *Service) UpdateLink(id int64, link Link) (Link, error) {
	_, err := s.db.Exec(
		`UPDATE nav_links
		SET group_id = ?, title = ?, url = ?, icon = ?, icon_type = ?, description = ?,
		    sort_order = ?, visible = ?, open_in_new_tab = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?`,
		link.GroupID,
		link.Title,
		link.URL,
		link.Icon,
		link.IconType,
		link.Description,
		link.SortOrder,
		link.Visible,
		link.OpenInNewTab,
		id,
	)
	if err != nil {
		return Link{}, fmt.Errorf("update link: %w", err)
	}
	link.ID = id
	return link, nil
}

func (s *Service) DeleteLink(id int64) error {
	_, err := s.db.Exec("DELETE FROM nav_links WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("delete link: %w", err)
	}
	return nil
}
