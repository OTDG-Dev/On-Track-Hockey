import { Routes } from '@angular/router';
import { Login } from './components/login/login';
import { Register } from './components/register/register';
import { CreatePlayer } from './components/create-player/create-player';
import { ViewPlayers } from './components/view-players/view-players';
import { CreateTeam } from './components/create-team/create-team';
import { CreateDivision } from './components/create-division/create-division';
import { CreateLeague } from './components/create-league/create-league';
import { ViewPlayer } from './components/view-player/view-player';
import { EditPlayer } from './components/edit-player/edit-player';

export const routes: Routes = [
    {
        path: '',
        component: Login
    },
    {
        path: 'register',
        component: Register
    },
    {
        path: 'create-player',
        component: CreatePlayer
    },
    {
        path: 'view-players',
        component: ViewPlayers
    },
    {
        path: 'create-team',
        component: CreateTeam
    },
    {
        path: 'create-division',
        component: CreateDivision
    },
    {
        path: 'create-league',
        component: CreateLeague
    },
    {
        path: 'view-player/:id',
        component: ViewPlayer
    },
    {
        path: 'edit-player/:id',
        component: EditPlayer
    }
];
