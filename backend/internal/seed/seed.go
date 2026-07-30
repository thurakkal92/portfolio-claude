package seed

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Run wipes content tables and inserts the canonical EN+DE content.
// contact_submissions and contact_rate_limit are NOT touched.
func Run(ctx context.Context, pool *pgxpool.Pool) error {
	tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if err := truncateContent(ctx, tx); err != nil {
		return fmt.Errorf("truncate: %w", err)
	}
	if err := insertLocales(ctx, tx); err != nil {
		return fmt.Errorf("locales: %w", err)
	}
	if err := insertSiteSettings(ctx, tx); err != nil {
		return fmt.Errorf("site_settings: %w", err)
	}
	if err := insertHero(ctx, tx); err != nil {
		return fmt.Errorf("hero: %w", err)
	}
	if err := insertAbout(ctx, tx); err != nil {
		return fmt.Errorf("about: %w", err)
	}
	if err := insertSkills(ctx, tx); err != nil {
		return fmt.Errorf("skills: %w", err)
	}
	if err := insertProjects(ctx, tx); err != nil {
		return fmt.Errorf("projects: %w", err)
	}
	if err := insertExperiences(ctx, tx); err != nil {
		return fmt.Errorf("experiences: %w", err)
	}
	if err := insertContact(ctx, tx); err != nil {
		return fmt.Errorf("contact: %w", err)
	}
	if err := insertSocialLinks(ctx, tx); err != nil {
		return fmt.Errorf("social_links: %w", err)
	}

	return tx.Commit(ctx)
}

func truncateContent(ctx context.Context, tx pgx.Tx) error {
	_, err := tx.Exec(ctx, `
		TRUNCATE TABLE
			social_links,
			contact,
			experience_translations, experiences,
			project_translations, projects,
			skill_items, skill_group_translations, skill_groups,
			about, hero, site_settings, locales
		RESTART IDENTITY CASCADE`)
	return err
}

func insertLocales(ctx context.Context, tx pgx.Tx) error {
	rows := []struct {
		Code      string
		Name      string
		IsDefault bool
	}{
		{"en", "English", true},
		{"de", "Deutsch", false},
	}
	for _, r := range rows {
		if _, err := tx.Exec(ctx,
			`INSERT INTO locales (code, name, is_default) VALUES ($1, $2, $3)`,
			r.Code, r.Name, r.IsDefault); err != nil {
			return err
		}
	}
	return nil
}

func insertSiteSettings(ctx context.Context, tx pgx.Tx) error {
	type s struct {
		locale, title, desc, og, cv string
	}
	rows := []s{
		{
			locale: "en",
			title:  "Nabeel Thurakkal | Senior Frontend Developer",
			desc:   "Senior Frontend Developer based in Ulm, Germany — 8+ years building scalable web solutions with React, TypeScript, and modern CSS.",
			og:     "/og/og-image-en.png",
			cv:     "/cv/nabeel-thurakkal-en.pdf",
		},
		{
			locale: "de",
			title:  "Nabeel Thurakkal | Senior Frontend-Entwickler",
			desc:   "Senior Frontend-Entwickler aus Ulm, Deutschland – über 8 Jahre Erfahrung in skalierbaren Weblösungen mit React, TypeScript und modernem CSS.",
			og:     "/og/og-image-de.png",
			cv:     "/cv/nabeel-thurakkal-de.pdf",
		},
	}
	for _, r := range rows {
		if _, err := tx.Exec(ctx, `
			INSERT INTO site_settings (locale, site_title, site_description, og_image_path, cv_path)
			VALUES ($1, $2, $3, $4, $5)`,
			r.locale, r.title, r.desc, r.og, r.cv); err != nil {
			return err
		}
	}
	return nil
}

