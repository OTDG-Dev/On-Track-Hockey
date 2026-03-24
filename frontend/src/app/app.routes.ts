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
import { ViewTeams } from './components/view-teams/view-teams';
import { ViewTeam } from './components/view-team/view-team';
import { EditTeam } from './components/edit-team/edit-team';
import { ViewDivisions } from './components/view-divisions/view-divisions';
import { EditDivision } from './components/edit-division/edit-division';
import { ViewDivision } from './components/view-division/view-division';
import { ViewLeagues } from './components/view-leagues/view-leagues';
import { ViewLeague } from './components/view-league/view-league';
import { EditLeague } from './components/edit-league/edit-league';

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
    },
    {
        path: 'view-teams',
        component: ViewTeams
    },
    {
        path: 'view-team/:id',
        component: ViewTeam
    },
    {
        path: 'edit-team/:id',
        component: EditTeam
    },
    {
        path: 'view-divisions',
        component: ViewDivisions
    },
    {
        path: 'view-division/:id',
        component: ViewDivision
    },
    {
        path: 'edit-division/:id',
        component: EditDivision
    },
    {
        path: 'view-leagues',
        component: ViewLeagues
    },
    {
        path: 'view-league/:id',
        component: ViewLeague
    },
    {
        path: 'edit-league/:id',
        component: EditLeague
    }
];
