import { RouteObject } from 'react-router-dom';

import { Path, PathValues } from './path';
import { MainLayout } from '../views/layout';
import BattlefieldView from '../views/battlefield/battlefield';
import { PageNotFoundView } from '../views/pageNotFound';

export type Route = RouteObject & {
    path: PathValues;
    element: React.ComponentType;
};

export const routes: Route[] = [
    {
        path: Path.Battlefield,
        // @ts-expect-error react-router-dom types are not correct
        element: MainLayout(BattlefieldView, {
            className: 'view--battlefield',
            i18n: { pageTitle: 'views.battlefield.pageTitle' }
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