func insertHero(ctx context.Context, tx pgx.Tx) error {
	type h struct {
		locale, eyebrow, name, subtitle, ctaP, ctaS string
	}
	rows := []h{
		{
			locale:   "en",
			eyebrow:  "Senior Frontend Developer · Ulm, Germany",
			name:     "Nabeel Thurakkal",
			subtitle: "8+ years of experience building scalable web solutions with React, TypeScript, and modern CSS.",
			ctaP:     "View Projects",
			ctaS:     "Download CV",
		},
		{
			locale:   "de",
			eyebrow:  "Senior Frontend-Entwickler · Ulm, Deutschland",
			name:     "Nabeel Thurakkal",
			subtitle: "Über 8 Jahre Erfahrung in der Entwicklung skalierbarer Weblösungen mit React, TypeScript und modernem CSS.",
			ctaP:     "Projekte ansehen",
			ctaS:     "Lebenslauf herunterladen",
		},
	}
	for _, r := range rows {
		if _, err := tx.Exec(ctx, `
			INSERT INTO hero (locale, eyebrow, name, subtitle, cta_primary_label, cta_secondary_label)
			VALUES ($1, $2, $3, $4, $5, $6)`,
			r.locale, r.eyebrow, r.name, r.subtitle, r.ctaP, r.ctaS); err != nil {
			return err
		}
	}
	return nil
}

func insertAbout(ctx context.Context, tx pgx.Tx) error {
	enFacts, _ := json.Marshal([]map[string]string{
		{"label": "Location", "value": "Ulm, Germany", "kind": "text"},
		{"label": "Languages", "value": "German (B1), English (Fluent)", "kind": "text"},
		{"label": "Experience", "value": "8+ Years", "kind": "text"},
		{"label": "Availability", "value": "Open to work", "kind": "status"},
	})
	deFacts, _ := json.Marshal([]map[string]string{
		{"label": "Standort", "value": "Ulm, Deutschland", "kind": "text"},
		{"label": "Sprachen", "value": "Deutsch (B1), Englisch (fließend)", "kind": "text"},
		{"label": "Erfahrung", "value": "8+ Jahre", "kind": "text"},
		{"label": "Verfügbarkeit", "value": "Offen für neue Möglichkeiten", "kind": "status"},
	})
	rows := []struct {
		locale, heading, body string
		facts                 []byte
	}{
		{
			"en", "About Me",
			"Experienced Senior Frontend Developer with over 8 years in the industry. Specialized in building design systems, component libraries, and data-intensive dashboards. Proven track record of promoting frontend best practices and delivering high-performance SaaS products.\n\nI focus on technical precision and minimal aesthetics, ensuring that every pixel serves a purpose.",
			enFacts,
		},
		{
			"de", "Über mich",
			"Erfahrener Senior Frontend-Entwickler mit über 8 Jahren Berufserfahrung. Spezialisiert auf den Aufbau von Designsystemen, Komponentenbibliotheken und datenintensiven Dashboards. Nachweislich erfolgreich darin, Frontend-Best-Practices zu fördern und leistungsstarke SaaS-Produkte zu liefern.\n\nMein Fokus liegt auf technischer Präzision und minimalistischer Ästhetik – jedes Pixel hat einen Zweck.",
			deFacts,
		},
	}
	for _, r := range rows {
		if _, err := tx.Exec(ctx, `
			INSERT INTO about (locale, heading, body_md, quick_facts) VALUES ($1, $2, $3, $4)`,
			r.locale, r.heading, r.body, r.facts); err != nil {
			return err
		}
	}
	return nil
}

func insertSkills(ctx context.Context, tx pgx.Tx) error {
	type group struct {
		slug, icon string
		titleEN    string
		titleDE    string
		items      []string
	}
	groups := []group{
		{
			slug: "languages", icon: "code",
			titleEN: "Languages", titleDE: "Sprachen",
			items: []string{"JavaScript (ES6+)", "TypeScript", "HTML5", "CSS3", "Sass", "Less", "MDX"},
		},
		{
			slug: "frameworks", icon: "layers",
			titleEN: "Frameworks & Libraries", titleDE: "Frameworks & Bibliotheken",
			items: []string{"React", "React Redux", "Next.js", "React Query", "Tailwind CSS"},
		},
		{
			slug: "tools", icon: "wrench",
			titleEN: "Tools", titleDE: "Tools",
			items: []string{"Figma", "Sketch", "Storybook", "GitHub Actions", "Datadog"},
		},
		{
			slug: "devops", icon: "terminal",
			titleEN: "Testing & DevOps", titleDE: "Testing & DevOps",
			items: []string{"Jest", "React Testing Library", "CI/CD Pipelines", "Agile"},
		},
	}
	for i, g := range groups {
		var groupID int64
		if err := tx.QueryRow(ctx, `
			INSERT INTO skill_groups (slug, icon, display_order) VALUES ($1, $2, $3) RETURNING id`,
			g.slug, g.icon, i).Scan(&groupID); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `
			INSERT INTO skill_group_translations (group_id, locale, title) VALUES ($1, 'en', $2), ($1, 'de', $3)`,
			groupID, g.titleEN, g.titleDE); err != nil {
			return err
		}
		for j, label := range g.items {
			if _, err := tx.Exec(ctx, `
				INSERT INTO skill_items (group_id, label, display_order) VALUES ($1, $2, $3)`,
				groupID, label, j); err != nil {
				return err
			}
		}
	}
	return nil
}

