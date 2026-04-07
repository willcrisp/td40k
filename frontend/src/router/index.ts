import {
  createRouter,
  createWebHistory,
  type NavigationGuardNext,
  type RouteLocationNormalized,
} from 'vue-router';
import HomeView from '@/views/HomeView.vue';
import LobbyView from '@/views/LobbyView.vue';
import GameView from '@/views/GameView.vue';
import { usePlayerStore } from '@/stores/usePlayerStore';
import { useRoomStore } from '@/stores/useRoomStore';

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      component: HomeView,
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

function requirePlayer(
  _to: RouteLocationNormalized,
  _from: RouteLocationNormalized,
  next: NavigationGuardNext
) {
  const player = usePlayerStore();
  if (!player.playerId) return next('/');
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
