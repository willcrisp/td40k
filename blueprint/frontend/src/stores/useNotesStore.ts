import { defineStore } from "pinia";
import { ref } from "vue";
import { useToast } from "primevue/usetoast";
import type { Note } from "@/types";
import { apiListNotes, apiCreateNote, apiDeleteNote } from "@/lib/api";

export const useNotesStore = defineStore("notes", () => {
  const notes = ref<Note[]>([]);
  const adding = ref(false);
  const toast = useToast();

  async function fetchNotes() {
    const res = await apiListNotes();
    notes.value = res.data;
  }

  async function createNote(content: string) {
    adding.value = true;
    try {
      await apiCreateNote(content);
      // list updated via WebSocket broadcast
    } catch {
      toast.add({
        severity: "error",
        summary: "Error",
        detail: "Failed to add note.",
        life: 3000,
      });
    } finally {
      adding.value = false;
    }
  }

  async function deleteNote(id: string) {
    try {
      await apiDeleteNote(id);
      // list updated via WebSocket broadcast
    } catch {
      toast.add({
        severity: "error",
        summary: "Error",
        detail: "Failed to delete note.",
        life: 3000,
      });
    }
  }

  function applyInsert(note: Note) {
    notes.value.unshift(note);
  }

  function applyDelete(id: string) {
    notes.value = notes.value.filter((n) => n.id !== id);
  }

  return { notes, adding, fetchNotes, createNote, deleteNote, applyInsert, applyDelete };
});
