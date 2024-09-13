import { RouteObject } from 'react-router-dom';

import { Path, PathValues } from './path';
import { MainLayout } from '../views/layout';
import HomeView from '../views/home/Home';
import { PageNotFoundView } from '../views/pageNotFound';

export type Route = RouteObject & {
    path: PathValues;
    element: React.ComponentType;
};

export const routes: Route[] = [
    {
        path: Path.Home,
        // @ts-expect-error react-router-dom types are not correct
        element: MainLayout(HomeView, {
            className: 'view--home',
            i18n: { pageTitle: 'views.home.pageTitle' }
        })
    },
    {
        path: Path.PageNotFound,
        // @ts-expect-error react-router-dom types are not correct
        element: MainLayout(PageNotFoundView, {
            className: 'view--page-not-found',
            i18n: { pageTitle: 'views.pageNotFound.pageTitle' }
        }),
    }
];