func insertProjects(ctx context.Context, tx pgx.Tx) error {
	type proj struct {
		slug, company, image string
		live, source         *string
		tags                 []string
		titleEN, descEN      string
		titleDE, descDE      string
		highlightsEN         []map[string]string
		highlightsDE         []map[string]string
	}
	ptr := func(s string) *string { return &s }
	projects := []proj{
		{
			slug:    "quillbot-com",
			company: "QuillBot",
			image:   "/images/projects/quillbot-com.svg",
			live:    ptr("https://quillbot.com"),
			tags:    []string{"React", "TypeScript", "Micro-frontend"},
			titleEN: "quillbot.com — Product Platform",
			descEN:  "Product-side work on QuillBot's flagship AI writing platform — micro-frontend migration of auth and monetisation surfaces, a Paraphraser Library refactor, and marketing UI feeding growth experiments.",
			titleDE: "quillbot.com — Produktplattform",
			descDE:  "Produktseitige Arbeit an QuillBots KI-Schreibplattform – Micro-Frontend-Migration der Login- und Monetarisierungsbereiche, Refactoring der Paraphraser-Library sowie Marketing-UI für Wachstumsexperimente.",
			highlightsEN: []map[string]string{
				{"icon": "check", "text": "Refactored the Paraphraser Library — custom React hooks, API integrations, performance and readability cleanup"},
				{"icon": "check", "text": "Micro-frontend migration of Login/Signup, Accounts, Premium and Upgrade surfaces"},
				{"icon": "check", "text": "Shipped sales popups and banners powering targeted marketing campaigns"},
				{"icon": "trending-up", "text": "Ran A/B tests via Amplitude to validate product hypotheses"},
			},
			highlightsDE: []map[string]string{
				{"icon": "check", "text": "Refactoring der Paraphraser-Library – Custom Hooks, API-Integrationen, Performance- und Lesbarkeits-Cleanup"},
				{"icon": "check", "text": "Micro-Frontend-Migration von Login/Signup-, Account-, Premium- und Upgrade-Bereichen"},
				{"icon": "check", "text": "Sales-Popups und Banner für zielgerichtete Marketingkampagnen"},
				{"icon": "trending-up", "text": "A/B-Tests via Amplitude zur Validierung von Produkthypothesen"},
			},
		},
		{
			slug:    "cleartrip",
			company: "Cleartrip",
			image:   "/images/projects/cleartrip.svg",
			live:    ptr("https://www.cleartrip.com/flights"),
			tags:    []string{"React", "styled-components", "PWA"},
			titleEN: "Cleartrip — Flights, Hotels & Mobile PWA",
			descEN:  "Frontend revamp across Cleartrip's flight and hotel booking flows including the mobile PWA — modernising the React architecture and unifying the booking experience across verticals.",
			titleDE: "Cleartrip — Flüge, Hotels & Mobile PWA",
			descDE:  "Frontend-Neugestaltung der Flug- und Hotelbuchungsflows von Cleartrip inkl. der mobilen PWA – Modernisierung der React-Architektur und Vereinheitlichung der Buchungsoberflächen über Bereiche hinweg.",
			highlightsEN: []map[string]string{
				{"icon": "check", "text": "Frontend revamp and refactor of the Cleartrip Flights division on React + Bento"},
				{"icon": "check", "text": "Migrated Cleartrip Hotels home and detail pages to a PWA desktop experience"},
				{"icon": "check", "text": "Built booking-page UI components for the Cleartrip mobile PWA (React + styled-components)"},
			},
			highlightsDE: []map[string]string{
				{"icon": "check", "text": "Frontend-Neugestaltung und Refactoring des Cleartrip-Flights-Bereichs mit React + Bento"},
				{"icon": "check", "text": "Migration der Cleartrip-Hotels-Start- und Detailseiten auf eine PWA-Desktop-Erfahrung"},
				{"icon": "check", "text": "UI-Komponenten für die Buchungsseiten der Cleartrip Mobile PWA (React + styled-components)"},
			},
		},
		{
			slug:    "uniqlo",
			company: "Fast Retailing (Uniqlo)",
			image:   "/images/projects/uniqlo.svg",
			live:    ptr("https://www.uniqlo.com/"),
			tags:    []string{"React", "Node.js", "JavaScript"},
			titleEN: "Uniqlo — Frontend Maintenance",
			descEN:  "Frontend maintenance and feature work for Uniqlo (Fast Retailing) during my QBurst engagement — a high-traffic global e-commerce storefront.",
			titleDE: "Uniqlo — Frontend-Wartung",
			descDE:  "Frontend-Wartung und Feature-Entwicklung für Uniqlo (Fast Retailing) im Rahmen meiner Tätigkeit bei QBurst – ein stark frequentierter globaler E-Commerce-Storefront.",
			highlightsEN: []map[string]string{
				{"icon": "check", "text": "Frontend maintenance and feature development on the Uniqlo storefront"},
				{"icon": "check", "text": "Delivered as part of the QBurst engineering team"},
			},
			highlightsDE: []map[string]string{
				{"icon": "check", "text": "Frontend-Wartung und Feature-Entwicklung im Uniqlo-Storefront"},
				{"icon": "check", "text": "Umgesetzt im Rahmen des QBurst-Engineering-Teams"},
			},
		},
		{
			slug:    "styleq",
			company: "QuillBot",
			image:   "/images/projects/styleq.svg",
			tags:    []string{"React", "TypeScript", "Storybook"},
			titleEN: "StyleQ Design System",
			descEN:  "A WCAG-compliant design system built with React and TypeScript, powering QuillBot's Paraphraser, Grammar Checker, Summarizer and other AI writing tools across their entire user base.",
			titleDE: "StyleQ Designsystem",
			descDE:  "Ein WCAG-konformes Designsystem mit React und TypeScript, das QuillBots Paraphraser, Grammar Checker, Summarizer und weitere KI-Schreibtools über die gesamte Nutzerbasis hinweg unterstützt.",
			highlightsEN: []map[string]string{
				{"icon": "check", "text": "Revamped Paraphraser, Grammar Checker, Summarizer and other product pages onto the new DS"},
				{"icon": "check", "text": "WCAG compliance improved accessibility and UX across platforms"},
				{"icon": "check", "text": "Storybook documentation for team-wide adoption"},
				{"icon": "trending-up", "text": "Improved conversion rate, interface aesthetics and revenue"},
			},
			highlightsDE: []map[string]string{
				{"icon": "check", "text": "Überarbeitung von Paraphraser, Grammar Checker, Summarizer und weiteren Produktseiten auf das neue DS"},
				{"icon": "check", "text": "WCAG-Konformität für bessere Barrierefreiheit und UX über Plattformen hinweg"},
				{"icon": "check", "text": "Storybook-Dokumentation für die teamweite Adoption"},
				{"icon": "trending-up", "text": "Verbesserte Conversion-Rate, Interface-Ästhetik und Umsatz"},
			},
		},
		{
			slug:    "styleq-icons",
			company: "QuillBot",
			image:   "/images/projects/styleq-icons.svg",
			tags:    []string{"React", "TypeScript", "SVG"},
			titleEN: "styleq-icons",
			descEN:  "An extensive icon library shipped alongside StyleQ, enhancing the visual language and usability of QuillBot's interface across product surfaces.",
			titleDE: "styleq-icons",
			descDE:  "Eine umfangreiche Icon-Bibliothek als Teil von StyleQ, die die Bildsprache und Usability der QuillBot-Oberfläche über alle Produktbereiche stärkt.",
			highlightsEN: []map[string]string{
				{"icon": "check", "text": "Icons packaged as tree-shakeable React components"},
				{"icon": "check", "text": "Unified iconography across QuillBot's family of AI writing tools"},
			},
			highlightsDE: []map[string]string{
				{"icon": "check", "text": "Icons als tree-shakeable React-Komponenten"},
				{"icon": "check", "text": "Einheitliche Bildsprache über QuillBots Familie von KI-Schreibtools hinweg"},
			},
		},
		{
			slug:    "bento",
			company: "Cleartrip",
			image:   "/images/projects/bento.svg",
			live:    ptr("https://ctbento-v1.netlify.app/"),
			tags:    []string{"React", "SCSS", "Gatsby", "Storybook"},
			titleEN: "Bento Design System",
			descEN:  "A component library and design system built from scratch for Cleartrip with React, SCSS, Gatsby and Storybook — the foundation the Cleartrip Flights revamp was built on.",
			titleDE: "Bento Designsystem",
			descDE:  "Eine von Grund auf für Cleartrip entwickelte Komponentenbibliothek und Designsystem mit React, SCSS, Gatsby und Storybook – die Basis der Cleartrip-Flights-Neugestaltung.",
			highlightsEN: []map[string]string{
				{"icon": "check", "text": "Built the component library from scratch — React, SCSS, Gatsby, Storybook"},
				{"icon": "check", "text": "Storybook documentation supporting team-wide adoption"},
				{"icon": "check", "text": "Powered the Cleartrip Flights frontend revamp"},
			},
			highlightsDE: []map[string]string{
				{"icon": "check", "text": "Aufbau der Komponentenbibliothek von Grund auf – React, SCSS, Gatsby, Storybook"},
				{"icon": "check", "text": "Storybook-Dokumentation für die teamweite Adoption"},
				{"icon": "check", "text": "Basis für die Frontend-Neugestaltung von Cleartrip Flights"},
			},
		},
		{
			slug:    "stormbreaker",
			company: "Open Source",
			image:   "/images/projects/stormbreaker.svg",
			live:    ptr("https://stormbreaker-v2.netlify.app/"),
			source:  ptr("https://github.com/thurakkal92"),
			tags:    []string{"React", "Emotion JS"},
			titleEN: "Stormbreaker",
			descEN:  "An MIT-licensed component library focused on developer experience and ease of customization via Emotion JS.",
			titleDE: "Stormbreaker",
			descDE:  "Eine MIT-lizenzierte Komponentenbibliothek mit Fokus auf Developer Experience und einfache Anpassbarkeit über Emotion JS.",
			highlightsEN: []map[string]string{
				{"icon": "check", "text": "Accessible and customizable components"},
				{"icon": "star", "text": "Community-driven development"},
			},
			highlightsDE: []map[string]string{
				{"icon": "check", "text": "Barrierefreie und anpassbare Komponenten"},
				{"icon": "star", "text": "Community-getriebene Entwicklung"},
			},
		},
	}
	for i, p := range projects {
		hlEN, _ := json.Marshal(p.highlightsEN)
		hlDE, _ := json.Marshal(p.highlightsDE)

		var id int64
		if err := tx.QueryRow(ctx, `
			INSERT INTO projects (slug, company, image_path, live_url, source_url, tags, display_order)
			VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
			p.slug, p.company, p.image, p.live, p.source, p.tags, i).Scan(&id); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `
			INSERT INTO project_translations (project_id, locale, title, description, highlights)
			VALUES ($1, 'en', $2, $3, $4), ($1, 'de', $5, $6, $7)`,
			id, p.titleEN, p.descEN, hlEN, p.titleDE, p.descDE, hlDE); err != nil {
			return err
		}
	}
	return nil
}

func insertExperiences(ctx context.Context, tx pgx.Tx) error {
	type xp struct {
		company           string
		start             time.Time
		end               *time.Time
		location          *string
		roleEN, summaryEN string
		roleDE, summaryDE string
	}
	ptrTime := func(y int, m time.Month) *time.Time {
		t := time.Date(y, m, 1, 0, 0, 0, 0, time.UTC)
		return &t
	}
	ptr := func(s string) *string { return &s }

	xps := []xp{
		{
			company:   "P19 GmbH",
			start:     time.Date(2025, 5, 1, 0, 0, 0, 0, time.UTC),
			end:       nil,
			location:  ptr("Ulm, Germany"),
			roleEN:    "Senior Frontend Developer",
			summaryEN: "Built a comprehensive super admin dashboard from scratch and implemented backend systems with Node.js and TypeScript.",
			roleDE:    "Senior Frontend-Entwickler",
			summaryDE: "Komplette Neuentwicklung eines umfassenden Super-Admin-Dashboards sowie Aufbau von Backend-Systemen mit Node.js und TypeScript.",
		},
		{
			company:   "QuillBot",
			start:     time.Date(2022, 3, 1, 0, 0, 0, 0, time.UTC),
			end:       ptrTime(2024, 8),
			location:  ptr("Remote"),
			roleEN:    "Senior Frontend Developer",
			summaryEN: "Developed the StyleQ design system and revamped all major product pages, significantly boosting conversion rates.",
			roleDE:    "Senior Frontend-Entwickler",
			summaryDE: "Entwicklung des StyleQ Designsystems und Überarbeitung aller wichtigen Produktseiten, mit deutlicher Steigerung der Conversion-Rate.",
		},
		{
			company:   "Cleartrip",
			start:     time.Date(2018, 8, 1, 0, 0, 0, 0, time.UTC),
			end:       ptrTime(2022, 3),
			location:  ptr("Remote"),
			roleEN:    "Senior UI Engineer",
			summaryEN: "Led the 'Bento' component library and the migration of hotels and the mobile PWA to modern React architectures.",
			roleDE:    "Senior UI Engineer",
			summaryDE: "Leitung der „Bento“-Komponentenbibliothek sowie der Migration der Hotelplattform und der mobilen PWA auf moderne React-Architekturen.",
		},
		{
			company:   "QBurst Technologies",
			start:     time.Date(2014, 11, 1, 0, 0, 0, 0, time.UTC),
			end:       ptrTime(2018, 7),
			location:  nil,
			roleEN:    "Senior Frontend Engineer",
			summaryEN: "Worked on large-scale applications including Millicom TIGO using React and Node.js.",
			roleDE:    "Senior Frontend Engineer",
			summaryDE: "Arbeit an Großprojekten, u. a. Millicom TIGO, mit React und Node.js.",
		},
	}

	for i, x := range xps {
		var id int64
		if err := tx.QueryRow(ctx, `
			INSERT INTO experiences (company, start_date, end_date, location, display_order)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			x.company, x.start, x.end, x.location, i).Scan(&id); err != nil {
			return err
		}
		if _, err := tx.Exec(ctx, `
			INSERT INTO experience_translations (experience_id, locale, role, summary)
			VALUES ($1, 'en', $2, $3), ($1, 'de', $4, $5)`,
			id, x.roleEN, x.summaryEN, x.roleDE, x.summaryDE); err != nil {
			return err
		}
	}
	return nil
}

