<script setup lang="ts">
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import Tabs from 'primevue/tabs';
import Tab from 'primevue/tab';
import TabList from 'primevue/tablist';
import TabPanel from 'primevue/tabpanel';
import TabPanels from 'primevue/tabpanels';
import InputText from 'primevue/inputtext';
import Password from 'primevue/password';
import Button from 'primevue/button';
import { usePlayerStore } from '@/stores/usePlayerStore';
import { apiLogin, apiRegister } from '@/lib/api';

const router = useRouter();
const playerStore = usePlayerStore();

// Login form
const loginUsername = ref('');
const loginPassword = ref('');
const loginError = ref('');
const loginLoading = ref(false);

// Register form
const regUsername = ref('');
const regNickname = ref('');
const regPassword = ref('');
const regError = ref('');
const regLoading = ref(false);

async function handleLogin() {
  loginError.value = '';
  if (!loginUsername.value || !loginPassword.value) return;
  loginLoading.value = true;
  try {
    const res = await apiLogin({
      username: loginUsername.value,
      password: loginPassword.value,
    });
    playerStore.setAuth(res.data);
    router.push('/');
  } catch {
    loginError.value = 'Invalid username or password.';
  } finally {
    loginLoading.value = false;
  }
}

async function handleRegister() {
  regError.value = '';
  if (!regUsername.value || !regNickname.value || !regPassword.value) return;
  regLoading.value = true;
  try {
    const res = await apiRegister({
      username: regUsername.value,
      nickname: regNickname.value,
      password: regPassword.value,
    });
    playerStore.setAuth(res.data);
    router.push('/');
  } catch (err: unknown) {
    const status =
      err && typeof err === 'object' && 'response' in err
        ? (err as { response: { status: number } }).response?.status
        : 0;
    regError.value =
      status === 409
        ? 'Username already taken. Choose another.'
        : 'Registration failed. Please try again.';
  } finally {
    regLoading.value = false;
  }
}
</script>

<template>
  <div class="layout-terminal flex items-center justify-center min-h-screen">
    <div class="w-full max-w-md">
      <div class="mb-8">
        <h1 class="text-3xl font-display text-primary">Game Tracker</h1>
        <p class="text-sm font-mono text-surface-variant">
          Status: Online // Sector 40K
        </p>
      </div>

      <div class="panel-low p-6 riveted">
        <Tabs value="login">
          <TabList>
            <Tab value="login">Login</Tab>
            <Tab value="register">Register</Tab>
          </TabList>

          <TabPanels>
            <!-- Login -->
            <TabPanel value="login">
              <form
                class="flex flex-col gap-4 mt-4"
                @submit.prevent="handleLogin"
              >
                <div class="flex flex-col gap-1">
                  <label class="text-xs font-mono text-surface-variant"
                    >Username</label
                  >
                  <InputText
                    v-model="loginUsername"
                    placeholder="username"
                    class="font-mono"
                    autocomplete="username"
                    autofocus
                  />
                </div>
                <div class="flex flex-col gap-1">
                  <label class="text-xs font-mono text-surface-variant"
                    >Password</label
                  >
                  <Password
                    v-model="loginPassword"
                    placeholder="password"
                    :feedback="false"
                    toggle-mask
                    fluid
                    autocomplete="current-password"
                  />
                </div>

                <p
                  v-if="loginError"
                  class="text-sm font-mono text-red-400"
                >
                  {{ loginError }}
                </p>

                <Button
                  type="submit"
                  label="Login"
                  :loading="loginLoading"
                  class="btn-tactical"
                />
              </form>
            </TabPanel>

            <!-- Register -->
            <TabPanel value="register">
              <form
                class="flex flex-col gap-4 mt-4"
                @submit.prevent="handleRegister"
              >
                <div class="flex flex-col gap-1">
                  <label class="text-xs font-mono text-surface-variant"
                    >Username</label
                  >
                  <InputText
                    v-model="regUsername"
                    placeholder="username"
                    class="font-mono"
                    autocomplete="username"
                  />
                </div>
                <div class="flex flex-col gap-1">
                  <label class="text-xs font-mono text-surface-variant"
                    >Nickname</label
                  >
                  <InputText
                    v-model="regNickname"
                    placeholder="display name"
                    class="font-mono"
                    autocomplete="nickname"
                  />
                </div>
                <div class="flex flex-col gap-1">
                  <label class="text-xs font-mono text-surface-variant"
                    >Password</label
                  >
                  <Password
                    v-model="regPassword"
                    placeholder="password"
                    :feedback="false"
                    toggle-mask
                    fluid
                    autocomplete="new-password"
                  />
                </div>

                <p v-if="regError" class="text-sm font-mono text-red-400">
                  {{ regError }}
                </p>

                <Button
                  type="submit"
                  label="Create Account"
                  :loading="regLoading"
                  class="btn-tactical"
                />
              </form>
            </TabPanel>
          </TabPanels>
        </Tabs>
      </div>
    </div>
  </div>
</template>
