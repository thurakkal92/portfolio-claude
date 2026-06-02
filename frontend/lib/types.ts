export type QuickFact = {
  label: string;
  value: string;
  kind: "text" | "status";
};

export type SiteSettings = {
  siteTitle: string;
  siteDescription: string;
  ogImagePath?: string;
  cvPath?: string;
};

export type Hero = {
  eyebrow: string;
  name: string;
  subtitle: string;
  ctaPrimaryLabel: string;
  ctaSecondaryLabel: string;
};

export type About = {
  heading: string;
  bodyMd: string;
  quickFacts: QuickFact[];
};

export type SkillGroup = {
  slug: string;
  icon: string;
  title: string;
  items: string[];
};

export type ProjectHighlight = { icon: string; text: string };

export type Project = {
  slug: string;
  company: string;
  title: string;
  description: string;
  imagePath: string;
  liveUrl?: string;
  sourceUrl?: string;
  tags: string[];
  highlights: ProjectHighlight[];
};

export type Experience = {
  company: string;
  role: string;
  startDate: string;
  endDate?: string;
  location?: string;
  summary: string;
};

export type ContactCopy = {
  heading: string;
  subheading: string;
  formNameLabel: string;
  formEmailLabel: string;
  formMessageLabel: string;
  formSubmitLabel: string;
  successMessage: string;
  errorMessage: string;
};

export type SocialLink = {
  kind: "email" | "github" | "linkedin" | string;
  href: string;
  displayLabel: string;
};

export type ContentPayload = {
  locale: string;
  siteSettings: SiteSettings;
  hero: Hero;
  about: About;
  skills: SkillGroup[];
  projects: Project[];
  experiences: Experience[];
  contact: ContactCopy;
  socialLinks: SocialLink[];
};
