import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

export default defineConfig({
  site: 'https://ska.gchiesa.dev',
  integrations: [
    starlight({
      title: 'SKA',
      description: 'Your scaffolding buddy - A powerful templating tool for Platform Engineers',
      logo: {
        src: './src/assets/ska-logo.png',
        replacesTitle: true,
      },
      social: {
        github: 'https://github.com/gchiesa/ska',
      },
      customCss: [
        './src/styles/tailwind.css',
        './src/styles/custom.css',
      ],
      sidebar: [
        {
          label: 'Getting Started',
          items: [
            { label: 'Introduction', slug: 'getting-started/introduction' },
            { label: 'Installation', slug: 'getting-started/installation' },
            { label: 'Quick Start', slug: 'getting-started/quick-start' },
          ],
        },
        {
          label: 'Concepts',
          items: [
            { label: 'Upstream Blueprints', slug: 'concepts/upstream-blueprints' },
            { label: 'Template Language', slug: 'concepts/template-language' },
            { label: 'Partial Sections', slug: 'concepts/partial-sections' },
            { label: 'Terminal UI', slug: 'concepts/terminal-ui' },
          ],
        },
        {
          label: 'Use Cases',
          items: [
            { label: 'Multiple Template Subfolders', slug: 'use-cases/multiple-template-subfolders' },
            { label: 'Partial File Management', slug: 'use-cases/partial-management' },
            { label: 'Ignore Files After First Run', slug: 'use-cases/ignore-files-after-first' },
            { label: 'Multiple Configs Same Folder', slug: 'use-cases/multiple-configs-same-folder' },
          ],
        },
      ],
      components: {
        Hero: './src/components/Hero.astro',
      },
    }),
  ],
});
