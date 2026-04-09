<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { usePlayerStore } from "@/stores/usePlayerStore";
import { apiLogin, apiRegister } from "@/lib/api";

const router = useRouter();
const playerStore = usePlayerStore();

const username = ref("");
const password = ref("");
const error = ref("");
const loading = ref(false);

async function handleLogin() {
  error.value = "";
  loading.value = true;
  try {
    const res = await apiLogin(username.value, password.value);
    playerStore.setAuth(res.data);
    router.push("/");
  } catch {
    error.value = "Invalid credentials.";
  } finally {
    loading.value = false;
  }
}

async function handleRegister() {
  error.value = "";
  loading.value = true;
  try {
    const res = await apiRegister(username.value, password.value);
    playerStore.setAuth(res.data);
    router.push("/");
  } catch (e: unknown) {
    const status = (e as { response?: { status?: number } })?.response?.status;
    error.value =
      status === 409 ? "Username already taken." : "Registration failed.";
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <div
    class="flex align-items-center justify-content-center"
    style="min-height: 100vh"
  >
    <div style="width: 360px">
      <Tabs value="login">
        <TabList>
          <Tab value="login">Login</Tab>
          <Tab value="register">Register</Tab>
        </TabList>
        <TabPanels>
          <TabPanel value="login">
            <div class="flex flex-column gap-3 pt-3">
              <InputText
                v-model="username"
                placeholder="Username"
                fluid
              />
              <Password
                v-model="password"
                placeholder="Password"
                :feedback="false"
                fluid
              />
              <Message v-if="error" severity="error">{{ error }}</Message>
              <Button
                label="Login"
                :loading="loading"
                fluid
                @click="handleLogin"
              />
            </div>
          </TabPanel>
          <TabPanel value="register">
            <div class="flex flex-column gap-3 pt-3">
              <InputText
                v-model="username"
                placeholder="Username"
                fluid
              />
              <Password
                v-model="password"
                placeholder="Password"
                :feedback="false"
                fluid
              />
              <Message v-if="error" severity="error">{{ error }}</Message>
              <Button
                label="Register"
                :loading="loading"
                fluid
                @click="handleRegister"
              />
            </div>
          </TabPanel>
        </TabPanels>
      </Tabs>
    </div>
  </div>
</template>
