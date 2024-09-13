export type MainLayoutProps = {
  className?: string;
  /**
   * Translations keys, not the actual translations.
   *
   * e.g.:
   * { pageTitle: 'views.browser.page-title' }
   *   vs
   * { pageTitle: 'Welcome, just don't expect much.' }
   */
  i18n?: { pageTitle: string };
};
