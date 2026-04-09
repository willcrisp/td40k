<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import { storeToRefs } from "pinia";
import { useRouter } from "vue-router";
import { useCounterStore } from "@/stores/useCounterStore";
import { useNotesStore } from "@/stores/useNotesStore";
import { useWebSocketStore } from "@/stores/useWebSocketStore";
import { useUserStore } from "@/stores/useUserStore";

const router = useRouter();
const counterStore = useCounterStore();
const notesStore = useNotesStore();
const wsStore = useWebSocketStore();
const userStore = useUserStore();

const { value, loading: counterLoading } = storeToRefs(counterStore);
const { notes, adding } = storeToRefs(notesStore);
const { username, userId } = storeToRefs(userStore);

const newNote = ref("");

onMounted(async () => {
  await Promise.all([counterStore.fetchCounter(), notesStore.fetchNotes()]);
  wsStore.connect();
});

onUnmounted(() => {
  wsStore.disconnect();
});

async function handleAddNote() {
  const content = newNote.value.trim();
  if (!content) return;
  await notesStore.createNote(content);
  newNote.value = "";
}

function handleLogout() {
  wsStore.disconnect();
  userStore.logout();
  router.push("/auth");
}
</script>

<template>
  <div class="flex flex-column align-items-center gap-4 p-6">
    <Toolbar style="width: 100%; max-width: 560px">
      <template #start>
        <span class="font-semibold">{{ username }}</span>
      </template>
      <template #end>
        <Button
          label="Logout"
          severity="secondary"
          size="small"
          @click="handleLogout"
        />
      </template>
    </Toolbar>

    <Card style="width: 100%; max-width: 560px">
      <template #title>Shared Counter</template>
      <template #content>
        <div class="flex flex-column align-items-center gap-4">
          <div class="text-8xl font-bold">{{ value }}</div>
          <Button
            label="Increment"
            :loading="counterLoading"
            size="large"
            @click="counterStore.increment()"
          />
        </div>
      </template>
    </Card>

    <Card style="width: 100%; max-width: 560px">
      <template #title>Shared Notes</template>
      <template #content>
        <div class="flex flex-column gap-3">
          <div class="flex gap-2">
            <Textarea
              v-model="newNote"
              rows="2"
              placeholder="Write a note..."
              class="flex-1"
              @keydown.ctrl.enter="handleAddNote"
            />
            <Button
              icon="pi pi-plus"
              :loading="adding"
              :disabled="!newNote.trim()"
              @click="handleAddNote"
            />
          </div>

          <div
            v-for="note in notes"
            :key="note.id"
            class="flex align-items-start gap-2"
          >
            <div class="flex-1">
              <span class="font-semibold text-sm">
                {{ note.user_id === userId ? "You" : note.username }}
              </span>
              <span class="text-sm ml-2">{{ note.content }}</span>
            </div>
            <Button
              v-if="note.user_id === userId"
              icon="pi pi-trash"
              severity="danger"
              size="small"
              text
              @click="notesStore.deleteNote(note.id)"
            />
          </div>

          <p v-if="notes.length === 0" class="text-center text-sm m-0">
            No notes yet.
          </p>
        </div>
      </template>
    </Card>
  </div>
</template>
