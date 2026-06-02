package content

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Service {
	return &Service{pool: pool}
}

var ErrLocaleNotFound = errors.New("locale not found")

// Get assembles the full content payload for a locale.
func (s *Service) Get(ctx context.Context, locale string) (*Payload, error) {
	p := &Payload{Locale: locale}

	if err := s.loadSiteSettings(ctx, locale, &p.SiteSettings); err != nil {
		return nil, err
	}
	if err := s.loadHero(ctx, locale, &p.Hero); err != nil {
		return nil, err
	}
	if err := s.loadAbout(ctx, locale, &p.About); err != nil {
		return nil, err
	}
	skills, err := s.loadSkills(ctx, locale)
	if err != nil {
		return nil, err
	}
	p.Skills = skills

	projects, err := s.loadProjects(ctx, locale)
	if err != nil {
		return nil, err
	}
	p.Projects = projects

	exps, err := s.loadExperiences(ctx, locale)
	if err != nil {
		return nil, err
	}
	p.Experiences = exps

	if err := s.loadContact(ctx, locale, &p.Contact); err != nil {
		return nil, err
	}

	socials, err := s.loadSocialLinks(ctx)
	if err != nil {
		return nil, err
	}
	p.SocialLinks = socials

	return p, nil
}

func (s *Service) loadSiteSettings(ctx context.Context, locale string, dst *SiteSettings) error {
	row := s.pool.QueryRow(ctx, `
		SELECT site_title, site_description, og_image_path, cv_path
		FROM site_settings WHERE locale = $1`, locale)
	if err := row.Scan(&dst.SiteTitle, &dst.SiteDescription, &dst.OGImagePath, &dst.CVPath); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrLocaleNotFound
		}
		return fmt.Errorf("site_settings: %w", err)
	}
	return nil
}

func (s *Service) loadHero(ctx context.Context, locale string, dst *Hero) error {
	row := s.pool.QueryRow(ctx, `
		SELECT eyebrow, name, subtitle, cta_primary_label, cta_secondary_label
		FROM hero WHERE locale = $1`, locale)
	if err := row.Scan(&dst.Eyebrow, &dst.Name, &dst.Subtitle, &dst.CTAPrimaryLabel, &dst.CTASecondaryLabel); err != nil {
		return fmt.Errorf("hero: %w", err)
	}
	return nil
}

func (s *Service) loadAbout(ctx context.Context, locale string, dst *About) error {
	var quickFactsJSON []byte
	row := s.pool.QueryRow(ctx, `
		SELECT heading, body_md, quick_facts
		FROM about WHERE locale = $1`, locale)
	if err := row.Scan(&dst.Heading, &dst.BodyMD, &quickFactsJSON); err != nil {
		return fmt.Errorf("about: %w", err)
	}
	if err := unmarshalJSONB(quickFactsJSON, &dst.QuickFacts); err != nil {
		return fmt.Errorf("about quick_facts: %w", err)
	}
	return nil
}

func (s *Service) loadSkills(ctx context.Context, locale string) ([]SkillGroup, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT g.id, g.slug, g.icon, t.title
		FROM skill_groups g
		JOIN skill_group_translations t ON t.group_id = g.id AND t.locale = $1
		ORDER BY g.display_order ASC`, locale)
	if err != nil {
		return nil, fmt.Errorf("skill_groups: %w", err)
	}
	defer rows.Close()

	type gRow struct {
		ID    int64
		Group SkillGroup
	}
	var groups []gRow
	for rows.Next() {
		var g gRow
		if err := rows.Scan(&g.ID, &g.Group.Slug, &g.Group.Icon, &g.Group.Title); err != nil {
			return nil, fmt.Errorf("skill scan: %w", err)
		}
		groups = append(groups, g)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	out := make([]SkillGroup, 0, len(groups))
	for _, g := range groups {
		itemRows, err := s.pool.Query(ctx, `
			SELECT label FROM skill_items WHERE group_id = $1 ORDER BY display_order ASC`, g.ID)
		if err != nil {
			return nil, fmt.Errorf("skill_items: %w", err)
		}
		items := []string{}
		for itemRows.Next() {
			var label string
			if err := itemRows.Scan(&label); err != nil {
				itemRows.Close()
				return nil, err
			}
			items = append(items, label)
		}
		itemRows.Close()
		g.Group.Items = items
		out = append(out, g.Group)
	}
	return out, nil
}

func (s *Service) loadProjects(ctx context.Context, locale string) ([]Project, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT p.slug, p.company, p.image_path, p.live_url, p.source_url, p.tags,
		       t.title, t.description, t.highlights
		FROM projects p
		JOIN project_translations t ON t.project_id = p.id AND t.locale = $1
		ORDER BY p.display_order ASC`, locale)
	if err != nil {
		return nil, fmt.Errorf("projects: %w", err)
	}
	defer rows.Close()

	out := []Project{}
	for rows.Next() {
		var p Project
		var highlightsJSON []byte
		if err := rows.Scan(&p.Slug, &p.Company, &p.ImagePath, &p.LiveURL, &p.SourceURL, &p.Tags,
			&p.Title, &p.Description, &highlightsJSON); err != nil {
			return nil, fmt.Errorf("project scan: %w", err)
		}
		if err := unmarshalJSONB(highlightsJSON, &p.Highlights); err != nil {
			return nil, fmt.Errorf("project highlights: %w", err)
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (s *Service) loadExperiences(ctx context.Context, locale string) ([]Experience, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT e.company, e.start_date, e.end_date, e.location, t.role, t.summary
		FROM experiences e
		JOIN experience_translations t ON t.experience_id = e.id AND t.locale = $1
		ORDER BY e.display_order ASC`, locale)
	if err != nil {
		return nil, fmt.Errorf("experiences: %w", err)
	}
	defer rows.Close()

	out := []Experience{}
	for rows.Next() {
		var (
			e         Experience
			startDate any
			endDate   any
		)
		if err := rows.Scan(&e.Company, &startDate, &endDate, &e.Location, &e.Role, &e.Summary); err != nil {
			return nil, fmt.Errorf("experience scan: %w", err)
		}
		e.StartDate = formatDate(startDate)
		if endDate != nil {
			s := formatDate(endDate)
			e.EndDate = &s
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (s *Service) loadContact(ctx context.Context, locale string, dst *Contact) error {
	row := s.pool.QueryRow(ctx, `
		SELECT heading, subheading, form_name_label, form_email_label,
		       form_message_label, form_submit_label, success_message, error_message
		FROM contact WHERE locale = $1`, locale)
	if err := row.Scan(&dst.Heading, &dst.Subheading, &dst.FormNameLabel, &dst.FormEmailLabel,
		&dst.FormMessageLabel, &dst.FormSubmitLabel, &dst.SuccessMessage, &dst.ErrorMessage); err != nil {
		return fmt.Errorf("contact: %w", err)
	}
	return nil
}

func (s *Service) loadSocialLinks(ctx context.Context) ([]SocialLink, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT kind, href, display_label FROM social_links ORDER BY display_order ASC`)
	if err != nil {
		return nil, fmt.Errorf("social_links: %w", err)
	}
	defer rows.Close()
	out := []SocialLink{}
	for rows.Next() {
		var l SocialLink
		if err := rows.Scan(&l.Kind, &l.Href, &l.DisplayLabel); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}