func insertContact(ctx context.Context, tx pgx.Tx) error {
	type c struct {
		locale, heading, sub, nameL, emailL, msgL, submitL, success, errMsg string
	}
	rows := []c{
		{
			locale:  "en",
			heading: "Get In Touch",
			sub:     "Open for new opportunities and collaborations. Let's build something precise.",
			nameL:   "Name", emailL: "Email", msgL: "Message", submitL: "Send Message",
			success: "Thanks — your message is on its way. I'll get back to you soon.",
			errMsg:  "Something went wrong. Please try again or email me directly.",
		},
		{
			locale:  "de",
			heading: "Kontakt aufnehmen",
			sub:     "Offen für neue Möglichkeiten und Kooperationen. Lass uns gemeinsam etwas Präzises bauen.",
			nameL:   "Name", emailL: "E-Mail", msgL: "Nachricht", submitL: "Nachricht senden",
			success: "Danke – deine Nachricht ist unterwegs. Ich melde mich bald.",
			errMsg:  "Etwas ist schiefgelaufen. Bitte versuche es erneut oder schreibe mir direkt eine E-Mail.",
		},
	}
	for _, r := range rows {
		if _, err := tx.Exec(ctx, `
			INSERT INTO contact (locale, heading, subheading, form_name_label, form_email_label,
				form_message_label, form_submit_label, success_message, error_message)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			r.locale, r.heading, r.sub, r.nameL, r.emailL, r.msgL, r.submitL, r.success, r.errMsg); err != nil {
			return err
		}
	}
	return nil
}

func insertSocialLinks(ctx context.Context, tx pgx.Tx) error {
	rows := []struct {
		kind, href, label string
		order             int
	}{
		{"email", "mailto:nabeel.thurakkal92@gmail.com", "Email", 0},
		{"linkedin", "https://linkedin.com/in/nabeelthurakkal", "LinkedIn", 1},
		{"github", "https://github.com/thurakkal92", "GitHub", 2},
	}
	for _, r := range rows {
		if _, err := tx.Exec(ctx, `
			INSERT INTO social_links (kind, href, display_label, display_order)
			VALUES ($1, $2, $3, $4)`,
			r.kind, r.href, r.label, r.order); err != nil {
			return err
		}
	}
	return nil
}
