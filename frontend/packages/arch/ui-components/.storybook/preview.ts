// frontend/packages/arch/ui-components/.storybook/preview.ts

import type { Preview } from '@storybook/react';
import '../src/styles/reset.css';

const preview: Preview = {
  parameters: {
    actions: { argTypesRegex: '^on[A-Z].*' },
    controls: {
      matchers: {
        color: /(background|color)$/i,
        date: /Date$/,
      },
    },
  },
};

export default preview;
