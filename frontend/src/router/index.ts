import {
  createRouter,
  createWebHistory,
  type NavigationGuardNext,
  type RouteLocationNormalized,
} from 'vue-router';
import HomeView from '@/views/HomeView.vue';
import LobbyView from '@/views/LobbyView.vue';
import GameView from '@/views/GameView.vue';
import LoginView from '@/views/LoginView.vue';
import { usePlayerStore } from '@/stores/usePlayerStore';
import { useRoomStore } from '@/stores/useRoomStore';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/auth',
      component: LoginView,
      beforeEnter: redirectIfAuthenticated,
    },
    {
      path: '/',
      component: HomeView,
      beforeEnter: requirePlayer,
    },
    {
      path: '/lobby/:id',
      component: LobbyView,
      beforeEnter: requirePlayer,
    },
    {
      path: '/game/:id',
      component: GameView,
      beforeEnter: [requirePlayer, requireParticipant],
    },
  ],
});

function redirectIfAuthenticated(
  _to: RouteLocationNormalized,
  _from: RouteLocationNormalized,
  next: NavigationGuardNext
) {
  const token = localStorage.getItem('token');
  if (token) return next('/');
  next();
}

function requirePlayer(
  _to: RouteLocationNormalized,
  _from: RouteLocationNormalized,
  next: NavigationGuardNext
) {
  const token = localStorage.getItem('token');
  if (!token) return next('/auth');
  next();
}

async function requireParticipant(
  to: RouteLocationNormalized,
  _from: RouteLocationNormalized,
  next: NavigationGuardNext
) {
  const player = usePlayerStore();
  const room = useRoomStore();
  await room.loadRoom(to.params.id as string);

  const isParticipant =
    room.gameMasterId === player.playerId ||
    room.attackerId === player.playerId ||
    room.defenderId === player.playerId;

  if (!isParticipant) return next(`/lobby/${to.params.id}`);
  if (room.status === 'lobby') return next(`/lobby/${to.params.id}`);
  next();
}

export default router;
